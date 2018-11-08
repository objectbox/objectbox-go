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

// Single header file for the ObjectBox C API
//
// Naming conventions
// ------------------
// * methods: obx_thing_action()
// * structs: OBX_thing {}
// * error codes: OBX_ERROR_REASON
// * enums: TODO
//

#ifndef OBJECTBOX_H
#define OBJECTBOX_H

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif

//----------------------------------------------
// ObjectBox version codes
//----------------------------------------------

// Note that you should use methods with prefix obx_version_ to check when linking against the dynamic library
#define OBX_VERSION_MAJOR 0
#define OBX_VERSION_MINOR 3
#define OBX_VERSION_PATCH 0

/// Returns the version of the library as ints. Pointers may be null
void obx_version(int* major, int* minor, int* patch);

/// Checks if the version of the library is equal to or higher than the given version ints.
/// @returns 1 if the condition is met and 0 otherwise
bool obx_version_is_at_least(int major, int minor, int patch);

/// Returns the version of the library to be printed.
/// The format may change; to query for version use the int based methods instead.
const char* obx_version_string();

/// Returns the version of the ObjectBox core to be printed.
/// The format may change, do not rely on its current form.
const char* obx_version_core_string();

//----------------------------------------------
// Return codes
//----------------------------------------------

/// Value returned when no error occurred (0)
#define OBX_SUCCESS 0

/// Returned by e.g. get operations if nothing was found for a specific ID.
/// This is NOT an error condition, and thus no last error info is set.
#define OBX_NOT_FOUND 404

// General errors
#define OBX_ERROR_ILLEGAL_STATE 10001
#define OBX_ERROR_ILLEGAL_ARGUMENT 10002
#define OBX_ERROR_ALLOCATION 10003
#define OBX_ERROR_NO_ERROR_INFO 10097
#define OBX_ERROR_GENERAL 10098
#define OBX_ERROR_UNKNOWN 10099

// Storage errors (often have a secondary error code)
#define OBX_ERROR_DB_FULL 10101
#define OBX_ERROR_MAX_READERS_EXCEEDED 10102
#define OBX_ERROR_STORE_MUST_SHUTDOWN 10103
#define OBX_ERROR_STORAGE_GENERAL 10199

// Data errors
#define OBX_ERROR_UNIQUE_VIOLATED 10201
#define OBX_ERROR_NON_UNIQUE_RESULT 10202
#define OBX_ERROR_PROPERTY_TYPE_MISMATCH 10203
#define OBX_ERROR_CONSTRAINT_VIOLATED 10299

// STD errors
#define OBX_ERROR_STD_ILLEGAL_ARGUMENT 10301
#define OBX_ERROR_STD_OUT_OF_RANGE 10302
#define OBX_ERROR_STD_LENGTH 10303
#define OBX_ERROR_STD_BAD_ALLOC 10304
#define OBX_ERROR_STD_RANGE 10305
#define OBX_ERROR_STD_OVERFLOW 10306
#define OBX_ERROR_STD_OTHER 10399

// Inconsistencies detected
#define OBX_ERROR_SCHEMA 10501
#define OBX_ERROR_FILE_CORRUPT 10502

//----------------------------------------------
// Common types
//----------------------------------------------
/// Schema entity & property identifiers
typedef uint32_t obx_schema_id;

/// Universal identifier used in schema for entities & properties
typedef uint64_t obx_uid;

/// ID of a single Object stored in the database
typedef uint64_t obx_id;

/// Error code returned by an obx_* function
typedef int obx_err;

//----------------------------------------------
// Error info
//----------------------------------------------

obx_err obx_last_error_code();

const char* obx_last_error_message();

obx_err obx_last_error_secondary();

void obx_last_error_clear();

//----------------------------------------------
// Model
//----------------------------------------------
typedef enum {
    PropertyType_Bool = 1,
    PropertyType_Byte = 2,
    PropertyType_Short = 3,
    PropertyType_Char = 4,
    PropertyType_Int = 5,
    PropertyType_Long = 6,
    PropertyType_Float = 7,
    PropertyType_Double = 8,
    PropertyType_String = 9,
    PropertyType_Date = 10,
    PropertyType_Relation = 11,
    PropertyType_ByteVector = 23,
} OBPropertyType;

