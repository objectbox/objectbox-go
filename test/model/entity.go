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

package model

//go:generate objectbox-gogen

// Tests all available GO & ObjectBox types
// TODO rename; e.g. TestEntity
type Entity struct {
	Id         uint64
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Bool       bool
	String     string
	Byte       byte
	ByteVector []byte
	Rune       rune
	Float32    float32
	Float64    float64
	Date       int64 `date`
	// TODO Complex64  complex64
	// TODO Complex128 complex128
}

type TestStringIdEntity struct {
	Id string `id`
}
