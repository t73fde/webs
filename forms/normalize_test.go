// -----------------------------------------------------------------------------
// Copyright (c) 2026-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2026-present Detlef Stern
// -----------------------------------------------------------------------------

package forms_test

import (
	"testing"

	"t73f.de/r/webs/forms"
)

func TestTrimNormalizer(t *testing.T) {
	fld := forms.TextAreaField("name", "label")
	if err := fld.SetValue("\tbla  fasel\n "); err != nil {
		t.Error("SetValue", err)
		return
	}
	trim := forms.TrimNormalizer{}
	if err := trim.Check(nil, fld); err != nil {
		t.Error("Check", err)
		return
	}
	exp := "bla  fasel"
	if got := fld.Value(); exp != got {
		t.Errorf("expected %q, but got %q", exp, got)
	}

	fld = forms.TextAreaField("name", "label", trim, forms.StringEqual("x", "msg"))
	if err := fld.SetValue("\t x \n "); err != nil {
		t.Error("SetValue2", err)
		return
	}
	form := forms.Define(fld)
	if !form.IsValid() {
		t.Error("IsValid()", form.Messages())
	}
}
