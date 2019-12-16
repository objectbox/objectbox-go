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

// Package fbutils provides utilities for the FlatBuffers in ObjectBox
package fbutils

import "github.com/google/flatbuffers/go"

// CreateStringOffset creates an offset in the FlatBuffers table
func CreateStringOffset(fbb *flatbuffers.Builder, value string) flatbuffers.UOffsetT {
	return fbb.CreateString(value)
}

// CreateByteVectorOffset creates an offset in the FlatBuffers table
func CreateByteVectorOffset(fbb *flatbuffers.Builder, value []byte) flatbuffers.UOffsetT {
	if value == nil {
		return 0
	}

	return fbb.CreateByteVector(value)
}

// CreateStringVectorOffset creates an offset in the FlatBuffers table
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
