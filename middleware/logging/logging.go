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

// Package logging provides middleware functors to log HTTP requests and
// responses.
package logging

import (
	"log/slog"
	"net/http"

	"t73f.de/r/webs/ip"
	"t73f.de/r/webs/middleware"
	"t73f.de/r/webs/middleware/reqid"
)

// DefaultRequestIDKey is the default name of the request id log attribute.
const DefaultRequestIDKey = "id"

// ReqConfig stores all configuration data to build a request logger.
type ReqConfig struct {
	Logger        *slog.Logger
	Level         slog.Level
	Message       string
	WithRequestID bool
	WithRemote    bool
	WithHeaders   bool
}

// Build the Functor from the configuration.
func (c *ReqConfig) Build() middleware.Functor {
	logger := c.Logger
	if logger == nil {
		return middleware.NilFunctor
	}
	level := c.Level
	msg := c.Message
	if msg == "" {
		msg = "REQ"
	}
	withRequestID, withRemote, withHeaders := c.WithRequestID, c.WithRemote, c.WithHeaders
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestIDAttr, remoteAttr, headerAttr slog.Attr
			if withRequestID {
				requestIDAttr = slog.Any(DefaultRequestIDKey, reqid.GetRequestID(r.Context()))
			}
			if withRemote {
				remoteValue := ip.GetRemoteAddr(r)
				if remoteValue != "" {
					remoteAttr = slog.String("remote", remoteValue)
				}
			}
			if withHeaders {
				headerAttr = slog.Any("header", r.Header)
			}

			logger.LogAttrs(r.Context(), level, msg, requestIDAttr,
				slog.String("method", r.Method), slog.Any("url", r.URL),
				remoteAttr, headerAttr)
			next.ServeHTTP(w, r)
		})
	}
}

// RespConfig stores all confguration data to build a response logger.
type RespConfig struct {
	Logger        *slog.Logger
	Level         slog.Level
	Message       string
	WithRequestID bool
	WithHeaders   bool
}

// Build the Functor from the configuration.
func (c *RespConfig) Build() middleware.Functor {
	logger := c.Logger
	if logger == nil {
		return middleware.NilFunctor
	}
	level := c.Level
	msg := c.Message
	if msg == "" {
		msg = "RSP"
	}
	withRequestID, withHeaders := c.WithRequestID, c.WithHeaders
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logw := logResponseWriter{w: w}
			next.ServeHTTP(&logw, r)

			var requestIDAttr, headerAttr slog.Attr
			if withRequestID {
				requestIDAttr = slog.Any(DefaultRequestIDKey, reqid.GetRequestID(r.Context()))
			}
			if withHeaders {
				headerAttr = slog.Any("header", logw.Header())
			}

			logger.LogAttrs(r.Context(), level, msg, requestIDAttr,
				slog.String("method", r.Method), slog.Any("url", r.URL),
				slog.Int("status", logw.code), slog.Int("length", logw.length),
				headerAttr)

		})
	}
}

type logResponseWriter struct {
	w      http.ResponseWriter
	code   int
	length int
}

func (lrw *logResponseWriter) Header() http.Header { return lrw.w.Header() }

func (lrw *logResponseWriter) Write(data []byte) (int, error) {
	length, err := lrw.w.Write(data)
	lrw.length += length
	return length, err
}
func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.code = code
	lrw.w.WriteHeader(code)
}
