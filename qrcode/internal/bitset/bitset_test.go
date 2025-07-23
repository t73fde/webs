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

package bitset

import (
	"math/rand"
	"testing"
)

const (
	x0 = false
	x1 = true
)

func TestNewBitset(t *testing.T) {
	tests := [][]bool{
		{},
		{x1},
		{x0},
		{x1, x0},
		{x1, x0, x1},
		{x0, x0, x1},
	}

	for _, v := range tests {
		result := New(v...)

		if !equal(result.Bits(), v) {
			t.Errorf("%s", result.String())
			t.Errorf("%v => %v, want %v", v, result.Bits(), v)
		}
	}
}

func TestAppend(t *testing.T) {
	randomBools := make([]bool, 128)
	rng := rand.New(rand.NewSource(1))
	for i := range randomBools {
		randomBools[i] = rng.Intn(2) == 1
	}

	for i := 0; i < len(randomBools)-1; i++ {
		a := New(randomBools[0:i]...)
		b := New(randomBools[i:]...)

		a.Append(b)
		if !equal(a.Bits(), randomBools) {
			t.Errorf("got %v, want %v", a.Bits(), randomBools)
		}
	}
}

func TestAppendByte(t *testing.T) {
	tests := []struct {
		initial  *Bitset
		value    byte
		numBits  int
		expected *Bitset
	}{
		{New(), 0x01, 1, New(x1)},
		{New(x1), 0x01, 1, New(x1, x1)},
		{New(x0), 0x01, 1, New(x0, x1)},
		{
			New(x1, x0, x1, x0, x1, x0, x1),
			0xAA, // 0b10101010
			2,
			New(x1, x0, x1, x0, x1, x0, x1, x1, x0),
		},
		{
			New(x1, x0, x1, x0, x1, x0, x1),
			0xAA, // 0b10101010
			8,
			New(x1, x0, x1, x0, x1, x0, x1, x1, x0, x1, x0, x1, x0, x1, x0),
		},
	}

	for _, test := range tests {
		test.initial.AppendByte(test.value, test.numBits)
		if !equal(test.initial.Bits(), test.expected.Bits()) {
			t.Errorf("Got %v, expected %v", test.initial.Bits(),
				test.expected.Bits())
		}
	}
}

func TestAppendUint32(t *testing.T) {
	tests := []struct {
		initial  *Bitset
		value    uint32
		numBits  int
		expected *Bitset
	}{
		{
			New(),
			0xAAAAAAAF,
			4,
			New(x1, x1, x1, x1),
		},
		{
			New(),
			0xFFFFFFFF,
			32,
			New(x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1,
				x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1, x1),
		},
		{
			New(),
			0x0,
			32,
			New(x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0,
				x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0, x0),
		},
		{
			New(),
			0xAAAAAAAA,
			32,
			New(x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1,
				x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0),
		},
		{
			New(),
			0xAAAAAAAA,
			31,
			New(x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1,
				x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0, x1, x0),
		},
	}

	for _, test := range tests {
		test.initial.AppendUint32(test.value, test.numBits)
		if !equal(test.initial.Bits(), test.expected.Bits()) {
			t.Errorf("Got %v, expected %v", test.initial.Bits(),
				test.expected.Bits())
		}
	}
}

func TestAppendBools(t *testing.T) {
	randomBools := make([]bool, 128)
	rng := rand.New(rand.NewSource(1))
	for i := range randomBools {
		randomBools[i] = rng.Intn(2) == 1
	}

	for i := 0; i < len(randomBools)-1; i++ {
		got := New(randomBools[0:i]...)
		got.AppendBools(randomBools[i:]...)

		if !equal(got.Bits(), randomBools) {
			t.Errorf("got %v, want %v", got.Bits(), randomBools)
		}
	}
}

func BenchmarkShortAppend(b *testing.B) {
	bitset := New()
	for b.Loop() {
		bitset.AppendBools(x0, x1, x0, x1, x0, x1, x0)
	}
}

func TestLen(t *testing.T) {
	randomBools := make([]bool, 128)
	rng := rand.New(rand.NewSource(1))
	for i := range randomBools {
		randomBools[i] = rng.Intn(2) == 1
	}

	for i := 0; i < len(randomBools)-1; i++ {
		if got := New(randomBools[0:i]...); got.Len() != i {
			t.Errorf("Len = %d, want %d", got.Len(), i)
		}
	}
}

func TestAt(t *testing.T) {
	test := []bool{x0, x1, x0, x1, x0, x1, x1, x0, x1}

	bitset := New(test...)
	for i, v := range test {
		if got := bitset.At(i); got != test[i] {
			t.Errorf("bitset[%d] => %t, want %t", i, got, v)
		}
	}
}

func equal(a []bool, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestExample(t *testing.T) {
	b := New()                       // {}
	b.AppendBools(true, true, false) // {1, 1, 0}
	b.AppendBools(true)              // {1, 1, 0, 1}
	b.AppendByte(0x02, 4)            // {1, 1, 0, 1, 0, 0, 1, 0}

	expected := []bool{x1, x1, x0, x1, x0, x0, x1, x0}
	if !equal(b.Bits(), expected) {
		t.Errorf("Got %v, expected %v", b.Bits(), expected)
	}
}

func TestByteAt(t *testing.T) {
	data := []bool{x0, x1, x0, x1, x0, x1, x1, x0, x1}

	tests := []struct {
		index    int
		expected byte
	}{
		{0, 0x56},
		{1, 0xad},
		{2, 0x2d},
		{5, 0x0d},
		{8, 0x01},
	}

	for _, test := range tests {
		b := New()
		b.AppendBools(data...)

		if got := b.ByteAt(test.index); got != test.expected {
			t.Errorf("Got %#x, expected %#x", got, test.expected)
		}
	}
}

func TestSubstr(t *testing.T) {
	data := []bool{x0, x1, x0, x1, x0, x1, x1, x0}

	tests := []struct {
		start    int
		end      int
		expected []bool
	}{
		{0, 8, []bool{x0, x1, x0, x1, x0, x1, x1, x0}},
		{0, 0, []bool{}},
		{0, 1, []bool{x0}},
		{2, 4, []bool{x0, x1}},
	}

	for _, test := range tests {
		b := New()
		b.AppendBools(data...)

		expected := New()
		expected.AppendBools(test.expected...)

		if got := b.Substr(test.start, test.end); !got.Equals(expected) {
			t.Errorf("Got %s, expected %s", got.String(), expected.String())
		}
	}
}
