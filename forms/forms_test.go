// -----------------------------------------------------------------------------
// Copyright (c) 2024-present Detlef Stern
//
// This file is part of sxwebs.
//
// sxwebs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2023-present Detlef Stern
// -----------------------------------------------------------------------------

package forms_test

import (
	"maps"
	"net/url"
	"slices"
	"strings"
	"testing"

	"t73f.de/r/webs/forms"
	"t73f.de/r/webs/htmls/render"
)

func TestSimpleRequiredForm(t *testing.T) {
	f := forms.Define(
		forms.TextField("username", "User name", forms.Required{"username"}),
		forms.PasswordField("password", "Password", forms.Required{"password"}),
		forms.SubmitField("submit", "Login"),
	)
	f.SetFormValues(nil, nil)
	if got := f.IsValid(); got {
		t.Error("empty form must not validate")
	}
	gotMsgs := f.Messages()
	if len(gotMsgs) == 0 {
		t.Error("form did not validate, but there are no messages")
	}
	expMsgs := forms.Messages{
		"password": {"password"},
		"username": {"username"},
	}
	if !maps.EqualFunc(expMsgs, gotMsgs, slices.Equal) {
		t.Errorf("expected errors: %v, but got %v", expMsgs, gotMsgs)
	}

	f.SetFormValues(url.Values{"username": nil, "password": nil}, nil)
	if got := f.IsValid(); got {
		t.Error("nil form must not validate")
	}

	f.SetFormValues(url.Values{"username": {"user"}, "password": {"pass"}}, nil)
	if got := f.IsValid(); !got {
		t.Error("normal form must validate")
	}
	expData := forms.Data{"password": "pass", "username": "user"}
	if gotData := f.Data(); !maps.Equal(expData, gotData) {
		t.Errorf("expected data %v, but got %v", expData, gotData)
	}
}

func renderForm(f *forms.Form) string {
	var sb strings.Builder
	if err := render.Render(&sb, f.Render()); err != nil {
		return "{[{" + err.Error() + "}]}"
	}
	return sb.String()
}

func TestRenderNilForm(t *testing.T) {
	var f *forms.Form
	if got := f.Render(); got != nil {
		t.Errorf("nil snippet expected, but got: %v", got)
	}
}
