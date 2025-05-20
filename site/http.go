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

package site

import (
	"net/http"
	"path"

	"t73f.de/r/webs/middleware"
)

// Registerer contains all methods need to register handler for HTTP.
type Registerer interface {
	GetHandler(string) (http.Handler, bool)
	GetMiddleware(string) (middleware.Middleware, bool)
	Handle(string, http.Handler)
}

// Handle registers all named handlers for the whole site.
func (st *Site) Handle(reg Registerer) {
	m, found := reg.GetMiddleware(st.Middleware)
	if !found {
		m = middleware.NewChain()
	}
	st.Root.handle(reg, st.Basepath, m)
}

// Handle registers all named handlers for the node and its children.
func (n *Node) handle(reg Registerer, basepath string, m middleware.Middleware) {
	upath := path.Join(basepath, n.Nodepath)

	var hPath string
	switch n.pathSpec {
	case pathSpecDir:
		hPath = path.Join(upath, "{$}")
	case pathSpecFull:
		hPath = upath + "/"
	case pathSpecItem:
		hPath = upath
	}

	if nm, found := reg.GetMiddleware(n.Middleware); found {
		lm := middleware.NewListFromMiddleware(m)
		lnm := middleware.NewListFromMiddleware(nm)
		m = lm.Extend(lnm)
	}

	methods := n.site.Methods
	for i, handlerName := range n.Handler {
		if len(methods) < i {
			break
		}
		method := methods[i]
		if method == "" || handlerName == "" {
			continue
		}
		handler, found := reg.GetHandler(handlerName)
		if !found {
			continue
		}
		hmw, found := reg.GetMiddleware(n.HandlerMW[i])
		if found {
			lm := middleware.NewListFromMiddleware(m)
			lnm := middleware.NewListFromMiddleware(hmw)
			m = lm.Extend(lnm)
		}
		handler = middleware.Apply(m, handler)
		reg.Handle(method+" "+hPath, handler)
	}

	for _, child := range n.Children {
		child.handle(reg, upath, m)
	}
}
