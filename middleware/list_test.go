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

	fts := slices.Collect(makeFunctors(6, &used))
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := http.NewServeMux()

	c1 := middleware.NewChain(fts[0], fts[1])
	l1 := middleware.NewListFromMiddleware(c1)
	m.Handle("GET /{$}", middleware.Apply(l1, hf))

	c2 := c1.Append(fts[2], fts[3])
	l2 := middleware.NewListFromMiddleware(c2)
	m.Handle("GET /foo", middleware.Apply(l2, hf))

	l3 := l2.Append(fts[4])
	m.Handle("GET /nested/foo", middleware.Apply(l3, hf))

	l4 := l1.Extend(nil)
	m.Handle("GET /bar", middleware.Apply(l4, hf))

	l5 := l1.Extend(middleware.NewList(fts[2], nil).Append(fts[3]))
	m.Handle("GET /ext", middleware.Apply(l5, hf))

	l6 := middleware.NewListFromMiddleware(l5)
	m.Handle("GET /lst", middleware.Apply(l6, hf))

	if l5 != l6 {
		t.Error("NewListFromMiddleware is not idempotent")
	}

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
		{method: "GET", path: "/lst", exp: ";0;1;2;3", status: http.StatusOK},
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

func TestListFunctors(t *testing.T) {
	var used string
	fts := slices.Collect(makeFunctors(2, &used))
	l := middleware.NewList(fts[0], middleware.NewList(fts[1], nil))
	var val middleware.Functor
	for f := range l.Functors() {
		val = f
		break
	}

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	m := http.NewServeMux()
	m.Handle("GET /foo", val(hf))
	m.Handle("GET /bar", fts[0](hf))

	used = ""
	r, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)
	valUsed := used

	used = ""
	r, err = http.NewRequest("GET", "/bar", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr = httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if used != valUsed {
		t.Errorf("%q expected, but got %v", used, valUsed)
	}
}
