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

// Package status provides a middleware that associates HTTP status codes with
// [http.Handler]s.
//
// One example is a 404 / [http.StatusNotFound] handler that redirects to
// another site, as proposed by https://4042307.org/:
//
//	cfg := status.Config{HandlerMap:HandlerMap{
//	           http.StatusNotFound: status.BaseRedirectHandler(
//	                                    "https://2017.4042307.org",
//	                                    http.StatusTemporaryRedirect)}}
//	f := cfg.Build()
package status

import (
	"net/http"
	"strings"

	"t73f.de/r/webs/middleware"
)

// Config stores the base data for the status redirect middleware functor.
type Config struct {
	// HandlerMap maps a HTTP status code to its handler.
	//
	// The provides status codes should be in the 4xx and 5xx range.
	HandlerMap HandlerMap

	// NoClearMap maps HTTP status codes to a boolean value that signals not
	// to clear the HTTP header before calling the handler.
	NoClearMap map[int]bool
}

// HandlerMap maps HTTP status codes to handler.
type HandlerMap map[int]http.Handler

// Build a middleware functor that will call a handler when the base handler
// results in a given status code.
func (c Config) Build() middleware.Functor {
	m := c.HandlerMap
	if m == nil {
		m = HandlerMap{}
	}
	nc := c.NoClearMap
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srw := statusRespWriter{m: m, nc: nc, w: w, r: r}
			next.ServeHTTP(&srw, r)
		})
	}
}

type statusRespWriter struct {
	m  HandlerMap
	nc map[int]bool
	w  http.ResponseWriter
	r  *http.Request

	found bool
}

func (srw *statusRespWriter) Header() http.Header {
	return srw.w.Header()
}
func (srw *statusRespWriter) WriteHeader(code int) {
	if h, found := srw.m[code]; found {
		srw.found = true
		if nc := srw.nc; nc == nil || !nc[code] {
			clear(srw.w.Header())
		}
		h.ServeHTTP(srw.w, srw.r)
		return
	}
	srw.w.WriteHeader(code)
}
func (srw *statusRespWriter) Write(data []byte) (int, error) {
	if srw.found {
		// Ignore data/body from original request as we started a new handler.
		return len(data), nil
	}
	return srw.w.Write(data)
}

// BaseRedirectHandler returns a handler that redirects each request it
// receives using the given status code. The redirect URL is calculated by
// appending the requests URL (a path, an optional query, and an optional
// fragment) to the given base URL.
//
// Base URL must not have a suffix "/" as `r.URL.Path` already has a "/"
// prefix.
//
// The provided code should be in the 3xx range.
func BaseRedirectHandler(baseURL string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		var sb strings.Builder
		sb.WriteString(baseURL)
		sb.WriteString(u.Path)
		if query := u.RawQuery; query != "" || u.ForceQuery {
			sb.WriteByte('?')
			sb.WriteString(query)
		}
		if fragment := u.Fragment; fragment != "" {
			sb.WriteByte('#')
			sb.WriteString(fragment)
		}
		http.Redirect(w, r, sb.String(), code)
	})
}
