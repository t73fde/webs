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

package forms

import "strings"

// ActNormalizersions are technically validators.
// They are used to modify field values.

// TrimNormalizer removes all leading and trailing white space.
type TrimNormalizer struct{}

// Check will trim the given field.
func (TrimNormalizer) Check(_ *Form, fld Field) error {
	return fld.SetValue(strings.TrimSpace(fld.Value()))
}
