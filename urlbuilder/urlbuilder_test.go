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
	"net/url"
	"strings"
	"testing"

	"t73f.de/r/webs/urlbuilder"
)

func TestVarURLBuilder(t *testing.T) {
	t.Parallel()

	var ub2 urlbuilder.URLBuilder
	ub2.AddPath("")
	if exp, got := "/", ub2.String(); exp != got {
		t.Errorf("empty path builder must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub2)
	ub2.AddPath("path")
	if exp, got := "/", ub2.String(); exp != got {
		t.Errorf("empty path builder must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub2)
	var ub3 urlbuilder.URLBuilder
	ub3.AddPath("path").AddPath("")
	if exp, got := "/path/", ub3.String(); exp != got {
		t.Errorf("URLBuilder.AddPath/DIR must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub3)

	var ub urlbuilder.URLBuilder
	checkCopy(t, &ub)
	if exp, got := "/", ub.String(); exp != got {
		t.Errorf("empty URLBuilder must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.AddPath("path")
	if exp, got := "/path", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.AddPath("/pf/ad")
	if exp, got := "/path/pf%2Fad", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath2 must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.AddPath("p/")
	if exp, got := "/path/pf%2Fad/p/", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddPath3 must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.SetFragment("frag")
	if exp, got := "/path/pf%2Fad/p/#frag", ub.String(); exp != got {
		t.Errorf("URLBuilder.SetFragment must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.AddQuery("k", "v")
	if exp, got := "/path/pf%2Fad/p/?k=v#frag", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddQuery must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.SetFragment("f")
	if exp, got := "/path/pf%2Fad/p/?k=v#f", ub.String(); exp != got {
		t.Errorf("URLBuilder.SetFragment2 must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.AddQuery("l", "w")
	if exp, got := "/path/pf%2Fad/p/?k=v&l=w#f", ub.String(); exp != got {
		t.Errorf("URLBuilder.AddQuery2 must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.RemoveQueries()
	if exp, got := "/path/pf%2Fad/p/#f", ub.String(); exp != got {
		t.Errorf("URLBuilder.RemoveQueries must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)
	ub.SetFragment(" ")
	if exp, got := "/path/pf%2Fad/p/", ub.String(); exp != got {
		t.Errorf("URLBuilder.SetFragment3 must result in string value %q, but got %q", exp, got)
		return
	}
	checkCopy(t, &ub)

	var ubCopy urlbuilder.URLBuilder
	ub = urlbuilder.URLBuilder{}
	ub.AddPath("a")
	ub.Copy(&ubCopy)
	ub.AddPath("b")
	if exp, got := ub.String(), ubCopy.String(); exp == got {
		t.Errorf("change after copy path should result in different builder, but got %q", got)
		return
	}

	ub.AddQuery("k", "v")
	ub.Copy(&ubCopy)
	ub.AddQuery("l", "w")
	if exp, got := ub.String(), ubCopy.String(); exp == got {
		t.Errorf("change after copy query should result in different builder, but got %q", got)
		return
	}
}

func checkCopy(t *testing.T, ub *urlbuilder.URLBuilder) {
	var ubCopy urlbuilder.URLBuilder
	ub.Copy(&ubCopy)
	exp := ub.String()
	if got := ubCopy.String(); got != exp {
		t.Errorf("copy of %q shoud not change, but got: %q", exp, got)
	}

	u, err := url.Parse(exp)
	if err != nil {
		t.Errorf("unable to parse as url.URL: %v", err)
	}
	if strings.Contains(u.Fragment, "=") {
		t.Errorf("fragment contains query: %q", u.Fragment)
	}
	if got := u.String(); got != exp {
		t.Errorf("parsed url.URL.String() differ from original, expected: %q, but got: %q", exp, got)
	}
}
