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

	"t73f.de/r/webs/htmls/render"
)

func TestEscape(t *testing.T) {
	testcases := []struct {
		name    string
		comment string
		exp     string
	}{
		{"nil", "", ""},
		{"space", " ", " "},
		{"amp", "&", "&amp;"},
		{"amp-s", "a&", "a&amp;"},
		{"amp-e", "&a", "&amp;a"},
		{"apos", "a'b", "a&#39;b"},
		{"lt", "a<b", "a&lt;b"},
		{"gt", "a>b", "a&gt;b"},
		{"quot", "a\"b", "a&quot;b"},
		{"abc", "abc", "abc"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := render.Escape(&sb, tc.comment)
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
func TestEscapeAttrKey(t *testing.T) {
	testcases := []struct {
		key string
		exp string
	}{
		{"", ""},
		{" ", ""}, {" a b c ", "abc"},
		{"=", ""}, {"=a", "a"}, {"a=b", "ab"}, {"a=", "a"},
	}

	for _, tc := range testcases {
		t.Run(tc.key, func(t *testing.T) {
			var sb strings.Builder
			err := render.EscapeAttrKey(&sb, tc.key)
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

func TestEscapeAttrValue(t *testing.T) {
	testcases := []struct {
		key string
		exp string
	}{
		{"", "\"\""},
		{" ", "\" \""}, {" a b c ", "\" a b c \""},
	}

	for _, tc := range testcases {
		t.Run(tc.key, func(t *testing.T) {
			var sb strings.Builder
			err := render.EscapeAttrValue(&sb, tc.key)
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

func TestEscapeComment(t *testing.T) {
	testcases := []struct {
		name    string
		comment string
		exp     string
	}{
		{"nil", "", ""},
		{"space", " ", " "},
		{"gt", ">abc", "&gt;abc"},
		{"gt-minus", "->abc", "-&gt;abc"},
		{"start", "abc<!--def", "abc&lt;!--def"},
		{"start-2", "<!--def", "&lt;!--def"},
		{"start-3", "abc<!--", "abc&lt;!--"},
		{"end", "abc-->def", "abc--&gt;def"},
		{"end-2", "-->def", "--&gt;def"},
		{"end-3", "abc-->", "abc--&gt;"},
		{"end!", "abc--!>def", "abc--!&gt;def"},
		{"end!-2", "--!>def", "--!&gt;def"},
		{"end!-3", "abc--!>", "abc--!&gt;"},
		{"ends", "abc<!-", "abc&lt;!-"},
		{"ends-ok", "abc<!-def", "abc<!-def"},
		{"amp", "abc&def;ghi", "abc&amp;def;ghi"},
		{"amp-2", "&def;ghi", "&amp;def;ghi"},
		{"amp-3", "abc&def;", "abc&amp;def;"},
		{"h>", "h>", "h>"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var sb strings.Builder
			err := render.EscapeComment(&sb, tc.comment)
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
