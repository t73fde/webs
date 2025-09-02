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

package tags_test

import (
	"testing"

	"t73f.de/r/webs/htmls/tags"
)

func BenchmarkVoid(b *testing.B) {
	var samples = []string{
		"a", "p", "div", "span", "em", "h1", "ul", "li", "li", "a",
		"area", "br", "img", "meta", "wbr"}

	for b.Loop() {
		for _, tag := range samples {
			_ = tags.IsVoid(tag)
		}
	}
}
