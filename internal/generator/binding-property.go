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

package generator

const (
	propertyFlagID = 1

	propertyFlagNonPrimitiveType = 2

	propertyFlagNotNull = 4

	propertyFlagIndexed = 8

	propertyFlagReserved = 16

	propertyFlagUnique = 32

	propertyFlaIDMonotonicSequence = 64

	propertyFlagIDSelfAssignable = 128

	propertyFlagIndexPartialSkipNull = 256

	propertyFlagIndexPartialSkipZero = 512

	propertyFlagVirtual = 1024

	propertyFlagIndexHash = 2048

	propertyFlagIndexHash64 = 4096

	propertyFlagUnsigned = 8192
)

const (
	propertyTypeBool         = 1
	propertyTypeByte         = 2
	propertyTypeShort        = 3
	propertyTypeChar         = 4
	propertyTypeInt          = 5
	propertyTypeLong         = 6
	propertyTypeFloat        = 7
	propertyTypeDouble       = 8
	propertyTypeString       = 9
	propertyTypeDate         = 10
	propertyTypeRelation     = 11
	propertyTypeByteVector   = 23
	propertyTypeStringVector = 30
)
