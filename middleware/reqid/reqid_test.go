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
	rqid := ""
	var reqidcfg reqid.Config
	reqidcfg.Generator = snow.New(5)
	reqidcfg.Initialize()
	reqidcfg.Generator = nil

	rmw := reqidcfg.Build()
	handler := func(w http.ResponseWriter, r *http.Request) {
		rqid = r.Header.Get(reqidcfg.HeaderKey)
	}
	mux := http.NewServeMux()
	mux.Handle("/foo", rmw(http.HandlerFunc(handler)))

	r, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}

	for range 10 {
		rqid = ""
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, r)
		res := rr.Result()
		if got := res.Header.Get(reqidcfg.HeaderKey); rqid != got {
			t.Errorf("request IDs differ: exp: %q, got: %q", rqid, got)
		}
	}
}
