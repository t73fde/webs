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

// Package reqid provides a middleware functor to enrich HTTP requests (and
// optionally HTTP responses) with a unique identifier.
package reqid

import (
	"context"
	"net/http"

	"t73f.de/r/zero/contexts"
	"t73f.de/r/zero/snow"

	"t73f.de/r/webs/middleware"
)

// DefaultHeaderKey specifies the HTTP header key, where the request ID should be stored.
const DefaultHeaderKey = "X-Request-Id"

// DefaultLoggingKey is the default value for [Config.LoggingKey].
const DefaultLoggingKey = "REQ-ID"

// Config stores all configuration to build a Functor.
type Config struct {
	HeaderKey    string // Key of HTTP request / response that sores the value
	Generator    *snow.Generator
	AppID        uint
	LoggingKey   string // Key for logging, see [Config.WithLogger].
	WithContext  bool
	WithResponse bool
}

// Build the Functor from the configuration.
func (c *Config) Build() middleware.Functor {
	headerKey := c.HeaderKey
	if c.HeaderKey == "" {
		headerKey = DefaultHeaderKey
	}
	gen, appID := c.Generator, c.AppID
	if gen == nil {
		gen = snow.New(0)
	}
	if m := gen.MaxAppID(); appID > m {
		appID = 0
	}
	withContext := c.WithContext
	withResponse := c.WithResponse
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := gen.Create(appID)
			if withContext {
				r = r.WithContext(withReqID(r.Context(), id))
			}
			s := id.String()
			r.Header.Set(headerKey, s)
			if withResponse {
				w.Header().Set(headerKey, s)
			}
			next.ServeHTTP(w, r)
		})
	}
}

type ctxKeyType struct{}

var withReqID, getReqID = contexts.WithAndValue[snow.Key](ctxKeyType{})

// GetRequestID returns the request identification injected by the middleware functor.
func GetRequestID(ctx context.Context) snow.Key {
	if id, ok := getReqID(ctx); ok {
		return id
	}
	return snow.Invalid
}
