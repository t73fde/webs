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

import (
	"fmt"
	"strconv"
	"strings"

	"t73f.de/r/webs/htmls"
)

// Field represents a HTTP form field.
type Field interface {
	Name() string
	Value() string
	Clear()
	SetValue(string) error
	Validators() Validators
	Disable()
	Render(string, []string) *htmls.Node
}

// ----- Submit input element

// SubmitElement represents an element <input type="submit" ...>
type SubmitElement struct {
	name           string
	label          string
	value          string
	prio           uint8
	disabled       bool
	noFormValidate bool
}

// SubmitField builds a new submit field.
func SubmitField(name, label string) *SubmitElement {
	return &SubmitElement{
		name:  name,
		label: label,
	}
}

// SetPriority sets the importance of the field. Only the values 0, 1, 2, and 3
// are allowed, with 0 being the highest priority.
func (se *SubmitElement) SetPriority(prio uint8) *SubmitElement {
	se.prio = min(prio, uint8(len(submitPrioClass)-1))
	return se
}

var submitPrioClass = map[uint8]string{
	0: "primary",
	1: "secondary",
	2: "tertiary",
	3: "cancel", // must always be the last, see se.SetCancel()
}

// NoFormValidate marks the submit field as an action that disables form
// validation, if this field causes the form to be sent.
func (se *SubmitElement) NoFormValidate() *SubmitElement {
	se.noFormValidate = true
	return se
}

// SetCancel marks the submit field to work as as a cancel button.
func (se *SubmitElement) SetCancel() *SubmitElement {
	se.prio = uint8(len(submitPrioClass) - 1)
	se.noFormValidate = true
	return se
}

// Name returns the name of this element.
func (se *SubmitElement) Name() string { return se.name }

// Value returns the value of this element.
func (se *SubmitElement) Value() string { return se.value }

// Clear the element.
func (se *SubmitElement) Clear() { se.value = "" }

// SetValue sets the value of this element.
func (se *SubmitElement) SetValue(value string) error { se.value = value; return nil }

// Validators return the currently active validators.
func (se *SubmitElement) Validators() Validators { return nil }

// Disable the submit element.
func (se *SubmitElement) Disable() { se.disabled = true }

// Render the submit element as SxHTML.
func (se *SubmitElement) Render(fieldID string, _ []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(se.Validators())
	attrs := makeAttributes(5, valAttrs, se.disabled, se.noFormValidate)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: se.name},
		htmls.Attribute{Key: "type", Value: "submit"},
		htmls.Attribute{Key: "value", Value: se.label},
		htmls.Attribute{Key: "class", Value: submitPrioClass[se.prio]},
	)
	attrs = addEnablingAttributes(attrs, se.disabled, valAttrs)
	attrs = addBoolAttribute(attrs, "formnovalidate", se.noFormValidate)
	return htmls.Elem("input", attrs)
}

// ----- Checkbox field

// CheckboxElement represents a checkbox.
type CheckboxElement struct {
	name     string
	label    string
	value    string
	disabled bool
}

// CheckboxField provides a checkbox.
func CheckboxField(name, label string) *CheckboxElement {
	return &CheckboxElement{
		name:  name,
		label: label,
	}
}

// Name returns the name of this element.
func (cbe *CheckboxElement) Name() string { return cbe.name }

// Value returns the value of this element.
func (cbe *CheckboxElement) Value() string { return cbe.value }

// Clear the element.
func (cbe *CheckboxElement) Clear() { cbe.value = "" }

// SetValue sets the value of this element.
func (cbe *CheckboxElement) SetValue(value string) error { cbe.value = value; return nil }

// SetChecked sets the value of the checkbox element.
func (cbe *CheckboxElement) SetChecked(val bool) {
	if val {
		cbe.value = cbe.name
	} else {
		cbe.value = ""
	}
}

// Validators return the currently active validators.
func (cbe *CheckboxElement) Validators() Validators { return nil }

// Disable the checkbox element.
func (cbe *CheckboxElement) Disable() { cbe.disabled = true }

// Render the checkbox element.
func (cbe *CheckboxElement) Render(fieldID string, _ []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(cbe.Validators())
	attrs := makeAttributes(5, valAttrs, cbe.value != "", cbe.disabled)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: cbe.name},
		htmls.Attribute{Key: "type", Value: "checkbox"},
		htmls.Attribute{Key: "value", Value: cbe.name},
	)
	attrs = addBoolAttribute(attrs, "checked", cbe.value != "")
	attrs = addEnablingAttributes(attrs, cbe.disabled, valAttrs)

	return htmls.Elem("div", nil,
		htmls.Elem("input", attrs),
		renderLabel(cbe, fieldID, cbe.label),
	)
}

