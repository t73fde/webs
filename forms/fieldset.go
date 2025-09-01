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

package forms

import (
	"t73f.de/r/webs/htmls"
)

// Fieldset represents an HTML <fieldset>
type Fieldset struct {
	form     *Form
	name     string
	legend   string
	fields   []Field
	disabled bool
}

func (fs *Fieldset) setForm(f *Form) {
	for _, fd := range fs.fields {
		f.addName(fd)
	}
	fs.form = f
}

// FieldsetField builds a Fieldset.
func FieldsetField(name, legend string, fields ...Field) *Fieldset {
	return &Fieldset{
		form:     nil,
		name:     name,
		legend:   legend,
		fields:   fields,
		disabled: false,
	}
}

// Name the Fieldset.
func (fs *Fieldset) Name() string { return fs.name }

// Value returns the value of the Fieldset: there is no value.
func (Fieldset) Value() string { return "" }

// Clear the Fieldset.
func (fs *Fieldset) Clear() {
	for _, f := range fs.fields {
		f.Clear()
	}
}

// SetValue resetturns the value of the Fieldset: there is no value -> ignore
func (Fieldset) SetValue(string) error { return nil }

// Validators returns the validators for this Fieldset: there are no validators.
func (Fieldset) Validators() Validators { return nil }

// Disable the Fieldset.
func (fs *Fieldset) Disable() {
	for _, f := range fs.fields {
		f.Disable()
	}
}

// Render the Fieldset.
func (fs *Fieldset) Render(fieldID string, messages []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(fs.Validators())
	attrs := makeAttributes(5, valAttrs, fs.disabled)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: fs.name},
	)
	attrs = addEnablingAttributes(attrs, fs.disabled, valAttrs)

	msgs := renderMessages(messages)
	numChildren := len(msgs) + len(fs.fields)
	if fs.legend != "" {
		numChildren++
	}

	fsNode := htmls.Elem("fieldset", attrs)
	fsNode.Children = make([]*htmls.Node, 0, numChildren)
	if legend := fs.legend; legend != "" {
		fsNode.Children = append(fsNode.Children, htmls.Elem("legend", nil, htmls.Text(legend)))
	}
	fsNode.Children = append(fsNode.Children, renderMessages(messages)...)
	form := fs.form
	for _, field := range fs.fields {
		fsNode.Children = append(fsNode.Children, field.Render(form.calcFieldID(field), form.messages[field.Name()]))
	}

	return fsNode
}
