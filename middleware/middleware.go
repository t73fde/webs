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

// Functor is a function that transforms an http.Handler into an http.Handler.
type Functor func(http.Handler) http.Handler

// NilFunctor is a Functor that does nothing.
func NilFunctor(h http.Handler) http.Handler { return h }

// Functors returns an iterator with only the given Functor. This allows to
// treat the functor as a Middleware.
func (f Functor) Functors() iter.Seq[Functor] {
	return func(yield func(Functor) bool) {
		_ = yield(f)
	}
}

// Middleware is a sequence of http.Handler transforming Functors.
type Middleware interface {

	// Functors returns an iterator of the Functors to apply.
	Functors() iter.Seq[Functor]
}

// Apply the Middleware sequence to the given handler, resulting in a modified handler.
func Apply(m Middleware, h http.Handler) http.Handler {
	for f := range m.Functors() {
		h = f(h)
	}
	return h
}

// Nil is a middleware that does nothing.
type Nil struct{}

// Functors returns an iterator of all Functors to apply. For this middleware:
// nothing is returned.
func (Nil) Functors() iter.Seq[Functor] { return func(func(Functor) bool) {} }
