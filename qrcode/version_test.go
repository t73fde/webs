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

func TestFormatInfo(t *testing.T) {
	tests := []struct {
		level       RecoveryLevel
		maskPattern int

		expected uint32
	}{ // L=01 M=00 Q=11 H=10
		{Low, 1, 0x72f3},
		{Medium, 2, 0x5e7c},
		{High, 3, 0x3a06},
		{Highest, 4, 0x0762},
		{Low, 5, 0x6318},
		{Medium, 6, 0x4f97},
		{High, 7, 0x2bed},
	}

	for i, test := range tests {
		expected := bitset.New()
		expected.AppendUint32(test.expected, formatInfoLengthBits)

		v := getQRCodeVersion(test.level, 1)
		got := v.formatInfo(test.maskPattern)
		if !expected.Equals(got) {
			t.Errorf("formatInfo test #%d got %s, expected %s", i, got.String(),
				expected.String())
		}
	}
}

func TestVersionInfo(t *testing.T) {
	tests := []struct {
		version  int
		expected uint32
	}{
		{7, 0x007c94},
		{10, 0x00a4d3},
		{20, 0x0149a6},
		{30, 0x01ed75},
		{40, 0x028c69},
	}

	for i, test := range tests {
		expected := bitset.New()
		expected.AppendUint32(test.expected, versionInfoLengthBits)

		v := getQRCodeVersion(Low, test.version)
		got := v.versionInfo()
		if !expected.Equals(got) {
			t.Errorf("versionInfo test #%d got %s, expected %s", i, got.String(),
				expected.String())
		}
	}
}

func TestNumBitsToPadToCodeoword(t *testing.T) {
	tests := []struct {
		level       RecoveryLevel
		version     int
		numDataBits int
		expected    int
	}{
		{Low, 1, 0, 0},
		{Low, 1, 1, 7},
		{Low, 1, 7, 1},
		{Low, 1, 8, 0},
	}

	for i, test := range tests {
		v := getQRCodeVersion(test.level, test.version)

		got := v.numBitsToPadToCodeword(test.numDataBits)
		if got != test.expected {
			t.Errorf("numBitsToPadToCodeword test %d (version=%d numDataBits=%d), got %d, expected %d",
				i, test.version, test.numDataBits, got, test.expected)
		}
	}
}

// getQRCodeVersion returns the QR Code version by version number and recovery
// level. Returns nil if the requested combination is not defined.
func getQRCodeVersion(level RecoveryLevel, version int) *qrCodeVersion {
	for _, v := range versions {
		if v.level == level && v.version == version {
			return &v
		}
	}

	return nil
}