// ----- <textarea ...>...</textarea> field

// TextAreaElement represents the corresponding textarea form element.
type TextAreaElement struct {
	name       string
	label      string
	rows       uint32
	cols       uint32
	value      string
	validators Validators
	disabled   bool
}

// TextAreaField creates a new text area element.
func TextAreaField(name, label string, validators ...Validator) *TextAreaElement {
	return &TextAreaElement{
		name:       name,
		label:      label,
		validators: validators,
	}
}

// SetRows sets the number of rows for the text area element.
func (tae *TextAreaElement) SetRows(rows uint32) *TextAreaElement {
	tae.rows = rows
	return tae
}

// SetCols sets the number of columns for the text area, i.e. the number of
// possibly visible lines.
func (tae *TextAreaElement) SetCols(cols uint32) *TextAreaElement {
	tae.cols = cols
	return tae
}

// Name returns the name of the text area element.
func (tae *TextAreaElement) Name() string { return tae.name }

// Value returns the value of the text area.
func (tae *TextAreaElement) Value() string { return tae.value }

// Clear the text area.
func (tae *TextAreaElement) Clear() { tae.value = "" }

// SetValue sets the value of the text area. Sequences of '\r\n' will be replaced by '\n'.
func (tae *TextAreaElement) SetValue(value string) error {
	tae.value = strings.ReplaceAll(value, "\r\n", "\n") // Unify Windows/Unix EOL handling
	return nil
}

// Validators returns the currently active validators for this text area.
func (tae *TextAreaElement) Validators() Validators {
	if tae.disabled {
		return nil
	}
	return tae.validators
}

// Disable the text area element.
func (tae *TextAreaElement) Disable() { tae.disabled = true }

// Render the text area.
func (tae *TextAreaElement) Render(fieldID string, messages []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(tae.Validators())
	attrs := makeAttributes(5, valAttrs, tae.rows > 0, tae.cols > 0, tae.disabled)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: tae.name},
	)
	if rows := tae.rows; rows > 0 {
		attrs = append(attrs, htmls.Attribute{Key: "rows", Value: strconv.FormatUint(uint64(rows), 10)})
	}
	if cols := tae.cols; cols > 0 {
		attrs = append(attrs, htmls.Attribute{Key: "cols", Value: strconv.FormatUint(uint64(cols), 10)})
	}
	attrs = addEnablingAttributes(attrs, tae.disabled, valAttrs)

	msgs := renderMessages(messages)
	divNode := htmls.Elem("div", nil)
	divNode.Children = make([]*htmls.Node, 2+len(msgs))
	divNode.AddChildren(renderLabel(tae, fieldID, tae.label))
	divNode.AddChildren(msgs...)
	divNode.AddChildren(htmls.Elem("textarea", attrs, htmls.Text(tae.value)))
	return divNode
}

// ----- <select ...>...</select> field

// SelectElement represents the corresponding select form element.
type SelectElement struct {
	name       string
	label      string
	choices    []string
	value      string
	validators Validators
	disabled   bool
}

// SelectField creates a new select element.
func SelectField(name, label string, choices []string, validators ...Validator) *SelectElement {
	se := &SelectElement{
		name:       name,
		label:      label,
		validators: validators,
	}
	se.SetChoices(choices)
	return se
}

// SetChoices allows to update the choices after field creation, e.g. for
// dynamically generated choices.
func (se *SelectElement) SetChoices(choices []string) {
	if len(choices) == 0 || len(choices) == 1 {
		se.choices = nil
	} else if len(choices)%2 != 0 {
		se.choices = choices[0 : len(choices)-2]
	} else {
		se.choices = choices
	}
}

// Name returns the element name.
func (se *SelectElement) Name() string { return se.name }

// Value returns the value of the select element.
func (se *SelectElement) Value() string { return se.value }

// Clear the select element.
func (se *SelectElement) Clear() { se.value = "" }

// SetValue sets the value of the select element.
func (se *SelectElement) SetValue(value string) error {
	se.value = value
	for i := 0; i < len(se.choices); i += 2 {
		if se.choices[i] == value {
			return nil
		}
	}
	return fmt.Errorf("no such choice: %q", value)
}

// Validators return the active validators for the select element.
func (se *SelectElement) Validators() Validators {
	if se.disabled {
		return nil
	}
	return se.validators
}

