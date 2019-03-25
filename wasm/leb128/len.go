// Copyright 2018 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leb128

// Uleb128Len calculates the length of v using unsigned LEB128 encoding.
func Uleb128Len(v uint64) int {
	var vlen int
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}

		vlen++
		if c&0x80 == 0 {
			break
		}
	}
	return vlen
}

// Sleb128Len calculates the length of v using signed LEB128 encoding.
func Sleb128Len(v int64) int {
	var vlen int
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		vlen++
		if c&0x80 == 0 {
			break
		}
	}
	return vlen
}
