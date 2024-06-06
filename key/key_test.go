// -----------------------------------------------------------------------------
// Copyright (c) 2023-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2023-present Detlef Stern
// -----------------------------------------------------------------------------

package key_test

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"t73f.de/r/webs/key"
)

func TestKeyString(t *testing.T) {
	var testcases = []struct {
		key key.Key
		exp string
	}{
		{0, "0"},
		{1, "1"},
		{0xffffffffffffffff, "FZZZZZZZZZZZZ"},
	}
	for _, tc := range testcases {
		t.Run(strconv.FormatUint(uint64(tc.key), 10), func(t *testing.T) {
			got := tc.key.String()
			if got != tc.exp {
				t.Errorf("%q expected, but got %q", tc.exp, got)
			}
			key, err := key.Parse(got)
			if err != nil {
				panic(err)
			}
			if key != tc.key {
				t.Errorf("key %d was printed as %q, but parsed as %d/%q", tc.key, got, key, key)
			}
		})
	}
}

func TestGenerator(t *testing.T) {
	var generator key.Generator
	var lastKey key.Key

	for i := 0; i < 1000000; i++ {
		k := generator.Make(0)
		if k <= lastKey {
			t.Errorf("key does not increase: %v -> %v", lastKey, k)
			return
		}
		lastKey = k
		checkParse(t, k)
	}
}

func checkParse(t *testing.T, k key.Key) {
	s := k.String()
	parsedKey, err := key.Parse(s)
	if err != nil {
		panic(err)
	}
	if parsedKey != k {
		t.Errorf("key %d/%q was parsed, but got %d/%v", k, s, parsedKey, parsedKey)
	}
}

func TestKeyID(t *testing.T) {
	for intBits := uint(0); intBits <= key.MaxAppBits; intBits++ {
		maxID := int32(1 << intBits)
		generator := key.NewGenerator(intBits)
		for i := 0; i < 512; i++ {
			exp := uint(rand.Int31n(maxID))
			k := generator.Make(exp)
			got := k.ID(intBits)
			if got != exp {
				t.Errorf("id of %v should be %d, but got %d", k, exp, got)
			}

			checkParse(t, k)
		}
	}
}

func TestKeyID2(t *testing.T) {
	var k key.Key
	if k.IsValid() {
		t.Errorf("key %v/%d is not invalid, but should be", k, k)
	}
	for intBits := uint(0); intBits <= key.MaxAppBits; intBits++ {
		got := k.ID(intBits)
		if got != 0 {
			t.Errorf("%v.ID() should be 0, but is %d", k, got)
		}
	}
}

func TestParseKey(t *testing.T) {
	var testcases = []struct {
		s   string
		r   int
		exp key.Key
	}{
		{"0000000000000", 0, 0},
		{"00-000-000-00-000", 0, 0},
		{"000-000-000-00-00", 0, 0},
		{"0-00-0-0-0-0-0-0-0-0-0-0", 0, 0},
		{"0000000000001", 0, 1},
		{"0E34NNFRTCQ15", 0, 507945423712181285},
		{"0DXZBE2D7TB04", 0, 502128752335858692},
		{"-0000000000000", 1, 0},
		{"0000000000000-", 1, 0},
		{"0DXZBE2D7<>04", 1, 0},
		{"1DXZBE2D7TB040", 2, 0},
		{"FZZZZZZZZZZZZ", 0, math.MaxUint64},
		{"F-zz-ZZZZZZZZ-zz", 0, math.MaxUint64},
	}

	for _, tc := range testcases {
		t.Run(tc.s, func(t *testing.T) {
			got, err := key.Parse(tc.s)
			if err != nil {
				switch tc.r {
				case 0:
					t.Errorf("error %v returned, but none expected", err)
				case 1:
					if !strings.HasPrefix(err.Error(), "non base-32 character ") {
						t.Errorf("error 'non base-32 character' expected, but got: %v", err)
					}
				case 2:
					if !strings.HasPrefix(err.Error(), "does not fit in uint64: \"") {
						t.Errorf("error 'string does not fit' expected, but got: %v", err)
					}
				default:
					t.Errorf("unknown result code %d in test case", tc.r)
				}
				return
			}
			if tc.r != 0 {
				t.Error("error expected, but got value:", got)
				return
			}
			if got != tc.exp {
				t.Errorf("external key %v/%d expected, but got %v/%d", tc.exp, tc.exp, got, got)
				return
			}
			checkParse(t, got)
		})
	}
}
