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

package check_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	zerocontext "t73f.de/r/zero/context"

	"t73f.de/r/webs/middleware/check"
)

type ctxType string

const ctxKey = ctxType("key")
const ctxVal = "123"

func TestChecker(t *testing.T) {
	used := ""
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		used, _ = getCtx(r.Context())
		w.WriteHeader(http.StatusNoContent)
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", check.Build(check.Func(checkFalse))(http.HandlerFunc(hf)))
	mux.Handle("/bar", check.Build(check.Func(checkTrue))(http.HandlerFunc(hf)))
	mux.Handle("/baz", check.Build(check.Func(checkTrueCtx))(http.HandlerFunc(hf)))

	r, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, r)
	res := rr.Result()
	if used != "" {
		t.Errorf("handler func was executed: %q", used)
	}
	if code := res.StatusCode; code != expErrCode {
		t.Errorf("status code %d expected, got: %d", expErrCode, code)
	}

	used = ""
	r, err = http.NewRequest("GET", "/bar", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, r)
	res = rr.Result()
	if used != "" {
		t.Errorf("context was modified: %q", used)
	}
	if code := res.StatusCode; code != expOKCode {
		t.Errorf("status code %d expected, got: %d", expOKCode, code)
	}

	used = ""
	r, err = http.NewRequest("GET", "/baz", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, r)
	res = rr.Result()
	if used != ctxVal {
		t.Errorf("context wrongly set, exp: %q, got: %q", ctxVal, used)
	}
	if code := res.StatusCode; code != expOKCode {
		t.Errorf("status code %d expected, got: %d", expOKCode, code)
	}
}

const expErrCode = http.StatusBadRequest
const expOKCode = http.StatusNoContent

func checkFalse(w http.ResponseWriter, _ *http.Request) (context.Context, bool) {
	w.WriteHeader(expErrCode)
	return nil, false
}
func checkTrue(_ http.ResponseWriter, r *http.Request) (context.Context, bool) {
	return r.Context(), true
}
func checkTrueCtx(_ http.ResponseWriter, r *http.Request) (context.Context, bool) {
	return withCtx(r.Context(), ctxVal), true
}

var withCtx, getCtx = zerocontext.WithAndValue[string](ctxKey)
