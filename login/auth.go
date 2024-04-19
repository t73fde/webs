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
)

// NoAuthenticator authenticates nothing.
type NoAuthenticator struct{}

func (*NoAuthenticator) Authenticate(context.Context, string, string) (UserInfo, error) {
	return UserInfo{}, ErrUsernamePassword
}

// TestAuthenticator is an Authenticator for testing purposes.
type TestAuthenticator struct {
	mx    sync.Mutex // protect the following map
	names map[string]UserInfo
}

func (ta *TestAuthenticator) Authenticate(_ context.Context, username, password string) (UserInfo, error) {
	if username[0] == 'x' {
		return UserInfo{}, ErrUsernamePassword
	}
	if username[0] == 'q' && username != password {
		return UserInfo{}, ErrUsernamePassword
	}
	ta.mx.Lock()
	defer ta.mx.Unlock()

	lenNames := len(ta.names)
	if lenNames == 0 {
		userinfo := UserInfo{UserID: 0, Username: username}
		ta.names = map[string]UserInfo{username: userinfo}
		return userinfo, nil
	}
	if userinfo, found := ta.names[username]; found {
		return userinfo, nil
	}
	if lenNames > 1024 {
		return UserInfo{}, ErrTooManyUsers
	}
	userinfo := UserInfo{UserID: int64(lenNames), Username: username}
	return userinfo, nil
}
