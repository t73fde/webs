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

package extkey_test

import (
	"testing"

	"t73f.de/r/webs/extkey"
)

func BenchmarkSnowflake(b *testing.B) {
	var generator extkey.Generator

	for i := 0; i < b.N; i++ {
		generator.Create(0)
	}
}
func BenchmarkSnowflakeX(b *testing.B) {
	bits := 7
	generator := extkey.NewGenerator(uint(bits))
	key := uint((1 << bits) - 1)
	for i := 0; i < b.N; i++ {
		generator.Create(key)
	}
}
