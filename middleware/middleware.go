//-----------------------------------------------------------------------------
// Copyright (c) 2024-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2024-present Detlef Stern
//-----------------------------------------------------------------------------

// Package middleware provides some net/http middleware.
package middleware

import (
	"iter"
	"net/http"
)

// Middleware is a function that transforms an http.Handler into an http.Handler.
type Middleware func(http.Handler) http.Handler

// Seq is a sequence of Middlewares.
type Seq interface {

	// Values returns an iterator of the Middlewares to apply.
	Values() iter.Seq[Middleware]
}

// Apply the Middleware sequence to the given handler, resulting in a modified handler.
func Apply(seq Seq, h http.Handler) http.Handler {
	for mw := range seq.Values() {
		h = mw(h)
	}
	return h
}
