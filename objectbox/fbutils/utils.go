/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package fbutils provides utilities for the FlatBuffers in ObjectBox
package fbutils

import "github.com/google/flatbuffers/go"

func CreateStringOffset(fbb *flatbuffers.Builder, value string) flatbuffers.UOffsetT {
	if len(value) > 0 {
		return fbb.CreateString(value)
	} else {
		return 0
	}
}

func CreateByteVectorOffset(fbb *flatbuffers.Builder, value []byte) flatbuffers.UOffsetT {
	if value == nil {
		return 0
	}

	return fbb.CreateByteVector(value)
}

func CreateStringVectorOffset(fbb *flatbuffers.Builder, values []string) flatbuffers.UOffsetT {
	if values == nil {
		return 0
	}

	var offsets = make([]flatbuffers.UOffsetT, len(values))
	for i, v := range values {
		offsets[i] = fbb.CreateString(v)
	}

	return createOffsetVector(fbb, offsets)
}

func createOffsetVector(fbb *flatbuffers.Builder, offsets []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	fbb.StartVector(int(flatbuffers.SizeUOffsetT), len(offsets), int(flatbuffers.SizeUOffsetT))
	for i := len(offsets) - 1; i >= 0; i-- {
		fbb.PrependUOffsetT(offsets[i])
	}
	return fbb.EndVector(len(offsets))
}

// Define some Get*Slot methods that are missing in the FlatBuffers table
// NOTE - don't use table.String because byteSliceToString is "unsafe" and doesn't play well
// with our []byte to C void* mapping. This leads to weird runtime errors because that string
// just points to a memory that has already been freed/reused by C.

func GetStringSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) string {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		return string(table.ByteVector(o + table.Pos))
	}
	return ""
}

func GetByteVectorSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) []byte {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		// we need to make a copy because the source bytes are directly mapped to a C void* (same as for GetStringSlot)
		var src = table.ByteVector(o + table.Pos)
		var result = make([]byte, len(src))
		copy(result, src)
		return result
	}
	return nil
}

func GetStringVectorSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) []string {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		var ln = table.VectorLen(o) // number of elements

		// prepare the result vector
		var values = make([]string, 0, ln)

		// iterate over the vector and read each element separately
		var start = table.Vector(o)
		var end = start + flatbuffers.UOffsetT(ln)*flatbuffers.SizeUOffsetT

		for pos := start; pos < end; pos += flatbuffers.SizeUOffsetT {
			values = append(values, string(table.ByteVector(pos)))
		}

		return values
	}

	return nil
}

// Setters always write values, regardless of default

func SetBoolSlot(fbb *flatbuffers.Builder, slot int, value bool) {
	if value {
		SetByteSlot(fbb, slot, 1)
	} else {
		SetByteSlot(fbb, slot, 0)
	}
}

func SetByteSlot(fbb *flatbuffers.Builder, slot int, value byte) {
	fbb.PrependByte(value)
	fbb.Slot(slot)
}

func SetUint8Slot(fbb *flatbuffers.Builder, slot int, value uint8) {
	fbb.PrependUint8(value)
	fbb.Slot(slot)
}

func SetUint16Slot(fbb *flatbuffers.Builder, slot int, value uint16) {
	fbb.PrependUint16(value)
	fbb.Slot(slot)
}

func SetUint32Slot(fbb *flatbuffers.Builder, slot int, value uint32) {
	fbb.PrependUint32(value)
	fbb.Slot(slot)
}

func SetUint64Slot(fbb *flatbuffers.Builder, slot int, value uint64) {
	fbb.PrependUint64(value)
	fbb.Slot(slot)
}

func SetInt8Slot(fbb *flatbuffers.Builder, slot int, value int8) {
	fbb.PrependInt8(value)
	fbb.Slot(slot)
}

func SetInt16Slot(fbb *flatbuffers.Builder, slot int, value int16) {
	fbb.PrependInt16(value)
	fbb.Slot(slot)
}

func SetInt32Slot(fbb *flatbuffers.Builder, slot int, value int32) {
	fbb.PrependInt32(value)
	fbb.Slot(slot)
}

func SetInt64Slot(fbb *flatbuffers.Builder, slot int, value int64) {
	fbb.PrependInt64(value)
	fbb.Slot(slot)
}

func SetFloat32Slot(fbb *flatbuffers.Builder, slot int, value float32) {
	fbb.PrependFloat32(value)
	fbb.Slot(slot)
}

func SetFloat64Slot(fbb *flatbuffers.Builder, slot int, value float64) {
	fbb.PrependFloat64(value)
	fbb.Slot(slot)
}

func SetUOffsetTSlot(fbb *flatbuffers.Builder, slot int, value flatbuffers.UOffsetT) {
	if value != 0 {
		fbb.PrependUOffsetT(value)
		fbb.Slot(slot)
	}
}
