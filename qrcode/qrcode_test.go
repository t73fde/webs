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
	"slices"
	"strings"
	"testing"
)

func TestQRCodeMaxCapacity(t *testing.T) {
	tests := []struct {
		string         string
		numRepetitions int
	}{
		{"0", 7089},
		{"A", 4296},
		{"#", 2953},
		// Alternate byte/numeric data types. Optimises to 2,952 bytes.
		{"#1", 1476},
	}

	for _, test := range tests {
		_, err := New(strings.Repeat(test.string, test.numRepetitions), Low)
		if err != nil {
			t.Errorf("%d x '%s' got %s expected success", test.numRepetitions,
				test.string, err.Error())
		}
	}

	for _, test := range tests {
		_, err := New(strings.Repeat(test.string, test.numRepetitions+1), Low)
		if err == nil {
			t.Errorf("%d x '%s' chars encodable, expected not encodable",
				test.numRepetitions+1, test.string)
		}
	}
}

func TestQRCodeVersionCapacity(t *testing.T) {
	tests := []struct {
		version         int
		level           RecoveryLevel
		maxNumeric      int
		maxAlphanumeric int
		maxByte         int
	}{
		{1, Low, 41, 25, 17},
		{2, Low, 77, 47, 32},
		{2, Highest, 34, 20, 14},
		{40, Low, 7089, 4296, 2953},
		{40, Highest, 3057, 1852, 1273},
	}

	for i, test := range tests {
		numericData := strings.Repeat("1", test.maxNumeric)
		alphanumericData := strings.Repeat("A", test.maxAlphanumeric)
		byteData := strings.Repeat("#", test.maxByte)

		n, err := New(numericData, test.level)
		if err != nil {
			t.Fatal(err.Error())
		}

		a, err := New(alphanumericData, test.level)
		if err != nil {
			t.Fatal(err.Error())
		}

		b, err := New(byteData, test.level)
		if err != nil {
			t.Fatal(err.Error())
		}

		if n.VersionNumber != test.version {
			t.Fatalf("Test #%d numeric has version #%d, expected #%d", i,
				n.VersionNumber, test.version)
		}

		if a.VersionNumber != test.version {
			t.Fatalf("Test #%d alphanumeric has version #%d, expected #%d", i,
				a.VersionNumber, test.version)
		}

		if b.VersionNumber != test.version {
			t.Fatalf("Test #%d byte has version #%d, expected #%d", i,
				b.VersionNumber, test.version)
		}
	}
}

func TestQRCodeISOAnnexIExample(t *testing.T) {
	q, err := New("01234567", Medium)
	if err != nil {
		t.Fatalf("Error producing ISO Annex I Example: %s, expected success",
			err.Error())
	}
	q.encode()

	const expectedMask int = 2
	if q.mask != expectedMask {
		t.Errorf("ISO Annex I example mask got %d, expected %d\n", q.mask,
			expectedMask)
	}
}

func BenchmarkQRCodeURLSize(b *testing.B) {
	for b.Loop() {
		_, _ = New("http://www.example.org", Medium)
	}
}

func BenchmarkQRCodeMaximumSize(b *testing.B) {
	// 7089 is the maximum encodable number of numeric digits.
	content := strings.Repeat("0", 7089)
	for b.Loop() {
		_, _ = New(content, Low)
	}
}

