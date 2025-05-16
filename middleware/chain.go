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

package middleware

// Based on the article "Organize your Go middleware without dependencies"
// (https://www.alexedwards.net/blog/organize-your-go-middleware-without-dependencies)
// by Alex Edwards and some ideas from https://github.com/justinas/alice

import (
	"iter"
	"slices"
)

// Chain is a immutable sequence of Middleware functors that encapsulate an handler.
type Chain struct {
	seq []Middleware
}

// NewChain creates a new Chain.
func NewChain(seq ...Middleware) Chain {
	return Chain{seq: slices.Clone(seq)}
}

// NewChainFromList builds a Chain from a given Middleware List.
func NewChainFromList(l *List) Chain {
	return Chain{seq: slices.Collect(l.Values())}
}

// Append middleware to the Chain, resulting in a new Chain.
func (chn Chain) Append(seq ...Middleware) Chain {
	return Chain{seq: slices.Concat(chn.seq, seq)}
}

// Extend a Chain by another one, resulting in a new Chain.
func (chn Chain) Extend(other Chain) Chain { return chn.Append(other.seq...) }

// Values return an iterator of the Middleware Chain, in order of application.
func (chn Chain) Values() iter.Seq[Middleware] {
	return func(yield func(Middleware) bool) {
		for _, mw := range slices.Backward(chn.seq) {
			if !yield(mw) {
				return
			}
		}
	}
}
