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

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"t73f.de/r/webs/middleware"
)

// Based on https://gist.github.com/alexedwards/219d88ebdb9c0c9e74715d243f5b2136

func TestList(t *testing.T) {
	used := ""

	mw := slices.Collect(makeMiddleware(6, &used))
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := http.NewServeMux()

	c1 := middleware.NewChain(mw[0], mw[1])
	l1 := middleware.NewListFromChain(c1)
	m.Handle("GET /{$}", middleware.Apply(l1, hf))

	c2 := c1.Append(mw[2], mw[3])
	l2 := middleware.NewListFromChain(c2)
	m.Handle("GET /foo", middleware.Apply(l2, hf))

	l3 := l2.Append(mw[4])
	m.Handle("GET /nested/foo", middleware.Apply(l3, hf))

	l4 := l1.Extend(nil)
	m.Handle("GET /bar", middleware.Apply(l4, hf))

	l5 := l1.Extend(middleware.NewList(mw[2], nil).Append(mw[3]))
	m.Handle("GET /ext", middleware.Apply(l5, hf))

	m.Handle("GET /baz", middleware.Apply(c1, hf))

	var tests = []struct {
		method string
		path   string
		exp    string
		status int
	}{
		{method: "GET", path: "/", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/foo", exp: ";0;1;2;3", status: http.StatusOK},
		{method: "GET", path: "/nested/foo", exp: ";0;1;2;3;4", status: http.StatusOK},
		{method: "GET", path: "/bar", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/baz", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/ext", exp: ";0;1;2;3", status: http.StatusOK},
		{method: "GET", path: "/boo", exp: "", status: http.StatusNotFound},
	}

	for _, test := range tests {
		used = ""

		r, err := http.NewRequest(test.method, test.path, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err)
		}

		rr := httptest.NewRecorder()
		m.ServeHTTP(rr, r)

		got := rr.Result()
		if status := got.StatusCode; status != test.status {
			t.Errorf("%s %s: expected status %d but was %d", test.method, test.path, test.status, status)
		}

		if used != test.exp {
			t.Errorf("%s %s: middleware used: expected %q; got %q", test.method, test.path, test.exp, used)
		}
	}
}
