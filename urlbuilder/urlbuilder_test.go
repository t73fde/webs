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

package urlbuilder_test

import (
	"testing"

	"t73f.de/r/webs/urlbuilder"
)

func TestNewURLBuilder(t *testing.T) {
	t.Parallel()
	checkURLBuilder(t, urlbuilder.New(""))
	checkURLBuilder(t, urlbuilder.New("/"))

	ub := urlbuilder.New("prefix/")
	if exp, got := "/prefix", ub.String(); exp != got {
		t.Errorf("prefix/ URLBuilder must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddPath("path")
	if exp, got := "/prefix/path", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath must result in string value %q, but got %q", exp, got)
		return
	}

	ub = urlbuilder.New("prefix")
	if exp, got := "/prefix", ub.String(); exp != got {
		t.Errorf("prefix/ URLBuilder must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddPath("path")
	if exp, got := "/prefix/path", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath must result in string value %q, but got %q", exp, got)
		return
	}
}

func TestVarURLBuilder(t *testing.T) {
	t.Parallel()
	var ub urlbuilder.URLBuilder
	checkURLBuilder(t, &ub)
}

func checkURLBuilder(t *testing.T, ub *urlbuilder.URLBuilder) {
	t.Helper()
	if exp, got := "/", ub.String(); exp != got {
		t.Errorf("empty URLBuilder must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddPath("path")
	if exp, got := "/path", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddPath("/pfad")
	if exp, got := "/path/pfad", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath2 must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddPath("p/")
	if exp, got := "/path/pfad/p/", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath3 must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.SetFragment("frag")
	if exp, got := "/path/pfad/p/#frag", ub.String(); exp != got {
		t.Errorf("URLBuilder.SetFragment must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.SetFragment("f")
	if exp, got := "/path/pfad/p/#f", ub.String(); exp != got {
		t.Errorf("URLBuilder.SetFragment2 must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddQuery("k", "v")
	if exp, got := "/path/pfad/p/#f?k=v", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddQuery must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.AddQuery("l", "w")
	if exp, got := "/path/pfad/p/#f?k=v&l=w", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddQuery2 must result in string value %q, but got %q", exp, got)
		return
	}
	ub = ub.RemoveQueries()
	if exp, got := "/path/pfad/p/#f", ub.String(); exp != got {
		t.Errorf("URLBuilder.RemoveQueries must result in string value %q, but got %q", exp, got)
		return
	}
}
