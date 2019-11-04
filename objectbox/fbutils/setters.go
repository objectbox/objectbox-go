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

// Setters always write values, regardless of the default

// SetBoolSlot sets the given value in the FlatBuffers table
func SetBoolSlot(fbb *flatbuffers.Builder, slot int, value bool) {
	if value {
		SetByteSlot(fbb, slot, 1)
	} else {
		SetByteSlot(fbb, slot, 0)
	}
}

// SetByteSlot sets the given value in the FlatBuffers table
func SetByteSlot(fbb *flatbuffers.Builder, slot int, value byte) {
	fbb.PrependByte(value)
	fbb.Slot(slot)
}

// SetUint8Slot sets the given value in the FlatBuffers table
func SetUint8Slot(fbb *flatbuffers.Builder, slot int, value uint8) {
	fbb.PrependUint8(value)
	fbb.Slot(slot)
}

// SetUint16Slot sets the given value in the FlatBuffers table
func SetUint16Slot(fbb *flatbuffers.Builder, slot int, value uint16) {
	fbb.PrependUint16(value)
	fbb.Slot(slot)
}

// SetUint32Slot sets the given value in the FlatBuffers table
func SetUint32Slot(fbb *flatbuffers.Builder, slot int, value uint32) {
	fbb.PrependUint32(value)
	fbb.Slot(slot)
}

// SetUint64Slot sets the given value in the FlatBuffers table
func SetUint64Slot(fbb *flatbuffers.Builder, slot int, value uint64) {
	fbb.PrependUint64(value)
	fbb.Slot(slot)
}

// SetInt8Slot sets the given value in the FlatBuffers table
func SetInt8Slot(fbb *flatbuffers.Builder, slot int, value int8) {
	fbb.PrependInt8(value)
	fbb.Slot(slot)
}

// SetInt16Slot sets the given value in the FlatBuffers table
func SetInt16Slot(fbb *flatbuffers.Builder, slot int, value int16) {
	fbb.PrependInt16(value)
	fbb.Slot(slot)
}

// SetInt32Slot sets the given value in the FlatBuffers table
func SetInt32Slot(fbb *flatbuffers.Builder, slot int, value int32) {
	fbb.PrependInt32(value)
	fbb.Slot(slot)
}

// SetInt64Slot sets the given value in the FlatBuffers table
func SetInt64Slot(fbb *flatbuffers.Builder, slot int, value int64) {
	fbb.PrependInt64(value)
	fbb.Slot(slot)
}

// SetFloat32Slot sets the given value in the FlatBuffers table
func SetFloat32Slot(fbb *flatbuffers.Builder, slot int, value float32) {
	fbb.PrependFloat32(value)
	fbb.Slot(slot)
}

// SetFloat64Slot sets the given value in the FlatBuffers table
func SetFloat64Slot(fbb *flatbuffers.Builder, slot int, value float64) {
	fbb.PrependFloat64(value)
	fbb.Slot(slot)
}

// SetUOffsetTSlot sets the given value in the FlatBuffers table
func SetUOffsetTSlot(fbb *flatbuffers.Builder, slot int, value flatbuffers.UOffsetT) {
	if value != 0 {
		fbb.PrependUOffsetT(value)
		fbb.Slot(slot)
	}
}
