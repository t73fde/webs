// -----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
// -----------------------------------------------------------------------------

package forms_test

import (
	"testing"

	"t73f.de/r/webs/forms"
)

func TestValidatorAnyNoneOf(t *testing.T) {
	valid, invalid := []string{"a", "c"}, []string{"b", "d"}
	anyOf := forms.AnyOf(valid...)
	noneOf := forms.NoneOf(valid...)

	fieldAny := forms.TextField("any", "any", anyOf)
	fieldNone := forms.TextField("none", "none", noneOf)
	form := forms.Define(fieldAny, fieldNone)

	for _, v := range valid {
		if err := fieldAny.SetValue(v); err != nil {
			t.Errorf("fieldAny.SetValue(%q) failed: %v", v, err)
		}
		if err := anyOf.Check(form, fieldAny); err != nil {
			t.Error("unexpected error for any:", err)
		}

		if err := fieldNone.SetValue(v); err != nil {
			t.Errorf("fieldNone.SetValue(%q) failed: %v", v, err)
		}
		if err := noneOf.Check(form, fieldNone); err == nil {
			t.Error("expected error for none, but got nil")
		}
	}

	for _, v := range invalid {
		if err := fieldAny.SetValue(v); err != nil {
			t.Errorf("fieldAny.SetValue(%q) failed: %v", v, err)
		}
		if err := anyOf.Check(form, fieldAny); err == nil {
			t.Error("expected error for any, but got nil")
		}

		if err := fieldNone.SetValue(v); err != nil {
			t.Errorf("fieldNone.SetValue(%q) failed: %v", v, err)
		}
		if err := noneOf.Check(form, fieldNone); err != nil {
			t.Error("unexpected error for none:", err)
		}
	}
}
