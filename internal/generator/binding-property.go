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
	/// One long property on an entity must be the ID
	PropertyFlagId = 1

	/// On languages like Java, a non-primitive type is used (aka wrapper types, allowing null)
	PropertyFlagNonPrimitiveType = 2

	/// Unused yet
	PropertyFlagNotNull = 4

	PropertyFlagIndexed = 8

	PropertyFlagReserved = 16

	/// Unused yet: Unique index
	PropertyFlagUnique = 32

	/// Unused yet: Use a persisted sequence to enforce ID to rise monotonic (no ID reuse)
	PropertyFlaIdMonotonicSequence = 64

	/// Allow IDs to be assigned by the developer
	PropertyFlagIdSelfAssignable = 128

	/// Unused yet
	PropertyFlagIndexPartialSkipNull = 256

	/// Unused yet, used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
	PropertyFlagIndexPartialSkipZero = 512

	/// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
	PropertyFlagVirtual = 1024

	/// Index uses a 32 bit hash instead of the value
	/// (32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions)
	PropertyFlagIndexHash = 2048

	/// Index uses a 64 bit hash instead of the value
	/// (recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash)
	PropertyFlagIndexHash64 = 4096

	/// The actual type of the variable is unsigned (used in combination with numeric OBXPropertyType_*)
	PropertyFlagUnsigned = 8192
)

const (
	PropertyTypeBool         = 1
	PropertyTypeByte         = 2
	PropertyTypeShort        = 3
	PropertyTypeChar         = 4
	PropertyTypeInt          = 5
	PropertyTypeLong         = 6
	PropertyTypeFloat        = 7
	PropertyTypeDouble       = 8
	PropertyTypeString       = 9
	PropertyTypeDate         = 10
	PropertyTypeRelation     = 11
	PropertyTypeByteVector   = 23
	PropertyTypeStringVector = 30
)
