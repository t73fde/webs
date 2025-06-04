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

package logging_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"t73f.de/r/webs/middleware/logging"
)

func TestRequestLogging(t *testing.T) {
	logh := testLoggingHandler{}
	logger := slog.New(&logh)

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mux := http.NewServeMux()

	tests := testcases{
		{"/foo", nil, false, false, "abc", nil},
		{"/bar", logger, false, false, "REQ", []string{"method", "GET", "url", "/bar"}},
		{"/baz", logger, false, true, "REQ", []string{"method", "GET", "url", "/baz", "header", "map[]"}},
		{"/rar", logger, true, false, "REQ", []string{"method", "GET", "url", "/rar", "remote", "127.0.0.1:54321"}},
		{"/raz", logger, true, true, "REQ", []string{"method", "GET", "url", "/raz", "remote", "127.0.0.1:54321", "header", "map[]"}},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			logh.records = nil
			cfg := logging.ReqConfig{Logger: tc.logger, WithRemote: tc.withRemote, WithHeaders: tc.withHeader}
			mux.Handle("GET "+tc.path, cfg.Build()(hf))

			r, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Errorf("NewRequest: %s", err)
			}
			r.RemoteAddr = "127.0.0.1:54321"
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, r)

			if tc.logger == nil {
				if got := len(logh.records); got != 0 {
					t.Errorf("expected no log record, got %d", got)
				}
				return
			}
			if got := len(logh.records); got != 1 {
				t.Errorf("expected one log record, got %d", got)
				return
			}
			rec := logh.records[0]
			if got := rec.Message; got != tc.expMsg {
				t.Errorf("message %q expected, got: %q", tc.expMsg, got)
			}
			attrs := []string{}
			rec.Attrs(func(a slog.Attr) bool {
				if !a.Equal(slog.Attr{}) {
					attrs = append(attrs, a.Key, a.Value.String())
				}
				return true
			})
			if !slices.Equal(tc.expAttrs, attrs) {
				t.Errorf("attrs %v expected, got %v", tc.expAttrs, attrs)
			}
		})
	}
}

func TestResponseLogging(t *testing.T) {
	logh := testLoggingHandler{}
	logger := slog.New(&logh)

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = io.WriteString(w, "Hello")
	})
	mux := http.NewServeMux()

	tests := testcases{
		{"/foo", nil, false, false, "abc", nil},
		{"/bar", logger, false, false, "RSP", []string{
			"method", "GET", "url", "/bar", "status", "200", "length", "5"}},
		{"/baz", logger, false, true, "RSP", []string{
			"method", "GET", "url", "/baz", "status", "200", "length", "5", "header", "map[]"}},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			logh.records = nil
			cfg := logging.RespConfig{Logger: tc.logger, WithHeaders: tc.withHeader}
			mux.Handle("GET "+tc.path, cfg.Build()(hf))

			r, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Errorf("NewRequest: %s", err)
			}
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, r)

			if tc.logger == nil {
				if got := len(logh.records); got != 0 {
					t.Errorf("expected no log record, got %d", got)
				}
				return
			}
			if got := len(logh.records); got != 1 {
				t.Errorf("expected one log record, got %d", got)
				return
			}
			rec := logh.records[0]
			if got := rec.Message; got != tc.expMsg {
				t.Errorf("message %q expected, got: %q", tc.expMsg, got)
			}
			attrs := []string{}
			rec.Attrs(func(a slog.Attr) bool {
				if !a.Equal(slog.Attr{}) {
					attrs = append(attrs, a.Key, a.Value.String())
				}
				return true
			})
			if !slices.Equal(tc.expAttrs, attrs) {
				t.Errorf("attrs expected:\n%v, got:\n%v", tc.expAttrs, attrs)
			}
		})
	}
}

type testcases []struct {
	path       string
	logger     *slog.Logger
	withRemote bool
	withHeader bool
	expMsg     string
	expAttrs   []string
}

type testLoggingHandler struct {
	records []slog.Record
}

func (h *testLoggingHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *testLoggingHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}

func (h *testLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // für Einfachheit, kein Attr-Support
}

func (h *testLoggingHandler) WithGroup(name string) slog.Handler {
	return h // keine Gruppen
}
