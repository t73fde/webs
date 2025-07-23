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

// Package qrcode allows to encode data as QR codes.
package qrcode

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"t73f.de/r/webs/qrcode/internal/bitset"
	"t73f.de/r/webs/qrcode/internal/reedsolomon"
)

// A QRCode represents a valid encoded QRCode.
type QRCode struct {
	content string // original content

	// QR Code type.
	recoveryLevel RecoveryLevel
	VersionNumber int

	// User settable drawing options.
	ForegroundColor color.Color
	BackgroundColor color.Color

	// Disable the QR Code border.
	DisableBorder bool

	encoder *dataEncoder
	version qrCodeVersion

	data   *bitset.Bitset
	symbol *symbol
	mask   int
}

// New constructs a QRCode.
//
// An error occurs if the content is too long.
func New(content string, level RecoveryLevel) (*QRCode, error) {
	var encoder *dataEncoder
	var encoded *bitset.Bitset
	var chosenVersion *qrCodeVersion
	var err error

	for i := range allDataEncoder {
		de := allDataEncoder[i] // we need a fresh copy
		encoder = &de

		encoded, err = encoder.encode([]byte(content))
		if err != nil {
			continue
		}

		chosenVersion = chooseQRCodeVersion(level, encoder, encoded.Len())
		if chosenVersion != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}
	if chosenVersion == nil {
		return nil, errors.New("content too long to encode")
	}

	q := &QRCode{
		content: content,

		recoveryLevel: level,
		VersionNumber: chosenVersion.version,

		ForegroundColor: color.Black,
		BackgroundColor: color.White,
		DisableBorder:   false,

		encoder: encoder,
		data:    encoded,
		version: *chosenVersion,
	}
	return q, nil
}

// Bitmap returns the QR Code as a 2D array of 1-bit pixels.
//
// bitmap[y][x] is true if the pixel at (x, y) is set.
//
// The bitmap includes the required "quiet zone" around the QR Code to aid
// decoding.
func (q *QRCode) Bitmap() [][]bool {
	q.encode()
	return q.symbol.bitmap()
}

// Image returns the QR Code as an image.Image.
//
// A positive size sets a fixed image width and height (e.g. 256 yields an
// 256x256px image).
//
// Depending on the amount of data encoded, fixed size images can have different
// amounts of padding (white space around the QR Code). As an alternative, a
// variable sized image can be generated instead:
//
// A negative size causes a variable sized image to be returned. The image
// returned is the minimum size required for the QR Code. Choose a larger
// negative number to increase the scale of the image. e.g. a size of -5 causes
// each module (QR Code "pixel") to be 5px in size.
func (q *QRCode) Image(size int) image.Image {
	q.encode()

	// Minimum pixels (both width and height) required.
	realSize := q.symbol.fullSize

	// Variable size support.
	if size < 0 {
		size = size * -1 * realSize
	}

	// Actual pixels available to draw the symbol. Automatically increase the
	// image size if it's not large enough.
	if size < realSize {
		size = realSize
	}

	// Output image.
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{size, size}}

	// Saves a few bytes to have them in this order
	p := color.Palette([]color.Color{q.BackgroundColor, q.ForegroundColor})
	img := image.NewPaletted(rect, p)
	fgClr := uint8(img.Palette.Index(q.ForegroundColor))

	bitmap := q.symbol.bitmap()

	// Map each image pixel to the nearest QR code module.
	modulesPerPixel := float64(realSize) / float64(size)
	for y := 0; y < size; y++ {
		y2 := int(float64(y) * modulesPerPixel)
		for x := 0; x < size; x++ {
			x2 := int(float64(x) * modulesPerPixel)
			if bitmap[y2][x2] {
				pos := img.PixOffset(x, y)
				img.Pix[pos] = fgClr
			}
		}
	}
	return img
}

// PNG returns the QR Code as a PNG image.
//
// size is both the image width and height in pixels. If size is too small then
// a larger image is silently returned. Negative values for size cause a
// variable sized image to be returned: See the documentation for Image().
func (q *QRCode) PNG(size int) ([]byte, error) {
	img := q.Image(size)

	var b bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&b, img); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// encode completes the steps required to encode the QR Code. These include
