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
	"fmt"
	"reflect"
	"testing"

	"t73f.de/r/webs/qrcode/internal/bitset"
)

func TestClassifyDataMode(t *testing.T) {
	tests := []struct {
		data   []byte
		actual []segment
	}{
		{
			[]byte{0x30},
			[]segment{
				{dataModeNumeric, []byte{0x30}},
			},
		},
		{
			[]byte{0x30, 0x41, 0x42, 0x43, 0x20, 0x00, 0xf0, 0xf1, 0xf2, 0x31},
			[]segment{
				{dataModeNumeric, []byte{0x30}},
				{dataModeAlphanumeric, []byte{0x41, 0x42, 0x43, 0x20}},
				{dataModeByte, []byte{0x00, 0xf0, 0xf1, 0xf2}},
				{dataModeNumeric, []byte{0x31}},
			},
		},
	}

	for _, test := range tests {
		encoder := allDataEncoder[0] //(dataEncoderType1To9)
		_, err := encoder.encode(test.data)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(test.actual, encoder.actual) {
			t.Errorf("Got %v, expected %v", encoder.actual, test.actual)
		}
	}
}

func TestByteModeLengthCalculations(t *testing.T) {
	tests := []struct {
		dataEncoderType int
		dataMode        dataMode
		numSymbols      int
		expectedLength  int
	}{}

	for i, test := range tests {
		encoder := allDataEncoder[test.dataEncoderType]
		var resultLength int

		resultLength, err := encoder.encodedLength(test.dataMode, test.numSymbols)

		if test.expectedLength == -1 {
			if err == nil {
				t.Errorf("Test %d: got length %d, expected error", i, resultLength)
			}
		} else if resultLength != test.expectedLength {
			t.Errorf("Test %d: got length %d, expected length %d", i, resultLength,
				test.expectedLength)
		}
	}
}

func TestSingleModeEncodings(t *testing.T) {
	tests := []struct {
		dataEncoderType int
		dataMode        dataMode
		data            string
		expected        *bitset.Bitset
	}{
		{
			0, // dataEncoderType1To9,
			dataModeNumeric, "01234567",
			bitset.NewFromBase2String("0001 0000001000 0000001100 0101011001 1000011"),
		},
		{
			0, // dataEncoderType1To9,
			dataModeAlphanumeric, "AC-42",
			bitset.NewFromBase2String("0010 000000101 00111001110 11100111001 000010"),
		},
		{
			0, // dataEncoderType1To9,
			dataModeByte, "123",
			bitset.NewFromBase2String("0100 00000011 00110001 00110010 00110011"),
		},
		{
			1, // dataEncoderType10To26,
			dataModeByte, "123",
			bitset.NewFromBase2String("0100 00000000 00000011 00110001 00110010 00110011"),
		},
		{
			1, // dataEncoderType27To40,
			dataModeByte, "123",
			bitset.NewFromBase2String("0100 00000000 00000011 00110001 00110010 00110011"),
		},
	}

	for _, test := range tests {
		encoder := allDataEncoder[test.dataEncoderType]
		encoded := bitset.New()

		encoder.encodeDataRaw([]byte(test.data), test.dataMode, encoded)

		if !test.expected.Equals(encoded) {
			t.Errorf("For %s got %s, expected %s", test.data, encoded.String(),
				test.expected.String())
		}
	}
}

type testModeSegment struct {
	dataMode dataMode
	numChars int
}

