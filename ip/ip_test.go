//-----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
//-----------------------------------------------------------------------------

package ip_test

import (
	"testing"

	"t73f.de/r/webs/ip"
)

func TestIsLoopbackAddr(t *testing.T) {
	localAddr := []string{
		"127.0.0.1", "127.0.0.1:80",
		"localhost", "localhost:80", "localhost%zone:80",
		"[::1]:23",
	}
	for _, addr := range localAddr {
		if !ip.IsLoopbackAddr(addr) {
			t.Errorf("%q should be a loopback addr", addr)
		}
	}
	remoteAddr := []string{
		"", ":80", "%zone:80",
		"::1", "::1:80",
	}
	for _, addr := range remoteAddr {
		if ip.IsLoopbackAddr(addr) {
			t.Errorf("%q should be a non-loopback addr", addr)
		}
	}
}