// Disable the field.
func (se *SelectElement) Disable() { se.disabled = true }

// Render the select element.
func (se *SelectElement) Render(fieldID string, messages []string) *htmls.Node {
	valAttrs := makeValidatorAttributes(se.Validators())
	attrs := makeAttributes(5, valAttrs, se.disabled)
	attrs = append(attrs,
		htmls.Attribute{Key: "id", Value: fieldID},
		htmls.Attribute{Key: "name", Value: se.name},
	)
	attrs = addEnablingAttributes(attrs, se.disabled, valAttrs)

	choiceNodes := make([]*htmls.Node, 0, len(se.choices)/2)
	for i := 0; i < len(se.choices); i += 2 {
		choice := se.choices[i]
		optAttrs := makeAttributes(1, nil, choice == "", se.value == choice)
		optAttrs = append(optAttrs, htmls.Attribute{Key: "value", Value: choice})
		optAttrs = addEnablingAttributes(optAttrs, se.disabled, nil)
		optAttrs = addBoolAttribute(optAttrs, "selected", se.value == choice)
		choiceNodes = append(choiceNodes, htmls.Elem("option", optAttrs, htmls.Text(se.choices[i+1])))
	}

	divElem := htmls.Elem("div", nil, renderLabel(se, fieldID, se.label))
	divElem.Children = append(divElem.Children, renderMessages(messages)...)
	divElem.Children = append(divElem.Children, htmls.Elem("select", attrs, choiceNodes...))
	return divElem
}

// ----- Flow Content -----

// FlowContentElement adds some flow content to the form.
type FlowContentElement struct {
	name    string
	content *htmls.Node
}

// FlowContentField allows to add some text (aka flow content) to the form.
func FlowContentField(name string, content *htmls.Node) *FlowContentElement {
	return &FlowContentElement{name: name, content: content}
}

// Name returns the element name.
func (fce *FlowContentElement) Name() string { return fce.name }

// Value returns the value of the select element.
func (*FlowContentElement) Value() string { return "" }

// Clear the select element.
func (*FlowContentElement) Clear() {}

// SetValue sets the value of the select element.
func (*FlowContentElement) SetValue(string) error {
	return fmt.Errorf("flow content has no specific value")
}

// Validators return the active validators for the select element.
func (*FlowContentElement) Validators() Validators { return nil }

// Disable the field.
func (*FlowContentElement) Disable() {}

// Render the flow content element.
func (fce *FlowContentElement) Render(string, []string) *htmls.Node {
	return fce.content
}

// ----- General utility functions for rendering etc.

func renderLabel(field Field, fieldID, label string) *htmls.Node {
	if label == "" {
		return nil
	}
	labelText := htmls.Text(label)
	if field.Validators().HasRequired() {
		labelText.Data += "*"
	}
	return htmls.Elem("label", []htmls.Attribute{{Key: "for", Value: fieldID}}, labelText)
}

func renderMessages(messages []string) []*htmls.Node {
	result := make([]*htmls.Node, 0, len(messages))
	for _, msg := range messages {
		result = append(result,
			htmls.Elem("span", []htmls.Attribute{{Key: "class", Value: "message"}}, htmls.Text(msg)))
	}
	return result
}

func addBoolAttribute(attrs []htmls.Attribute, key string, val bool) []htmls.Attribute {
	if val {
		return append(attrs, htmls.Attribute{Key: key})
	}
	return attrs
}

// addEnablingAttributes adds some attributes, depending whether the field is
// disabled or not. If it is disabled, the "disabled" attribute will be added,
// and no validator attributes are added too.
// Otherwise, the field is enable and therefore the attributes of an validator
// will be added.
func addEnablingAttributes(attrs []htmls.Attribute, disabled bool, validatorAttributes []htmls.Attribute) []htmls.Attribute {
	if disabled {
		return append(attrs, htmls.Attribute{Key: "disabled"})
	}
	return append(attrs, validatorAttributes...)
}

func makeValidatorAttributes(validators []Validator) []htmls.Attribute {
	if len(validators) == 0 {
		return nil
	}
	result := make([]htmls.Attribute, 0, len(validators))
	for _, val := range validators {
		result = append(result, val.Attributes()...)
	}
	return result
}

func makeAttributes(minLen int, valAttrs []htmls.Attribute, opt ...bool) []htmls.Attribute {
	length := minLen + len(valAttrs)
	for _, b := range opt {
		if b {
			length++
		}
	}
	return make([]htmls.Attribute, 0, length)
}
