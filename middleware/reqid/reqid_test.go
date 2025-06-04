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

package reqid_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"t73f.de/r/webs/middleware/reqid"
	"t73f.de/r/zero/snow"
)

func TestSimpleReqID(t *testing.T) {
	rqid, rqidctx := "", snow.Invalid
	var reqidcfg reqid.Config
	reqidcfg.WithResponse = true
	reqidcfg.AppID = 65535

	rmw := reqidcfg.Build()
	handler := func(w http.ResponseWriter, r *http.Request) {
		rqid = r.Header.Get(reqid.DefaultHeaderKey)
		rqidctx = reqid.GetRequestID(r.Context())
	}
	mux := http.NewServeMux()
	mux.Handle("/foo", rmw(http.HandlerFunc(handler)))

	r, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	for range 10 {
		rqid, rqidctx = "", snow.Invalid
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, r)
		res := rr.Result()
		if rqid == "" {
			t.Error("no header set")
			break
		}
		if rqidctx.IsValid() {
			t.Error("context id wrongly set")
		}
		if got := res.Header.Get(reqid.DefaultHeaderKey); rqid != got {
			t.Errorf("request IDs differ: exp: %q, got: %q", rqid, got)
			break
		}
	}

	reqidcfg.WithResponse = false
	rmw = reqidcfg.Build()
	mux.Handle("/bar", rmw(http.HandlerFunc(handler)))
	r, err = http.NewRequest("GET", "/bar", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	for range 10 {
		rqid = ""
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, r)
		res := rr.Result()
		if rqid == "" {
			t.Error("no header set")
			break
		}
		if got := res.Header.Get(reqid.DefaultHeaderKey); got != "" {
			t.Errorf("no response key expected, got: %q", got)
			break
		}
	}

	reqidcfg.WithContext = true
	rmw = reqidcfg.Build()
	mux.Handle("/baz", rmw(http.HandlerFunc(handler)))
	r, err = http.NewRequest("GET", "/baz", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	for range 10 {
		rqid, rqidctx = "", snow.Invalid
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, r)
		res := rr.Result()
		if rqid == "" {
			t.Error("no header set")
			break
		}
		if got := res.Header.Get(reqid.DefaultHeaderKey); got != "" {
			t.Errorf("no response key expected, got: %q", got)
			break
		}
		if rqid != rqidctx.String() {
			t.Errorf("header and context id differ: %v vs %v", rqid, rqidctx)
		}
	}
}
