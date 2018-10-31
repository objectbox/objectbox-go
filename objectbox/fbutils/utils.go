// Package fbutils provides utilities for the FlatBuffers in ObjectBox
package fbutils

import "github.com/google/flatbuffers/go"

type Table struct {
	_tab flatbuffers.Table
}

func CreateStringOffset(fbb *flatbuffers.Builder, value string) flatbuffers.UOffsetT {
	if len(value) > 0 {
		return fbb.CreateString(value)
	} else {
		return 0
	}
}

func GetRootAsTable(buf []byte, offset flatbuffers.UOffsetT) *Table {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Table{}
	x._tab.Bytes = buf
	x._tab.Pos = n + offset
	return x
}

func (table *Table) OffsetAsUint64(vtableOffset flatbuffers.VOffsetT) uint64 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetUint64(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsInt64(vtableOffset flatbuffers.VOffsetT) int64 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt64(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsString(vtableOffset flatbuffers.VOffsetT) string {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return string(table._tab.ByteVector(o + table._tab.Pos))
	}
	return ""
}
