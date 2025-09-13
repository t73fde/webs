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

package status_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"t73f.de/r/webs/middleware/status"
)

func TestStatusBuilder(t *testing.T) {
	const headerKey = "X-Handler"
	data200, data404 := "", ""
	hf200 := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data200 = "ok"
		w.Header().Set(headerKey, "OK")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Okay"))
	})
	hf404 := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data404 = "set"
		w.Header().Set(headerKey, "Not Found")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Missing"))
	})
	cfg := status.Config{}
	mux := http.NewServeMux()

	mux.Handle("GET /200/base", cfg.Build()(hf200))
	check200(t, mux, "base", &data200)

	cfg.HandlerMap = status.HandlerMap{
		http.StatusNotFound: http.RedirectHandler(
			"/foo", http.StatusTemporaryRedirect)}
	mux.Handle("GET /200/set", cfg.Build()(hf200))
	check200(t, mux, "set", &data200)

	mux.Handle("GET /404", cfg.Build()(hf404))
	r, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
		return
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, r)
	if got := rr.Code; got != http.StatusTemporaryRedirect {
		t.Errorf("code %d expected, got: %d", http.StatusTemporaryRedirect, got)
	}
	if data404 != "set" {
		t.Error("handler not called")
	}
	data404 = ""
	h := rr.Header()
	if h.Get("Location") != "/foo" {
		t.Errorf("header should contain location '/foo', but got: %v", h)
	}
	if got := h.Get(headerKey); got != "" {
		t.Errorf("%s must not be set, got: %q", headerKey, got)
	}

	cfg.NoClearMap = map[int]bool{404: true}
	mux.Handle("GET /200/nc", cfg.Build()(hf200))
	check200(t, mux, "nc", &data200)

	mux.Handle("GET /404c", cfg.Build()(hf404))
	r, err = http.NewRequest("GET", "/404c", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
		return
	}
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, r)
	if got := rr.Code; got != http.StatusTemporaryRedirect {
		t.Errorf("code %d expected, got: %d", http.StatusTemporaryRedirect, got)
	}
	if data404 != "set" {
		t.Error("handler not called")
	}
	data404 = ""
	h = rr.Header()
	if h.Get("Location") != "/foo" {
		t.Errorf("header should contain location '/foo', but got: %v", h)
	}
	if got := h.Get(headerKey); got == "" {
		t.Errorf("%s must be set, got: %v", headerKey, h)
	}
}

func check200(t *testing.T, mux *http.ServeMux, name string, data200 *string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		r, err := http.NewRequest("GET", "/200/"+name, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err)
			return
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, r)
		if *data200 != "ok" {
			t.Error("handler not called")
		}
		*data200 = ""
		if got := rr.Code; got != 200 {
			t.Errorf("code 200 expected, but got: %d", got)
		}
		if h := rr.Header(); len(h) != 1 || h.Get("X-Handler") != "OK" {
			t.Errorf("only X-Handler:OK expected, but got %v", h)
		}
	})
}

func TestBaseRedirectHandler(t *testing.T) {
	const newSite = "https://2017.4042307.org"
	const code = http.StatusTemporaryRedirect
	const prefix = "/foo"
	const sitePrefix = newSite + prefix
	mux := http.NewServeMux()
	mux.Handle("GET "+prefix, status.BaseRedirectHandler(newSite, code))

	testcases := []struct {
		name string
		path string
		exp  string
	}{
		{"root", "", sitePrefix},
		{"q1", "?a", sitePrefix + "?a"},
		{"q1a", "?a=1", sitePrefix + "?a=1"},
		{"q2", "?a&b", sitePrefix + "?a&b"},
		{"q2a", "?b=1&a=2", sitePrefix + "?b=1&a=2"},
		{"qs", "?a=1&a=2", sitePrefix + "?a=1&a=2"},
		{"f1", "#frag", sitePrefix + "#frag"},
		{"qf", "?s=404#frag", sitePrefix + "?s=404#frag"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", prefix+tc.path, nil)
			if err != nil {
				t.Errorf("NewRequest: %s", err)
				return
			}
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, r)

			if got := rr.Code; got != code {
				t.Errorf("status code %d expected, got: %d", code, got)
			}
			h := rr.Header()
			if got := h.Get("Location"); got != tc.exp {
				t.Errorf("Location %q expected, but got: %q", tc.exp, got)
			}
		})
	}
}
