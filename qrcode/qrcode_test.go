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
		84, 69, 255, 255, 255, 0, 0, 0, 85, 194, 211, 126, 0, 0, 0, 156, 73,
		68, 65, 84, 120, 218, 116, 206, 175, 170, 2, 65, 24, 64, 241, 115, 195,
		93, 190, 101, 192, 175, 46, 98, 16, 20, 17, 12, 162, 70, 131, 73, 48,
		90, 4, 31, 195, 199, 176, 105, 20, 22, 255, 4, 65, 116, 124, 128, 101,
		147, 209, 36, 198, 197, 104, 26, 5, 139, 65, 38, 10, 246, 61, 225, 228,
		31, 121, 233, 189, 227, 65, 170, 109, 1, 179, 47, 58, 48, 219, 219,
		239, 111, 7, 82, 73, 5, 244, 105, 61, 16, 10, 80, 240, 171, 50, 252,
		45, 47, 2, 161, 61, 36, 16, 4, 143, 57, 104, 183, 55, 6, 214, 181, 1,
		132, 159, 126, 2, 255, 113, 52, 4, 51, 153, 157, 0, 205, 60, 232, 43,
		93, 128, 212, 93, 6, 198, 54, 167, 96, 118, 173, 13, 152, 227, 232, 12,
		210, 40, 93, 65, 93, 28, 229, 98, 191, 3, 0, 131, 33, 40, 6, 21, 27,
		225, 66, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130,
	}
	if !slices.Equal(got, expPNG) {
		t.Error("unexpected PNG")
		// t.Error(got)
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
	if !slices.EqualFunc(bm, expBM, slices.Equal) {
		t.Error("unexpected bitmap")
		// t.Error(bm)
	}
}
