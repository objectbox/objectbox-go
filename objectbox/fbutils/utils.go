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
	if len(value) > 0 {
		return fbb.CreateByteVector(value)
	} else {
		return 0
	}
}

// Define some Get*Slot methods that are missing in the FlatBuffers table

func GetStringSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) string {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		return table.String(o + table.Pos)
	}
	return ""
}

func GetByteVectorSlot(table *flatbuffers.Table, slot flatbuffers.VOffsetT) []byte {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		return table.ByteVector(o + table.Pos)
	}
	return nil
}
