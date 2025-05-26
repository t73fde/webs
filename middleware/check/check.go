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

// Package check provides a middleware functor to checks HTTP requests.
package check

import (
	"context"
	"net/http"

	"t73f.de/r/webs/middleware"
)

// Checker test an HTTP request for some precondition.
type Checker interface {
	// Check the precondition of an request. If the precondition is satisfied,
	// return a Context and a true value. The request may be modified too, e.g.
	// header data. Otherwise return any Context and a false value. In
	// addition, the Checker must write an error message to the ResponseWriter.
	Check(http.ResponseWriter, *http.Request) (context.Context, bool)
}

// Func is a Checker inside a function.
type Func func(http.ResponseWriter, *http.Request) (context.Context, bool)

// Check the request.
func (cf Func) Check(w http.ResponseWriter, r *http.Request) (context.Context, bool) {
	return cf(w, r)
}

// Build a Checker middleware.
func Build(c Checker) middleware.Functor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ctx, ok := c.Check(w, r); ok {
				if ctx != r.Context() {
					r = r.WithContext(ctx)
				}
				next.ServeHTTP(w, r)
			}
		})
	}
}
