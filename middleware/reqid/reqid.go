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

package reqid

import (
	"net/http"

	"t73f.de/r/webs/middleware"
	"t73f.de/r/zero/snow"
)

// DefaultHeaderKey specifies the HTTP header key, where the request ID should be stored.
const DefaultHeaderKey = "X-Request-Id"

// New creates a new Middleware to inject a unique request ID.
func New() middleware.Middleware {
	gen := snow.NewGenerator(0)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := gen.Create(0)
			s := id.String()
			r.Header.Set(DefaultHeaderKey, s)
			w.Header().Set(DefaultHeaderKey, s)
			next.ServeHTTP(w, r)
		})
	}
}
