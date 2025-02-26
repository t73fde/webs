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

// Package site allows to define the structure of a web site.
package site

import (
	"fmt"
	"net/http"
	"path"
	"slices"
	"strings"

	"t73f.de/r/webs/urlbuilder"
)

// Site contains information about the web site.
type Site struct {
	Name     string // Name of the site.
	Basepath string // Base path, where the site is available.
	Language string
	Methods  []string // HTTP methods to be used by node handler. Default: GET, POST.
	Root     Node     // Root note of the site.

	baked     bool
	basepaths []string
	nodes     map[string]*Node
}

// DefaultLanguage is the language value used as a default.
const DefaultLanguage = "en"

// Bake the internal data of the Site.
func (st *Site) Bake() error {
	if st.baked {
		return nil
	}

	st.Name = strings.TrimSpace(st.Name)

	basepath := strings.TrimSpace(st.Basepath)
	for p := range strings.SplitSeq(basepath, "/") {
		if p != "" {
			st.basepaths = append(st.basepaths, p)
		}
	}
	st.Basepath = "/" + path.Join(st.basepaths...)

	st.Language = strings.TrimSpace(st.Language)
	if st.Language == "" {
		st.Language = DefaultLanguage
	}

	if len(st.Methods) == 0 {
		st.Methods = []string{
			http.MethodGet,
			http.MethodPost,
		}
	}

	err := st.Root.bake(st, nil)
	st.baked = (err == nil)
	return err
}

// Basepaths returns the base path of the application as a string slice.
func (st *Site) Basepaths() []string { return st.basepaths }

// Node returns the Node with the given identification.
func (st *Site) Node(id string) *Node {
	if nodes := st.nodes; nodes != nil {
		return nodes[id]
	}
	return nil
}

// BestNode returns the node that matches the given path at best. If an
// absolute path (starting with '/') is given, a nil result indicates
func (st *Site) BestNode(path string) *Node {
	if path == "" {
		return &st.Root
	}
	relpath := path
	if relpath[0] == '/' {
		relpath = relpath[1:]
	}
	return st.Root.BestNode(relpath)
}

// MakeURLBuilder creates and configures a new URL builder for the web user interface.
func (st *Site) MakeURLBuilder() *urlbuilder.URLBuilder {
	var ub urlbuilder.URLBuilder
	for _, p := range st.basepaths {
		ub.AddPath(p)
	}
	return &ub
}

// BuilderFor returns an URL builder initialized with the full path of the
// node with the given identifier.
func (st *Site) BuilderFor(nodeID string, args ...string) *urlbuilder.URLBuilder {
	n := st.Node(nodeID)
	if n == nil {
		return nil
	}
	return n.BuilderFor(args...)
}

// Node stores information about one element of the web site, i.e. a web page.
type Node struct {
	ID       string            // Unique identification
	Nodepath string            // Path element
	Title    string            // Title of the node: <title>{TITLE}</title>, <h1>{TITLE}</h1>
	Language string            // Language of the node
	Extra    map[string]string // Some extra information, to be defined by application
	Handler  []string          // 0=GET, 1=POST (see Site.Methods)
	Children []*Node           // Child nodes

	site     *Site
	parent   *Node
	pathSpec pathSpec
	hmap     map[string]string
}

// pathSpec determines the type of the nodepath.
type pathSpec uint8

// pathSpec constants / enumeration
const (
	pathSpecDir  pathSpec = iota // Path ends with '/'
	pathSpecFull                 // Node has no children, it should match full path
	pathSpecItem                 // Path does not end with '/'
)

// Path returns the full path to this node.
func (n *Node) Path() string {
	ancestors := []string{}
	for a := n; a != nil; a = a.parent {
		if p := a.Nodepath; p != "" {
			ancestors = append(ancestors, a.Nodepath)
		}
	}
	ancestors = append(ancestors, n.site.Basepath)
	slices.Reverse(ancestors)
	result := path.Join(ancestors...)
	if n.pathSpec == pathSpecDir && result[len(result)-1] != '/' {
		return result + "/"
	}
	return result
}

// GetTitle returns the title of the node. If no title is stored, its ID is returned.
func (n *Node) GetTitle() string {
	if title := n.Title; title != "" {
		return title
	}
	return n.ID
}

// SetHandler set the given handler name for the given method.
func (n *Node) SetHandler(method, handler string) {
	if hm := n.hmap; hm != nil {
		hm[method] = handler
		return
	}
	n.hmap = map[string]string{method: handler}
}

