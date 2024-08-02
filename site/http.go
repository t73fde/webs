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

import "net/http"

// Registerer contains all methods need to register handler for HTTP.
type Registerer interface {
	GetHandler(string) (http.Handler, bool)
	AddMiddleware(*Node, http.Handler) http.Handler
	Handle(string, http.Handler)
}

// Handle registers all named handlers for the whole site.
func (st *Site) Handle(reg Registerer) {
	st.Root.Handle(reg, st.Basepath)
}

// Handle registers all named handlers for the node and its children.
func (n *Node) Handle(reg Registerer, basepath string) {
	rawpath := basepath + n.Nodepath
	path := rawpath + "/"

	var hPath string
	switch n.PathSpec() {
	case pathSpecDir:
		hPath = rawpath + "/{$}"
	case pathSpecFull:
		hPath = path
	case pathSpecItem:
		hPath = rawpath
	}

	for i, handlerName := range n.Handler {
		st := n.site
		if len(st.Methods) < i {
			break
		}
		handler, found := reg.GetHandler(handlerName)
		if !found {
			continue
		}
		handler = reg.AddMiddleware(n, handler)

		method := st.Methods[i]
		reg.Handle(method+" "+hPath, handler)
	}

	for _, child := range n.Children {
		child.Handle(reg, path)
	}
}