/// Not really an enum, but binary flags to use across languages
typedef enum {
    /// One long property on an entity must be the ID
    PropertyFlags_ID = 1,

    /// On languages like Java, a non-primitive type is used (aka wrapper types, allowing null)
    PropertyFlags_NON_PRIMITIVE_TYPE = 2,

    /// Unused yet
    PropertyFlags_NOT_NULL = 4,
    PropertyFlags_INDEXED = 8,
    PropertyFlags_RESERVED = 16,
    /// Unused yet: Unique index
    PropertyFlags_UNIQUE = 32,
    /// Unused yet: Use a persisted sequence to enforce ID to rise monotonic (no ID reuse)
    PropertyFlags_ID_MONOTONIC_SEQUENCE = 64,
    /// Allow IDs to be assigned by the developer
    PropertyFlags_ID_SELF_ASSIGNABLE = 128,
    /// Unused yet
    PropertyFlags_INDEX_PARTIAL_SKIP_NULL = 256,
    /// Unused yet, used by References for 1) back-references and 2) to clear references to deleted objects (required
    /// for ID reuse)
    PropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512,
    /// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
    PropertyFlags_VIRTUAL = 1024,
    /// Index uses a 32 bit hash instead of the value
    /// (32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions)

    PropertyFlags_INDEX_HASH = 2048,
    /// Index uses a 64 bit hash instead of the value
    /// (recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit
    /// hash)
    PropertyFlags_INDEX_HASH64 = 4096

} OBPropertyFlags;

struct OBX_model;
typedef struct OBX_model OBX_model;

OBX_model* obx_model_create();

/// Only call when not calling obx_store_open (which will free it internally)
obx_err obx_model_free(OBX_model* model);

obx_err obx_model_entity(OBX_model* model, const char* name, obx_schema_id entity_id, obx_uid entity_uid);

obx_err obx_model_property(OBX_model* model, const char* name, OBPropertyType type, obx_schema_id property_id,
                           obx_uid property_uid);

obx_err obx_model_property_flags(OBX_model* model, OBPropertyFlags flags);

obx_err obx_model_property_relation(OBX_model* model, const char* target_entity, obx_schema_id index_id,
                                    obx_uid index_uid);

obx_err obx_model_property_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

void obx_model_last_entity_id(OBX_model*, obx_schema_id entity_id, obx_uid entity_uid);

void obx_model_last_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

void obx_model_last_relation_id(OBX_model* model, obx_schema_id relation_id, obx_uid relation_uid);

obx_err obx_model_entity_last_property_id(OBX_model* model, obx_schema_id property_id, obx_uid property_uid);

//----------------------------------------------
// Store
//----------------------------------------------

struct OBX_store;
typedef struct OBX_store OBX_store;

typedef struct OBX_store_options {
    /// Use NULL for default value ("objectbox")
    char* directory;

    /// Use zero for default value
    uint64_t maxDbSizeInKByte;

    /// Use zero for default value
    unsigned int fileMode;

    /// Use zero for default value
    unsigned int maxReaders;
} OBX_store_options;

typedef enum {
    DebugFlags_LOG_TRANSACTIONS_READ = 1,
    DebugFlags_LOG_TRANSACTIONS_WRITE = 2,
    DebugFlags_LOG_QUERIES = 4,
    DebugFlags_LOG_QUERY_PARAMETERS = 8,
    DebugFlags_LOG_ASYNC_QUEUE = 16,
} OBDebugFlags;

typedef struct OBX_bytes {
    void* data;
    size_t size;
} OBX_bytes;

typedef struct OBX_bytes_array {
    OBX_bytes* bytes;
    size_t count;
} OBX_bytes_array;

typedef struct OBX_id_array {
    obx_id* ids;
    size_t count;
} OBX_id_array;

OBX_store* obx_store_open_bytes(const void* model_bytes, size_t model_size, const OBX_store_options* options);

/// Note: the model is freed by calling this method
OBX_store* obx_store_open(OBX_model* model, const OBX_store_options* options);

obx_schema_id obx_store_entity_id(OBX_store* store, const char* entity_name);

obx_schema_id obx_store_entity_property_id(OBX_store* store, obx_schema_id entity_id, const char* property_name);

obx_err obx_store_await_async_completion(OBX_store* store);

obx_err obx_store_debug_flags(OBX_store* store, OBDebugFlags flags);

obx_err obx_store_close(OBX_store* store);

//----------------------------------------------
// Transaction
//----------------------------------------------

struct OBX_txn;
typedef struct OBX_txn OBX_txn;

OBX_txn* obx_txn_begin(OBX_store* store);

OBX_txn* obx_txn_begin_read(OBX_store* store);

obx_err obx_txn_close(OBX_txn* txn);

obx_err obx_txn_abort(OBX_txn* txn);

obx_err obx_txn_commit(OBX_txn* txn);

//----------------------------------------------
// Cursor
//----------------------------------------------

struct OBX_cursor;
typedef struct OBX_cursor OBX_cursor;

