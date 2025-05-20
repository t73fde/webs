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
	f    Functor
	next *List
}

// NewList creates a new list, based on the previous list and a Middleware.
func NewList(f Functor, lst *List) *List {
	return &List{f: f, next: lst}
}

// NewListFromMiddleware build a new list from a given Middleware.
func NewListFromMiddleware(m Middleware) *List {
	var sentinel List
	curr := &sentinel
	for f := range m.Functors() {
		l := NewList(f, nil)
		curr.next = l
		curr = l
	}
	return sentinel.next
}

// Append builds a new list by adding the given middleware to the list.
func (l *List) Append(f Functor) *List {
	return NewList(f, l)
}

// Extend the list by some other list.
func (l *List) Extend(other *List) *List {
	if other == nil {
		return l
	}
	first := NewList(other.f, nil)
	prev := first
	for curr := other; ; {
		curr = curr.next
		if curr == nil {
			prev.next = l
			return first
		}
		prev.next = NewList(curr.f, nil)
		prev = prev.next
	}
}

// Functors returns an iterator of Middleware to apply.
func (l *List) Functors() iter.Seq[Functor] {
	return func(yield func(Functor) bool) {
		for e := l; e != nil; e = e.next {
			if !yield(e.f) {
				return
			}
		}
	}
}
