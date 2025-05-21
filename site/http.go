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
	st.Root.handle(reg, st.Basepath, middleware.Nil{})
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

	m = extendMiddleware(reg, m, n.Middleware)

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
		hmw := extendMiddleware(reg, m, n.HandlerMW[i])
		handler = middleware.Apply(hmw, handler)
		reg.Handle(method+" "+hPath, handler)
	}

	for _, child := range n.Children {
		child.handle(reg, upath, m)
	}
}

func extendMiddleware(reg Registerer, baseMW middleware.Middleware, name string) middleware.Middleware {
	if mw, found := reg.GetMiddleware(name); found {
		lb := middleware.NewListFromMiddleware(baseMW)
		lmw := middleware.NewListFromMiddleware(mw)
		return lb.Extend(lmw)
	}
	return baseMW
}
