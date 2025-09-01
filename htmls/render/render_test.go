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

package render_test

import (
	"strings"
	"testing"

	"t73f.de/r/webs/htmls"
	"t73f.de/r/webs/htmls/render"
)

func TestRender(t *testing.T) {
	testcases := []struct {
		name string
		node *htmls.Node
		exp  string
	}{
		{"nil", nil, ""},
		{"a<b", htmls.Text("a<b"), "a&lt;b"},
		{"h1", htmls.Elem("h1", nil, htmls.Text("heading")), "<h1>heading</h1>"},
		{"h1class",
			htmls.Elem("h1", htmls.Attrs("class", "mini"), htmls.Text("MinHead")),
			"<h1 class=\"mini\">MinHead</h1>",
		},
		{"ol",
			htmls.Elem("ol", []htmls.Attribute{{Key: "reversed"}},
				htmls.Elem("li", nil, htmls.Text("1")),
				htmls.Elem("li", htmls.Attrs("value", "two"), htmls.Text("2")),
				htmls.Elem("li", nil, htmls.Text("3")),
			),
			"<ol reversed=\"\"><li>1</li><li value=\"two\">2</li><li>3</li></ol>"},
		{"br", htmls.Elem("br", nil), "<br>"},
		{"br-child",
			htmls.Elem("br", nil, htmls.Text("error")),
			"{[{void tag \"br\" contains children}]}"},
		{"comment",
			&htmls.Node{Type: htmls.CommentNode, Data: "comment"},
			"<-- comment -->"},
		{"raw",
			&htmls.Node{Type: htmls.RawNode, Data: "<h1>"},
			"<h1>"},
		{"script",
			htmls.Elem("script", nil, htmls.Text("a<b")),
			"<script>a<b</script>"},
		{"pre-nl",
			htmls.Elem("pre", nil, htmls.Text("\nabc\n")),
			"<pre>\n\nabc\n</pre>",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := render.Render(&sb, tc.node)
			var got string
			if err != nil {
				got = "{[{" + err.Error() + "}]}"
			} else {
				got = sb.String()
			}
			if got != tc.exp {
				t.Errorf("\nexpected: %q\n but got: %q", tc.exp, got)
			}
		})
	}
}