// adding the terminator bits and padding, splitting the data into blocks and
// applying the error correction, and selecting the best data mask.
func (q *QRCode) encode() {
	numTerminatorBits := q.version.numTerminatorBitsRequired(q.data.Len())

	q.addTerminatorBits(numTerminatorBits)
	q.addPadding()

	encoded := q.encodeBlocks()

	const numMasks int = 8
	penalty := 0

	for mask := range numMasks {
		s := buildRegularSymbol(q.version, mask, encoded, !q.DisableBorder)

		numEmptyModules := s.numEmptyModules()
		if numEmptyModules != 0 {
			panic(fmt.Sprintf("BUG: numEmptyModules is %d (expected 0) (version=%d)",
				numEmptyModules, q.VersionNumber))
		}

		p := s.penaltyScore()
		if q.symbol == nil || p < penalty {
			q.symbol = s
			q.mask = mask
			penalty = p
		}
	}
}

// addTerminatorBits adds final terminator bits to the encoded data.
//
// The number of terminator bits required is determined when the QR Code version
// is chosen (which itself depends on the length of the data encoded). The
// terminator bits are thus added after the QR Code version
// is chosen, rather than at the data encoding stage.
func (q *QRCode) addTerminatorBits(numTerminatorBits int) {
	q.data.AppendNumBools(numTerminatorBits, false)
}

// encodeBlocks takes the completed (terminated & padded) encoded data, splits
// the data into blocks (as specified by the QR Code version), applies error
// correction to each block, then interleaves the blocks together.
//
// The QR Code's final data sequence is returned.
func (q *QRCode) encodeBlocks() *bitset.Bitset {
	// Split into blocks.
	type dataBlock struct {
		data          *bitset.Bitset
		ecStartOffset int
	}

	block := make([]dataBlock, q.version.numBlocks())

	start, end, blockID := 0, 0, 0
	for _, b := range q.version.block {
		for j := 0; j < b.numBlocks; j++ {
			start = end
			end = start + b.numDataCodewords*8

			// Apply error correction to each block.
			numErrorCodewords := b.numCodewords - b.numDataCodewords
			block[blockID].data = reedsolomon.Encode(q.data.Substr(start, end), numErrorCodewords)
			block[blockID].ecStartOffset = end - start

			blockID++
		}
	}

	// Interleave the blocks.

	result := bitset.New()

	// Combine data blocks.
	working := true
	for i := 0; working; i += 8 {
		working = false

		for j, b := range block {
			if i >= block[j].ecStartOffset {
				continue
			}

			result.Append(b.data.Substr(i, i+8))
			working = true
		}
	}

	// Combine error correction blocks.
	working = true
	for i := 0; working; i += 8 {
		working = false

		for j, b := range block {
			offset := i + block[j].ecStartOffset
			if offset >= block[j].data.Len() {
				continue
			}

			result.Append(b.data.Substr(offset, offset+8))
			working = true
		}
	}

	// Append remainder bits.
	result.AppendNumBools(q.version.numRemainderBits, false)
	return result
}

// addPadding pads the encoded data upto the full length required.
func (q *QRCode) addPadding() {
	numDataBits := q.version.numDataBits()
	if q.data.Len() == numDataBits {
		return
	}

	// Pad to the nearest codeword boundary.
	q.data.AppendNumBools(q.version.numBitsToPadToCodeword(q.data.Len()), false)

	// Pad codewords 0b11101100 and 0b00010001.
	padding := [2]*bitset.Bitset{
		bitset.New(b1, b1, b1, b0, b1, b1, b0, b0),
		bitset.New(b0, b0, b0, b1, b0, b0, b0, b1),
	}

	// Insert pad codewords alternately.
	i := 0
	for numDataBits-q.data.Len() >= 8 {
		q.data.Append(padding[i])

		i = 1 - i // Alternate between 0 and 1.
	}

	if q.data.Len() != numDataBits {
		panic(fmt.Sprintf("BUG: got len %d, expected %d", q.data.Len(), numDataBits))
	}
}
