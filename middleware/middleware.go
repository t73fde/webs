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

import "net/http"

// Middleware is a function that wraps a http.Handler around a http.Handler.
type Middleware func(http.Handler) http.Handler

// Func is a middleware that encapsulates a HandlerFunc.
type Func func(http.HandlerFunc) http.HandlerFunc
