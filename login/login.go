//-----------------------------------------------------------------------------
// Copyright (c) 2024-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2024-present Detlef Stern
//-----------------------------------------------------------------------------

// Package login provides a mechanism for simple login / logout use cases of web sites.
package login

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Provider is an object that handles everything w.r.t authentication.
// It is the main element to log in / log out.
type Provider struct {
	logger *slog.Logger
	auth   Authenticator
	sess   SessionManager
	redir  Redirector

	passlen      int // max length of username and password
	authlen      int // max length of cookie value
	cookiePath   string
	maxCookieAge int
	secureCookie bool

	UsernameKey string
	PasswordKey string
	cookieName  string

	mxAuthProgress sync.Mutex
	authProgress   map[string]struct{}
	authWait       time.Duration
}

// MakeProvider make a new authenticator. Typically, you only need one
// authenticator for an application.
func MakeProvider(logger *slog.Logger, auth Authenticator, sess SessionManager, redir Redirector) *Provider {
	provider := Provider{
		logger: logger,
		auth:   auth,
		sess:   sess,
		redir:  redir,

		passlen:      127,
		authlen:      32,
		cookiePath:   "/", // TODO: should be set-able
		maxCookieAge: 366 * 24 * 3600,
		secureCookie: false,

		UsernameKey: "username",
		PasswordKey: "password",
		cookieName:  "auth",

		authProgress: map[string]struct{}{},
		authWait:     2 * time.Second, // wait time for multiple logins
	}
	return &provider
}

// Authenticator allows to authenticate a human user.
type Authenticator interface {
	// Authenticate with the given user name and password, giving some data
	// about the user.
	Authenticate(ctx context.Context, username, password string) (UserInfo, error)
}

var ErrUsernamePassword = errors.New("username and password do not match")
var ErrTooManyUsers = errors.New("too many users")

// UserInfo gives some information about a user, w.r.t. authentication.
// Other data must be handled separately.
type UserInfo interface {
	// Unique user name
	Name() string
}

type (
	SessionInfo struct {
		SessionID SessionID
		User      UserInfo
	}
	SessionID string
)

// SessionManager handles the set of logged-in users.
type SessionManager interface {
	// Associate an user info with a session identifier.
	SetUserAuth(context.Context, UserInfo, SessionID) error

	// Retrieve the user info based on the session identifier.
	UserAuth(context.Context, SessionID) (UserInfo, error)

	// Remove session. May remove all sessions of the associated user.
	Remove(context.Context, SessionID) error
}

var ErrNoSuchSession = errors.New("no such session")
var ErrTooManySessions = errors.New("too many open sessions")

// Redirector will redirect the user to an appropriate URL.
type Redirector interface {
	// Redirect to login page.
	LoginRedirect(http.ResponseWriter, *http.Request)

	// Redirect after a successful login.
	SuccessRedirect(http.ResponseWriter, *http.Request, UserInfo)

	// Redirect after a login with errors.
	ErrorRedirect(http.ResponseWriter, *http.Request, UserInfo, error)

	// Redirect after logout.
	LogoutRedirect(http.ResponseWriter, *http.Request)
}

// Login creates a handler to implement a POST request from the login web page.
func (lp *Provider) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := strings.TrimSpace(r.FormValue(lp.UsernameKey))
		password := strings.TrimSpace(r.FormValue(lp.PasswordKey))

		if l := lp.passlen; username == "" || len(username) > l || password == "" || len(password) > l {
			lp.logger.Info("invalid password attempt")
			lp.loginRedirect(w, r)
			return
		}

		ctx := r.Context()
		if !lp.rateAndWait(username) {
			lp.logger.InfoContext(ctx, "login rated", "username", username)
			lp.loginRedirect(w, r)
			return
		}

		userinfo, err := lp.auth.Authenticate(ctx, username, password)
		if err != nil {
			lp.logger.InfoContext(ctx, "login failed", "error", err)
			lp.loginRedirect(w, r)
			return
		}

		lp.LoginUser(w, r, userinfo)
	})
}
func (lp *Provider) rateAndWait(username string) bool {
	lp.mxAuthProgress.Lock()
	defer lp.mxAuthProgress.Unlock()
	if _, found := lp.authProgress[username]; found {
		return false
	}
	lp.authProgress[username] = struct{}{}
	go func(name string) {
		time.Sleep(lp.authWait)
		lp.mxAuthProgress.Lock()
		delete(lp.authProgress, name)
		lp.mxAuthProgress.Unlock()
	}(username)

	return true
}
func (lp *Provider) loginRedirect(w http.ResponseWriter, r *http.Request) {
	lp.clearAuthCookie(w)
	lp.redir.LoginRedirect(w, r)
}

