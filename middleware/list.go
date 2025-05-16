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

import "iter"

// List is a single linked list of Middleware.
type List struct {
	mw   Middleware
	next *List
}

// NewList creates a new list, based on the previous list and a Middleware.
func NewList(mw Middleware, lst *List) *List {
	return &List{mw: mw, next: lst}
}

// NewListFromChain build a new list from a Middleware Chain.
func NewListFromChain(chn Chain) (l *List) {
	for _, mw := range chn.seq {
		l = NewList(mw, l)
	}
	return l
}

// Append builds a new list by adding the given middleware to the list.
func (l *List) Append(mw Middleware) *List {
	return NewList(mw, l)
}

// Extend the list by some other list.
func (l *List) Extend(other *List) *List {
	if other == nil {
		return l
	}
	first := NewList(other.mw, nil)
	prev := first
	for curr := other; ; {
		curr = curr.next
		if curr == nil {
			prev.next = l
			return first
		}
		prev.next = NewList(curr.mw, nil)
		prev = prev.next
	}
}

// Values returns an iterator of Middleware to apply.
func (l *List) Values() iter.Seq[Middleware] {
	return func(yield func(Middleware) bool) {
		for e := l; e != nil; e = e.next {
			if !yield(e.mw) {
				return
			}
		}
	}
}
