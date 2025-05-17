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
	"fmt"
	"iter"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"t73f.de/r/webs/middleware"
)

// Based on https://gist.github.com/alexedwards/219d88ebdb9c0c9e74715d243f5b2136

func TestMiddleware(t *testing.T) {
	used := ""

	mw := slices.Collect(makeMiddleware(3, &used))
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := http.NewServeMux()

	m.Handle("GET /{$}", middleware.Apply(mw[0], hf))
	m.Handle("GET /foo", middleware.Apply(mw[1], hf))
	m.Handle("GET /baz", middleware.Apply(mw[2], hf))

	var tests = []struct {
		method string
		path   string
		exp    string
		status int
	}{
		{method: "GET", path: "/", exp: ";0", status: http.StatusOK},
		{method: "GET", path: "/foo", exp: ";1", status: http.StatusOK},
		{method: "GET", path: "/baz", exp: ";2", status: http.StatusOK},
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

func makeMiddleware(n int, used *string) iter.Seq[middleware.Middleware] {
	return func(yield func(middleware.Middleware) bool) {
		for i := range n {
			m := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					*used += fmt.Sprintf(";%d", i)
					next.ServeHTTP(w, r)
				})
			}
			if !yield(m) {
				return
			}
		}
	}
}
