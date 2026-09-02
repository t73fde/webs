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
	"strings"
	"testing"

	"t73f.de/r/webs/forms"
	"t73f.de/r/webs/htmls"
	"t73f.de/r/webs/htmls/render"
)

func TestFlowContent(t *testing.T) {
	form := forms.Define(forms.FlowContentField("fce1", htmls.Elem("p", nil, htmls.Text("Test"))))

	exp := "<form action=\"\" method=\"POST\"><p>Test</p></form>"
	if got := renderForm(form); got != exp {
		t.Errorf("expected: %q, but got: %q", exp, got)
	}
}

func TestSubmitElement(t *testing.T) {
	se := forms.SubmitField("Edit", "edit")
	var sb strings.Builder
	if err := render.Render(&sb, se.Render("se", nil)); err != nil {
		t.Error(err)
		return
	}
	exp := "<input id=\"se\" name=\"Edit\" type=\"submit\" value=\"edit\" class=\"primary\">"
	if got := sb.String(); got != exp {
		t.Errorf("render should be %q, but got %q", exp, got)
	}

	se.SetPriority(17)
	sb.Reset()
	if err := render.Render(&sb, se.Render("se", nil)); err != nil {
		t.Error(err)
		return
	}
	exp = "<input id=\"se\" name=\"Edit\" type=\"submit\" value=\"edit\" class=\"level-17\">"
	if got := sb.String(); got != exp {
		t.Errorf("render should be %q, but got %q", exp, got)
	}

	se.SetCancel()
	sb.Reset()
	if err := render.Render(&sb, se.Render("se", nil)); err != nil {
		t.Error(err)
		return
	}
	exp = "<input id=\"se\" name=\"Edit\" type=\"submit\" value=\"edit\" class=\"cancel\" formnovalidate=\"\">"
	if got := sb.String(); got != exp {
		t.Errorf("render should be %q, but got %q", exp, got)
	}
}
