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

// Package key provides a generic key to be used both as an URI element and as a
// primary key in a database.
package key

import (
	"fmt"
	"sync"
	"time"
)

// Key is the generic, external primary key for all data.
//
// Snowflake/TSID:
// * 42 bit timestamp, enough to be used in the year 2160.
// * 22 bit application / sequence number
//   - 0-20 bit application defined data, e.g. for tables, nodes, ...
//   - 2-22 bit sequence number
//
// The timestamp value starts at 2024-06-01. 42 bits of milliseconds allow an
// end date in the year 2163.
//
// The rest of the 64 bits, ie 22 bits can be splitted at the users demand.
// Two bits are always reserved to be used as a sequence number if two keys are
// generated within the same millisecond. This can be enlarged up 22 bits.
// A maximum of 20 bits can be used for the application. The more bits are
// used, the less bits are available for the sequence number. An application
// can use the bits to store the number of a database table, or the number of a
// computing node.
type Key uint64

// Invalid is the default invalid key.
const Invalid Key = 0

const (
	timestampBits = 42
	randomBits    = 22

	maxTimeStamp = 1<<timestampBits - 1
)

// MaxAppBits states the maximum number of bits reserved for the application
// defined part of the key.
const MaxAppBits = 20

// Parse will parse a string into an external key.
func Parse(s string) (Key, error) {
	result := Key(0)
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '-' && i > 0 && i < len(s)-1 {
			continue
		}
		if '0' <= ch && ch <= 128 {
			val := decode32map[ch-'0']
			if 0 <= val && val <= 31 {
				if result&0xF800000000000000 != 0 {
					return Invalid, fmt.Errorf("does not fit in uint64: %q / %x", s, uint64(result))
				}
				result = (result << 5) | Key(val)
				continue
			}
		}
		return result, fmt.Errorf("non base-32 character %c/%v found", ch, ch)
	}
	return result, nil
}

// IsInvalid returns true if the key is definitely an invalid key.
func (extkey Key) IsInvalid() bool { return extkey == Invalid }

// IsValid returns true if the key is definitely an invalid key.
func (extkey Key) IsValid() bool { return extkey != Invalid }

var decode32map = [...]int8{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -1, -1, -1, -1, -1, // 0x30 .. 0x3f
	-1, 10, 11, 12, 13, 14, 15, 16, 17, 1, 18, 19, 1, 20, 21, 0, // 0x40 .. 0x4f
	22, 23, 24, 25, 26, 36, 27, 28, 29, 30, 31, -1, -1, -1, -1, -1, // 0x50 .. 0x5f
	-1, 10, 11, 12, 13, 14, 15, 16, 17, 1, 18, 19, 1, 20, 21, 0, // 0x60 .. 0x6f
	22, 23, 24, 25, 26, -1, 27, 28, 29, 30, 31, -1, -1, -1, -1, -1, // 0x70 .. 0x7f
}

// Time returns the timestamp value, when the key was generated.
func (pk Key) Time() time.Time {
	return time.UnixMilli(int64(pk>>randomBits) + epochAdjust)
}

// ID returns the application defined part of the key.
func (pk Key) ID(appBits uint) uint {
	if appBits == 0 || appBits > MaxAppBits {
		return 0
	}
	return uint((pk & 0x3fffff) >> (randomBits - appBits))
}

// String returns a base-32 representation of the key as a string.
// It contains at most 13 characters.
func (pk Key) String() string {
	if pk == 0 {
		return "0"
	}
	u64 := uint64(pk)
	temp := [13]byte{}
	tpos := 0
	for u64 > 0 {
		temp[tpos] = base32chars[u64%32]
		tpos++
		u64 = u64 >> 5
	}
	var result [13]byte
	for i := 0; i < tpos; i++ {
		result[i] = temp[tpos-i-1]
	}
	return string(result[:tpos])
}

const base32chars = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// Generator is a generator for unique keys as int64.
type Generator struct {
	mx      sync.Mutex // Protects the next two fields
	lastTS  uint64     // Last timestamp
	nextSeq uint64     // Next sequence number for lastTS
	appBits uint       // number of bits for application use. range: 0-MaxAppBits
	appMax  uint       // 1 << appBits (if appBits > 0; else: 0)
}

// NewGenerator creates a new key generator with a given number of bits for
// application use.
func NewGenerator(appBits uint) Generator {
	if appBits > MaxAppBits {
		panic(fmt.Sprintf("key generator need too many bits (max %d): %v", appBits, MaxAppBits))
	}
	return Generator{
		appBits: appBits,
		appMax:  1 << appBits,
	}
}

// epochAdjust is used to make the timestamp values smaller, so they better fit
// in 42 bits.
//
// Its value is time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
const epochAdjust = 1717200000000

// Make builds a new Key.
func (kg *Generator) Make(appId uint) Key {
	if appId > 0 && appId >= kg.appMax {
		panic(fmt.Errorf("application value out of range: %v (max: %v)", appId, kg.appMax))
	}
	for {
		milli := uint64(time.Now().UnixMilli())
		var seq uint64

		kg.mx.Lock()
		if milli > kg.lastTS {
			kg.lastTS = milli
			kg.nextSeq = 1
			seq = 0
		} else {
			seq = kg.nextSeq
			kg.nextSeq++
		}
		kg.mx.Unlock()

		if seq < (1 << (randomBits - kg.appBits)) {
			ts := milli - epochAdjust
			if ts > maxTimeStamp {
				panic(fmt.Sprintf("timestamp %v exceeds largest possible value %v", ts, maxTimeStamp))
			}

			// 42bit=ts, kg.intBits=appId, 22-kg.intBits=seq
			k := (ts << randomBits) | (uint64(appId) << (randomBits - kg.appBits)) | seq
			return Key(k)
		}

		time.Sleep(1 * time.Millisecond)
	}
}
