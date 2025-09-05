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
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/slogtest"

	"t73f.de/r/webs/middleware/reqid"
)

func TestReqIDLogging(t *testing.T) {
	var sb strings.Builder
	th := slog.NewTextHandler(&sb, nil)
	baseLogger := slog.New(th)

	reqConfig := reqid.Config{WithContext: true}
	reqLogger := reqConfig.WithLogger(baseLogger)
	reqLogger = reqLogger.With("key", "val")

	handler := func(_ http.ResponseWriter, r *http.Request) {
		reqLogger.Info("NOCO")
		reqLogger.InfoContext(r.Context(), "CTX!")
	}
	rmw := reqConfig.Build()
	mux := http.NewServeMux()
	mux.Handle("/foo", rmw(http.HandlerFunc(handler)))
	r, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, r)

	lines := strings.Split(sb.String(), "\n")
	if len(lines) < 2 {
		t.Error("at least two lines expected, but got:", lines)
		return
	}
	nocoLine, ctxLine := lines[0], lines[1]
	if exp := " level=INFO msg=NOCO"; !strings.Contains(nocoLine, exp) {
		t.Errorf("expected %q in first log, but got: %q", exp, nocoLine)
	}
	if exp := " level=INFO msg=CTX! key=val REQ-ID="; !strings.Contains(ctxLine, exp) {
		t.Errorf("expected %q in second log, but got: %q", exp, ctxLine)
	}
}

func TestReqIDHandler(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := slog.New(slog.NewJSONHandler(&buf, nil))

	reqConfig := reqid.Config{WithContext: true}
	reqLogger := reqConfig.WithLogger(baseLogger)

	results := func() []map[string]any {
		var ms []map[string]any
		for line := range bytes.SplitSeq(buf.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				panic(err) // In a real test, use t.Fatal.
			}
			ms = append(ms, m)
		}
		return ms
	}
	err := slogtest.TestHandler(reqLogger.Handler(), results)
	if err != nil {
		t.Error(err)
	}
}

func TestReqIDLoggerDefaul(t *testing.T) {
	reqConfig := reqid.Config{WithContext: true}
	logger := reqConfig.WithLogger(nil)
	if logger != nil {
		t.Errorf("logger must be nil, but got: %v", logger)
	}

	var buf bytes.Buffer
	baseLogger := slog.New(slog.NewJSONHandler(&buf, nil))
	reqConfig.WithContext = false
	logger = reqConfig.WithLogger(baseLogger)
	if logger != baseLogger {
		t.Errorf("logger must be base: %v, but got: %v", baseLogger, logger)
	}
}
