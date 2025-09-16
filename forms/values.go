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

// Package forms handles HTML form data.
package forms

import (
	"strconv"
	"time"
)

// Time layouts of data coming from HTML forms.
const (
	htmlDateLayout     = "2006-01-02"
	htmlDatetimeLayout = "2006-01-02T15:04"
)

// DateValue returns the date as a string suitable for a HTML date field value.
func DateValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDateLayout)
}

// DatetimeValue returns the time as a string suitable for a HTML datetime-local field value.
func DatetimeValue(t time.Time) string {
	if t.Equal(time.Time{}) {
		return ""
	}
	return t.Format(htmlDatetimeLayout)
}

// IntValue returns the value as a string to be stored in a field.
func IntValue(i int) string { return strconv.Itoa(i) }

// UintValue returns the value as a string to be stored in a field.
func UintValue(i uint64) string { return strconv.FormatUint(i, 10) }

// CheckboxValue returns a value for a checkbox.
// The value should be the name for the [CheckboxField].
func CheckboxValue(b bool, val string) string {
	if b {
		return val
	}
	return ""
}

// Data contains all form data, as a map of field names to field values.
type Data map[string]string

// Get string data of a field. Return empty string for unknown field.
func (d Data) Get(fieldName string) string {
	if len(d) == 0 {
		return ""
	}
	if value, found := d[fieldName]; found {
		return value
	}
	return ""
}

// GetDate returns the value of the given field as a time.Time, but only
// as a real date, with time 00:00:00.
func (d Data) GetDate(fieldName string) time.Time {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := time.Parse(htmlDateLayout, value); err == nil {
				return result
			}
		}
	}
	return time.Time{}
}

// GetDatetime returns the value of the given field as a time.Time.
func (d Data) GetDatetime(fieldName string) time.Time {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := time.ParseInLocation(htmlDatetimeLayout, value, time.Local); err == nil {
				return result
			}
		}
	}
	return time.Time{}
}

// GetInt returns the value of the given field as an int.
func (d Data) GetInt(fieldName string, defaultValue int) int {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := strconv.Atoi(value); err == nil {
				return result
			}
		}
	}
	return defaultValue
}

// GetUint returns the value of the given field as a number.
func (d Data) GetUint(fieldName string, defaultValue uint64) uint64 {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := strconv.ParseUint(value, 10, 64); err == nil {
				return result
			}
		}
	}
	return defaultValue
}

// GetFloat returns the value of the given field as a number.
func (d Data) GetFloat(fieldName string, defaultValue float64) float64 {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			if result, err := strconv.ParseFloat(value, 64); err == nil {
				return result
			}
		}
	}
	return defaultValue
}
