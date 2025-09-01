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

// Package tags provides some helper functions for HTML tags.
package tags

// IsVoid returns true, if the given tag is a "void tag", according to section
// 13.1.2 of the HTML5 spec.
func IsVoid(tag string) bool {
	_, found := voidTags[tag]
	return found
}

var voidTags = map[string]struct{}{
	"area":   {},
	"base":   {},
	"br":     {},
	"col":    {},
	"embed":  {},
	"hr":     {},
	"img":    {},
	"input":  {},
	"link":   {},
	"meta":   {},
	"source": {},
	"track":  {},
	"wbr":    {},
}

// IsLiteralChildTextTag tests agains HTML5 spec, section 13.3, if children
// text nodes must be written unescaped.
func IsLiteralChildTextTag(tag string) bool {
	switch tag {
	case "style", "script", "xmp", "iframe", "noembed", "noframes", "plaintext":
		return true
	}
	return false
}
