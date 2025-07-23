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

package reedsolomon

import "testing"

func TestGFMultiplicationIdentities(t *testing.T) {
	for i := range 256 {
		value := gfElement(i)
		if gfMultiply(gfZero, value) != gfZero {
			t.Errorf("0 . %d != 0", value)
		}

		if gfMultiply(value, gfOne) != value {
			t.Errorf("%d . 1 == %d, want %d", value, gfMultiply(value, gfOne), value)
		}
	}
}

func TestGFMultiplicationAndDivision(t *testing.T) {
	// a * b == result
	var tests = []struct {
		a   gfElement
		b   gfElement
		exp gfElement
	}{
		{0, 29, 0},
		{1, 1, 1},
		{1, 32, 32},
		{2, 4, 8},
		{16, 128, 232},
		{17, 17, 28},
		{27, 9, 195},
	}

	for _, test := range tests {
		if got := gfMultiply(test.a, test.b); got != test.exp {
			t.Errorf("%d * %d = %d, want %d", test.a, test.b, got, test.exp)
		}

		if test.b != gfZero && test.exp != gfZero {
			if b := gfDivide(test.exp, test.a); b != test.b {
				t.Errorf("%d / %d = %d, want %d", test.exp, test.a, b, test.b)
			}
		}
	}
}

func TestGFInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		a := gfElement(i)
		inverse := gfInverse(a)

		if got := gfMultiply(a, inverse); got != gfOne {
			t.Errorf("%d * %d^-1 == %d, want %d", a, inverse, got, gfOne)
		}
	}
}

func TestGFDivide(t *testing.T) {
	for i := 1; i < 256; i++ {
		for j := 1; j < 256; j++ {
			// a * b == product
			a := gfElement(i)
			b := gfElement(j)
			product := gfMultiply(a, b)

			// product / b == a
			if got := gfDivide(product, b); got != a {
				t.Errorf("%d / %d == %d, want %d", product, b, got, a)
			}
		}
	}
}