func TestPNGBitmap(t *testing.T) {
	qr, err := New("http://example.org", Low)
	if err != nil {
		t.Fatal(err)
	}
	if exp := 2; qr.VersionNumber != exp {
		t.Errorf("expected version %d, but got %d", exp, qr.VersionNumber)
		return
	}
	got, err := qr.PNG(1)
	if err != nil {
		t.Fatal(err)
	}
	expPNG := []byte{
		137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0,
		33, 0, 0, 0, 33, 1, 3, 0, 0, 0, 109, 42, 80, 44, 0, 0, 0, 6, 80, 76,
		84, 69, 255, 255, 255, 0, 0, 0, 85, 194, 211, 126, 0, 0, 0, 157, 73,
		68, 65, 84, 120, 218, 116, 142, 63, 11, 1, 97, 28, 128, 31, 131, 235,
		119, 253, 202, 187, 94, 50, 40, 146, 50, 8, 163, 193, 164, 140, 22,
		229, 99, 248, 24, 54, 70, 117, 249, 51, 40, 113, 62, 192, 117, 147,
		209, 36, 227, 101, 52, 189, 148, 197, 160, 119, 20, 251, 61, 195, 179,
		62, 15, 89, 152, 123, 199, 129, 84, 219, 2, 186, 47, 90, 208, 237,
		237, 239, 183, 5, 169, 36, 2, 230, 25, 57, 192, 23, 160, 224, 86, 101,
		200, 45, 47, 2, 126, 116, 136, 193, 243, 30, 115, 48, 221, 222, 24,
		88, 215, 6, 224, 127, 250, 49, 228, 195, 96, 8, 58, 153, 157, 126,
		153, 212, 129, 121, 37, 11, 144, 186, 77, 65, 163, 230, 20, 116, 215,
		218, 128, 30, 71, 103, 144, 70, 233, 10, 198, 134, 65, 230, 236, 55,
		0, 0, 255, 255, 131, 33, 40, 6, 168, 148, 4, 108, 0, 0, 0, 0, 73, 69,
		78, 68, 174, 66, 96, 130}
	if !slices.Equal(got, expPNG) {
		t.Error("unexpected PNG")
		// t.Error((data))
	}

	bm := qr.Bitmap()
	expBM := [][]bool{
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b1, b1, b1, b1, b0, b0, b0, b1, b0, b0, b0, b1, b1, b0, b0, b1, b1, b1, b1, b1, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b0, b0, b0, b1, b0, b0, b1, b0, b0, b0, b0, b1, b1, b0, b0, b1, b0, b0, b0, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b1, b0, b0, b1, b0, b0, b0, b1, b1, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b1, b0, b0, b1, b1, b0, b1, b1, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b1, b0, b0, b1, b1, b1, b1, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b0, b0, b0, b1, b0, b0, b0, b1, b1, b1, b0, b1, b1, b1, b0, b1, b0, b0, b0, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b1, b1, b1, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b1, b1, b1, b1, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b1, b0, b0, b1, b0, b0, b0, b0, b1, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b0, b1, b1, b1, b1, b1, b0, b0, b0, b1, b0, b0, b1, b1, b1, b1, b0, b0, b0, b1, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b1, b1, b0, b0, b1, b1, b1, b0, b1, b1, b1, b0, b0, b1, b1, b1, b0, b0, b0, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b1, b1, b0, b1, b0, b1, b0, b1, b1, b1, b0, b1, b0, b1, b0, b1, b0, b1, b0, b1, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b1, b1, b0, b0, b0, b0, b0, b0, b1, b1, b0, b1, b1, b1, b0, b1, b0, b0, b1, b1, b0, b0, b1, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b1, b0, b0, b1, b1, b1, b0, b1, b1, b0, b0, b1, b1, b1, b1, b1, b1, b0, b1, b0, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b1, b0, b0, b1, b1, b1, b1, b1, b0, b0, b1, b0, b0, b1, b1, b0, b0, b1, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b1, b1, b1, b1, b1, b0, b1, b1, b0, b0, b1, b0, b0, b0, b1, b0, b1, b1, b0, b1, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b1, b0, b1, b1, b0, b0, b1, b1, b0, b1, b1, b0, b0, b0, b1, b1, b0, b0, b0, b0, b1, b0, b1, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b0, b1, b1, b1, b0, b0, b0, b0, b1, b0, b0, b0, b1, b1, b1, b1, b1, b1, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b1, b1, b1, b1, b1, b1, b0, b1, b1, b0, b0, b0, b1, b1, b1, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b1, b1, b1, b1, b0, b1, b1, b0, b1, b1, b0, b1, b1, b1, b0, b1, b0, b1, b0, b0, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b0, b0, b0, b1, b0, b1, b0, b0, b0, b1, b1, b1, b0, b1, b0, b0, b0, b1, b1, b0, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b1, b0, b1, b1, b0, b0, b1, b0, b1, b1, b1, b1, b1, b0, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b0, b1, b1, b1, b0, b0, b1, b1, b0, b0, b0, b1, b1, b0, b1, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b1, b1, b1, b0, b1, b0, b1, b1, b0, b0, b0, b1, b0, b1, b0, b1, b0, b1, b1, b1, b0, b0, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b0, b0, b0, b0, b0, b1, b0, b1, b0, b1, b1, b0, b0, b0, b1, b1, b1, b0, b0, b1, b1, b0, b1, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b1, b1, b1, b1, b1, b1, b1, b0, b1, b0, b0, b0, b1, b0, b0, b1, b1, b0, b1, b1, b0, b0, b0, b1, b1, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
		{b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0, b0},
	}
	if !slices.EqualFunc(bm, expBM, func(l1, l2 []bool) bool { return slices.Equal(l1, l2) }) {
		t.Error("unexpected bitmap")
		// t.Error(bm)
	}
}
