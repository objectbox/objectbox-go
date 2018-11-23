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
