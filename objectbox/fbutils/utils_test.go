package fbutils

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/test/assert"
)

// This test is here to make sure our flatbuffers integration is correct and mainly focuses on memory management
// because of how the c-api integration is implemented (not copying C void* but mapping it to []byte).
// The test covers all supported types (see entity.go in the test directory).
// It creates an FB object, creates an unsafe copy (similar to the c-arrays.go) and than clears the memory behind that
// unsafe copy. Afterwards, it tests that the "parsed" flatbuffers object is still the same as the written one.

// this is just a test of the auxiliary functions in this file, so that we can concentrate on actual tests elsewhere
func TestAuxUnsafeBytes(t *testing.T) {
	var objBytesManaged = createObjectBytes(object)
	assert.Eq(t, 216, len(objBytesManaged))
	assert.Eq(t, 52, int(objBytesManaged[0]))
	assert.Eq(t, 51, int(objBytesManaged[57]))

	// get an unsafe copy (really just a pointer), as the the cursor would do
	var unsafeBytes = getUnsafeBytes(objBytesManaged)

	// get a safe copy
	var safeBytes = make([]byte, len(objBytesManaged))
	copy(safeBytes, objBytesManaged)

	// at this point, they should all be the same
	assert.Eq(t, objBytesManaged, unsafeBytes)
	assert.NotEq(t, unsafe.Pointer(&objBytesManaged), unsafe.Pointer(&unsafeBytes))
	assert.Eq(t, unsafe.Pointer(&objBytesManaged[0]), unsafe.Pointer(&unsafeBytes[0]))

	assert.Eq(t, objBytesManaged, safeBytes)
	assert.NotEq(t, unsafe.Pointer(&objBytesManaged), unsafe.Pointer(&safeBytes))

	// now let's clear the object bytes, and check the copies
	clearBytes(&objBytesManaged)
	assert.Eq(t, 216, len(objBytesManaged))
	assert.Eq(t, 0, int(objBytesManaged[0]))

	// now we assert the unsafe bytes has changed if it wasn't supposed to
	assert.Eq(t, objBytesManaged, unsafeBytes)

	// but the safe copy is still the same
	assert.Eq(t, 52, int(safeBytes[0]))
}

func TestObjectRead(t *testing.T) {
	var objBytesManaged = createObjectBytes(object)

	// get an unsafe copy (really just a pointer), as the the cursor would do
	var unsafeBytes = getUnsafeBytes(objBytesManaged)

	// read the object value
	var read = readObject(unsafeBytes)
	assert.Eq(t, object, read)

	// clear the source bytes
	clearBytes(&objBytesManaged)

	// the read object should be still correct
	assert.Eq(t, object, read)
	fmt.Println(read)
}

// this simulates what cVoidPtrToByteSlice is doing, i.e. mapping an unmanaged pointer to a new []byte slice
func getUnsafeBytes(source []byte) []byte {
	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = uintptr(unsafe.Pointer(&source[0]))
	header.Len = len(source)
	header.Cap = header.Len
	return bytes
}

func clearBytes(data *[]byte) {
	for i := range *data {
		(*data)[i] = 0
	}
}

// this is a copy of test/model/entity.go
type entity struct {
	Id           uint64
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Bool         bool
	String       string
	StringVector []string
	Byte         byte
	ByteVector   []byte
	Rune         rune
	Float32      float32
	Float64      float64
}

// a prototype object
var object = entity{
	Id:           1,
	Int:          2,
	Int8:         3,
	Int16:        4,
	Int32:        5,
	Int64:        6,
	Uint:         7,
	Uint8:        8,
	Uint16:       9,
	Uint32:       10,
	Uint64:       11,
	Bool:         true,
	String:       "Str",
	StringVector: []string{"First", "second", ""},
	Byte:         12,
	ByteVector:   []byte{13, 14, 15},
	Rune:         16,
	Float32:      17.18,
	Float64:      19.20,
}

func createObjectBytes(obj entity) []byte {
	// this is a copy of test/model/entity.obx.go Flatten()
	var fbb = flatbuffers.NewBuilder(512)
	var offsetString = CreateStringOffset(fbb, obj.String)
	var offsetStringVector = CreateStringVectorOffset(fbb, obj.StringVector)
	var offsetByteVector = CreateByteVectorOffset(fbb, obj.ByteVector)

	// build the FlatBuffers object
	fbb.StartObject(21)
	SetUint64Slot(fbb, 0, obj.Id)
	SetInt64Slot(fbb, 1, int64(obj.Int))
	SetInt8Slot(fbb, 2, obj.Int8)
	SetInt16Slot(fbb, 3, obj.Int16)
	SetInt32Slot(fbb, 4, obj.Int32)
	SetInt64Slot(fbb, 5, obj.Int64)
	SetUint64Slot(fbb, 6, uint64(obj.Uint))
	SetUint8Slot(fbb, 7, obj.Uint8)
	SetUint16Slot(fbb, 8, obj.Uint16)
	SetUint32Slot(fbb, 9, obj.Uint32)
	SetUint64Slot(fbb, 10, obj.Uint64)
	SetBoolSlot(fbb, 11, obj.Bool)
	SetUOffsetTSlot(fbb, 12, offsetString)
	SetUOffsetTSlot(fbb, 20, offsetStringVector)
	SetByteSlot(fbb, 13, obj.Byte)
	SetUOffsetTSlot(fbb, 14, offsetByteVector)
	SetInt32Slot(fbb, 15, obj.Rune)
	SetFloat32Slot(fbb, 16, obj.Float32)
	SetFloat64Slot(fbb, 17, obj.Float64)
	fbb.Finish(fbb.EndObject())
	return fbb.FinishedBytes()
}

func readObject(data []byte) entity {
	// this is a copy of test/model/entity.obx.go Load()
	table := &flatbuffers.Table{
		Bytes: data,
		Pos:   flatbuffers.GetUOffsetT(data),
	}

	return entity{
		Id:           table.GetUint64Slot(4, 0),
		Int:          int(table.GetUint64Slot(6, 0)),
		Int8:         table.GetInt8Slot(8, 0),
		Int16:        table.GetInt16Slot(10, 0),
		Int32:        table.GetInt32Slot(12, 0),
		Int64:        table.GetInt64Slot(14, 0),
		Uint:         uint(table.GetUint64Slot(16, 0)),
		Uint8:        table.GetUint8Slot(18, 0),
		Uint16:       table.GetUint16Slot(20, 0),
		Uint32:       table.GetUint32Slot(22, 0),
		Uint64:       table.GetUint64Slot(24, 0),
		Bool:         table.GetBoolSlot(26, false),
		String:       GetStringSlot(table, 28),
		StringVector: GetStringVectorSlot(table, 44),
		Byte:         table.GetByteSlot(30, 0),
		ByteVector:   GetByteVectorSlot(table, 32),
		Rune:         rune(table.GetInt32Slot(34, 0)),
		Float32:      table.GetFloat32Slot(36, 0),
		Float64:      table.GetFloat64Slot(38, 0),
	}
}
