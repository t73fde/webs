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

// Package header provides middleware functors to set response header to specific values.
package header

import (
	"maps"
	"net/http"

	"t73f.de/r/webs/middleware"
)

// Config stores all configuration data to build a header setting functor.
type Config struct {
	Constants map[string]string
}

// Build the Functor from the configuration.
func (c *Config) Build() middleware.Functor {
	if len(c.Constants) == 0 {
		return middleware.NilFunctor
	}
	constMap := maps.Clone(c.Constants)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := w.Header()
			for k, v := range constMap {
				if _, found := header[k]; !found {
					header.Add(k, v)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
