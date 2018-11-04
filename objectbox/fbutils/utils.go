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

func CreateByteVectorOffset(fbb *flatbuffers.Builder, value []byte) flatbuffers.UOffsetT {
	if len(value) > 0 {
		return fbb.CreateByteVector(value)
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

func (table *Table) OffsetAsBool(vtableOffset flatbuffers.VOffsetT) bool {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetBool(o + table._tab.Pos)
	}
	return false
}

func (table *Table) OffsetAsByte(vtableOffset flatbuffers.VOffsetT) byte {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetByte(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsRune(vtableOffset flatbuffers.VOffsetT) rune {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt32(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsInt(vtableOffset flatbuffers.VOffsetT) int {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return int(table._tab.GetInt32(o + table._tab.Pos))
	}
	return 0
}

func (table *Table) OffsetAsInt8(vtableOffset flatbuffers.VOffsetT) int8 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt8(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsInt16(vtableOffset flatbuffers.VOffsetT) int16 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt16(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsInt32(vtableOffset flatbuffers.VOffsetT) int32 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt32(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsInt64(vtableOffset flatbuffers.VOffsetT) int64 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetInt64(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsUint(vtableOffset flatbuffers.VOffsetT) uint {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return uint(table._tab.GetUint32(o + table._tab.Pos))
	}
	return 0
}

func (table *Table) OffsetAsUint8(vtableOffset flatbuffers.VOffsetT) uint8 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetUint8(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsUint16(vtableOffset flatbuffers.VOffsetT) uint16 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetUint16(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsUint32(vtableOffset flatbuffers.VOffsetT) uint32 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetUint32(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsUint64(vtableOffset flatbuffers.VOffsetT) uint64 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetUint64(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsFloat32(vtableOffset flatbuffers.VOffsetT) float32 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetFloat32(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsFloat64(vtableOffset flatbuffers.VOffsetT) float64 {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.GetFloat64(o + table._tab.Pos)
	}
	return 0
}

func (table *Table) OffsetAsString(vtableOffset flatbuffers.VOffsetT) string {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return string(table._tab.ByteVector(o + table._tab.Pos))
	}
	return ""
}

func (table *Table) OffsetAsByteVector(vtableOffset flatbuffers.VOffsetT) []byte {
	if o := flatbuffers.UOffsetT(table._tab.Offset(vtableOffset)); o != 0 {
		return table._tab.ByteVector(o + table._tab.Pos)
	}
	return nil
}
