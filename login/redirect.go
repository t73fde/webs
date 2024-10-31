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

import "net/http"

// SimpleRedirector provides some static URLs.
type SimpleRedirector struct {
	LoginURL   string
	SuccessURL string
	ErrorURL   string
	LogoutURL  string
}

// LoginRedirect performs a redirection if user must authenticate itself.
func (sr *SimpleRedirector) LoginRedirect(w http.ResponseWriter, r *http.Request) {
	if sr.LoginURL == "" {
		sr.LoginURL = "/login/"
	}
	http.Redirect(w, r, sr.LoginURL, http.StatusSeeOther)
}

// SuccessRedirect performs a redirection after the user was successfully authenticated.
func (sr *SimpleRedirector) SuccessRedirect(w http.ResponseWriter, r *http.Request, _ UserInfo) {
	if sr.SuccessURL == "" {
		sr.SuccessURL = "/"
	}
	http.Redirect(w, r, sr.SuccessURL, http.StatusSeeOther)
}

// ErrorRedirect performs a redirection if user was not authenticated during login.
func (sr *SimpleRedirector) ErrorRedirect(w http.ResponseWriter, r *http.Request, _ UserInfo, _ error) {
	if sr.ErrorURL == "" {
		sr.ErrorURL = "/"
	}
	http.Redirect(w, r, sr.ErrorURL, http.StatusSeeOther)
}

// LogoutRedirect performs a rediration when user logs out.
func (sr *SimpleRedirector) LogoutRedirect(w http.ResponseWriter, r *http.Request) {
	if sr.LogoutURL == "" {
		sr.LogoutURL = "/"
	}
	http.Redirect(w, r, sr.LogoutURL, http.StatusSeeOther)
}
