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

// Package flash allow to display flash messages on web sites.
package flash

import (
	"context"
	"sync"
	"time"

	"t73f.de/r/webs/login"
)

// Flasher allows to set key-based flash messages, and to retrieve them.
type Flasher interface {
	// Add a flash message with the given key. If the key is the empty string,
	// it is assumed a "global" message.
	//
	// Add can be called multiple times, if multiple messages must be displayed.
	Add(ctx context.Context, key string, message string)

	// Messages returns all messages as a map.
	//
	// A second call will return a nil value, i.e. messages are removed.
	Messages(context.Context) map[string][]string
}

type memoryFlasher struct {
	mx       sync.Mutex
	sessions map[login.SessionID]*memMessages
}
type memMessages struct {
	messages map[string][]string
	expiry   time.Time
}

func MakeMemoryFlasher() Flasher {
	return &memoryFlasher{sessions: make(map[login.SessionID]*memMessages, 128)}
}

func (mf *memoryFlasher) Add(ctx context.Context, key, message string) {
	session := login.Session(ctx)
	if session == nil {
		return
	}
	sessid := session.SessionID
	if sessid == "" {
		return
	}
	now := time.Now()
	expiry := now.Add(5 * time.Second)
	mf.mx.Lock()
	defer mf.mx.Unlock()
	sessions := mf.sessions
	if sess, hasSession := sessions[sessid]; hasSession {
		sess.messages[key] = append(sess.messages[key], message)
		sess.expiry = expiry
		return
	}

	sessions[sessid] = &memMessages{
		messages: map[string][]string{key: {message}},
		expiry:   expiry,
	}

	// Check other sessions for outdates messages.
	for sessid, sessMsgs := range sessions {
		if sessMsgs.expiry.Before(now) {
			delete(sessions, sessid)
		}
	}
}

func (mf *memoryFlasher) Messages(ctx context.Context) map[string][]string {
	session := login.Session(ctx)
	if session == nil {
		return nil
	}
	sessid := session.SessionID
	if sessid == "" {
		return nil
	}
	mf.mx.Lock()
	defer mf.mx.Unlock()

	sessions := mf.sessions
	if sess, hasSession := sessions[sessid]; hasSession {
		delete(sessions, sessid)
		if sess.expiry.Before(time.Now()) {
			return nil
		}
		return sess.messages
	}
	return nil
}
