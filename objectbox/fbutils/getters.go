/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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

package fbutils

import flatbuffers "github.com/google/flatbuffers/go"

// Define some Get*Slot methods that are missing in the FlatBuffers table
// NOTE - don't use table.String because byteSliceToString is "unsafe" and doesn't play well
// with our []byte to C void* mapping. This leads to weird runtime errors because that string
// just points to a memory that has already been freed/reused by C.

func GetStringSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) string {
	if o := table.Offset(slot); o != 0 {
		return string(table.ByteVector(flatbuffers.UOffsetT(o) + table.Pos))
	}
	return ""
}

func GetStringPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *string {
	if o := table.Offset(slot); o != 0 {
		var value = string(table.ByteVector(flatbuffers.UOffsetT(o) + table.Pos))
		return &value
	}
	return nil
}

func GetByteVectorSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) []byte {
	if vector := GetByteVectorPtrSlot(table, slot); vector != nil {
		return *vector
	}
	return nil
}

func GetByteVectorPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *[]byte {
	if o := table.Offset(slot); o != 0 {
		// we need to make a copy because the source bytes are directly mapped to a C void* (same as for GetStringSlot)
		var src = table.ByteVector(flatbuffers.UOffsetT(o) + table.Pos)
		var result = make([]byte, len(src))
		copy(result, src)
		return &result
	}
	return nil
}

func GetStringVectorSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) []string {
	if vector := GetStringVectorPtrSlot(table, slot); vector != nil {
		return *vector
	}
	return nil
}

func GetStringVectorPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *[]string {
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

		return &values
	}

	return nil
}

func GetBoolSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) bool {
	return table.GetBoolSlot(slot, false)
}

func GetBoolPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *bool {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetBool(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetByteSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) byte {
	return table.GetByteSlot(slot, 0)
}

func GetBytePtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *byte {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetByte(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetRuneSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) rune {
	return table.GetInt32Slot(slot, 0)
}

func GetRunePtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *rune {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetInt32(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetIntSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) int {
	return int(table.GetInt64Slot(slot, 0))
}

func GetIntPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *int {
	if o := table.Offset(slot); o != 0 {
		var value = int(table.GetInt64(flatbuffers.UOffsetT(o) + table.Pos))
		return &value
	}
	return nil
}

func GetInt8Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) int8 {
	return table.GetInt8Slot(slot, 0)
}

func GetInt8PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *int8 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetInt8(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetInt16Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) int16 {
	return table.GetInt16Slot(slot, 0)
}

func GetInt16PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *int16 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetInt16(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetInt32Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) int32 {
	return table.GetInt32Slot(slot, 0)
}

func GetInt32PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *int32 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetInt32(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetInt64Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) int64 {
	return table.GetInt64Slot(slot, 0)
}

func GetInt64PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *int64 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetInt64(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetUintSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) uint {
	return uint(table.GetUint64Slot(slot, 0))
}

func GetUintPtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *uint {
	if o := table.Offset(slot); o != 0 {
		var value = uint(table.GetUint64(flatbuffers.UOffsetT(o) + table.Pos))
		return &value
	}
	return nil
}

func GetUint8Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) uint8 {
	return table.GetUint8Slot(slot, 0)
}

func GetUint8PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *uint8 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetUint8(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetUint16Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) uint16 {
	return table.GetUint16Slot(slot, 0)
}

func GetUint16PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *uint16 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetUint16(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetUint32Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) uint32 {
	return table.GetUint32Slot(slot, 0)
}

func GetUint32PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *uint32 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetUint32(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetUint64Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) uint64 {
	return table.GetUint64Slot(slot, 0)
}

func GetUint64PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *uint64 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetUint64(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetFloat32Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) float32 {
	return table.GetFloat32Slot(slot, 0)
}

func GetFloat32PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *float32 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetFloat32(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}

func GetFloat64Slot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) float64 {
	return table.GetFloat64Slot(slot, 0)
}

func GetFloat64PtrSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) *float64 {
	if o := table.Offset(slot); o != 0 {
		var value = table.GetFloat64(flatbuffers.UOffsetT(o) + table.Pos)
		return &value
	}
	return nil
}
