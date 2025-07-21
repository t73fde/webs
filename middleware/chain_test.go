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

func TestChain(t *testing.T) {
	used := ""

	fts := slices.Collect(makeFunctors(6, &used))
	hf := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m := http.NewServeMux()

	c1 := middleware.NewChain(fts[0], fts[1])
	m.Handle("GET /{$}", middleware.Apply(c1, hf))

	c2 := c1.Append(fts[2], fts[3])
	m.Handle("GET /foo", middleware.Apply(c2, hf))

	c3 := c2.Append(fts[4])
	m.Handle("GET /nested/foo", middleware.Apply(c3, hf))

	c4 := c1.Extend(middleware.NewChain(fts[5]))
	m.Handle("GET /bar", middleware.Apply(c4, hf))

	c5 := middleware.NewChainFromMiddleware(
		middleware.NewList(fts[1], middleware.NewList(fts[0], nil)))
	m.Handle("GET /lst", middleware.Apply(c5, hf))

	c6 := middleware.NewChainFromMiddleware(c3)
	m.Handle("GET /chn", middleware.Apply(c6, hf))

	m.Handle("GET /baz", middleware.Apply(c1, hf))

	var tests = Testcases{
		{method: "GET", path: "/", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/foo", exp: ";0;1;2;3", status: http.StatusOK},
		{method: "GET", path: "/nested/foo", exp: ";0;1;2;3;4", status: http.StatusOK},
		{method: "GET", path: "/bar", exp: ";0;1;5", status: http.StatusOK},
		{method: "GET", path: "/baz", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/boo", exp: "", status: http.StatusNotFound},
		{method: "GET", path: "/lst", exp: ";0;1", status: http.StatusOK},
		{method: "GET", path: "/chn", exp: ";0;1;2;3;4", status: http.StatusOK},
	}
	tests.Run(t, &used, m)
}

func TestChainFunctors(t *testing.T) {
	var used string
	fts := slices.Collect(makeFunctors(2, &used))
	l := middleware.NewChain(fts...)
	var val middleware.Functor
	for f := range l.Functors() {
		val = f
		break
	}

	hf := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	m := http.NewServeMux()
	m.Handle("GET /foo", val(hf))
	m.Handle("GET /bar", fts[1](hf))

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
