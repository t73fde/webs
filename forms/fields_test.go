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
	"t73f.de/r/webs/htmls"
)

func TestFlowContent(t *testing.T) {
	form := forms.Define(forms.FlowContentField("fce1", htmls.Elem("p", nil, htmls.Text("Test"))))

	exp := "<form action=\"\" method=\"POST\"><p>Test</p></form>"
	if got := renderForm(form); got != exp {
		t.Errorf("expected: %q, but got: %q", exp, got)
	}
}
