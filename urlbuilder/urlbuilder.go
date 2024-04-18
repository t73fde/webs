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

package urlbuilder

import (
	"net/url"
	"slices"
	"strings"
)

// URLBuilder helps to build (absolute) URLs.
type URLBuilder struct {
	path     []string
	fragment string
	query    []urlQuery
}
type urlQuery struct{ key, val string }

// Copy an URLBuilder.
func (ub *URLBuilder) Copy(dest *URLBuilder) {
	dest.path = slices.Clone(ub.path)
	dest.fragment = ub.fragment
	dest.query = slices.Clone(ub.query)
}

// AddPath adds a new path element.
func (ub *URLBuilder) AddPath(p string) *URLBuilder {
	for len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	if p != "" {
		ub.path = append(ub.path, p)
	}
	return ub
}

// SetFragment stores the fragment
func (ub *URLBuilder) SetFragment(frag string) *URLBuilder {
	ub.fragment = strings.TrimSpace(frag)
	return ub
}

// AddQuery adds a new key/value query parameter
func (ub *URLBuilder) AddQuery(key, value string) *URLBuilder {
	ub.query = append(ub.query, urlQuery{key, value})
	return ub
}

// RemoveQueries removes all previously added key/value query parameter.
// This allows to recycle an URLBuilder, to be used for various query
// parameter values, where the path (and the fragment) stays constant.
func (ub *URLBuilder) RemoveQueries() *URLBuilder {
	ub.query = nil
	return ub
}

// String constructs a string representation of the URL.
func (ub *URLBuilder) String() string {
	var sb strings.Builder

	if len(ub.path) == 0 {
		sb.WriteByte('/')
	}
	for _, p := range ub.path {
		sb.WriteByte('/')
		if pl := len(p); pl > 0 && p[pl-1] == '/' {
			sb.WriteString(url.PathEscape(p[0 : pl-1]))
			sb.WriteByte('/')
		} else {
			sb.WriteString(url.PathEscape(p))
		}
	}

	for i, q := range ub.query {
		if i == 0 {
			sb.WriteByte('?')
		} else {
			sb.WriteByte('&')
		}
		sb.WriteString(q.key)
		if val := q.val; val != "" {
			sb.WriteByte('=')
			sb.WriteString(url.QueryEscape(val))
		}
	}

	if ub.fragment != "" {
		sb.WriteByte('#')
		sb.WriteString(ub.fragment)
	}

	return sb.String()
}
