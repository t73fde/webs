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

// Package ip provides methods to work with IP addresses in a web server context.
package ip

import (
	"net"
	"net/http"
	"strings"
)

// GetRemoteAddr returns the network address that sent the request.
// The format is "IP:port".
//
// In contrast to `http.Request.RemoteAddr` it works with intermediate server,
// i.e. it returns the network address of the first client, not the address
// from an intermediate server that acts as a client.
func GetRemoteAddr(r *http.Request) string {
	if r == nil {
		return ""
	}
	if from := r.Header.Get("X-Forwarded-For"); from != "" {
		return from
	}
	return r.RemoteAddr
}

// IsLoopbackAddr returns true if the address is an address of the local computer.
// In most cases, it is the client address, if the loopback network interface is used.
// Invalid addresses are treated as non-loopback / remote addresses.
func IsLoopbackAddr(addr string) bool {
	var host string
	if !strings.ContainsRune(addr, ':') {
		addr = addr + ":80"
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if pos := strings.IndexRune(host, '%'); pos >= 0 {
		host = host[0:pos]
	}
	ip := net.ParseIP(host)
	if ip.IsLoopback() {
		return true
	}
	return strings.ToLower(host) == "localhost"
}