OBX_cursor* obx_cursor_create(OBX_txn* txn, obx_schema_id entity_id);

OBX_cursor* obx_cursor_create2(OBX_txn* txn, const char* entity_name);

obx_err obx_cursor_close(OBX_cursor* cursor);

obx_id obx_cursor_id_for_put(OBX_cursor* cursor, obx_id id_or_zero);

/// ATTENTION: ensure that the given value memory is allocated to the next 4 bytes boundary.
/// ObjectBox needs to store bytes with sizes dividable by 4 for internal reasons.
/// Use obx_cursor_put_padded otherwise.
obx_err obx_cursor_put(OBX_cursor* cursor, obx_id id, const void* data, size_t size, bool checkForPreviousValueFlag);

/// Prefer obx_cursor_put (non-padded) if possible, as this does a memcpy if the size is not dividable by 4.
obx_err obx_cursor_get(OBX_cursor* cursor, obx_id id, void** data, size_t* size);

/// Gets all objects as bytes.
/// For bigger quantities, it's recommended to iterate using obx_cursor_first and obx_cursor_next.
/// However, if the calling overhead is high (e.g. for language bindings), this method helps.
OBX_bytes_array* obx_cursor_get_all(OBX_cursor* cursor);

obx_err obx_cursor_first(OBX_cursor* cursor, void** data, size_t* size);

obx_err obx_cursor_next(OBX_cursor* cursor, void** data, size_t* size);

obx_err obx_cursor_remove(OBX_cursor* cursor, obx_id id);

obx_err obx_cursor_remove_all(OBX_cursor* cursor);

obx_err obx_cursor_count(OBX_cursor* cursor, uint64_t* count);

OBX_bytes_array* obx_cursor_backlink_bytes(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                           obx_id id);

OBX_id_array* obx_cursor_backlink_ids(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                      obx_id id);

//----------------------------------------------
// Box
//----------------------------------------------

/// A box may be used across threads
struct OBX_box;
typedef struct OBX_box OBX_box;

OBX_box* obx_box_create(OBX_store* store, obx_schema_id entity_id);

obx_err obx_box_close(OBX_box* box);

obx_id obx_box_id_for_put(OBX_box* box, obx_id id_or_zero);

obx_err obx_box_put_async(OBX_box* box, obx_id id, const void* data, size_t size, bool checkForPreviousValueFlag);

//----------------------------------------------
// Query Builder
//----------------------------------------------
struct OBX_query_builder;
typedef struct OBX_query_builder OBX_query_builder;

/// Query Builder condition identifier
/// - returned by condition creating functions,
/// - used to combine conditions with any/all, thus building more complex conditions
typedef int obx_qb_cond;

OBX_query_builder* obx_qb_create(OBX_store* store, obx_schema_id entity_id);
obx_err obx_qb_close(OBX_query_builder* builder);
obx_err obx_qb_error_code(OBX_query_builder* builder);
const char* obx_qb_error_message(OBX_query_builder* builder);

obx_qb_cond obx_qb_null(OBX_query_builder* builder, obx_schema_id property_id);
obx_qb_cond obx_qb_not_null(OBX_query_builder* builder, obx_schema_id property_id);

obx_qb_cond obx_qb_string_equal(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                bool case_sensitive);

obx_qb_cond obx_qb_string_not_equal(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                    bool case_sensitive);
obx_qb_cond obx_qb_string_contains(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                   bool case_sensitive);
obx_qb_cond obx_qb_string_starts_with(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                      bool case_sensitive);
obx_qb_cond obx_qb_string_ends_with(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                    bool case_sensitive);
obx_qb_cond obx_qb_string_greater(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                  bool case_sensitive, bool with_equal);
obx_qb_cond obx_qb_string_less(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                               bool case_sensitive, bool with_equal);
obx_qb_cond obx_qb_string_in(OBX_query_builder* builder, obx_schema_id property_id, const char* values[], int count,
                             bool case_sensitive);

