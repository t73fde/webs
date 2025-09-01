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

// Package htmls allows to create structured HTML snippets.
//
// These snippets do not form a full DOM tree, they are not designed to build
// a full HTML document.
package htmls

// A Node consists of a NodeType and some Data (tag name for element nodes,
// content for text nodes, comment nodes, and some more). An element node
// may also contain some Attributes and Children. Data is always not escaped,
// i.e. it stores "a<b" and not "a&lt;b".
type Node struct {
	Data       string
	Attributes []Attribute
	Children   []*Node
	Type       NodeType
}

// AddChildren adds some more children to the Node.
func (node *Node) AddChildren(children ...*Node) {
	for _, child := range children {
		if child != nil {
			node.Children = append(node.Children, child)
		}
	}
}

// NodeType is the type of a [Node].
type NodeType uint8

const (
	_ NodeType = iota

	// TextNode signals that only unescaped text is stored in [Node.Data].
	// [Node.Attributes] and [Node.Children] are not used.
	TextNode

	// ElementNode represents an HTML element, i.e. something with a tag,
	// optionally some [Node.Attributes] and [Node.Children]
	ElementNode

	// CommentNode signals that only unescaped comment is stored in [Node.Data].
	// [Node.Attributes] and [Node.Children] are not used.
	CommentNode

	// RawNode signals that already process HTML text is stored in [Node.Data].
	// [Node.Attributes] and [Node.Children] are not used.
	RawNode
)

// An Attribute is a key-value pair to be used in a [Node].
//
// Both Key and Value store unescaped values, i.e. "a<b" instead of "a&lt;b".
type Attribute struct {
	Key   string
	Value string
}

// Elem returns an element node.
func Elem(tag string, attrs []Attribute, children ...*Node) *Node {
	nilChildren := 0
	for _, child := range children {
		if child == nil {
			nilChildren++
		}
	}
	if nilChildren == len(children) {
		children = nil
	} else {
		newChildren := make([]*Node, 0, len(children)-nilChildren)
		for _, child := range children {
			if child != nil {
				newChildren = append(newChildren, child)
			}
		}
		children = newChildren
	}
	return &Node{
		Data:       tag,
		Attributes: attrs,
		Children:   children,
		Type:       ElementNode,
	}
}

// Text returns a text node.
func Text(data string) *Node {
	return &Node{
		Data:       data,
		Attributes: nil,
		Children:   nil,
		Type:       TextNode,
	}
}

// Attrs returns a slice of [Attribute] values.
func Attrs(keyval ...string) []Attribute {
	if len(keyval)%2 == 1 {
		keyval = append(keyval, "")
	}
	result := make([]Attribute, 0, len(keyval)/2)
	for i := 0; i < len(keyval); i += 2 {
		result = append(result, Attribute{Key: keyval[i], Value: keyval[i+1]})
	}
	return result
}
