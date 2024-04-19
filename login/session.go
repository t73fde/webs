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

package login

import (
	"context"
	"sync"
	"time"
)

// RAMSessions is a SessionManager that stores its sessions in main memory.
type RAMSessions struct {
	mx       sync.Mutex // procted following map
	sessions map[string]*sessionData
}
type sessionData struct {
	UserInfo
	expires time.Time
}

func (rs *RAMSessions) SetUserAuth(_ context.Context, userinfo UserInfo, auth string) error {
	session := sessionData{
		UserInfo: userinfo,
		expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	rs.mx.Lock()
	numSessions := len(rs.sessions)
	if numSessions == 0 {
		rs.sessions = map[string]*sessionData{auth: &session}
	} else if numSessions > 1024 {
		return ErrTooManySessions
	} else {
		rs.sessions[auth] = &session
	}
	rs.mx.Unlock()
	return nil
}

func (rs *RAMSessions) UserAuth(_ context.Context, auth string) (UserInfo, error) {
	var session *sessionData
	var found bool
	rs.mx.Lock()
	defer rs.mx.Unlock()

	if len(rs.sessions) > 0 {
		session, found = rs.sessions[auth]
	}
	if !found || session == nil {
		return UserInfo{}, ErrNoSuchSession
	}
	now := time.Now()
	if now.After(session.expires) {
		delete(rs.sessions, auth)
		return UserInfo{}, ErrNoSuchSession
	}
	if session.expires.Before(now.Add(3 * 24 * time.Hour)) {
		session.expires = now.Add(7 * 24 * time.Hour)
	}
	return session.UserInfo, nil
}

func (rs *RAMSessions) Remove(_ context.Context, auth string) error {
	rs.mx.Lock()
	if len(rs.sessions) > 0 {
		delete(rs.sessions, auth)
	}
	rs.mx.Unlock()
	return nil
}