obx_qb_cond obx_qb_int_equal(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
obx_qb_cond obx_qb_int_not_equal(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
obx_qb_cond obx_qb_int_greater(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
obx_qb_cond obx_qb_int_less(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
obx_qb_cond obx_qb_int_between(OBX_query_builder* builder, obx_schema_id property_id, int64_t value_a, int64_t value_b);

obx_qb_cond obx_qb_int64_in(OBX_query_builder* builder, obx_schema_id property_id, const int64_t values[], int count);
obx_qb_cond obx_qb_int64_not_in(OBX_query_builder* builder, obx_schema_id property_id, const int64_t values[],
                                int count);

obx_qb_cond obx_qb_int32_in(OBX_query_builder* builder, obx_schema_id property_id, const int32_t values[], int count);
obx_qb_cond obx_qb_int32_not_in(OBX_query_builder* builder, obx_schema_id property_id, const int32_t values[],
                                int count);

obx_qb_cond obx_qb_double_greater(OBX_query_builder* builder, obx_schema_id property_id, double value);
obx_qb_cond obx_qb_double_less(OBX_query_builder* builder, obx_schema_id property_id, double value);
obx_qb_cond obx_qb_double_between(OBX_query_builder* builder, obx_schema_id property_id, double value_a,
                                  double value_b);

obx_qb_cond obx_qb_bytes_equal(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size);
obx_qb_cond obx_qb_bytes_greater(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size,
                                 bool with_equal);
obx_qb_cond obx_qb_bytes_less(OBX_query_builder* builder, obx_schema_id property_id, const void* value, size_t size,
                              bool with_equal);

/// Combines conditions[] to a new condition using operator AND (all) or OR (any)
/// Note that these functions remove original conditions from the condition list and thus affect indices of remaining
/// conditions in the list
obx_qb_cond obx_qb_all(OBX_query_builder* builder, const obx_qb_cond conditions[], int count);
obx_qb_cond obx_qb_any(OBX_query_builder* builder, const obx_qb_cond conditions[], int count);

obx_err obx_qb_parameter_alias(OBX_query_builder* builder, const char* alias);

//----------------------------------------------
// Query
//----------------------------------------------
struct OBX_query;
typedef struct OBX_query OBX_query;

OBX_query* obx_query_create(OBX_query_builder* builder);
obx_err obx_query_close(OBX_query* query);

OBX_bytes_array* obx_query_find(OBX_query* query, OBX_cursor* cursor);
OBX_id_array* obx_query_find_ids(OBX_query* query, OBX_cursor* cursor);
obx_err obx_query_count(OBX_query* query, OBX_cursor* cursor, uint64_t* count);

/// Removes (deletes!) all matching entities.
obx_err obx_query_remove(OBX_query* query, OBX_cursor* cursor, uint64_t* count);

// TODO either introduce other group of "param" functions with entityId, or require entityId in each call
obx_err obx_query_string_param(OBX_query* query, obx_schema_id property_id, const char* value);
obx_err obx_query_string_params_in(OBX_query* query, obx_schema_id property_id, const char* values[], int count);
obx_err obx_query_int_param(OBX_query* query, obx_schema_id property_id, int64_t value);
obx_err obx_query_int_params(OBX_query* query, obx_schema_id property_id, int64_t value_a, int64_t value_b);
obx_err obx_query_int64_params_in(OBX_query* query, obx_schema_id property_id, const int64_t values[], int count);
obx_err obx_query_int32_params_in(OBX_query* query, obx_schema_id property_id, const int32_t values[], int count);
obx_err obx_query_double_param(OBX_query* query, obx_schema_id property_id, double value);
obx_err obx_query_double_params(OBX_query* query, obx_schema_id property_id, double value_a, double value_b);
obx_err obx_query_bytes_param(OBX_query* query, obx_schema_id property_id, const void* value, size_t size);

obx_err obx_query_string_param_alias(OBX_query* query, const char* alias, const char* value);
obx_err obx_query_string_params_in_alias(OBX_query* query, const char* alias, const char* values[], int count);
obx_err obx_query_int_param_alias(OBX_query* query, const char* alias, int64_t value);
obx_err obx_query_int_params_alias(OBX_query* query, const char* alias, int64_t value_a, int64_t value_b);
obx_err obx_query_int64_params_in_alias(OBX_query* query, const char* alias, const int64_t values[], int count);
obx_err obx_query_int32_params_in_alias(OBX_query* query, const char* alias, const int32_t values[], int count);
obx_err obx_query_double_param_alias(OBX_query* query, const char* alias, double value);
obx_err obx_query_double_params_alias(OBX_query* query, const char* alias, double value_a, double value_b);
obx_err obx_query_bytes_param_alias(OBX_query* query, const char* alias, const void* value, size_t size);

/// the resulting char* is valid until another call on describe_parameters is made on the same query or until the query
/// is freed
const char* obx_query_describe_parameters(OBX_query* query);

/// the resulting char* is valid until another call on to_string is made on the same query or until the query is freed
const char* obx_query_to_string(OBX_query* query);

void obx_bytes_free(OBX_bytes* bytes);
void obx_bytes_array_free(OBX_bytes_array* bytes_array);
void obx_id_array_free(OBX_id_array* id_array);

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_H
