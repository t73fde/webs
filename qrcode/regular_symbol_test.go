//-----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// This file was originally created by Tom Harwood under an MIT license, but
// later changed to fulfil the needs of webs. The text of the original license
// can be found in file ORIG_LICENSE. The following statements affects the
// original code as found on https://github.com/skip2/go-qrcode (Commit:
// da1b6568686e89143e94f980a98bc2dbd5537f13, 2020-06-17):
//
// go-qrcode
// Copyright 2014 Tom Harwood
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
//-----------------------------------------------------------------------------

package qrcode

import (
	"testing"

	"t73f.de/r/webs/qrcode/internal/bitset"
)

func TestBuildRegularSymbol(_ *testing.T) {
	for k := range 8 {
		v := getQRCodeVersion(Low, 1)

		data := bitset.New()
		for range 26 {
			data.AppendNumBools(8, false)
		}

		s := buildRegularSymbol(*v, k, data, false)
		_ = s
		//fmt.Print(m.string())
	}
}