// LoginUser performs the login session handling for an already authenticated user.
func (lp *Provider) LoginUser(w http.ResponseWriter, r *http.Request, userinfo UserInfo) {
	ctx := r.Context()

	hasher := sha512.New512_256()
	_, _ = io.CopyN(hasher, rand.Reader, 32)
	auth := lp.asHex(hasher)
	lp.setAuthCookie(w, auth)

	hasher.Reset()
	hasher.Write([]byte(auth))
	sessid := SessionID(lp.asHex(hasher))
	if err := lp.sess.SetUserAuth(ctx, userinfo, sessid); err != nil {
		lp.logger.Error("set-session", "error", err)
		lp.redir.ErrorRedirect(w, r, userinfo, err)
		return
	}
	lp.logger.Info("Login", "user", userinfo.Name())
	r = r.WithContext(setSession(ctx, &SessionInfo{SessionID: sessid, User: userinfo}))
	lp.redir.SuccessRedirect(w, r, userinfo)
}

// Logout creates a handler that implements a logout.
func (lp *Provider) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userinfo, auth, err := lp.checkCookie(r)
		if err != nil {
			lp.logger.Info("invalid cookie", "error", err)
		} else {
			ctx := r.Context()
			err = lp.sess.Remove(ctx, auth)
			if err != nil {
				lp.logger.Error("unable to remove auth", "error", err)
			}
			lp.logger.Info("Logout", "user", userinfo.Name())
		}
		lp.clearAuthCookie(w)
		lp.redir.LogoutRedirect(w, r)
	})
}

type sessionKeytype struct{}

var sessionKey sessionKeytype

func setSession(ctx context.Context, session *SessionInfo) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// EnrichUserInfo is a middleware that retrieves the user info based on the
// cookie and stores it in the request context.
//
// Function User() will provide the actual user info for handlers.
func (lp *Provider) EnrichUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userinfo, sessid, err := lp.checkCookie(r); err == nil {
			ctx := setSession(r.Context(), &SessionInfo{SessionID: sessid, User: userinfo})
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// Required ensures a logged-in user. Otherwise the anonymous user is
// redirected to the login page.
//
// Required implies EnrichUserInfo, i.e. there is no need to wrap a handler
// function with EnrichUserInfo.
//
// Function User() can be used to retrieve the actual user inside a handler.
func (lp *Provider) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userinfo, sessid, err := lp.checkCookie(r); err == nil {
			ctx := setSession(r.Context(), &SessionInfo{SessionID: sessid, User: userinfo})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			lp.loginRedirect(w, r)
		}
	})
}

// Session returns a reference to the current user session, or nil if there is
// no session.
func Session(ctx context.Context) *SessionInfo {
	if session, ok := ctx.Value(sessionKey).(*SessionInfo); ok {
		return session
	}
	return nil
}

var errInvalidCookie = errors.New("invalid cookie")

func (lp *Provider) checkCookie(r *http.Request) (UserInfo, SessionID, error) {
	cookie := lp.getAuthCookie(r)
	if cookie == "" {
		return nil, "", errInvalidCookie
	}
	hasher := sha512.New512_256()
	hasher.Write([]byte(cookie))
	auth := SessionID(lp.asHex(hasher))
	ctx := r.Context()
	userinfo, err := lp.sess.UserAuth(ctx, auth)
	return userinfo, auth, err
}

func (lp *Provider) getAuthCookie(r *http.Request) string {
	cookie, err := r.Cookie(lp.cookieName)
	if err != nil {
		return ""
	}
	auth := cookie.Value
	if len(auth) != lp.authlen {
		lp.logger.Info("bad authentication", "auth", auth)
		return ""
	}
	return auth
}

func (lp *Provider) setAuthCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     lp.cookieName,
		Value:    value,
		Path:     lp.cookiePath,
		MaxAge:   lp.maxCookieAge,
		Secure:   lp.secureCookie,
		HttpOnly: true, // TODO: "false" possibly needed for htmx
		SameSite: http.SameSiteLaxMode,
	})
}

func (lp *Provider) clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     lp.cookieName,
		Value:    "",
		Path:     lp.cookiePath,
		MaxAge:   -1,
		Secure:   lp.secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (lp *Provider) asHex(hasher hash.Hash) string {
	return fmt.Sprintf("%x", hasher.Sum(nil))[0:lp.authlen]
}
