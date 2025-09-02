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
	"io"

	"t73f.de/r/zero/runes"
)

// Escape writes the text, where some characters are replaced by HTML entities.
func Escape(w io.Writer, text string) error {
	if mw, ok := w.(myWriter); ok {
		return escape(mw, text)
	}
	buf := bufio.NewWriter(w)
	if err := escape(buf, text); err != nil {
		return err
	}
	return buf.Flush()
}

func escape(w myWriter, s string) error {
	pos := 0
	lenS := len(s)
	for i := range lenS {
		var escaped string

		switch s[i] {
		case '&':
			escaped = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			escaped = "&#39;"
		case '<':
			escaped = "&lt;"
		case '>':
			escaped = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			escaped = "&#34;"
		default:
			continue
		}

		if pos < i {
			if _, err := w.WriteString(s[pos:i]); err != nil {
				return nil
			}
		}
		if _, err := w.WriteString(escaped); err != nil {
			return err
		}
		pos = i + 1
	}

	if pos < lenS {
		if _, err := w.WriteString(s[pos:]); err != nil {
			return err
		}
	}
	return nil
}

// EscapeAttrKey writes an attribute key. Illegal characters, as specified in
// https://html.spec.whatwg.org/multipage/syntax.html#syntax-attribute-name
// are ignored.
func EscapeAttrKey(w io.Writer, key string) error {
	if mw, ok := w.(myWriter); ok {
		return escapeAttrKey(mw, key)
	}
	buf := bufio.NewWriter(w)
	if err := escapeAttrKey(buf, key); err != nil {
		return err
	}
	return buf.Flush()
}

func escapeAttrKey(w myWriter, key string) error {
	pos := 0
	for i, r := range key {
		if runes.IsAttributeName(r) {
			continue
		}
		if pos < i {
			if _, err := w.WriteString(key[pos:i]); err != nil {
				return err
			}
		}
		pos = i + 1
	}
	if pos < len(key) {
		if _, err := w.WriteString(key[pos:]); err != nil {
			return err
		}
	}
	return nil
}

// EscapeAttrValue writes an attribute value.
func EscapeAttrValue(w io.Writer, value string) error {
	if mw, ok := w.(myWriter); ok {
		return escapeAttrValue(mw, value)
	}
	buf := bufio.NewWriter(w)
	if err := escapeAttrValue(buf, value); err != nil {
		return err
	}
	return buf.Flush()
}

func escapeAttrValue(w myWriter, value string) error {
	if err := w.WriteByte('"'); err != nil {
		return err
	}
	if err := escape(w, value); err != nil {
		return err
	}
	if err := w.WriteByte('"'); err != nil {
		return err
	}
	return nil
}

// EscapeComment writes the string as a valid HTML5 comment.
func EscapeComment(w io.Writer, s string) error {
	if mw, ok := w.(myWriter); ok {
		return escapeComment(mw, s)
	}
	buf := bufio.NewWriter(w)
	if err := escapeComment(buf, s); err != nil {
		return err
	}
	return buf.Flush()
}

func escapeComment(w myWriter, s string) error {
	start := 0
	lenS := len(s)
	lenSm3 := lenS - 3
	for i := range lenS {
		var escaped string

		switch s[i] {
		case '&':
			escaped = "&amp;"

		case '>':
			if i > 0 {
				if prev := s[i-1]; prev != '!' && prev != '-' {
					continue
				}
			}
			escaped = "&gt;"

		case '<':
			if i < lenSm3 && s[i+1] == '!' && s[i+2] == '-' && s[i+3] == '-' {
				escaped = "&lt;"
			} else if i == lenSm3 && s[i+1] == '!' && s[i+2] == '-' {
				escaped = "&lt;"
			} else {
				continue
			}

		default:
			continue
		}

		if start < i {
			if _, err := w.WriteString(s[start:i]); err != nil {
				return nil
			}
		}
		if _, err := w.WriteString(escaped); err != nil {
			return err
		}
		start = i + 1
	}

	if start < lenS {
		if _, err := w.WriteString(s[start:]); err != nil {
			return err
		}
	}
	return nil
}