func TestOptimiseEncoding(t *testing.T) {
	tests := []struct {
		dataEncoderType int
		actual          []testModeSegment
		optimised       []testModeSegment
	}{
		// Coalescing multiple segments.
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1}, // length = 4 + 9 + 6 = 19 bits
				{dataModeNumeric, 1},      // length = 4 + 10 + 4 = 18 bits
				{dataModeAlphanumeric, 1}, // 19 bits.
				{dataModeNumeric, 1},      // 18 bits.
				{dataModeAlphanumeric, 1}, // 19 bits.
				// total = 93 bits.
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 5}, // length = 4 + 9 + 22 + 6 = 41.
			},
		},
		// Coalesing not necessary.
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 20},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 20},
			},
		},
		// Switch to more general dataMode.
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 100},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 100},
				{dataModeByte, 2},
			},
		},
		// Sometimes encoding everything as bytes is more efficient.
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeByte, 3},
			},
		},
		// https://www.google.com/123456789012345678901234567890
		// BBBBBAAABBBABBBBBBABBBANNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeByte, 5},
				{dataModeAlphanumeric, 3},
				{dataModeByte, 3},
				{dataModeAlphanumeric, 1},
				{dataModeByte, 6},
				{dataModeAlphanumeric, 1},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 30},
			},
			[]testModeSegment{
				{dataModeByte, 23},
				{dataModeNumeric, 30},
			},
		},
		// https://www.google.com/123
		// BBBBBAAABBBABBBBBBABBBANNN
		// Small segments are inefficient because of additional metadata.
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeByte, 5},
				{dataModeAlphanumeric, 3},
				{dataModeByte, 3},
				{dataModeAlphanumeric, 1},
				{dataModeByte, 6},
				{dataModeAlphanumeric, 1},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 3},
			},
			[]testModeSegment{
				{dataModeByte, 26},
			},
		},
		// HTTPS://WWW.GOOGLE.COM/123
		// AAAAAAAAAAAAAAAAAAAAAAANNN
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 23},
				{dataModeNumeric, 3},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 26},
			},
		},
		{
			2, // dataEncoderType27To40,
			[]testModeSegment{
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
				{dataModeByte, 1},
				{dataModeNumeric, 1},
			},
			[]testModeSegment{
				{dataModeByte, 8},
			},
		},
		// HTTPS://ABC.DE/Q/393AABB6998877XYZ0518AUQCRVJN25
		// AAAAAAAAAAAAAAAAANNNAAAANNNNNNNAAANNNNAAAAAAAANN
		// different to below---------^--------------------
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 17},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 7},
				{dataModeAlphanumeric, 3},
				{dataModeNumeric, 4},
				{dataModeAlphanumeric, 8},
				{dataModeNumeric, 2},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 48},
			},
		},
		// HTTPS://ABC.DE/Q/393AABB699E877XYZ0518AUQCRVJN25
		// AAAAAAAAAAAAAAAAANNNAAAANNNANNNAAANNNNAAAAAAAANN
		// different to above---------^--------------------
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeAlphanumeric, 17},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 4},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 1},
				{dataModeNumeric, 3},
				{dataModeAlphanumeric, 3},
				{dataModeNumeric, 4},
				{dataModeAlphanumeric, 8},
				{dataModeNumeric, 2},
			},
			[]testModeSegment{
				{dataModeAlphanumeric, 48},
			},
		},
		// 0123456789
		// NNNNNNNNNN
		{
			0, // dataEncoderType1To9,
			[]testModeSegment{
				{dataModeNumeric, 10},
			},
			[]testModeSegment{
				{dataModeNumeric, 10},
			},
		},
	}

	for _, test := range tests {
		numTotalChars := 0
		for _, v := range test.actual {
			numTotalChars += v.numChars
		}

		data := make([]byte, numTotalChars)

		i := 0
		for _, v := range test.actual {
			for j := 0; j < v.numChars; j++ {
				switch v.dataMode {
				case dataModeNumeric:
					data[i] = '1'
				case dataModeAlphanumeric:
					data[i] = 'A'
				case dataModeByte:
					data[i] = '#'
				default:
					t.Fatal("Unrecognised data mode")
				}

				i++
			}
		}

		encoder := allDataEncoder[test.dataEncoderType]

		_, err := encoder.encode(data)

		if err != nil {
			t.Errorf("Got %s, expected valid encoding", err.Error())
		} else {
			ok := true

			if len(encoder.optimised) != len(test.optimised) {
				ok = false
			} else {
				for i, s := range test.optimised {
					if encoder.optimised[i].dataMode != s.dataMode ||
						len(encoder.optimised[i].data) != s.numChars {
						ok = false
						break
					}
				}
			}

			if !ok {
				t.Errorf("got %s, expected %s", segmentsString(encoder.optimised),
					testModeSegmentsString(test.optimised))
			}
		}
	}
}

func testModeSegmentsString(segments []testModeSegment) string {
	result := "["
	for i, segment := range segments {
		if i > 0 {
			result += ", "
		}

		result += fmt.Sprintf("%d*%s", segment.numChars,
			dataModeString(segment.dataMode))
	}
	result += "]"
	return result
}

func segmentsString(segments []segment) string {
	result := "["
	for i, segment := range segments {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d*%s", len(segment.data),
			dataModeString(segment.dataMode))
	}
	result += "]"
	return result
}
func dataModeString(d dataMode) string {
	switch d {
	case dataModeNone:
		return "none"
	case dataModeNumeric:
		return "numeric"
	case dataModeAlphanumeric:
		return "alphanumeric"
	case dataModeByte:
		return "byte"
	}
	return "unknown"
}
