// Copyright 2018 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/ray-lang/wagon/wasm/leb128"
)

const currentVersion = 0x01

// EncodeModule writes a provided module to w using WASM binary encoding.
func EncodeModule(w io.Writer, m *Module) error {
	var offset uint32
	if err := writeU32(w, Magic); err != nil {
		return err
	}
	if err := writeU32(w, currentVersion); err != nil {
		return err
	}

	offset += 8 // 4 for magic and 4 for version
	sections := m.Sections
	buf := new(bytes.Buffer)
	for _, s := range sections {
		var err error
		var size int
		if size, err = leb128.WriteVarUint32(w, uint32(s.SectionID())); err != nil {
			return err
		}
		offset += uint32(size)
		buf.Reset()

		rawSec := s.GetRawSection()
		rawSec.Start = int64(offset)
		if err := s.WritePayload(buf); err != nil {
			return err
		}

		if size, err = leb128.WriteVarUint32(w, uint32(buf.Len())); err != nil {
			return err
		}

		offset += uint32(size + buf.Len())
		if _, err := buf.WriteTo(w); err != nil {
			return err
		}

		rawSec.End = int64(offset)
	}
	return nil
}

func writeStringUint(w io.Writer, s string) error {
	return writeBytesUint(w, []byte(s))
}

func writeBytesUint(w io.Writer, p []byte) error {
	_, err := leb128.WriteVarUint32(w, uint32(len(p)))
	if err != nil {
		return err
	}
	_, err = w.Write(p)
	return err
}

func writeU32(w io.Writer, n uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], n)
	_, err := w.Write(buf[:])
	return err
}

func writeU64(w io.Writer, n uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], n)
	_, err := w.Write(buf[:])
	return err
}
