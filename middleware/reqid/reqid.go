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

// Package reqid provides a middleware to enrich HTTP requests (and optionally
// HTTP responses) with a unique identifier.
package reqid

import (
	"net/http"

	"t73f.de/r/webs/middleware"
	"t73f.de/r/zero/snow"
)

// DefaultHeaderKey specifies the HTTP header key, where the request ID should be stored.
const DefaultHeaderKey = "X-Request-Id"

// Config stores all configutation to build a Middleware.
type Config struct {
	HeaderKey    string
	Generator    *snow.Generator
	AppID        uint
	WithResponse bool
}

// Initialize the configuration data with useful data.
func (c *Config) Initialize() {
	if c.HeaderKey == "" {
		c.HeaderKey = DefaultHeaderKey
	}
	if gen := c.Generator; gen != nil {
		c.AppID = min(c.AppID, gen.MaxAppID())
	}
}

// Build the Middleware from the configuration.
func (c *Config) Build() middleware.Middleware {
	headerKey := c.HeaderKey
	if c.HeaderKey == "" {
		headerKey = DefaultHeaderKey
	}
	gen, appID := c.Generator, c.AppID
	if gen == nil {
		gen = snow.New(0)
		appID = 0
	}
	withResponse := c.WithResponse
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := gen.Create(appID)
			s := id.String()
			r.Header.Set(headerKey, s)
			if withResponse {
				w.Header().Set(headerKey, s)
			}
			next.ServeHTTP(w, r)
		})
	}
}
