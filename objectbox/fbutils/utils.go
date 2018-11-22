// Package fbutils provides utilities for the FlatBuffers in ObjectBox
package fbutils

import "github.com/google/flatbuffers/go"

type Table struct {
	*flatbuffers.Table
}

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

func GetRootAsTable(buf []byte, offset flatbuffers.UOffsetT) *Table {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	return &Table{
		&flatbuffers.Table{
			Bytes: buf,
			Pos:   n + offset,
		},
	}
}

// Define some Get*Slot methods that are missing in the FlatBuffers table

func (table *Table) GetStringSlot(slot flatbuffers.VOffsetT) string {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		return table.String(o + table.Pos)
	}
	return ""
}

func (table *Table) GetByteVectorSlot(slot flatbuffers.VOffsetT) []byte {
	if o := flatbuffers.UOffsetT(table.Offset(slot)); o != 0 {
		return table.ByteVector(o + table.Pos)
	}
	return nil
}
