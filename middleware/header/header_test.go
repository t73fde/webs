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

package header_test

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"t73f.de/r/webs/middleware/header"
)

func TestHeader(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mux := http.NewServeMux()

	var cnst map[string]string
	tests := []struct {
		path string
		inp  map[string]string
		exp  http.Header
	}{
		{"/foo", cnst, http.Header{}},
		{"/bar", map[string]string{}, http.Header{}},
		{"/baz", map[string]string{"server": "DAS"}, http.Header{"Server": {"DAS"}}},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			cfg := header.Config{Constants: tc.inp}
			mux.Handle("GET "+tc.path, cfg.Build()(hf))

			r, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Errorf("NewRequest: %s", err)
			}
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, r)

			if got := rr.Header(); !maps.EqualFunc(tc.exp, got, slices.Equal) {
				t.Errorf("expected: %v, but got %v", tc.exp, got)
			}
		})
	}
}