// SetExtra set a key to a value.
func (n *Node) SetExtra(key, val string) {
	if extra := n.Extra; len(extra) > 0 {
		extra[key] = val
		return
	}
	n.Extra = map[string]string{key: val}
}

// GetExtra returns the stored value of a specific key.
func (n *Node) GetExtra(key string) (string, bool) {
	if extra := n.Extra; len(extra) > 0 {
		val, found := extra[key]
		return val, found
	}
	return "", false
}

// Parent returns the superior / parent node (or nil, if root node).
func (n *Node) Parent() *Node { return n.parent }

// BestNode returns the node that matches the given relative path the best.
// It never returns nil.
func (n *Node) BestNode(relpath string) *Node {
	for _, child := range n.Children {
		childpath := child.Nodepath
		if len(childpath) > 1 && childpath[0] == '{' && childpath[len(childpath)-1] == '}' {
			// child path is a placeholder
			sepPos := strings.IndexByte(relpath, '/')
			if sepPos < 0 || sepPos == len(relpath)-1 {
				return child
			}
			return child.BestNode(relpath[sepPos+1:])
		}
		if strings.TrimSuffix(relpath, "/") == childpath {
			return child
		}
		if len(child.Children) > 0 {
			childpath += "/"
			if relpath == childpath {
				return child
			}
			if len(relpath) >= len(childpath) && childpath == relpath[0:len(childpath)] {
				return child.BestNode(relpath[len(childpath):])
			}
		}
	}
	return n
}

// bake the node data.
func (n *Node) bake(st *Site, p *Node) error {
	if id := strings.TrimSpace(n.ID); id != "" {
		if st.nodes == nil {
			st.nodes = map[string]*Node{id: n}
		} else if st.nodes[id] == nil {
			st.nodes[id] = n
		} else {
			return fmt.Errorf("duplicate id %q for node %v", id, n.Nodepath)
		}
		n.ID = id
	}

	nodepath := strings.TrimSuffix(n.Nodepath, "/")
	if len(nodepath) > 0 {
		switch nodepath[0] {
		case '/':
			n.pathSpec = pathSpecDir
			nodepath = nodepath[1:]
		case '>':
			n.pathSpec = pathSpecFull
			nodepath = nodepath[1:]
		case '*':
			n.pathSpec = pathSpecItem
			nodepath = nodepath[1:]
		default:
			n.pathSpec = pathSpecDir
		}
	}
	n.Nodepath = nodepath

	n.Title = strings.TrimSpace(n.Title)

	n.Language = strings.TrimSpace(n.Language)
	if n.Language == "" {
		if p != nil {
			n.Language = p.Language
		} else {
			n.Language = st.Language
		}
	}

	n.site = st
	n.parent = p

	if hm := n.hmap; hm != nil {
		hsl := make([]string, len(st.Methods))
		for m, h := range hm {
			if pos := slices.Index(st.Methods, m); pos >= 0 {
				hsl[pos] = h
			}
		}
		n.Handler = hsl
	} else if numHandler := len(n.Handler); numHandler > 0 {
		hm = make(map[string]string, numHandler)
		numMethods := len(st.Methods)
		for i, h := range n.Handler {
			if i >= numMethods {
				break
			}
			hm[st.Methods[i]] = h
		}
	}

	children := make([]*Node, 0, len(n.Children))
	for _, child := range n.Children {
		if child == nil {
			continue
		}
		err := child.bake(st, n)
		if err != nil {
			return err
		}
		children = append(children, child)
	}
	n.Children = slices.Clip(children)
	return nil
}

// BuilderFor returns an URL builder for a specific node.
func (n *Node) BuilderFor(args ...string) *urlbuilder.URLBuilder {
	pos := 0
	ancestors := []string{}
	for a := n; a != nil; a = a.parent {
		if pe := a.Nodepath; pe != "" {
			if pe[0] == '{' && pe[len(pe)-1] == '}' {
				if pos < len(args) {
					pe = args[pos]
				} else {
					pe = fmt.Sprintf("missing-arg-%d", pos)
				}
				pos++
			}
			ancestors = append(ancestors, pe)
		}
	}
	ub := n.site.MakeURLBuilder()
	for i := len(ancestors) - 1; i >= 0; i-- {
		ub = ub.AddPath(ancestors[i])
	}

	// Add extra args that were not consumed by key values
	for ; pos < len(args); pos++ {
		ub = ub.AddPath(args[pos])
	}

	if n.pathSpec == pathSpecDir {
		ub = ub.AddPath("")
	}
	return ub
}
