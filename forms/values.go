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

package forms

import (
	"reflect"
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
	return t.Local().Format(htmlDatetimeLayout)
}

// SignedInt is the super-type of all signed integers.
type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInt is the super-type of all unsigned integers.
type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// IntValue returns the value as a string to be stored in a field.
func IntValue[T SignedInt](i T) string { return strconv.FormatInt(int64(i), 10) }

// UintValue returns the value as a string to be stored in a field.
func UintValue[T UnsignedInt](i T) string { return strconv.FormatUint(uint64(i), 10) }

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

// GetDate returns the value of the given field as a UTC-based time.Time,
// but only as a real date, with time 00:00:00.
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
				return result.UTC()
			}
		}
	}
	return time.Time{}
}

// GetInt returns the value of the given field as an int.
func (d Data) GetInt[T SignedInt](fieldName string, defaultValue T) T {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			bits := reflect.TypeOf(defaultValue).Bits()
			if result, err := strconv.ParseInt(value, 10, bits); err == nil {
				return T(result)
			}
		}
	}
	return defaultValue
}

// GetUint returns the value of the given field as a number.
func (d Data) GetUint[T UnsignedInt](fieldName string, defaultValue T) T {
	if len(d) > 0 {
		if value, found := d[fieldName]; found {
			bits := reflect.TypeOf(defaultValue).Bits()
			if result, err := strconv.ParseUint(value, 10, bits); err == nil {
				return T(result)
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
