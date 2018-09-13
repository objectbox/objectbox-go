package binding

import (
	"github.com/google/flatbuffers/go"
	. "github.com/objectbox/objectbox-go/objectbox"
	. "github.com/objectbox/objectbox-go/test/model/iot/flat"
	"github.com/objectbox/objectbox-go/test/model/iot/object"
	"strconv"
)

// This file could be generated in the future
type ReadingBinding struct {
}

func (ReadingBinding) AddToModel(model *Model) {
	model.Entity("Reading", 2, 10002)
	model.Property("id", PropertyType_Long, 1, 10002001)
	model.PropertyFlags(PropertyFlags_ID)
	model.Property("eventId", PropertyType_Long, 2, 10002002)
	//model.Property("eventId", PropertyType_Relation, 2, 10002002)
	//model.PropertyFlags(PropertyFlags_INDEXED)
	model.Property("date", PropertyType_Date, 3, 10002003)
	model.Property("valueName", PropertyType_String, 4, 10002004)
	model.Property("valueString", PropertyType_String, 5, 10002005)
	model.Property("valueInteger", PropertyType_Long, 6, 10002006)
	model.Property("valueFloating", PropertyType_Double, 7, 10002007)
	model.EntityLastPropertyId(7, 10002007)
}

func (ReadingBinding) GetId(entity interface{}) (id uint64, err error) {
	reading, ok := entity.(*object.Reading)
	if !ok {
		// Programming error, OK to panic
		panic("Object has wrong type")
	}
	idString := string(reading.Id)
	if idString == "" {
		return 0, nil
	}
	return strconv.ParseUint(idString, 10, 64)
}

func (ReadingBinding) Flatten(entity interface{}, fbb *flatbuffers.Builder, id uint64) {
	flattenModelReading(entity.(*object.Reading), fbb, id)
}

func flattenModelReading(reading *object.Reading, fbb *flatbuffers.Builder, id uint64) {
	offsetValueName := Unavailable
	if reading.ValueName != "" {
		offsetValueName = fbb.CreateString(reading.ValueName)
	}

	offsetValueString := Unavailable
	if reading.ValueString != "" {
		offsetValueString = fbb.CreateString(reading.ValueString)
	}

	ReadingStart(fbb)

	ReadingAddId(fbb, id)
	if reading.Date != 0 {
		ReadingAddDate(fbb, reading.Date)
	}

	if offsetValueName != Unavailable {
		ReadingAddValueName(fbb, offsetValueName)
	}
	if offsetValueString != Unavailable {
		ReadingAddValueString(fbb, offsetValueString)
	}
}

func (ReadingBinding) ToObject(bytes []byte) interface{} {
	flatReading := GetRootAsReading(bytes, flatbuffers.UOffsetT(0))
	return toModelReading(flatReading)
}

func (ReadingBinding) MakeSlice(capacity int) interface{} {
	return make([]object.Reading, 0, 16)
}

func (ReadingBinding) AppendToSlice(slice interface{}, entity interface{}) (sliceNew interface{}) {
	return append(slice.([]object.Reading), *entity.(*object.Reading))
}

func toModelReadingFromBytes(bytesData []byte) *object.Reading {
	flatReading := GetRootAsReading(bytesData, flatbuffers.UOffsetT(0))
	return toModelReading(flatReading)
}

func toModelReading(src *Reading) *object.Reading {
	return &object.Reading{
		Id:            src.Id(),
		Date:          src.Date(),
		ValueName:     string(src.ValueString()),
		ValueString:   string(src.ValueString()),
		ValueInteger:  src.ValueInteger(),
		ValueFloating: src.ValueFloating(),
	}

}
