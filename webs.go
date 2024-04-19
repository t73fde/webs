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

// Package webs provides some common definitions
package webs

import "net/http"

// FuncMiddleware is a functions that encapsulates a HandlerFunc.
type FuncMiddleware func(http.HandlerFunc) http.HandlerFunc
