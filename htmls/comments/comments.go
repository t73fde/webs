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

// Package comments provides some helper functions for HTML comments.
package comments

import (
	"bufio"
	"io"
)

// Escape writes the string as a valid HTML5 comment.
func Escape(w io.Writer, s string) error {
	if mw, ok := w.(myWriter); ok {
		return escape(mw, s)
	}
	buf := bufio.NewWriter(w)
	if err := escape(buf, s); err != nil {
		return err
	}
	return buf.Flush()
}

func escape(w myWriter, s string) error {
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

type myWriter interface {
	WriteString(string) (int, error)
}
