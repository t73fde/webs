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
	return nil, ErrUsernamePassword
}

// TestAuthenticator is an Authenticator for testing purposes.
type TestAuthenticator struct {
	mx    sync.Mutex // protect the following map
	names map[string]testUserInfo
}

type testUserInfo string

func (u testUserInfo) Name() string { return string(u) }

func (ta *TestAuthenticator) Authenticate(_ context.Context, username, password string) (UserInfo, error) {
	if username[0] == 'x' {
		return nil, ErrUsernamePassword
	}
	if username[0] == 'q' && username != password {
		return nil, ErrUsernamePassword
	}
	ta.mx.Lock()
	defer ta.mx.Unlock()

	lenNames := len(ta.names)
	if lenNames == 0 {
		userinfo := testUserInfo(username)
		ta.names = map[string]testUserInfo{username: userinfo}
		return userinfo, nil
	}
	if userinfo, found := ta.names[username]; found {
		return userinfo, nil
	}
	if lenNames > 1024 {
		return nil, ErrTooManyUsers
	}
	userinfo := testUserInfo(username)
	return userinfo, nil
}
