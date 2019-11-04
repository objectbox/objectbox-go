/*
 * Copyright 2018-2019 ObjectBox Ltd. All rights reserved.
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

/// When using ObjectBox as a dynamic library, you should verify that a compatible version was linked using obx_version() or obx_version_is_at_least().
#define OBX_VERSION_MAJOR 0
#define OBX_VERSION_MINOR 7
#define OBX_VERSION_PATCH 1  // values >= 100 are reserved for dev releases leading to the next minor/major increase

/// Returns the version of the library as ints. Pointers may be null.
void obx_version(int* major, int* minor, int* patch);

/// Checks if the version of the library is equal to or higher than the given version ints.
bool obx_version_is_at_least(int major, int minor, int patch);

/// Returns the version of the library to be printed.
/// The format may change; to query for version use the int based methods instead.
const char* obx_version_string(void);

/// Returns the version of the ObjectBox core to be printed.
/// The format may change, do not rely on its current form.
const char* obx_version_core_string(void);

//----------------------------------------------
// Utilities
//----------------------------------------------

/// delete the store files from the given directory
int obx_remove_db_files(char const* directory);

/// checks whether functions returning OBX_bytes_array are fully supported (depends on build, invariant during runtime)
bool obx_supports_bytes_array(void);

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
#define OBX_ERROR_ID_ALREADY_EXISTS 10210
#define OBX_ERROR_ID_NOT_FOUND 10211
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

/// A requested schema object (e.g. entity or property) was not found in the schema
#define OBX_ERROR_SCHEMA_OBJECT_NOT_FOUND 10503

//----------------------------------------------
// Common types
//----------------------------------------------
/// Schema entity & property identifiers
typedef uint32_t obx_schema_id;

/// Universal identifier used in schema for entities & properties
typedef uint64_t obx_uid;

/// ID of a single Object stored in the database
typedef uint64_t obx_id;

/// Error/success code returned by an obx_* function; see defines OBX_SUCCESS, OBX_NOT_FOUND, and OBX_ERROR_*
typedef int obx_err;

/// The callback for reading data one-by-one
/// @param user_data is a pass-through argument passed to the called API
/// @param data is the read data buffer
/// @param size specifies the length of the read data
/// @return true to keep going, false to cancel.
typedef bool obx_data_visitor(void* user_data, const void* data, size_t size);

//----------------------------------------------
// Error info; obx_last_error_*
//----------------------------------------------

/// Returns the error status on the current thread and clears the error state.
/// The buffer returned in out_message is valid only until the next call into ObjectBox.
/// @param out_error receives the error code; may be NULL
/// @param out_message receives the pointer to the error messages; may be NULL
/// @returns true if an error was pending
bool obx_last_error_pop(obx_err* out_error, const char** out_message);

/// The last error raised by an ObjectBox API call on the current thread, or OBX_SUCCESS if no error occured yet.
/// Note that API calls do not clear this error code (also true for this method).
/// Thus, if you receive an error from this, it's usually a good idea to call obx_last_error_clear() to clear the error
/// state (or use obx_last_error_pop()) for future API calls.
obx_err obx_last_error_code(void);

/// The error message string attached to the error returned by obx_last_error_code().
/// Like obx_last_error_code(), this is bound to the current thread, and this call does not clear the error state.
/// The buffer returned is valid only until the next call into ObjectBox.
const char* obx_last_error_message(void);

/// The underlying error for the error returned by obx_last_error_code(). Where obx_last_error_code() may be a generic
/// error like OBX_ERROR_STORAGE_GENERAL, this will give a further underlying and possibly platform-specific error code.
obx_err obx_last_error_secondary(void);

/// Clears the error state on the current thread; e.g. obx_last_error_code() will now return OBX_SUCCESS.
/// Note that clearing the error state does not happen automatically;
/// API calls set the error state when they produce an error, but do not clear it on success.
/// See also: obx_last_error_pop() to retrieve the error state and clear it.
void obx_last_error_clear(void);

//----------------------------------------------
// Model
//----------------------------------------------

typedef enum {
    OBXPropertyType_Bool = 1,   ///< 1 byte
    OBXPropertyType_Byte = 2,   ///< 1 byte
    OBXPropertyType_Short = 3,  ///< 2 bytes
    OBXPropertyType_Char = 4,   ///< 1 byte
    OBXPropertyType_Int = 5,    ///< 4 bytes
    OBXPropertyType_Long = 6,   ///< 8 bytes
    OBXPropertyType_Float = 7,  ///< 4 bytes
    OBXPropertyType_Double = 8, ///< 8 bytes
    OBXPropertyType_String = 9,
    OBXPropertyType_Date = 10,  ///< Unix timestamp (milliseconds since 1970) in 8 bytes
    OBXPropertyType_Relation = 11,
    OBXPropertyType_ByteVector = 23,
    OBXPropertyType_StringVector = 30,
} OBXPropertyType;

/// Bit-flags defining the behavior of properties.
/// Note: Numbers indicate the bit position
typedef enum {
    /// 64 bit long property (internally unsigned) representing the ID of the entity.
    /// May be combined with: NON_PRIMITIVE_TYPE, ID_MONOTONIC_SEQUENCE, ID_SELF_ASSIGNABLE.
    OBXPropertyFlags_ID = 1,

    /// On languages like Java, a non-primitive type is used (aka wrapper types, allowing null)
    OBXPropertyFlags_NON_PRIMITIVE_TYPE = 2,

    /// Unused yet
    OBXPropertyFlags_NOT_NULL = 4,

    OBXPropertyFlags_INDEXED = 8,

    /// Unused yet
    OBXPropertyFlags_RESERVED = 16,

    /// Unique index
    OBXPropertyFlags_UNIQUE = 32,

    /// Unused yet: Use a persisted sequence to enforce ID to rise monotonic (no ID reuse)
    OBXPropertyFlags_ID_MONOTONIC_SEQUENCE = 64,

    /// Allow IDs to be assigned by the developer
    OBXPropertyFlags_ID_SELF_ASSIGNABLE = 128,

    /// Unused yet
    OBXPropertyFlags_INDEX_PARTIAL_SKIP_NULL = 256,

    /// used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
    OBXPropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512,

    /// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
    OBXPropertyFlags_VIRTUAL = 1024,

    /// Index uses a 32 bit hash instead of the value
    /// 32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions
    OBXPropertyFlags_INDEX_HASH = 2048,

    /// Index uses a 64 bit hash instead of the value
    /// recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash
    OBXPropertyFlags_INDEX_HASH64 = 4096,

    /// Unused yet: While our default are signed ints, queries & indexes need do know signing info.
    /// Note: Don't combine with ID (IDs are always unsigned internally).
    /// Used in combination with integer types defined in OBXPropertyType_*.
    OBXPropertyFlags_UNSIGNED = 8192,
} OBXPropertyFlags;

struct OBX_model;
typedef struct OBX_model OBX_model;

/// Create an (empty) data meta model which is to be consumed by obx_opt_model().
/// Model creation is usually done by language bindings, which automatically build the model based on parsed source code
/// (for examples, see ObjectBox Go or Swift, which also use this C API).
///
/// For manual creation, these are the basic steps:
/// - define entity types using obx_model_entity() and obx_model_property()
/// - Pass the last ever used IDs with obx_model_last_entity_id(), obx_model_last_index_id(), obx_model_last_relation_id()
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details.
///               Note that obx_model_* functions handle OBX_model NULL pointers (will indicate an error but not crash).
OBX_model* obx_model(void);

/// Only call when not calling obx_store_open() (which will free it internally)
/// @param model NULL-able; returns OBX_SUCCESS if model is NULL
obx_err obx_model_free(OBX_model* model);

/// @param model NULL-able; returns OBX_ERROR_ILLEGAL_ARGUMENT if model is NULL
obx_err obx_model_error_code(OBX_model* model);

/// @param model NULL-able; returns NULL if model is NULL
const char* obx_model_error_message(OBX_model* model);

/// Starts the definition of a new entity type for the meta data model.
/// After this, call obx_model_property() to add properties to the entity type.
obx_err obx_model_entity(OBX_model* model, const char* name, obx_schema_id entity_id, obx_uid entity_uid);

/// Starts the definition of a new property for the entity type of the last obx_model_entity() call.
obx_err obx_model_property(OBX_model* model, const char* name, OBXPropertyType type, obx_schema_id property_id,
                           obx_uid property_uid);

/// Refines the property definition of the property of the last obx_model_property() call.
obx_err obx_model_property_flags(OBX_model* model, OBXPropertyFlags flags);

/// Refines the property definition of the property of the last obx_model_property() call.
obx_err obx_model_property_relation(OBX_model* model, const char* target_entity, obx_schema_id index_id,
                                    obx_uid index_uid);

/// Refines the property definition of the property of the last obx_model_property() call.
obx_err obx_model_property_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

/// Add a standalone relation between the active entity and the target entity to the model
obx_err obx_model_relation(OBX_model* model, obx_schema_id relation_id, obx_uid relation_uid, obx_schema_id target_id,
                           obx_uid target_uid);

void obx_model_last_entity_id(OBX_model*, obx_schema_id entity_id, obx_uid entity_uid);

void obx_model_last_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

void obx_model_last_relation_id(OBX_model* model, obx_schema_id relation_id, obx_uid relation_uid);

obx_err obx_model_entity_last_property_id(OBX_model* model, obx_schema_id property_id, obx_uid property_uid);

//----------------------------------------------
// Store
//----------------------------------------------

struct OBX_store;
typedef struct OBX_store OBX_store;

struct OBX_store_options;
typedef struct OBX_store_options OBX_store_options;

typedef enum {
    OBXDebugFlags_LOG_TRANSACTIONS_READ = 1,
    OBXDebugFlags_LOG_TRANSACTIONS_WRITE = 2,
    OBXDebugFlags_LOG_QUERIES = 4,
    OBXDebugFlags_LOG_QUERY_PARAMETERS = 8,
    OBXDebugFlags_LOG_ASYNC_QUEUE = 16,
} OBXDebugFlags;

typedef struct OBX_bytes {
    const void* data;
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

typedef struct OBX_string_array {
    const char** items;
    size_t count;
} OBX_string_array;

typedef struct OBX_int64_array {
    const int64_t* items;
    size_t count;
} OBX_int64_array;

typedef struct OBX_int32_array {
    const int32_t* items;
    size_t count;
} OBX_int32_array;

typedef struct OBX_int16_array {
    const int16_t* items;
    size_t count;
} OBX_int16_array;

typedef struct OBX_int8_array {
    const int8_t* items;
    size_t count;
} OBX_int8_array;

typedef struct OBX_double_array {
    const double* items;
    size_t count;
} OBX_double_array;

typedef struct OBX_float_array {
    const float* items;
    size_t count;
} OBX_float_array;

/// Create a default set of store options.
/// @returns NULL on failure, a default set of options on success
OBX_store_options* obx_opt();

/// Set the store directory on the options. The default is "objectbox".
obx_err obx_opt_directory(OBX_store_options* opt, const char* dir);

/// Set the maximum db size on the options. The default is 1Gb.
void obx_opt_max_db_size_in_kb(OBX_store_options* opt, size_t size_in_kb);

/// Set the file mode on the options. The default is 0755 (unix-style)
void obx_opt_file_mode(OBX_store_options* opt, int file_mode);

/// Set the maximum number of readers on the options.
/// "Readers" are an finite resource for which we need to define a maximum number upfront.
/// The default value is enough for most apps and usually you can ignore it completely.
/// However, if you get the OBX_ERROR_MAX_READERS_EXCEEDED error, you should verify your threading.
/// For each thread, ObjectBox uses multiple readers.
/// Their number (per thread) depends on number of types, relations, and usage patterns.
/// Thus, if you are working with many threads (e.g. in a server-like scenario), it can make sense to increase the
/// maximum number of readers.
/// Note: The internal default is currently around 120. So when hitting this limit, try values around 200-500.
void obx_opt_max_readers(OBX_store_options* opt, int max_readers);

/// Set the model on the options. The default is no model.
/// NOTE: the model is always freed by this function, including when an error occurs.
obx_err obx_opt_model(OBX_store_options* opt, OBX_model* model);

/// Set the model on the options copying the given bytes. The default is no model.
obx_err obx_opt_model_bytes(OBX_store_options* opt, const void* bytes, size_t size);

/// Like obx_opt_model_bytes BUT WITHOUT copying the given bytes.
/// Thus, you must keep the bytes available until the store was created.
obx_err obx_opt_model_bytes_direct(OBX_store_options* opt, const void* bytes, size_t size);

/// Free the options.
/// Note: Only free *unused* options, obx_store_open() frees the options internally
void obx_opt_free(OBX_store_options* opt);

/// Note: the given options are always freed by this function, including when an error occurs.
/// @param opt required parameter holding the data model (obx_opt_model()) and optional options (see obx_opt_*())
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_store* obx_store_open(OBX_store_options* opt);

obx_schema_id obx_store_entity_id(OBX_store* store, const char* entity_name);

obx_schema_id obx_store_entity_property_id(OBX_store* store, obx_schema_id entity_id, const char* property_name);

/// Awaits for all (including future) async submissions to be completed (the async queue becomes idle for a moment).
/// @returns true if all submissions were completed or async processing was not started; false if shutting down
/// @returns false if shutting down or an error occurred
bool obx_store_await_async_completion(OBX_store* store);

/// Awaits for previously submitted async operations to be completed (the async queue does not have to become idle).
/// @returns true if all submissions were completed or async processing was not started
/// @returns false if shutting down or an error occurred
bool obx_store_await_async_submitted(OBX_store* store);

obx_err obx_store_debug_flags(OBX_store* store, OBXDebugFlags flags);

obx_err obx_store_close(OBX_store* store);

//----------------------------------------------
// Transaction
//----------------------------------------------

struct OBX_txn;
typedef struct OBX_txn OBX_txn;

/// Creates a write transaction (read and write).
/// Transaction creation can be nested (recursive), however only the outermost transaction is relevant on the DB level.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if called when inside a read transaction.
OBX_txn* obx_txn_write(OBX_store* store);

/// Creates a read transaction (read only).
/// Transaction creation can be nested (recursive), however only the outermost transaction is relevant on the DB level.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_txn* obx_txn_read(OBX_store* store);

/// "Finishes" this write transaction successfully and closes it; performs a commit if this is the top level
/// transaction and all inner transactions (if any) were also successful (obx_txn_success() were called on them).
/// Because this also closes the given transaction, the given OBX_txn pointer must not be used afterwards.
/// @return OBX_ERROR_ILLEGAL_STATE if the given transaction is not a write transaction.
obx_err obx_txn_success(OBX_txn* txn);

/// Closes (deletes) the transaction (read or write); the given OBX_txn pointer must not be used afterwards.
/// While this is the only way to release read transactions, this call is also an alternative to call obx_txn_success()
/// on write transactions.
/// In combination with obx_txn_mark_success(), this potentially commits or aborts a write transaction on the DB:
/// 1) If it's an outermost TX and all (inner) TXs were marked successful, this commits the transaction.
/// 2) If this transaction was not marked successful, this aborts the transaction (even if it's an inner TX).
/// If an error is returned (e.g. commit failed because DB is full), you can assume that the transaction was closed.
obx_err obx_txn_close(OBX_txn* txn);

/// Aborts the underlying transaction immediately and thus frees DB resources.
/// Only obx_txn_close() is allowed to be called on the transaction after calling this.
obx_err obx_txn_abort(OBX_txn* txn);  // Internal note: will make more sense once we introduce obx_txn_reset

/// Marks the given write transaction as successful or failed.
/// You can call this method multiple times with different values before calling obx_txn_close() on the transaction.
/// @return OBX_ERROR_ILLEGAL_STATE if the given transaction is not a write transaction.
obx_err obx_txn_mark_success(OBX_txn* txn, bool wasSuccessful);

//------------------------------------------------------------------
// Cursor (lower level API, check also the more convenient Box API)
//------------------------------------------------------------------

struct OBX_cursor;
typedef struct OBX_cursor OBX_cursor;

typedef enum {
    /// Standard put ("insert or update")
    OBXPutMode_PUT = 1,

    /// Put succeeds only if the entity does not exist yet.
    OBXPutMode_INSERT = 2,

    /// Put succeeds only if the entity already exist.
    OBXPutMode_UPDATE = 3,

    // Not used yet (does not make sense for asnyc puts)
    // The given ID (non-zero) is guaranteed to be new; don't use unless you know exactly what you are doing!
    // This is primarily used internally. Wrong usage leads to inconsistent data (e.g. index data not updated)!
    // OBXPutMode_PUT_ID_GUARANTEED_TO_BE_NEW = 4

} OBXPutMode;

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_cursor* obx_cursor(OBX_txn* txn, obx_schema_id entity_id);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_cursor* obx_cursor2(OBX_txn* txn, const char* entity_name);

obx_err obx_cursor_close(OBX_cursor* cursor);

/// Call this when putting an object to generate an ID for it.
/// @param id_or_zero The ID of the entity. If you pass 0, this will generate a new one.
/// @seealso obx_box_id_for_put()
obx_id obx_cursor_id_for_put(OBX_cursor* cursor, obx_id id_or_zero);

/// ATTENTION: ensure that the given value memory is allocated to the next 4 bytes boundary.
/// ObjectBox needs to store bytes with sizes dividable by 4 for internal reasons.
/// Use obx_cursor_put_padded() otherwise.
/// @param id non-zero
obx_err obx_cursor_put(OBX_cursor* cursor, obx_id id, const void* data, size_t size, bool checkForPreviousValue);

/// Prefer obx_cursor_put (non-padded) if possible, as this does a memcpy if the size is not dividable by 4.
obx_err obx_cursor_put_padded(OBX_cursor* cursor, obx_id id, const void* data, size_t size, bool checkForPreviousValue);

obx_err obx_cursor_get(OBX_cursor* cursor, obx_id id, void** data, size_t* size);

/// Gets all objects as bytes.
/// For bigger quantities, it's recommended to iterate using obx_cursor_first and obx_cursor_next.
/// However, if the calling overhead is high (e.g. for language bindings), this method helps.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_bytes_array* obx_cursor_get_all(OBX_cursor* cursor);

obx_err obx_cursor_first(OBX_cursor* cursor, void** data, size_t* size);

obx_err obx_cursor_next(OBX_cursor* cursor, void** data, size_t* size);

obx_err obx_cursor_seek(OBX_cursor* cursor, obx_id id);

obx_err obx_cursor_current(OBX_cursor* cursor, void** data, size_t* size);

obx_err obx_cursor_remove(OBX_cursor* cursor, obx_id id);

obx_err obx_cursor_remove_all(OBX_cursor* cursor);

/// Count the number of available objects
obx_err obx_cursor_count(OBX_cursor* cursor, uint64_t* count);

/// Count the number of available objects up to the specified maximum
obx_err obx_cursor_count_max(OBX_cursor* cursor, uint64_t max_count, uint64_t* out_count);

/// Results true if there is no object available (false if at least one object is available)
obx_err obx_cursor_is_empty(OBX_cursor* cursor, bool* out_is_empty);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_bytes_array* obx_cursor_backlink_bytes(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                           obx_id id);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_id_array* obx_cursor_backlink_ids(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                      obx_id id);

obx_err obx_cursor_rel_put(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id, obx_id target_id);
obx_err obx_cursor_rel_remove(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id, obx_id target_id);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_id_array* obx_cursor_rel_ids(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id);

//----------------------------------------------
// Box
//----------------------------------------------

struct OBX_box;
typedef struct OBX_box OBX_box;

/// Gets the box for the given entity type. A box may be used across threads.
/// Boxes are shared instances and managed by the store; so there's no need to close/free boxes manually.
OBX_box* obx_box(OBX_store* store, obx_schema_id entity_id);

/// Checks whether a given object exists in the box.
obx_err obx_box_contains(OBX_box* box, obx_id id, bool* out_contains);

/// Checks whether this box contains objects with all of the IDs given
/// @param out_contains is set to true if all of the IDs are present, otherwise false
obx_err obx_box_contains_many(OBX_box* box, const OBX_id_array* ids, bool* out_contains);

/// Fetch a single object from the box; must be called inside a (reentrant) transaction.
/// The exposed data comes directly from the OS to allow zero-copy access, which limits the data lifetime:
/// \attention The exposed data is only valid as long as the (top) transaction is still active and no write
/// \attention operation (e.g. put/remove) was executed. Accessing data after this is undefined behavior.
/// @returns OBX_ERROR_ILLEGAL_STATE if not inside of an active transaction (see obx_txn_read() and obx_txn_write())
obx_err obx_box_get(OBX_box* box, obx_id id, void** data, size_t* size);

/// Fetch multiple objects for the given IDs from the box; must be called inside a (reentrant) transaction.
/// \attention See obx_box_get() for important notes on the limited lifetime of the exposed data.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if not inside of an active transaction
///               (see obx_txn_read() and obx_txn_write())
OBX_bytes_array* obx_box_get_many(OBX_box* box, const OBX_id_array* ids);

/// Fetch all objects from the box; must be called inside a (reentrant) transaction.
/// NOTE: don't call this in 32 bit mode! Use obx_box_visit_all() instead.
/// \attention See obx_box_get() for important notes on the limited lifetime of the exposed data.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if not inside of an active transaction
///               (see obx_txn_read() and obx_txn_write())
OBX_bytes_array* obx_box_get_all(OBX_box* box);

/// Read given objects from the database in a single transaction.
/// Call the visitor() on each object, passing user_data, object data & size as arguments.
/// The given visitor must return true to keep receiving results, false to cancel.
/// If an object is not found, the visitor() is still called, passing NULL as data and a 0 as size.
obx_err obx_box_visit_many(OBX_box* box, const OBX_id_array* ids, obx_data_visitor* visitor, void* user_data);

/// Read all objects in a single transaction.
/// Calls the visitor() on each object, passing visitor_arg, object data & size as arguments.
/// The given visitor must return true to keep receiving results, false to cancel.
obx_err obx_box_visit_all(OBX_box* box, obx_data_visitor* visitor, void* user_data);

/// Reserve an ID for insertion
/// @param id_or_zero The ID of the entity. If you pass 0, this will generate a new one.
/// @seealso obx_cursor_id_for_put
obx_id obx_box_id_for_put(OBX_box* box, obx_id id_or_zero);

/// Reserve the given number of IDs for insertion.
/// @param count number of IDs to reserve, max 10000
/// @param out_first_id the first ID of the sequence as
/// @returns an error in case the required number of IDs could not be reserved.
obx_err obx_box_ids_for_put(OBX_box* box, uint64_t count, obx_id* out_first_id);

/// Put object synchronously (in a single write transaction)
/// ATTENTION: ensure that the given value memory is allocated to the next 4 bytes boundary.
/// ObjectBox needs to store bytes with sizes dividable by 4 for internal reasons.
obx_err obx_box_put(OBX_box* box, obx_id id, const void* data, size_t size, OBXPutMode mode);

/// Put all given objects in the database in a single transaction
/// @param ids Previously allocated IDs for the given given objects (e.g. using obx_box_ids_for_put)
obx_err obx_box_put_many(OBX_box* box, const OBX_bytes_array* objects, const obx_id* ids, OBXPutMode mode);

/// Remove a single object
/// will return OBX_NOT_FOUND if an object with the given ID doesn't exist
obx_err obx_box_remove(OBX_box* box, obx_id id);

/// Remove all given objects from the database in a single transaction.
/// Note that this method will not fail if the object is not found (e.g. already removed).
/// In case you need to strictly check whether all of the objects exist before removing them,
///  execute obx_box_contains_ids() and obx_box_remove_ids() inside a single write transaction.
/// @param out_count Pointer to retrieve the number of removed objects; may be NULL.
obx_err obx_box_remove_many(OBX_box* box, const OBX_id_array* ids, uint64_t* out_count);

/// Remove all objects and set the out_count the the number of removed objects.
/// @param out_count Pointer to retrieve the number of removed objects; may be NULL.
obx_err obx_box_remove_all(OBX_box* box, uint64_t* out_count);

/// Checks whether there are any objects for this entity and updates the out_is_empty accordingly
obx_err obx_box_is_empty(OBX_box* box, bool* out_is_empty);

/// Count the number of objects in the box, up to the given maximum.
/// You can pass limit=0 to count all objects without any limitation.
obx_err obx_box_count(OBX_box* box, uint64_t limit, uint64_t* out_count);

/// Fetch IDs of all objects that link back to the given object (ID) using the given relation property (ID).
/// Note: This method refers to "property based relations" unlike the "stand-alone relations" (see obx_box_rel_*).
/// @param property_id the relation property, which must belong to the entity type represented by this box
/// @param id object ID; the object's type is the target of the relation property (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_id_array* obx_box_get_backlink_ids(OBX_box* box, obx_schema_id property_id, obx_id id);

//----------------------------------------------
// Box + stand-alone relation; obx_box_rel_*
//----------------------------------------------

/// Insert a standalone relation entry between two objects.
/// @param relation_id must be a standalone relation ID with source entity belonging to this box
/// @param source_id identifies an object from this box
/// @param target_id identifies an object from the target box (as per the relation definition)
obx_err obx_box_rel_put(OBX_box* box, obx_schema_id relation_id, obx_id source_id, obx_id target_id);

/// Remove a standalone relation entry between two objects.
/// See obx_box_rel_put() for parameters documentation.
obx_err obx_box_rel_remove(OBX_box* box, obx_schema_id relation_id, obx_id source_id, obx_id target_id);

/// Fetch IDs of all objects in this Box related to the given object (typically from another Box).
/// Used for a stand-alone relation and its "regular" direction; this Box represents the target of the relation.
/// @param relation_id ID of a standalone relation, whose target type matches this Box
/// @param id object ID of the relation source type (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_id_array* obx_box_rel_get_ids(OBX_box* box, obx_schema_id relation_id, obx_id id);

/// Fetch IDs of all objects in this Box related to the given object (typically from another Box).
/// Used for a stand-alone relation and its "backlink" direction; this Box represents the source of the relation.
/// @param relation_id ID of a standalone relation, whose source type matches this Box
/// @param id object ID of the relation target type (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_id_array* obx_box_rel_get_backlink_ids(OBX_box* box, obx_schema_id relation_id, obx_id id);

/// Created by obx_box_async, used for async operations like obx_async_put.
struct OBX_async;
typedef struct OBX_async OBX_async;

//----------------------------------------------
// Async
//----------------------------------------------

/// Note: DO NOT close this OBX_async; its lifetime is tied to the OBX_box instance.
OBX_async* obx_async(OBX_box* box);

/// Puts asynchronously using the given mode.
obx_err obx_async_put_mode(OBX_async* async, obx_id id, const void* data, size_t size, OBXPutMode mode);

/// Puts asynchronously with standard put semantics (insert or update).
obx_err obx_async_put(OBX_async* async, obx_id id, const void* data, size_t size);

/// Puts asynchronously with inserts semantics (won't put if object already exists).
obx_err obx_async_insert(OBX_async* async, obx_id id, const void* data, size_t size);

/// Puts asynchronously with update semantics (won't put if object is not yet present).
obx_err obx_async_update(OBX_async* async, obx_id id, const void* data, size_t size);

/// Reserves an ID, which is returned immediately for future reference, and puts asynchronously.
/// Note: of course, it can NOT be guaranteed that the entity will actually be put successfully in the DB.
/// @param data the given bytes are mutated to update the contained ID data.
obx_id obx_async_id_put(OBX_async* async, void* data, size_t size);

/// Reserves an ID, which is returned immediately for future reference, and inserts asynchronously.
/// Note: of course, it can NOT be guaranteed that the entity will actually be inserted successfully in the DB.
/// @param data the given bytes are mutated to update the contained ID data.
obx_id obx_async_id_insert(OBX_async* async, void* data, size_t size);

/// Removes asynchronously.
obx_err obx_async_remove(OBX_async* async, obx_id id);

/// Note: for standard tasks, prefer obx_box_async() giving you a shared instance that does not have to be closed.
/// Creates a custom OBX_async instance that has to be closed using obx_async_close().
OBX_async* obx_async_create(OBX_box* box, uint64_t enqueueTimeoutMillis);

/// Closes a custom OBX_async instance created with obx_async_create().
/// @return OBX_ERROR_ILLEGAL_ARGUMENT if you pass the shared instance from obx_box_async()
obx_err obx_async_close(OBX_async* async);

//----------------------------------------------
// Query Builder
//----------------------------------------------

/// Not really an enum, but binary flags to use across languages
typedef enum {
    /// Reverts the order from ascending (default) to descending.
    OBXOrderFlags_DESCENDING = 1,

    /// Makes upper case letters (e.g. "Z") be sorted before lower case letters (e.g. "a").
    /// If not specified, the default is case insensitive for ASCII characters.
    OBXOrderFlags_CASE_SENSITIVE = 2,

    /// For scalars only: changes the comparison to unsigned (default is signed).
    OBXOrderFlags_UNSIGNED = 4,

    /// null values will be put last.
    /// If not specified, by default null values will be put first.
    OBXOrderFlags_NULLS_LAST = 8,

    /// null values should be treated equal to zero (scalars only).
    OBXOrderFlags_NULLS_ZERO = 16,
} OBXOrderFlags;

struct OBX_query_builder;
typedef struct OBX_query_builder OBX_query_builder;

/// Query Builder condition identifier
/// - returned by condition creating functions,
/// - used to combine conditions with any/all, thus building more complex conditions
typedef int obx_qb_cond;

/// Create a query builder which is used to collect conditions using the obx_qb_* functions.
/// Once all conditions are applied, use obx_query() to build a OBX_query that is used to actually retrieve data.
/// Use obx_qb_close() to close (free) the query builder.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_query_builder* obx_query_builder(OBX_store* store, obx_schema_id entity_id);

/// Closes the query builder; note that OBX_query objects outlive their builder and thus are not affected by this call.
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

obx_qb_cond obx_qb_strings_contain(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
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

obx_err obx_qb_param_alias(OBX_query_builder* builder, const char* alias);

obx_err obx_qb_order(OBX_query_builder* builder, obx_schema_id property_id, OBXOrderFlags flags);

/// Create a link based on a property-relation (many-to-one)
OBX_query_builder* obx_qb_link_property(OBX_query_builder* builder, obx_schema_id property_id);

/// Create a backlink based on a property-relation used in reverse (one-to-many)
OBX_query_builder* obx_qb_backlink_property(OBX_query_builder* builder, obx_schema_id source_entity_id,
                                            obx_schema_id source_property_id);

/// Create a link based on a standalone relation (many-to-many)
OBX_query_builder* obx_qb_link_standalone(OBX_query_builder* builder, obx_schema_id relation_id);

/// Create a backlink based on a standalone relation (many-to-many, reverse direction)
OBX_query_builder* obx_qb_backlink_standalone(OBX_query_builder* builder, obx_schema_id relation_id);

//----------------------------------------------
// Query
//----------------------------------------------
struct OBX_query;
typedef struct OBX_query OBX_query;

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_query* obx_query(OBX_query_builder* builder);
obx_err obx_query_close(OBX_query* query);

/// Finds entities matching the query; NOTE: the returned data is only valid as long the transaction is active!
OBX_bytes_array* obx_query_find(OBX_query* query, uint64_t offset, uint64_t limit);

/// Walks over matching objects using the given data visitor
obx_err obx_query_visit(OBX_query* query, obx_data_visitor* visitor, void* user_data, uint64_t offset,
                        uint64_t limit);

/// Returns IDs of all matching objects
OBX_id_array* obx_query_find_ids(OBX_query* query, uint64_t offset, uint64_t limit);

/// Returns the number of matching objects
obx_err obx_query_count(OBX_query* query, uint64_t* count);

/// Removes all matching objects from the database & returns the number of deleted objects
obx_err obx_query_remove(OBX_query* query, uint64_t* count);

/// the resulting char* is valid until another call on to_string is made on the same query or until the query is freed
const char* obx_query_describe(OBX_query* query);

/// the resulting char* is valid until another call on describe_parameters is made on the same query or until the query
/// is freed
const char* obx_query_describe_params(OBX_query* query);

//----------------------------------------------
// Query using Cursor (lower level API)
//----------------------------------------------
obx_err obx_query_cursor_visit(OBX_query* query, OBX_cursor* cursor, obx_data_visitor* visitor, void* user_data,
                               uint64_t offset, uint64_t limit);

/// Finds entities matching the query; NOTE: the returned data is only valid as long the transaction is active!
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_bytes_array* obx_query_cursor_find(OBX_query* query, OBX_cursor* cursor, uint64_t offset, uint64_t limit);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_id_array* obx_query_cursor_find_ids(OBX_query* query, OBX_cursor* cursor, uint64_t offset, uint64_t limit);
obx_err obx_query_cursor_count(OBX_query* query, OBX_cursor* cursor, uint64_t* count);

/// Removes (deletes!) all matching objects.
obx_err obx_query_cursor_remove(OBX_query* query, OBX_cursor* cursor, uint64_t* count);

//----------------------------------------------
// Query parameters (obx_query_{type}_param(s))
//----------------------------------------------
obx_err obx_query_string_param(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, const char* value);
obx_err obx_query_string_params_in(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                   const char* values[], int count);
obx_err obx_query_int_param(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, int64_t value);
obx_err obx_query_int_params(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, int64_t value_a,
                             int64_t value_b);
obx_err obx_query_int64_params_in(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                  const int64_t values[], int count);
obx_err obx_query_int32_params_in(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                  const int32_t values[], int count);
obx_err obx_query_double_param(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, double value);
obx_err obx_query_double_params(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, double value_a,
                                double value_b);
obx_err obx_query_bytes_param(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id, const void* value,
                              size_t size);

obx_err obx_query_string_param_alias(OBX_query* query, const char* alias, const char* value);
obx_err obx_query_string_params_in_alias(OBX_query* query, const char* alias, const char* values[], int count);
obx_err obx_query_int_param_alias(OBX_query* query, const char* alias, int64_t value);
obx_err obx_query_int_params_alias(OBX_query* query, const char* alias, int64_t value_a, int64_t value_b);
obx_err obx_query_int64_params_in_alias(OBX_query* query, const char* alias, const int64_t values[], int count);
obx_err obx_query_int32_params_in_alias(OBX_query* query, const char* alias, const int32_t values[], int count);
obx_err obx_query_double_param_alias(OBX_query* query, const char* alias, double value);
obx_err obx_query_double_params_alias(OBX_query* query, const char* alias, double value_a, double value_b);
obx_err obx_query_bytes_param_alias(OBX_query* query, const char* alias, const void* value, size_t size);

//----------------------------------------------
// Property-Query - getting a single property instead of whole objects
//----------------------------------------------
struct OBX_query_prop;
typedef struct OBX_query_prop OBX_query_prop;

/// Creates a "property query" with results referring to single property (not complete objects).
/// Also provides aggregates like for example obx_query_prop_avg().
OBX_query_prop* obx_query_prop(OBX_query* query, obx_schema_id property_id);

obx_err obx_query_prop_close(OBX_query_prop* query);

/// Configures the property query to work only on distinct values
/// @note not all methods support distinct, those that don't will return an error
obx_err obx_query_prop_distinct(OBX_query_prop* query, bool distinct);

/// Configures the property query to work only on distinct values
/// @note not all methods support distinct, those that don't will return an error
obx_err obx_query_prop_distinct_string(OBX_query_prop* query, bool distinct, bool case_sensitive);

/// Count the number of non-NULL values of the given property across all objects matching the query
obx_err obx_query_prop_count(OBX_query_prop* query, uint64_t* out_count);

/// Calculate an average value for the given numeric property across all objects matching the query
obx_err obx_query_prop_avg(OBX_query_prop* query, double* out_average);

/// Find the minimum value of the given floating-point property across all objects matching the query
obx_err obx_query_prop_min(OBX_query_prop* query, double* out_minimum);

/// Find the maximum value of the given floating-point property across all objects matching the query
obx_err obx_query_prop_max(OBX_query_prop* query, double* out_maximum);

/// Calculate the sum of the given floating-point property across all objects matching the query
obx_err obx_query_prop_sum(OBX_query_prop* query, double* out_sum);

/// Find the minimum value of the given property across all objects matching the query
obx_err obx_query_prop_min_int(OBX_query_prop* query, int64_t* out_minimum);

/// Find the maximum value of the given property across all objects matching the query
obx_err obx_query_prop_max_int(OBX_query_prop* query, int64_t* out_maximum);

/// Calculate the sum of the given property across all objects matching the query
obx_err obx_query_prop_sum_int(OBX_query_prop* query, int64_t* out_sum);

/// Returns an array of strings stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified field are skipped
OBX_string_array* obx_query_prop_string_find(OBX_query_prop* query, const char* value_if_null);

/// Returns an array of ints stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_int64_array* obx_query_prop_int64_find(OBX_query_prop* query, const int64_t* value_if_null);

/// Returns an array of ints stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_int32_array* obx_query_prop_int32_find(OBX_query_prop* query, const int32_t* value_if_null);

/// Returns an array of ints stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_int16_array* obx_query_prop_int16_find(OBX_query_prop* query, const int16_t* value_if_null);

/// Returns an array of ints stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_int8_array* obx_query_prop_int8_find(OBX_query_prop* query, const int8_t* value_if_null);

/// Returns an array of doubles stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_double_array* obx_query_prop_double_find(OBX_query_prop* query, const double* value_if_null);

/// Returns an array of int stored as the given property across all objects matching the query
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_float_array* obx_query_prop_float_find(OBX_query_prop* query, const float* value_if_null);

//===================================================================================================================
struct OBX_observer;
//===================================================================================================================
/// Observers are called back when data has changed in the database.
/// See obx_observe(), or obx_observe_single_type() to listen to a changes that affect a single entity type
typedef struct OBX_observer OBX_observer;

/// Callback for obx_observe()
/// @param user_data user data given to obx_observe()
/// @param type_ids IDs of object types that had changes
/// @param type_ids_count number of IDs of type_ids
typedef void obx_observer(void* user_data, const obx_schema_id* type_ids, int type_ids_count);

/// Callback for obx_observe_single_type()
typedef void obx_observer_single_type(void* user_data);

/// Create an observer (callback) to be notified about all data changes (for all object types).
/// The callback is invoked right after a successful commit.
/// \par Threading note
/// The given callback is called on the thread issuing the commit for the data change, e.g. via obx_txn_success().
/// Future versions might change that to a background thread, so be careful with threading assumptions.
/// Also, it's a usually good idea to make the callback return quickly to let the calling thread continue.
/// \attention Currently, you can not call any data operations from inside the call back.
/// \attention More accurately, no transaction may be strated. (This restriction may be removed in a later version.)
/// @param user_data any value you want to be forwarded to the given observer callback (usually some context info).
/// @param callback pointer to be called when the observed data changes.
OBX_observer* obx_observe(OBX_store* store, obx_observer* callback, void* user_data);

/// Create an observer (callback) to be notified about data changes for a given object type.
/// The callback is invoked right after a successful commit.
/// \note  If you intend to observe more than one type, it is more efficient to use obx_observe().
/// \par Threading note
/// The given callback is called on the thread issuing the commit for the data change, e.g. via obx_txn_success().
/// Future versions might change that to a background thread, so be careful with threading assumptions.
/// Also, it's a usually good idea to make the callback return quickly to let the calling thread continue.
/// \attention Currently, you can not call any data operations from inside the call back.
/// \attention More accurately, no transaction may be strated. (This restriction may be removed in a later version.)
/// @param type_id ID of the object type to be observer.
/// @param user_data any value you want to be forwarded to the given observer callback (usually some context info).
/// @param callback pointer to be called when the observed data changes.
OBX_observer* obx_observe_single_type(OBX_store* store, obx_schema_id type_id, obx_observer_single_type* callback,
                                      void* user_data);

/// Free the memory used by the given observer and unsubscribe it from its box or query.
void obx_observer_close(OBX_observer* observer);

//----------------------------------------------
// Utilities for bytes/ids/arrays
//----------------------------------------------
void obx_bytes_free(OBX_bytes* bytes);

/// Allocates a bytes array struct of the given size, ready for the data to be pushed
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_bytes_array* obx_bytes_array(size_t count);

/// Sets the given data as the index in the bytes array. The data is not copied, just referenced through the pointer
obx_err obx_bytes_array_set(OBX_bytes_array* array, size_t index, const void* data, size_t size);

/// Frees the bytes array struct
void obx_bytes_array_free(OBX_bytes_array* array);

/// Creates an ID array struct, copying the given IDs as the contents
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_id_array* obx_id_array(const obx_id* ids, size_t count);

/// Frees the array struct
void obx_id_array_free(OBX_id_array* array);

/// Frees the array struct
void obx_string_array_free(OBX_string_array* array);

/// Frees the array struct
void obx_int64_array_free(OBX_int64_array* array);

/// Frees the array struct
void obx_int32_array_free(OBX_int32_array* array);

/// Frees the array struct
void obx_int16_array_free(OBX_int16_array* array);

/// Frees the array struct
void obx_int8_array_free(OBX_int8_array* array);

/// Frees the array struct
void obx_double_array_free(OBX_double_array* array);

/// Frees the array struct
void obx_float_array_free(OBX_float_array* array);

/// Only for Apple platforms: set the prefix to use for mutexes based on POSIX semaphores.
/// You must supply the application group identifier for sand-boxed macOS apps here; see also:
/// https://developer.apple.com/library/archive/documentation/Security/Conceptual/AppSandboxDesignGuide/AppSandboxInDepth/AppSandboxInDepth.html#//apple_ref/doc/uid/TP40011183-CH3-SW24
void obx_posix_sem_prefix_set(const char* prefix);

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_H
