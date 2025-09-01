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
// SPDX-FileCopyrightText: 2024-present Detlef Stern
// -----------------------------------------------------------------------------

package forms

// ----- <input ...> fields

import (
	"time"

	"t73f.de/r/webs/htmls"
)

// InputElement represents a HTTP <input> field.
type InputElement struct {
	name       string
	label      string
	value      string
	validators Validators
	disabled   bool
	itype      inputType
}

type inputType uint

// Constants for inputType
const (
	_ inputType = iota
	itypeCheckbox
	itypeDate
	itypeDatetime
	itypeEmail
	itypeNumber
	itypePassword
	itypeText
)

// Name returns the name of this element.
func (fd *InputElement) Name() string { return fd.name }

// Value returns the value of the input element.
func (fd *InputElement) Value() string { return fd.value }

// Clear the input element.
func (fd *InputElement) Clear() { fd.value = "" }

// Time layouts of data coming from HTML forms.
const (
	htmlDateLayout     = "2006-01-02"
	htmlDatetimeLayout = "2006-01-02T15:04"
)

// SetValue sets the value of this input element.
func (fd *InputElement) SetValue(value string) (err error) {
	fd.value = value
	switch fd.itype {
	case itypeDate:
		if value != "" {
			_, err = time.Parse(htmlDateLayout, value)
		}
	case itypeDatetime:
		if value != "" {
			_, err = time.Parse(htmlDatetimeLayout, value)
		}
	}
	return err
}

// Validators returns all currently active Validators.
func (fd *InputElement) Validators() Validators {
	if fd.disabled {
		return nil
	}
	return fd.validators
}

// Disable the input element.
func (fd *InputElement) Disable() { fd.disabled = true }

// Render the form input element as SxHTML.
func (fd *InputElement) Render(fieldID string, messages []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(fd.Validators())
	attrs := makeAttributes(5, valAttrs, fd.disabled)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: fd.name},
		htmls.Attribute{Key: "type", Value: inputTypeString[fd.itype]},
		htmls.Attribute{Key: "value", Value: fd.value},
	)
	attrs = addEnablingAttributes(attrs, fd.disabled, valAttrs)

	divNode := htmls.Elem("div", nil, renderLabel(fd, fieldID, fd.label))
	divNode.Children = append(divNode.Children, renderMessages(messages)...)
	divNode.Children = append(divNode.Children, htmls.Elem("input", attrs))
	return divNode
}

var inputTypeString = map[inputType]string{
	itypeCheckbox: "checkbox",
	itypeDate:     "date",
	itypeDatetime: "datetime-local",
	itypeEmail:    "email",
	itypeNumber:   "number",
	itypePassword: "password",
	itypeText:     "text",
}

// TextField builds a new text field.
func TextField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeText,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateField builds a new field to enter dates.
func DateField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeDate,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DateValue returns the date as a string suitable for a HTML date field value.
func DateValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDateLayout)
}

// DatetimeField builds a new field to enter a local date/time.
func DatetimeField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeDatetime,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// DatetimeValue returns the time as a string suitable for a HTML datetime-local field value.
func DatetimeValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDatetimeLayout)
}

// PasswordField builds a new password field.
func PasswordField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypePassword,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// EmailField builds a new e-mail field.
func EmailField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeEmail,
		name:       name,
		label:      label,
		validators: validators,
	}
}

// NumberField builds a new number field.
func NumberField(name, label string, validators ...Validator) *InputElement {
	return &InputElement{
		itype:      itypeNumber,
		name:       name,
		label:      label,
		validators: validators,
	}
}
