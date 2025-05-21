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

	"t73f.de/r/webs/middleware"
)

// ReqConfig stores all configuration data to build a request logger.
type ReqConfig struct {
	Logger      *slog.Logger
	Level       slog.Level
	Message     string
	WithHeaders bool
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
	if c.WithHeaders {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger.Log(r.Context(), level, msg, "method", r.Method, "url", r.URL, "header", r.Header)
				next.ServeHTTP(w, r)
			})
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log(r.Context(), level, msg, "method", r.Method, "url", r.URL)
			next.ServeHTTP(w, r)
		})
	}
}

// RespConfig stores all confguration data to build a response logger.
type RespConfig struct {
	Logger      *slog.Logger
	Level       slog.Level
	Message     string
	WithHeaders bool
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
		msg = "RESP"
	}
	if c.WithHeaders {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logw := logResponseWriter{w: w}
				next.ServeHTTP(&logw, r)
				logger.Log(r.Context(), level, msg,
					"method", r.Method, "url", r.URL,
					"status", logw.code, "length", logw.length,
					"header", logw.Header())
			})
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logw := logResponseWriter{w: w}
			next.ServeHTTP(&logw, r)
			logger.Log(r.Context(), level, msg,
				"method", r.Method, "url", r.URL,
				"status", logw.code, "length", logw.length)
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
