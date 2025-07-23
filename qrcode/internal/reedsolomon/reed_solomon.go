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

// Package reedsolomon provides error correction encoding for QR Code 2005.
//
// QR Code 2005 uses a Reed-Solomon error correcting code to detect and correct
// errors encountered during decoding.
//
// The generated RS codes are systematic, and consist of the input data with
// error correction bytes appended.
package reedsolomon

import "t73f.de/r/webs/qrcode/internal/bitset"

// Encode data for QR Code 2005 using the appropriate Reed-Solomon code.
//
// numECBytes is the number of error correction bytes to append, and is
// determined by the target QR Code's version and error correction level.
//
// ISO/IEC 18004 table 9 specifies the numECBytes required. e.g. a 1-L code has
// numECBytes=7.
func Encode(data *bitset.Bitset, numECBytes int) *bitset.Bitset {
	// Create a polynomial representing |data|.
	//
	// The bytes are interpreted as the sequence of coefficients of a polynomial.
	// The last byte's value becomes the x^0 coefficient, the second to last
	// becomes the x^1 coefficient and so on.
	ecpoly := newGFPolyFromData(data)
	ecpoly = gfPolyMultiply(ecpoly, newGFPolyMonomial(gfOne, numECBytes))

	generator := rsGeneratorPoly(numECBytes)        // pick generator polynomial
	remainder := gfPolyRemainder(ecpoly, generator) // generate error correction bytes

	// Combine the data & error correcting bytes.
	// The mathematically correct answer is:
	//
	//	result := gfPolyAdd(ecpoly, remainder).
	//
	// The encoding used by QR Code 2005 is slightly different this result: To
	// preserve the original |data| bit sequence exactly, the data and remainder
	// are combined manually below. This ensures any most significant zero bits
	// are preserved (and not optimised away).
	result := bitset.Clone(data)
	result.AppendBytes(remainder.data(numECBytes))
	return result
}

// rsGeneratorPoly returns the Reed-Solomon generator polynomial with |degree|.
//
// The generator polynomial is calculated as:
// (x + a^0)(x + a^1)...(x + a^degree-1)
func rsGeneratorPoly(degree int) gfPoly {
	if degree < 2 {
		panic("degree < 2")
	}

	generator := gfPoly{term: []gfElement{1}}
	for i := range degree {
		nextPoly := gfPoly{term: []gfElement{gfExpTable[i], 1}}
		generator = gfPolyMultiply(generator, nextPoly)
	}
	return generator
}
