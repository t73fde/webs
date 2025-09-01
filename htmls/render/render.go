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

// Package render provides a function to render htmls.Node into HTML.
package render

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"strings"

	"t73f.de/r/webs/htmls"
	"t73f.de/r/webs/htmls/comments"
	"t73f.de/r/webs/htmls/tags"
)

// Render writes the given node as simplified HTML5 to the provided writer.
//
// Note: This implementation does not fully comply with HTML5. Escaping is
// minimal and many special rules are ignored. The function is intended for
// testing purposes only.
func Render(w io.Writer, node *htmls.Node) error {
	if mw, ok := w.(myWriter); ok {
		return render(mw, node)
	}
	buf := bufio.NewWriter(w)
	if err := render(buf, node); err != nil {
		return err
	}
	return buf.Flush()
}

func render(w myWriter, node *htmls.Node) error {
	if node == nil {
		return nil
	}
	switch node.Type {
	case htmls.TextNode:
		_, err := w.WriteString(html.EscapeString(node.Data))
		return err
	case htmls.ElementNode:
		// no-op, fall through
	case htmls.CommentNode:
		if _, err := w.WriteString("<-- "); err != nil {
			return err
		}
		if err := comments.Escape(w, node.Data); err != nil {
			return err
		}
		if _, err := w.WriteString(" -->"); err != nil {
			return err
		}
		return nil
	case htmls.RawNode:
		_, err := w.WriteString(node.Data)
		return err
	default:
		return fmt.Errorf("unknown node type: %v", node.Type)
	}

	tag := node.Data
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(tag); err != nil {
		return err
	}
	for _, attr := range node.Attributes {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		if _, err := w.WriteString(html.EscapeString(attr.Key)); err != nil {
			return err
		}
		if _, err := w.WriteString("=\""); err != nil {
			return err
		}
		if _, err := w.WriteString(html.EscapeString(attr.Value)); err != nil {
			return err
		}
		if err := w.WriteByte('"'); err != nil {
			return err
		}
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	if tags.IsVoid(tag) {
		if len(node.Children) > 0 {
			return fmt.Errorf("void tag %q contains children", tag)
		}
		return nil
	}

	// Add initial newline, when it is possible that a newline in text child will be ignored.
	if len(node.Children) > 0 {
		if child := node.Children[0]; child.Type == htmls.TextNode && strings.HasPrefix(child.Data, "\n") {
			switch tag {
			case "pre", "textarea":
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
		}
	}

	if tags.IsLiteralChildTextTag(tag) {
		for _, child := range node.Children {
			if child.Type == htmls.TextNode {
				if _, err := w.WriteString(child.Data); err != nil {
					return err
				}
			} else {
				if err := render(w, child); err != nil {
					return err
				}
			}
		}
	} else {
		for _, child := range node.Children {
			if err := render(w, child); err != nil {
				return err
			}
		}
	}

	if _, err := w.WriteString("</"); err != nil {
		return err
	}
	if _, err := w.WriteString(tag); err != nil {
		return err
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}
	return nil
}

type myWriter interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}
