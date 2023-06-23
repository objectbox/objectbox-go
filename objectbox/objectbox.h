/*
 * Copyright 2018-2023 ObjectBox Ltd. All rights reserved.
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

/**
 * @defgroup c ObjectBox C API
 * @{
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

// NOLINTBEGIN(modernize-use-using)

#ifndef OBX_C_API
#define OBX_C_API
#endif

//----------------------------------------------
// ObjectBox version codes
//----------------------------------------------

/// When using ObjectBox as a dynamic library, you should verify that a compatible version was linked using
/// obx_version() or obx_version_is_at_least().
#define OBX_VERSION_MAJOR 0
#define OBX_VERSION_MINOR 18
#define OBX_VERSION_PATCH 1  // values >= 100 are reserved for dev releases leading to the next minor/major increase

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
/// @param data is the read data buffer
/// @param size specifies the length of the read data
/// @param user_data is a pass-through argument passed to the called API
/// @return true to keep going, false to cancel.
typedef bool obx_data_visitor(const void* data, size_t size, void* user_data);

//----------------------------------------------
// Runtime library information
//
// Functions in this group provide information about the loaded ObjectBox library.
// Their return values are invariable during runtime - they depend solely on the loaded library and its build settings.
//----------------------------------------------

/// Return the (runtime) version of the library as ints. Pointers may be null
OBX_C_API void obx_version(int* major, int* minor, int* patch);

/// Check if the (runtime) version of the library is equal to or higher than the given version ints.
OBX_C_API bool obx_version_is_at_least(int major, int minor, int patch);

/// Return the (runtime) version of the library to be printed.
/// The current format is "major.minor.patch" (e.g. "1.0.0") but may change in any future release.
/// Thus, only use for information purposes.
/// @see obx_version() and obx_version_is_at_least() for integer based versions
OBX_C_API const char* obx_version_string(void);

/// Return the version of the ObjectBox core to be printed (currently also contains a version date and features).
/// The format may change in any future release; only use for information purposes.
OBX_C_API const char* obx_version_core_string(void);

typedef enum {
    /// Functions that are returning multiple results (e.g. multiple objects) can be only used if this is available.
    /// This is only available for 64-bit OSes and is the opposite of "chunked mode", which forces to consume results
    /// in chunks (e.g. one by one).
    /// Since chunked mode consumes a bit less RAM, ResultArray style functions are typically only preferable if
    /// there's an additional overhead per call, e.g. caused by a higher level language abstraction like CGo.
    OBXFeature_ResultArray = 1,

    /// TimeSeries support (date/date-nano companion ID and other time-series functionality).
    OBXFeature_TimeSeries = 2,

    /// Sync client availability. Visit https://objectbox.io/sync for more details.
    OBXFeature_Sync = 3,

    /// Check whether debug log can be enabled during runtime.
    OBXFeature_DebugLog = 4,

    /// Administration interface (HTTP server) with a database browser.
    OBXFeature_Admin = 5,

    /// Tree & GraphQL support
    OBXFeature_Tree = 6,

    /// Sync server availability. Visit https://objectbox.io/sync for more details.
    OBXFeature_SyncServer = 7,
} OBXFeature;

/// Checks whether the given feature is available in the currently loaded library.
OBX_C_API bool obx_has_feature(OBXFeature feature);

//----------------------------------------------
// Utilities
//----------------------------------------------

/// Log level as passed to obx_log_callback.
typedef enum {
    /// Log level for verbose messages (not emitted at the moment)
    OBXLogLevel_Verbose = 10,

    /// Log level for debug messages (may be limited to special debug builds)
    OBXLogLevel_Debug = 20,

    /// Log level for info messages
    OBXLogLevel_Info = 30,

    /// Log level for warning messages
    OBXLogLevel_Warn = 40,

    /// Log level for error messages
    OBXLogLevel_Error = 50,
} OBXLogLevel;

/// Callback for logging, which can be provided to store creation via options.
typedef void obx_log_callback(OBXLogLevel log_level, const char* message, size_t message_size, void* user_data);

/// To be used for putting objects with prepared ID slots, e.g. obx_cursor_put_object().
#define OBX_ID_NEW 0xFFFFFFFFFFFFFFFF

/// Delete the store files from the given directory
OBX_C_API obx_err obx_remove_db_files(char const* directory);

/// @returns the file size of the main database file, or 0 if the file does not exist or some error occurred.
OBX_C_API size_t obx_db_file_size(char const* directory);

/// Enable (or disable) debug logging for ObjectBox internals.
/// This requires a version of the library with the DebugLog feature.
/// You can check if the feature is available with obx_has_feature(OBXFeature_DebugLog).
OBX_C_API obx_err obx_debug_log(bool enabled);

/// Checks if debug logs are enabled for ObjectBox internals. This depends on the availability of the DebugLog feature.
/// If the feature is available, it returns the current state, which is adjustable via obx_debug_log().
/// Otherwise, it always returns false for standard release builds (or true if you are having a special debug version).
OBX_C_API bool obx_debug_log_enabled();

/// Gets the number, as used by ObjectBox, of the current thread.
/// This e.g. allows to "associate" the thread with ObjectBox logs (each log entry contains the thread number).
OBX_C_API int obx_thread_number();

//----------------------------------------------
// Return and error codes
//----------------------------------------------

/// Value returned when no error occurred (0)
#define OBX_SUCCESS 0

/// Returned by, e.g., get operations if nothing was found for a specific ID.
/// This is NOT an error condition, and thus no "last error" info (code/text) is set.
#define OBX_NOT_FOUND 404

/// Indicates that a function had "no success", which is typically a likely outcome and not a "hard error".
/// This is NOT an error condition, and thus no "last error" info is set.
#define OBX_NO_SUCCESS 1001

/// Indicates that a function reached a time out, which is typically a likely outcome and not a "hard error".
/// This is NOT an error condition, and thus no "last error" info is set.
#define OBX_TIMEOUT 1002

// --- General errors --------------------------------------------------------

/// A request does not make sense in the current state. For example, doing actions on a closed object.
#define OBX_ERROR_ILLEGAL_STATE 10001

/// The passed arguments were illegal; e.g. passed null to a non-null argument.
#define OBX_ERROR_ILLEGAL_ARGUMENT 10002

/// A resource could not be allocated. This usually indicates a severe system status, e.g running out of memory.
#define OBX_ERROR_ALLOCATION 10003

/// Thrown when a numeric value did overflow/underflow, e.g. an addition would result in a value of the type's limits.
#define OBX_ERROR_NUMERIC_OVERFLOW 10004

/// Tried to use a special feature that is not part of this lib edition.
/// The ObjectBox library comes in different flavors, and may not contain all features.
#define OBX_ERROR_FEATURE_NOT_AVAILABLE 10005

/// A required resource (e.g. a store or some internal queue) is shutting down, and thus, the current request could not
/// be fulfilled. Typically occurs in asynchronous operations.
#define OBX_ERROR_SHUTTING_DOWN 10006

/// An error occurred but error information could not be extracted.
/// This falls in the "should not happen" category and should be very rare;
/// please report if you ever encounter this in a reproducible fashion.
#define OBX_ERROR_NO_ERROR_INFO 10097

/// A general error occurred that did not fit in a more specific category.
#define OBX_ERROR_GENERAL 10098

/// An error of unknown type occurred.
/// This falls in the "should not happen" category and should be very rare;
/// please report if you ever encounter this in a reproducible fashion.
#define OBX_ERROR_UNKNOWN 10099

// --- Storage errors (often have a secondary error code) --------------------

/// The database has reached its given storage limit.
/// Note that is limit is user defined and thus you can set it as an option when opening a database via
/// obx_opt_max_db_size_in_kb(). Further note that this limit is in place to protect the system from unreasonable
/// storage consumption due to unexpected usage and user errors (e.g. inserting data in an infinite loop).
#define OBX_ERROR_DB_FULL 10101

/// Thrown when the maximum of readers (read transactions) was exceeded.
/// Verify that you run a reasonable amount of threads only.
/// If you actually intend to work with a very high number of threads (>100), consider increasing the number of maximum
/// readers via obx_opt_max_readers().
#define OBX_ERROR_MAX_READERS_EXCEEDED 10102

/// Can not continue, ObjectBox store must be closed.
/// This falls in the "should not happen" category and should be very rare;
/// please report if you ever encounter this in a reproducible fashion.
#define OBX_ERROR_STORE_MUST_SHUTDOWN 10103

/// Occurs when a transaction is about to commit but it would exceed the user-defined data size limit.
/// See obx_opt_max_data_size_in_kb() for details.
#define OBX_ERROR_MAX_DATA_SIZE_EXCEEDED 10104

/// A general database related error occurred that did not fit in a more specific category.
/// No secondary error code is available.
#define OBX_ERROR_DB_GENERAL 10198

/// A storage error occurred that did not fit in a more specific (storage) category.
/// A secondary error code is available via obx_last_error_secondary().
#define OBX_ERROR_STORAGE_GENERAL 10199

// --- Data errors -----------------------------------------------------------

/// A unique constraint of was violated when trying to put.
#define OBX_ERROR_UNIQUE_VIOLATED 10201

/// A result was supposed to be unique, but was not. E.g. a query results in multiple results where one one is expected.
#define OBX_ERROR_NON_UNIQUE_RESULT 10202

/// A Property did not match an expected type; e.g. when defining a string-based query condition for an integer type.
#define OBX_ERROR_PROPERTY_TYPE_MISMATCH 10203

/// A given ID already existed, e.g. when trying to insert with a pre-existing ID.
/// Should happen only with user assigned IDs.
#define OBX_ERROR_ID_ALREADY_EXISTS 10210

/// A given ID was not found (does not exist), e.g. when trying to update an object that does not exist.
/// Typically happens when an object was removed before while trying to update it.
#define OBX_ERROR_ID_NOT_FOUND 10211

/// A Time Series operation failed, e.g. a time value was out of range
#define OBX_ERROR_TIME_SERIES 10212

/// A constraint violation occurred that did not fit in a more specific (constraint violation) category.
#define OBX_ERROR_CONSTRAINT_VIOLATED 10299

// --- STD errors --------------------------------------------------------

// These are low-level errors falling in the "should not happen" category and should be very rare;
// please report if you ever encounter one of these in a reproducible fashion.

#define OBX_ERROR_STD_ILLEGAL_ARGUMENT 10301
#define OBX_ERROR_STD_OUT_OF_RANGE 10302
#define OBX_ERROR_STD_LENGTH 10303
#define OBX_ERROR_STD_BAD_ALLOC 10304
#define OBX_ERROR_STD_RANGE 10305
#define OBX_ERROR_STD_OVERFLOW 10306
#define OBX_ERROR_STD_OTHER 10399

// --- Schema and file errors ------------------------------------------------

/// Data schema (aka data model) related error;
/// typically when a provided schema is inconsistent/incompatible to an existing one (on disk).
#define OBX_ERROR_SCHEMA 10501

/// DB file has errors, e.g. illegal values or structural inconsistencies were detected.
#define OBX_ERROR_FILE_CORRUPT 10502

/// DB file has errors related to pages, e.g. bad page refs outside of the file.
#define OBX_ERROR_FILE_PAGES_CORRUPT 10503

/// A requested schema object (e.g. an entity or a property) was not found in the schema
#define OBX_ERROR_SCHEMA_OBJECT_NOT_FOUND 10504

// --- Tree errors: 106xx ----------------------------------------------------

/// Could not setup tree because the data model (schema) of the database does not contain all required tree
/// types/properties.
#define OBX_ERROR_TREE_MODEL_INVALID 10601

/// Accessing a tree value failed because the stored value was put using a different value type.
#define OBX_ERROR_TREE_VALUE_TYPE_MISMATCH 10602

/// At a given path, a leaf was expected but a branch node was found; e.g. values can only be stored in leaves.
#define OBX_ERROR_TREE_PATH_NON_UNIQUE 10603

/// When trying to put or get using an illegal tree path. Examples for illegal paths with delimiter '/':
/// "", "/", "a", "/a", "a/", "a//b", " ", " / ", "a/ /b", " a/b", "a/b ".
#define OBX_ERROR_TREE_PATH_ILLEGAL 10604

/// An error related to tree functionality was detected;
/// please check the error message for details (see obx_last_error_message()).
#define OBX_ERROR_TREE_OTHER 10699

//----------------------------------------------
// Error info; obx_last_error_*
//----------------------------------------------

/// Return the error status on the current thread and clear the error state.
/// The buffer returned in out_message is valid only until the next call into ObjectBox.
/// @param out_error receives the error code; optional: may be NULL
/// @param out_message receives the pointer to the error messages; optional: may be NULL
/// @returns true if an error was pending
OBX_C_API bool obx_last_error_pop(obx_err* out_error, const char** out_message);

/// The last error raised by an ObjectBox API call on the current thread, or OBX_SUCCESS if no error occurred yet.
/// Note that API calls do not clear this error code (also true for this method).
/// Thus, if you receive an error from this, it's usually a good idea to call obx_last_error_clear() to clear the error
/// state (or use obx_last_error_pop()) for future API calls.
OBX_C_API obx_err obx_last_error_code(void);

/// The error message string attached to the error returned by obx_last_error_code().
/// Like obx_last_error_code(), this is bound to the current thread, and this call does not clear the error state.
/// The buffer returned is valid only until the next call into ObjectBox.
OBX_C_API const char* obx_last_error_message(void);

/// The underlying error for the error returned by obx_last_error_code(). Where obx_last_error_code() may be a generic
/// error like OBX_ERROR_STORAGE_GENERAL, this will give a further underlying and possibly platform-specific error code.
OBX_C_API obx_err obx_last_error_secondary(void);

/// Clear the error state on the current thread; e.g. obx_last_error_code() will now return OBX_SUCCESS.
/// Note that clearing the error state does not happen automatically;
/// API calls set the error state when they produce an error, but do not clear it on success.
/// See also: obx_last_error_pop() to retrieve the error state and clear it.
OBX_C_API void obx_last_error_clear(void);

/// Set the last error code and test - reserved for internal use from generated code.
OBX_C_API bool obx_last_error_set(obx_err code, obx_err secondary, const char* message);

//----------------------------------------------
// Model
//----------------------------------------------

typedef enum {
    OBXPropertyType_Bool = 1,    ///< 1 byte
    OBXPropertyType_Byte = 2,    ///< 1 byte
    OBXPropertyType_Short = 3,   ///< 2 bytes
    OBXPropertyType_Char = 4,    ///< 1 byte
    OBXPropertyType_Int = 5,     ///< 4 bytes
    OBXPropertyType_Long = 6,    ///< 8 bytes
    OBXPropertyType_Float = 7,   ///< 4 bytes
    OBXPropertyType_Double = 8,  ///< 8 bytes
    OBXPropertyType_String = 9,
    OBXPropertyType_Date = 10,  ///< Unix timestamp (milliseconds since 1970) in 8 bytes
    OBXPropertyType_Relation = 11,
    OBXPropertyType_DateNano = 12,  ///< Unix timestamp (nanoseconds since 1970) in 8 bytes
    OBXPropertyType_Flex = 13,  ///< Flexible" type, which may contain scalars (integers, floating points), strings or
                                ///< containers (lists and maps). Note: a flex map must use string keys.
    OBXPropertyType_ByteVector = 23,
    OBXPropertyType_StringVector = 30,
} OBXPropertyType;

/// Bit-flags defining the behavior of entities.
/// Note: Numbers indicate the bit position
typedef enum {
    /// Enable "data synchronization" for this entity type: objects will be synced with other stores over the network.
    /// It's possible to have local-only (non-synced) types and synced types in the same store (schema/data model).
    OBXEntityFlags_SYNC_ENABLED = 2,

    /// Makes object IDs for a synced types (SYNC_ENABLED is set) global.
    /// By default (not using this flag), the 64 bit object IDs have a local scope and are not unique globally.
    /// This flag tells ObjectBox to treat object IDs globally and thus no ID mapping (local <-> global) is performed.
    /// Often this is used with assignable IDs (ID_SELF_ASSIGNABLE property flag is set) and some special ID scheme.
    /// Note: typically you won't do this with automatically assigned IDs, set by the local ObjectBox store.
    ///       Two devices would likely overwrite each other's object during sync as object IDs are prone to collide.
    ///       It might be OK if you can somehow ensure that only a single device will create new IDs.
    OBXEntityFlags_SHARED_GLOBAL_IDS = 4,
} OBXEntityFlags;

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

    /// Used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
    OBXPropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512,

    /// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
    OBXPropertyFlags_VIRTUAL = 1024,

    /// Index uses a 32 bit hash instead of the value
    /// 32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions
    OBXPropertyFlags_INDEX_HASH = 2048,

    /// Index uses a 64 bit hash instead of the value
    /// recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash
    OBXPropertyFlags_INDEX_HASH64 = 4096,

    /// The actual type of the variable is unsigned (used in combination with numeric OBXPropertyType_*).
    /// While our default are signed ints, queries & indexes need do know signing info.
    /// Note: Don't combine with ID (IDs are always unsigned internally).
    OBXPropertyFlags_UNSIGNED = 8192,

    /// By defining an ID companion property, a special ID encoding scheme is activated involving this property.
    ///
    /// For Time Series IDs, a companion property of type Date or DateNano represents the exact timestamp.
    OBXPropertyFlags_ID_COMPANION = 16384,

    /// Unique on-conflict strategy: the object being put replaces any existing conflicting object (deletes it).
    OBXPropertyFlags_UNIQUE_ON_CONFLICT_REPLACE = 32768,

    /// If a date property has this flag (max. one per entity type), the date value specifies the time by which
    /// the object expires, at which point it MAY be removed (deleted), which can be triggered by an API call.
    OBXPropertyFlags_EXPIRATION_TIME = 65536,
} OBXPropertyFlags;

struct OBX_model;  // doxygen (only) picks up the typedef struct below

/// Model represents a database schema and must be provided when opening the store.
/// Model initialization is usually done by language bindings, which automatically build the model based on parsed
/// source code (for examples, see ObjectBox Go or Swift, which also use this C API).
///
/// For manual creation, these are the basic steps:
/// - define entity types using obx_model_entity() and obx_model_property()
/// - Pass the last ever used IDs with obx_model_last_entity_id(), obx_model_last_index_id(),
///   obx_model_last_relation_id()
typedef struct OBX_model OBX_model;

/// Create an (empty) data meta model which is to be consumed by obx_opt_model().
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details.
///               Note that obx_model_* functions handle OBX_model NULL pointers (will indicate an error but not crash).
OBX_C_API OBX_model* obx_model(void);

/// Only call when not calling obx_store_open() (which will free it internally)
/// @param model NULL-able; returns OBX_SUCCESS if model is NULL
OBX_C_API obx_err obx_model_free(OBX_model* model);

/// To minimise the amount of error handling code required when building a model, the first error is stored and can be
/// obtained here. All the obx_model_XXX functions are null operations after the first model error has occurred.
/// @param model NULL-able; returns OBX_ERROR_ILLEGAL_ARGUMENT if model is NULL
OBX_C_API obx_err obx_model_error_code(OBX_model* model);

/// To minimise the amount of error handling code required when building a model, the first error is stored and can be
/// obtained here. All the obx_model_XXX functions are null operations after the first model error has occurred.
/// @param model NULL-able; returns NULL if model is NULL
OBX_C_API const char* obx_model_error_message(OBX_model* model);

/// Starts the definition of a new entity type for the meta data model.
/// After this, call obx_model_property() to add properties to the entity type.
/// @param name A human readable name for the entity. Must be unique within the model
/// @param entity_id Must be unique within this version of the model
/// @param entity_uid Used to identify entities between versions of the model. Must be globally unique.
OBX_C_API obx_err obx_model_entity(OBX_model* model, const char* name, obx_schema_id entity_id, obx_uid entity_uid);

/// Refine the definition of the entity declared by the most recent obx_model_entity() call, specifying flags.
OBX_C_API obx_err obx_model_entity_flags(OBX_model* model, OBXEntityFlags flags);

/// Starts the definition of a new property for the entity type of the last obx_model_entity() call.
/// @param name A human readable name for the property. Must be unique within the entity
/// @param type The type of property required
/// @param property_id Must be unique within the entity
/// @param property_uid Used to identify properties between versions of the entity. Must be global unique.
OBX_C_API obx_err obx_model_property(OBX_model* model, const char* name, OBXPropertyType type,
                                     obx_schema_id property_id, obx_uid property_uid);

/// Refine the definition of the property declared by the most recent obx_model_property() call, specifying flags.
OBX_C_API obx_err obx_model_property_flags(OBX_model* model, OBXPropertyFlags flags);

/// Refine the definition of the property declared by the most recent obx_model_property() call, declaring it a
/// relation.
/// @param target_entity The name of the entity linked to by the relation
/// @param index_id Must be unique within this version of the model
/// @param index_uid Used to identify relations between versions of the model. Must be globally unique.
OBX_C_API obx_err obx_model_property_relation(OBX_model* model, const char* target_entity, obx_schema_id index_id,
                                              obx_uid index_uid);

/// Refine the definition of the property declared by the most recent obx_model_property() call, adding an index.
/// @param index_id Must be unique within this version of the model
/// @param index_uid Used to identify relations between versions of the model. Must be globally unique.
OBX_C_API obx_err obx_model_property_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

/// Add a standalone relation between the active entity and the target entity to the model
/// @param relation_id Must be unique within this version of the model
/// @param relation_uid Used to identify relations between versions of the model. Must be globally unique.
/// @param target_id The id of the target entity of the relation
/// @param target_uid The uid of the target entity of the relation
OBX_C_API obx_err obx_model_relation(OBX_model* model, obx_schema_id relation_id, obx_uid relation_uid,
                                     obx_schema_id target_id, obx_uid target_uid);

/// Set the highest ever known entity id in the model. Should always be equal to or higher than the
/// last entity id of the previous version of the model
OBX_C_API void obx_model_last_entity_id(OBX_model*, obx_schema_id entity_id, obx_uid entity_uid);

/// Set the highest ever known index id in the model. Should always be equal to or higher than the
/// last index id of the previous version of the model
OBX_C_API void obx_model_last_index_id(OBX_model* model, obx_schema_id index_id, obx_uid index_uid);

/// Set the highest every known relation id in the model. Should always be equal to or higher than the
/// last relation id of the previous version of the model.
OBX_C_API void obx_model_last_relation_id(OBX_model* model, obx_schema_id relation_id, obx_uid relation_uid);

/// Set the highest ever known property id in the entity. Should always be equal to or higher than the
/// last property id of the previous version of the entity.
OBX_C_API obx_err obx_model_entity_last_property_id(OBX_model* model, obx_schema_id property_id, obx_uid property_uid);

//----------------------------------------------
// Store
//----------------------------------------------

struct OBX_store;  // doxygen (only) picks up the typedef struct below

/// \brief A ObjectBox store represents a database storing data in a given directory on a local file system.
///
/// Once opened using obx_store_open(), it's an entry point to data access APIs such as box, query, cursor, transaction.
/// After your work is done, you must close obx_store_close() to safely release all the handles and avoid data loss.
///
/// It's possible open multiple stores in different directories.
typedef struct OBX_store OBX_store;

struct OBX_store_options;  // doxygen (only) picks up the typedef struct below

/// Store options customize the behavior of ObjectBox before opening a store. Options can't be changed once the store is
/// open but of course you can close the store and open it again with the changed options.
/// Some of the notable options are obx_opt_directory() and obx_opt_max_db_size_in_kb().
typedef struct OBX_store_options OBX_store_options;

/// Debug flags typically enable additional "debug logging" that can be helpful to better understand what is going on
/// internally. These are intended for the development process only; typically one does not enable them for releases.
typedef enum {
    OBXDebugFlags_LOG_TRANSACTIONS_READ = 1,
    OBXDebugFlags_LOG_TRANSACTIONS_WRITE = 2,
    OBXDebugFlags_LOG_QUERIES = 4,
    OBXDebugFlags_LOG_QUERY_PARAMETERS = 8,
    OBXDebugFlags_LOG_ASYNC_QUEUE = 16,
    OBXDebugFlags_LOG_CACHE_HITS = 32,
    OBXDebugFlags_LOG_CACHE_ALL = 64,
    OBXDebugFlags_LOG_TREE = 128,
    /// For a limited number of error conditions, this will try to print stack traces.
    /// Note: this is Linux-only, experimental, and has several limitations:
    /// The usefulness of these stack traces depends on several factors and might not be helpful at all.
    OBXDebugFlags_LOG_EXCEPTION_STACK_TRACE = 256,
    /// Run a quick self-test to verify basic threading; somewhat paranoia to check the platform and the library setup.
    OBXDebugFlags_RUN_THREADING_SELF_TEST = 512,
} OBXDebugFlags;

/// Defines a padding mode for putting data bytes.
/// Depending on how that data is created, this mode may optimize data handling by avoiding copying memory.
/// Internal background: data buffers used by put operations are required to have a size divisible by 4 for an
///                      efficient data layout.
typedef enum {
    /// Adds a padding when needed (may require a memory copy): this is the safe option and also the default.
    /// The extra memory copy may impact performance, however this is usually not noticeable.
    OBXPutPaddingMode_PaddingAutomatic = 1,

    /// Indicates that data buffers are safe to be extended for padding (adding up to 3 bytes to size is OK).
    /// Typically, it depends on the used FlatBuffers builder; e.g. the official C++ seems to ensure it, but
    /// flatcc (3rd party implementation for plain C) may not.
    OBXPutPaddingMode_PaddingAllowedByBuffer = 2,

    /// The caller ensures that all data bytes are already padded.
    /// ObjectBox will verify the buffer size and returns an error if it's not divisible by 4.
    OBXPutPaddingMode_PaddingByCaller = 3,
} OBXPutPaddingMode;

/// Bytes struct is an input/output wrapper typically used for a single object data (represented as FlatBuffers).
typedef struct OBX_bytes {
    const void* data;
    size_t size;
} OBX_bytes;

/// Bytes array struct is an input/output wrapper for multiple FlatBuffers object data representation.
typedef struct OBX_bytes_array {
    OBX_bytes* bytes;
    size_t count;
} OBX_bytes_array;

struct OBX_bytes_lazy;  // doxygen (only) picks up the typedef struct below

/// Bytes, which must be resolved "lazily" via obx_bytes_lazy_get() and released via obx_bytes_lazy_free().
/// Unlike OBX_bytes, this may represent allocated resources and/or bytes that are only produced on demand.
typedef struct OBX_bytes_lazy OBX_bytes_lazy;

/// Get the actual bytes from the given OBX_bytes_lazy.
/// Because of the potential lazy creation of bytes, the given bytes are not const as it may be resolved internally.
/// For the same reason, this function is not thread-safe, at least for the first call on a OBX_bytes_lazy instance.
/// @param out_bytes The pointer to a data pointer may be null, e.g. if you are only interested in the size.
OBX_C_API obx_err obx_bytes_lazy_get(OBX_bytes_lazy* bytes, const void** out_bytes, size_t* out_size);

/// Releases any resources associated with the given OBX_bytes_lazy.
OBX_C_API void obx_bytes_lazy_free(OBX_bytes_lazy* bytes);

/// ID array struct is an input/output wrapper for an array of object IDs.
typedef struct OBX_id_array {
    obx_id* ids;
    size_t count;
} OBX_id_array;

/// String array struct is an input/output wrapper for an array of character strings.
typedef struct OBX_string_array {
    const char** items;
    size_t count;
} OBX_string_array;

/// Int64 array struct is an input/output wrapper for an array of int64 numbers.
typedef struct OBX_int64_array {
    const int64_t* items;
    size_t count;
} OBX_int64_array;

/// Int32 array struct is an input/output wrapper for an array of int32 numbers.
typedef struct OBX_int32_array {
    const int32_t* items;
    size_t count;
} OBX_int32_array;

/// Int16 array struct is an input/output wrapper for an array of int16 numbers.
typedef struct OBX_int16_array {
    const int16_t* items;
    size_t count;
} OBX_int16_array;

/// Int8 array struct is an input/output wrapper for an array of int8 numbers.
typedef struct OBX_int8_array {
    const int8_t* items;
    size_t count;
} OBX_int8_array;

/// Double array struct is an input/output wrapper for an array of double precision floating point numbers.
typedef struct OBX_double_array {
    const double* items;
    size_t count;
} OBX_double_array;

/// Float array struct is an input/output wrapper for an array of single precision floating point numbers.
typedef struct OBX_float_array {
    const float* items;
    size_t count;
} OBX_float_array;

//----------------------------------------------
// Store Options
//----------------------------------------------

/// Create a default set of store options.
/// @returns NULL on failure, a default set of options on success
OBX_C_API OBX_store_options* obx_opt();

/// Set the store directory on the options. The default is "objectbox".
OBX_C_API obx_err obx_opt_directory(OBX_store_options* opt, const char* dir);

/// Set the maximum db size on the options. The default is 1Gb.
OBX_C_API void obx_opt_max_db_size_in_kb(OBX_store_options* opt, uint64_t size_in_kb);

/// Data size tracking is more involved than DB size tracking, e.g. it stores an internal counter.
/// Thus only use it if a stricter, more accurate limit is required (it's off by default).
/// It tracks the size of actual data bytes of objects (system and metadata is not considered).
/// On the upside, reaching the data limit still allows data to be removed (assuming DB limit is not reached).
/// Max data and DB sizes can be combined; data size must be below the DB size.
OBX_C_API void obx_opt_max_data_size_in_kb(OBX_store_options* opt, uint64_t size_in_kb);

/// Set the file mode on the options. The default is 0644 (unix-style)
OBX_C_API void obx_opt_file_mode(OBX_store_options* opt, unsigned int file_mode);

/// Set the maximum number of readers (related to read transactions) on the given options.
/// "Readers" are an finite resource for which we need to define a maximum number upfront.
/// The default value is enough for most apps and usually you can ignore it completely.
/// However, if you get the OBX_ERROR_MAX_READERS_EXCEEDED error, you should verify your threading.
/// For each thread, ObjectBox uses multiple readers.
/// Their number (per thread) depends on number of types, relations, and usage patterns.
/// Thus, if you are working with many threads (e.g. in a server-like scenario), it can make sense to increase the
/// maximum number of readers.
///
/// \note The internal default is currently 126. So when hitting this limit, try values around 200-500.
///
/// \attention Each thread that performed a read transaction and is still alive holds on to a reader slot.
///       These slots only get vacated when the thread ends. Thus be mindful with the number of active threads.
///       Alternatively, you can opt to try the experimental noReaderThreadLocals option flag.
OBX_C_API void obx_opt_max_readers(OBX_store_options* opt, unsigned int max_readers);

/// Disables the usage of thread locals for "readers" related to read transactions.
/// This can make sense if you are using a lot of threads that are kept alive.
/// \note This is still experimental, as it comes with subtle behavior changes at a low level and may affect
///       corner cases with e.g. transactions, which may not be fully tested at the moment.
OBX_C_API void obx_opt_no_reader_thread_locals(OBX_store_options* opt, bool flag);

/// Set the model on the options. The default is no model.
/// NOTE: the model is always freed by this function, including when an error occurs.
OBX_C_API obx_err obx_opt_model(OBX_store_options* opt, OBX_model* model);

/// Set the model on the options copying the given bytes. The default is no model.
OBX_C_API obx_err obx_opt_model_bytes(OBX_store_options* opt, const void* bytes, size_t size);

/// Like obx_opt_model_bytes BUT WITHOUT copying the given bytes.
/// Thus, you must keep the bytes available until after the store is created.
OBX_C_API obx_err obx_opt_model_bytes_direct(OBX_store_options* opt, const void* bytes, size_t size);

/// When the DB is opened initially, ObjectBox can do a consistency check on the given amount of pages.
/// Reliable file systems already guarantee consistency, so this is primarily meant to deal with unreliable
/// OSes, file systems, or hardware. Thus, usually a low number (e.g. 1-20) is sufficient and does not impact
/// startup performance significantly. To completely disable this you can pass 0, but we recommend a setting of
/// at least 1.
/// Note: ObjectBox builds upon ACID storage, which guarantees consistency given that the file system is working
/// correctly (in particular fsync).
/// @param page_limit limits the number of checked pages (currently defaults to 0, but will be increased in the future)
/// @param leaf_level enable for visiting leaf pages (defaults to false)
OBX_C_API void obx_opt_validate_on_open(OBX_store_options* opt, size_t page_limit, bool leaf_level);

/// Don't touch unless you know exactly what you are doing:
/// Advanced setting typically meant for language bindings (not end users). See OBXPutPaddingMode description.
OBX_C_API void obx_opt_put_padding_mode(OBX_store_options* opt, OBXPutPaddingMode mode);

/// Advanced setting meant only for special scenarios: setting to false causes opening the database in a limited,
/// schema-less mode. If you don't know what this means exactly: ignore this flag. Defaults to true.
OBX_C_API void obx_opt_read_schema(OBX_store_options* opt, bool value);

/// Advanced setting recommended to be used together with read-only mode to ensure no data is lost.
/// Ignores the latest data snapshot (committed transaction state) and uses the previous snapshot instead.
/// When used with care (e.g. backup the DB files first), this option may also recover data removed by the latest
/// transaction. Defaults to false.
OBX_C_API void obx_opt_use_previous_commit(OBX_store_options* opt, bool value);

/// Open store in read-only mode: no schema update, no write transactions. Defaults to false.
OBX_C_API void obx_opt_read_only(OBX_store_options* opt, bool value);

/// Configure debug flags (OBXDebugFlags); e.g. to influence logging. Defaults to NONE.
/// Combine multiple flags using bitwise OR.
OBX_C_API void obx_opt_debug_flags(OBX_store_options* opt, uint32_t flags);

/// Adds debug flags (OBXDebugFlags) to potentially existing ones.
/// Combine multiple flags using bitwise OR.
OBX_C_API void obx_opt_add_debug_flags(OBX_store_options* opt, uint32_t flags);

/// Maximum of async elements in the queue before new elements will be rejected.
/// Hitting this limit usually hints that async processing cannot keep up;
/// data is produced at a faster rate than it can be persisted in the background.
/// In that case, increasing this value is not the only alternative; other values might also optimize throughput.
/// For example, increasing maxInTxDurationMicros may help too.
OBX_C_API void obx_opt_async_max_queue_length(OBX_store_options* opt, size_t value);

/// Producers (AsyncTx submitter) is throttled when the queue size hits this
OBX_C_API void obx_opt_async_throttle_at_queue_length(OBX_store_options* opt, size_t value);

/// Sleeping time for throttled producers on each submission
OBX_C_API void obx_opt_async_throttle_micros(OBX_store_options* opt, uint32_t value);

/// Maximum duration spent in a transaction before AsyncQ enforces a commit.
/// This becomes relevant if the queue is constantly populated at a high rate.
OBX_C_API void obx_opt_async_max_in_tx_duration(OBX_store_options* opt, uint32_t micros);

/// Maximum operations performed in a transaction before AsyncQ enforces a commit.
/// This becomes relevant if the queue is constantly populated at a high rate.
OBX_C_API void obx_opt_async_max_in_tx_operations(OBX_store_options* opt, uint32_t value);

/// Before the AsyncQ is triggered by a new element in queue to starts a new run, it delays actually starting the
/// transaction by this value.
/// This gives a newly starting producer some time to produce more than one a single operation before AsyncQ starts.
/// Note: this value should typically be low to keep latency low and prevent accumulating too much operations.
OBX_C_API void obx_opt_async_pre_txn_delay(OBX_store_options* opt, uint32_t delay_micros);

/// Before the AsyncQ is triggered by a new element in queue to starts a new run, it delays actually starting the
/// transaction by this value.
/// This gives a newly starting producer some time to produce more than one a single operation before AsyncQ starts.
/// Note: this value should typically be low to keep latency low and prevent accumulating too much operations.
OBX_C_API void obx_opt_async_pre_txn_delay4(OBX_store_options* opt, uint32_t delay_micros, uint32_t delay2_micros,
                                            size_t min_queue_length_for_delay2);

/// Similar to preTxDelay but after a transaction was committed.
/// One of the purposes is to give other transactions some time to execute.
/// In combination with preTxDelay this can prolong non-TX batching time if only a few operations are around.
OBX_C_API void obx_opt_async_post_txn_delay(OBX_store_options* opt, uint32_t delay_micros);

/// Similar to preTxDelay but after a transaction was committed.
/// One of the purposes is to give other transactions some time to execute.
/// In combination with preTxDelay this can prolong non-TX batching time if only a few operations are around.
/// @param subtract_processing_time If set, delay_micros is interpreted from the start of TX processing.
///        In other words, the actual delay is delay_micros minus the TX processing time including the commit.
///        This can make timings more accurate (e.g. when fixed batching interval are given).
OBX_C_API void obx_opt_async_post_txn_delay5(OBX_store_options* opt, uint32_t delay_micros, uint32_t delay2_micros,
                                             size_t min_queue_length_for_delay2, bool subtract_processing_time);

/// Numbers of operations below this value are considered "minor refills"
OBX_C_API void obx_opt_async_minor_refill_threshold(OBX_store_options* opt, size_t queue_length);

/// If non-zero, this allows "minor refills" with small batches that came in (off by default).
OBX_C_API void obx_opt_async_minor_refill_max_count(OBX_store_options* opt, uint32_t value);

/// Default value: 10000, set to 0 to deactivate pooling
OBX_C_API void obx_opt_async_max_tx_pool_size(OBX_store_options* opt, size_t value);

/// Total cache size; default: ~ 0.5 MB
OBX_C_API void obx_opt_async_object_bytes_max_cache_size(OBX_store_options* opt, uint64_t value);

/// Maximal size for an object to be cached (only cache smaller ones)
OBX_C_API void obx_opt_async_object_bytes_max_size_to_cache(OBX_store_options* opt, uint64_t value);

/// Registers a log callback, which is called for a selection of log events.
/// Note: this does not replace the default logging, which is much more extensive (at least at this point).
OBX_C_API void obx_opt_log_callback(OBX_store_options* opt, obx_log_callback* callback, void* user_data);

/// Gets the option for "directory"; this is either the default, or, the value set by obx_opt_directory().
/// The returned value must not be modified and is only valid for the lifetime of the options or until the value is
/// changed.
/// @returns null if an error occurred, e.g. the given options were null.
OBX_C_API const char* obx_opt_get_directory(OBX_store_options* opt);

/// Gets the option for "max DB size"; this is either the default, or, the value set by obx_opt_max_db_size_in_kb().
/// @returns 0 if an error occurred, e.g. the given options were null.
OBX_C_API uint64_t obx_opt_get_max_db_size_in_kb(OBX_store_options* opt);

/// Gets the option for "max data size"; this is either the default, or, the value set by obx_opt_max_data_size_in_kb().
/// @returns 0 if an error occurred, e.g. the given options were null.
OBX_C_API uint64_t obx_opt_get_max_data_size_in_kb(OBX_store_options* opt);

/// Gets the option for "debug flags"; this is either the default, or, the value set by obx_opt_debug_flags().
/// @returns 0 if an error occurred, e.g. the given options were null.
OBX_C_API uint32_t obx_opt_get_debug_flags(OBX_store_options* opt);

/// Free the options.
/// Note: Only free *unused* options, obx_store_open() frees the options internally
OBX_C_API void obx_opt_free(OBX_store_options* opt);

//----------------------------------------------
// Store
//----------------------------------------------

/// Opens (creates) a "store", which represents an ObjectBox database instance in a given directory.
/// The store is an entry point to data access APIs such as box (obx_box_*), query (obx_qb_* and obx_query_*),
/// and transaction (obx_txn_*).
/// It's possible open multiple stores in different directories, e.g. at the same time.
/// See also obx_store_close() to close a previously opened store.
/// Note: the given options are always freed by this function, including when an error occurs.
/// @param opt required parameter holding the data model (obx_opt_model()) and optional options (see obx_opt_*())
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_store* obx_store_open(OBX_store_options* opt);

/// Check if an open store was found for the given path (i.e. opened before and not yet closed).
OBX_C_API bool obx_store_is_open(const char* path);

/// Attach to a previously opened store matching the path of the DB directory, which was used for opening the store.
/// The returned store is a new instance (e.g. different pointer value) and must also be closed via obx_store_close().
/// The actual underlying store is only closed when the last store OBX_store instance is closed.
/// @returns nullptr if no open store was found (i.e. not opened before or already closed)
/// @see obx_store_clone() for "attaching" to a available store instance.
OBX_C_API OBX_store* obx_store_attach(const char* path);

/// Attach to a previously opened store matching the given store ID.
/// The returned store is a new instance (e.g. different pointer value) and must also be closed via obx_store_close().
/// The actual underlying store is only closed when the last store OBX_store instance is closed.
/// @param store_id
/// @returns nullptr if no open store was found (i.e. not opened before or already closed)
/// @see obx_store_clone() for "attaching" to a available store instance.
OBX_C_API OBX_store* obx_store_attach_id(uint64_t store_id);

/// Combines the functionality of obx_store_attach() and obx_store_open() in a thread-safe way.
/// @param check_matching_options if true, some effort will be taken to check if the given options match an existing
///        store. Usually this should be set to true as options should not diverge when opening/attaching.
///        Note that additional checks will be added in the future.
/// @param out_attached (optional) if given a pointer to a flag that telling if the function attached to an existing
///        store (true) or a new store was created (false).
OBX_C_API OBX_store* obx_store_attach_or_open(OBX_store_options* opt, bool check_matching_options, bool* out_attached);

/// Store IDs can be used to attach to a store later.
/// The IDs are stable and unique during the lifetime of the process.
/// E.g. these IDs can be shared across threads efficiently and can serve a similar purpose as weak pointers do.
OBX_C_API uint64_t obx_store_id(OBX_store* store);

/// Clone a previously opened store; while a store instance is usable from multiple threads, situations may exist
/// in which cloning a store simplifies the overall lifecycle.
/// E.g. when a store is used for multiple threads and it may only be fully released once the last thread completes.
/// The returned store is a new instance (e.g. different pointer value) and must also be closed via obx_store_close().
/// The actual underlying store is only closed when the last store OBX_store instance is closed.
/// @returns nullptr if the store could not be cloned
/// @see obx_store_attach() for "cloning" using the store's path.
OBX_C_API OBX_store* obx_store_clone(OBX_store* store);

/// For stores created outside of this C API, e.g. via C++ or Java, this is how you can use it via C too.
/// Like this, it is OK to use the same store instance (same database) from multiple languages in parallel.
/// Note: the store's life time will still be managed outside of the C API;
/// thus ensure that store is not closed while calling any C function on it.
/// Once you are done with the C specific OBX_store, call obx_store_close() to free any C related resources.
/// This, however, will not close the "core store".
/// @param core_store A pointer to the core C++ ObjectStore, or the native JNI handle for a BoxStore.
OBX_C_API OBX_store* obx_store_wrap(void* core_store);

/// Look for an entity with the given name in the model and return its Entity ID.
OBX_C_API obx_schema_id obx_store_entity_id(OBX_store* store, const char* entity_name);

/// Return the property id from the property name or 0 if the name is not found
OBX_C_API obx_schema_id obx_store_entity_property_id(OBX_store* store, obx_schema_id entity_id,
                                                     const char* property_name);

/// Await for all (including future) async submissions to be completed (the async queue becomes empty).
/// @returns true if all submissions were completed or async processing was not started
/// @returns false if shutting down or an error occurred
OBX_C_API bool obx_store_await_async_completion(OBX_store* store);

/// Await async operations that have been submitted up to this point to be completed
/// (the async queue may still contain elements).
/// @returns true if all submissions were completed (or async processing was not started)
/// @returns false if shutting down or an error occurred
OBX_C_API bool obx_store_await_async_submitted(OBX_store* store);

/// Configure debug logging
OBX_C_API obx_err obx_store_debug_flags(OBX_store* store, OBXDebugFlags flags);

/// @returns true if the store was opened with a previous commit
/// @see obx_opt_use_previous_commit()
OBX_C_API bool obx_store_opened_with_previous_commit(OBX_store* store);

/// Prepares the store to close by setting its internal state to "closing".
/// Functions relying on the store will result in OBX_ERROR_SHUTTING_DOWN status once closing is initiated.
/// Unlike obx_store_close(), this method will return immediately and does not free resources just yet.
/// This is typically used in a multi-threaded context to allow an orderly shutdown in stages which go through a
/// "not accepting new requests" state.
OBX_C_API obx_err obx_store_prepare_to_close(OBX_store* store);

/// Closes a previously opened store and thus freeing all resources associated with the store.
/// \note This waits for write transactions to finish before returning from this call.
/// @param store may be NULL
OBX_C_API obx_err obx_store_close(OBX_store* store);

//----------------------------------------------
// Transaction
//----------------------------------------------

struct OBX_txn;  // doxygen (only) picks up the typedef struct below

/// Transaction provides the mean to use explicit database transactions, grouping several operations into a single unit
/// of work that either executes completely or not at all. If you are looking for a more detailed introduction to
/// transactions in general, please consult other resources, e.g., https://en.wikipedia.org/wiki/Database_transaction
///
/// You may not notice it, but almost all interactions with ObjectBox involve transactions. For example, if you call
/// obx_box_put() a write transaction is used. Also if you call obx_box_count(), a read transaction is used. All of this
/// is done under the hood and transparent to you.
/// However, there are situations where an explicit read transaction is necessary, e.g. obx_box_get(). Also, it’s
/// usually worth learning transaction basics to make your app more consistent and efficient, especially for writes.
typedef struct OBX_txn OBX_txn;

/// Create a write transaction (read and write).
/// Transaction creation can be nested (recursive), however only the outermost transaction is relevant on the DB level.
/// Once you are done tithe the transaction, you must call obx_txn_success() or obx_txn_close().
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if called when inside a read transaction.
OBX_C_API OBX_txn* obx_txn_write(OBX_store* store);

/// Create a read transaction (read only).
/// Transaction creation can be nested (recursive), however only the outermost transaction is relevant on the DB level.
/// Once you are done tithe the transaction, you must call obx_txn_close().
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_txn* obx_txn_read(OBX_store* store);

/// "Finish" this write transaction successfully and close it, performing a commit if this is the top level
/// transaction and all inner transactions (if any) were also successful (obx_txn_success() was called on them).
/// Because this also closes the given transaction along with all resources, the given OBX_txn pointer becomes invalid.
/// Thus, the OBX_txn pointer must not be used afterwards; do not even call obx_txn_close().
/// @return OBX_ERROR_ILLEGAL_STATE if the given transaction is not a write transaction.
OBX_C_API obx_err obx_txn_success(OBX_txn* txn);

/// Close the transaction (read or write) and free all of its resources; the given OBX_txn pointer must not be used
/// afterwards. If the given transaction is a write transaction, it is aborted (not committed) right away (without
/// waiting for the top level transaction).
/// If an error is returned (e.g., a commit failed because DB is full), you can assume that the transaction was closed.
/// @param txn may be NULL
OBX_C_API obx_err obx_txn_close(OBX_txn* txn);

/// Note: At the moment, you typically want to use only obx_txn_close() instead.
/// Abort the underlying transaction immediately and thus frees DB resources.
/// Only obx_txn_close() is allowed to be called on the transaction after calling this.
OBX_C_API obx_err obx_txn_abort(OBX_txn* txn);  // Internal note: will make more sense once we introduce obx_txn_reset

OBX_C_API obx_err obx_txn_data_size(OBX_txn* txn, uint64_t* out_committed_size, uint64_t* out_size_change);

//------------------------------------------------------------------
// Cursor (lower level API, check also the more convenient Box API)
//------------------------------------------------------------------

struct OBX_cursor;  // doxygen (only) picks up the typedef struct below

/// Cursor provides fine-grained (lower level API) access to the stored objects. Check also the more convenient Box API.
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
OBX_C_API OBX_cursor* obx_cursor(OBX_txn* txn, obx_schema_id entity_id);

/// @param cursor may be NULL
OBX_C_API obx_err obx_cursor_close(OBX_cursor* cursor);

/// Call this when putting an object to generate/prepare an ID for it.
/// @param id_or_zero The ID of the entity. If you pass 0, this will generate a new one.
/// @seealso obx_box_id_for_put()
OBX_C_API obx_id obx_cursor_id_for_put(OBX_cursor* cursor, obx_id id_or_zero);

/// Puts the given object data using the given ID.
/// A "put" in ObjectBox follows "insert or update" semantics;
/// New objects (no pre-existing object for given ID) are inserted while existing objects are replaced/updated.
/// @param id non-zero
OBX_C_API obx_err obx_cursor_put(OBX_cursor* cursor, obx_id id, const void* data, size_t size);

/// Like put obx_cursor_put(), but takes an additional parameter (4th parameter) for choosing a put mode.
/// @param id non-zero
/// @param mode Changes the put semantics to the given mode, e.g. OBXPutMode_INSERT or OBXPutMode_UPDATE.
/// @returns OBX_SUCCESS if the put operation was successful
/// @returns OBX_ERROR_ID_ALREADY_EXISTS OBXPutMode_INSERT was used, but an existing object was found using the given ID
/// @returns OBX_ERROR_ID_NOT_FOUND OBXPutMode_UPDATE was used, but no object was found for the given ID
OBX_C_API obx_err obx_cursor_put4(OBX_cursor* cursor, obx_id id, const void* data, size_t size, OBXPutMode mode);

/// An optimized version of obx_cursor_put() if you can ensure that the given ID is not used yet.
/// Typically used right after getting a new ID via obx_cursor_id_for_put().
/// WARNING: using this incorrectly (an object with the given ID already exists) may result in inconsistent data
/// (e.g. indexes do not get updated).
/// @param id non-zero
OBX_C_API obx_err obx_cursor_put_new(OBX_cursor* cursor, obx_id id, const void* data, size_t size);

/// Convenience for obx_cursor_put4() with OBXPutMode_INSERT.
/// @param id non-zero
/// @returns OBX_ERROR_ID_ALREADY_EXISTS if an insert fails because of a colliding ID
OBX_C_API obx_err obx_cursor_insert(OBX_cursor* cursor, obx_id id, const void* data, size_t size);

/// Convenience for obx_cursor_put4() with OBXPutMode_UPDATE.
/// @param id non-zero
/// @returns OBX_ERROR_ID_NOT_FOUND  if an update fails because the given ID does not represent any object
OBX_C_API obx_err obx_cursor_update(OBX_cursor* cursor, obx_id id, const void* data, size_t size);

/// FB ID slot must be present; new entities must prepare the slot using the special value OBX_ID_NEW.
/// Alternatively, you may also pass 0 to indicate a new entity if you are aware that FlatBuffers builders typically
/// skip zero values by default. Thus, you have to "force" writing the zero in FlatBuffers.
/// @param data object data, non-const because the ID slot will be written (mutated) for new entites (see above)
/// @returns id if the object could be put, or 0 in case of an error
OBX_C_API obx_id obx_cursor_put_object(OBX_cursor* cursor, void* data, size_t size);

/// @overload obx_id obx_cursor_put_object(OBX_cursor* cursor, void* data, size_t size)
OBX_C_API obx_id obx_cursor_put_object4(OBX_cursor* cursor, void* data, size_t size, OBXPutMode mode);

OBX_C_API obx_err obx_cursor_get(OBX_cursor* cursor, obx_id id, const void** data, size_t* size);

/// Get all objects as bytes.
/// For larger quantities, it's recommended to iterate using obx_cursor_first and obx_cursor_next.
/// However, if the calling overhead is high (e.g., for language bindings), this method helps.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_bytes_array* obx_cursor_get_all(OBX_cursor* cursor);

OBX_C_API obx_err obx_cursor_first(OBX_cursor* cursor, const void** data, size_t* size);

OBX_C_API obx_err obx_cursor_next(OBX_cursor* cursor, const void** data, size_t* size);

OBX_C_API obx_err obx_cursor_seek(OBX_cursor* cursor, obx_id id);

OBX_C_API obx_err obx_cursor_current(OBX_cursor* cursor, const void** data, size_t* size);

OBX_C_API obx_err obx_cursor_remove(OBX_cursor* cursor, obx_id id);

OBX_C_API obx_err obx_cursor_remove_all(OBX_cursor* cursor);

/// Count the number of available objects
OBX_C_API obx_err obx_cursor_count(OBX_cursor* cursor, uint64_t* count);

/// Count the number of available objects up to the specified maximum
OBX_C_API obx_err obx_cursor_count_max(OBX_cursor* cursor, uint64_t max_count, uint64_t* out_count);

/// Return true if there is no object available (false if at least one object is available)
OBX_C_API obx_err obx_cursor_is_empty(OBX_cursor* cursor, bool* out_is_empty);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_bytes_array* obx_cursor_backlinks(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                                obx_id id);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_id_array* obx_cursor_backlink_ids(OBX_cursor* cursor, obx_schema_id entity_id, obx_schema_id property_id,
                                                obx_id id);

OBX_C_API obx_err obx_cursor_rel_put(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id, obx_id target_id);
OBX_C_API obx_err obx_cursor_rel_remove(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id,
                                        obx_id target_id);

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_id_array* obx_cursor_rel_ids(OBX_cursor* cursor, obx_schema_id relation_id, obx_id source_id);

//----------------------------------------------
// Time series
//----------------------------------------------

/// Time series: get the limits (min/max time values) over all objects
/// @param out_min_id pointer to receive an output (may be NULL)
/// @param out_min_value pointer to receive an output (may be NULL)
/// @param out_max_id pointer to receive an output (may be NULL)
/// @param out_max_value pointer to receive an output (may be NULL)
/// @returns OBX_NOT_FOUND if no objects are stored
OBX_C_API obx_err obx_cursor_ts_min_max(OBX_cursor* cursor, obx_id* out_min_id, int64_t* out_min_value,
                                        obx_id* out_max_id, int64_t* out_max_value);

/// Time series: get the limits (min/max time values) over objects within the given time range
/// @param out_min_id pointer to receive an output (may be NULL)
/// @param out_min_value pointer to receive an output (may be NULL)
/// @param out_max_id pointer to receive an output (may be NULL)
/// @param out_max_value pointer to receive an output (may be NULL)
/// @returns OBX_NOT_FOUND if no objects are stored in the given range
OBX_C_API obx_err obx_cursor_ts_min_max_range(OBX_cursor* cursor, int64_t range_begin, int64_t range_end,
                                              obx_id* out_min_id, int64_t* out_min_value, obx_id* out_max_id,
                                              int64_t* out_max_value);

//----------------------------------------------
// Box
//----------------------------------------------

struct OBX_box;  // doxygen (only) picks up the typedef struct below

/// \brief A Box offers database operations for objects of a specific type.
///
/// Given an OBX_store, obx_box() gives you Box instances to interact with object data (e.g. get and put operations).
/// A Box instance is associated with a specific object type (data class) and gives you a high level API to interact
/// with data objects of that type.
///
/// Box operations automatically start an implicit transaction when accessing the database.
/// And because transactions offered by this C API are always reentrant, you can set your own transaction boundary
/// using obx_txn_read() or obx_txn_write().
/// Using these explicit transactions is very much encouraged for calling multiple write operations that
/// logically belong together for better consistency (ACID) and performance.
///
/// Box instances are thread-safe and cached internally (see obx_box()).
typedef struct OBX_box OBX_box;

/// Get access to the box for the given entity. A box may be used across threads.
/// Boxes are shared instances and managed by the store so there's no need to close/free them manually.
OBX_C_API OBX_box* obx_box(OBX_store* store, obx_schema_id entity_id);

/// Get access to the store this box belongs to - utility for when you only have access to the `box` variable but need
/// some store method, such as starting a transaction.
/// This doesn't produce a new instance of OBX_store, just gives you back the same pointer you've created this box with.
/// In other words, don't close the returned store separately.
OBX_C_API OBX_store* obx_box_store(OBX_box* box);

/// Check whether a given object exists in the box.
OBX_C_API obx_err obx_box_contains(OBX_box* box, obx_id id, bool* out_contains);

/// Check whether this box contains objects with all of the IDs given.
/// @param out_contains is set to true if all of the IDs are present, otherwise false
OBX_C_API obx_err obx_box_contains_many(OBX_box* box, const OBX_id_array* ids, bool* out_contains);

/// Fetch a single object from the box; must be called inside a (reentrant) transaction.
/// The exposed data comes directly from the OS to allow zero-copy access, which limits the data lifetime:
/// \attention The exposed data is only valid as long as the (top) transaction is still active and no write
/// \attention operation (e.g. put/remove) was executed. Accessing data after this is undefined behavior.
/// @returns OBX_ERROR_ILLEGAL_STATE if not inside of an active transaction (see obx_txn_read() and obx_txn_write())
OBX_C_API obx_err obx_box_get(OBX_box* box, obx_id id, const void** data, size_t* size);

/// Fetch multiple objects for the given IDs from the box; must be called inside a (reentrant) transaction.
/// \attention See obx_box_get() for important notes on the limited lifetime of the exposed data.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if not inside of an active transaction
///               (see obx_txn_read() and obx_txn_write())
OBX_C_API OBX_bytes_array* obx_box_get_many(OBX_box* box, const OBX_id_array* ids);

/// Fetch all objects from the box; must be called inside a (reentrant) transaction.
/// NOTE: don't call this in 32 bit mode! Use obx_box_visit_all() instead.
/// \attention See obx_box_get() for important notes on the limited lifetime of the exposed data.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details; e.g. code
///               OBX_ERROR_ILLEGAL_STATE will be set if not inside of an active transaction
///               (see obx_txn_read() and obx_txn_write())
OBX_C_API OBX_bytes_array* obx_box_get_all(OBX_box* box);

/// Read given objects from the database in a single transaction.
/// Call the visitor() on each object, passing user_data, object data & size as arguments.
/// The given visitor must return true to keep receiving results, false to cancel.
/// If an object is not found, the visitor() is still called, passing NULL as data and a 0 as size.
OBX_C_API obx_err obx_box_visit_many(OBX_box* box, const OBX_id_array* ids, obx_data_visitor* visitor, void* user_data);

/// Read all objects in a single transaction.
/// Calls the visitor() on each object, passing visitor_arg, object data & size as arguments.
/// The given visitor must return true to keep receiving results, false to cancel.
OBX_C_API obx_err obx_box_visit_all(OBX_box* box, obx_data_visitor* visitor, void* user_data);

/// Prepares an ID for insertion: pass in 0 (zero) to reserve a new ID or an existing ID to check/prepare it.
/// @param id_or_zero The ID of the entity. If you pass 0, this will generate a new one.
/// @seealso obx_cursor_id_for_put()
OBX_C_API obx_id obx_box_id_for_put(OBX_box* box, obx_id id_or_zero);

/// Reserve the given number of (new) IDs for insertion; a bulk version of obx_box_id_for_put().
/// @param count number of IDs to reserve, max 10000
/// @param out_first_id the first ID of the sequence as
/// @returns an error in case the required number of IDs could not be reserved.
OBX_C_API obx_err obx_box_ids_for_put(OBX_box* box, uint64_t count, obx_id* out_first_id);

/// Put the given object using the given ID synchronously; note that the ID also must match the one present in data.
/// @param id An ID usually reserved via obx_box_id_for_put().
/// @see obx_box_put5() to additionally provide a put mode
/// @see obx_box_put_object() for a variant not requiring reserving IDs
OBX_C_API obx_err obx_box_put(OBX_box* box, obx_id id, const void* data, size_t size);

/// Convenience for obx_box_put5() with OBXPutMode_INSERT.
/// @param id non-zero
/// @returns OBX_ERROR_ID_ALREADY_EXISTS if an insert fails because of a colliding ID
OBX_C_API obx_err obx_box_insert(OBX_box* box, obx_id id, const void* data, size_t size);

/// Convenience for obx_cursor_put4() with OBXPutMode_UPDATE.
/// @param id non-zero
/// @returns OBX_ERROR_ID_NOT_FOUND  if an update fails because the given ID does not represent any object
OBX_C_API obx_err obx_box_update(OBX_box* box, obx_id id, const void* data, size_t size);

/// Put the given object using the given ID synchronously; note that the ID also must match the one present in data.
/// @param id An ID usually reserved via obx_box_id_for_put().
/// @see obx_box_put() for standard put mode
/// @see obx_box_put_object() for a variant not requiring reserving IDs
OBX_C_API obx_err obx_box_put5(OBX_box* box, obx_id id, const void* data, size_t size, OBXPutMode mode);

/// FB ID slot must be present in the given data; new entities must have an ID value of zero or OBX_ID_NEW.
/// @param data writable data buffer, which may be updated for the ID
/// @returns id if the object could be put, or 0 in case of an error
OBX_C_API obx_id obx_box_put_object(OBX_box* box, void* data, size_t size);

/// FB ID slot must be present in the given data; new entities must have an ID value of zero or OBX_ID_NEW
/// @param data writable data buffer, which may be updated for the ID
/// @returns id if the object, or 0 in case of an error, e.g. the entity was not put according to OBXPutMode
OBX_C_API obx_id obx_box_put_object4(OBX_box* box, void* data, size_t size, OBXPutMode mode);

/// Put all given objects in the database in a single transaction. If any of the individual objects failed to put,
/// none are put and an error is returned, equivalent to calling obx_box_put_many5() with fail_on_id_failure=true.
/// @param ids Previously allocated IDs for the given given objects (e.g. using obx_box_ids_for_put)
OBX_C_API obx_err obx_box_put_many(OBX_box* box, const OBX_bytes_array* objects, const obx_id* ids, OBXPutMode mode);

/// Like obx_box_put_many(), but with an additional flag indicating how to treat ID failures with OBXPutMode_INSERT and
/// OBXPutMode_UPDATE.
/// @param fail_on_id_failure if set to true, an ID failure (OBX_ERROR_ID_ALREADY_EXISTS and OBX_ERROR_ID_NOT_FOUND)
///        will fail the transaction, and none of the objects are put/inserted/updated.
/// Note 1: If this function is run inside a managed TX (created by obx_txn_write()) with fail_on_id_failure=true and
///         a failure occurs, the whole outer TX is also aborted.
/// Note 2: ID failure errors are returned even if fail_on_id_failure=false and the TX wasn't aborted.
OBX_C_API obx_err obx_box_put_many5(OBX_box* box, const OBX_bytes_array* objects, const obx_id* ids, OBXPutMode mode,
                                    bool fail_on_id_failure);

/// Remove a single object
/// will return OBX_NOT_FOUND if an object with the given ID doesn't exist
OBX_C_API obx_err obx_box_remove(OBX_box* box, obx_id id);

/// Remove all given objects from the database in a single transaction.
/// Note that this method will not fail if the object is not found (e.g. already removed).
/// In case you need to strictly check whether all of the objects exist before removing them,
/// execute obx_box_contains_ids() and obx_box_remove_ids() inside a single write transaction.
/// @param out_count Pointer to retrieve the number of removed objects; optional: may be NULL.
OBX_C_API obx_err obx_box_remove_many(OBX_box* box, const OBX_id_array* ids, uint64_t* out_count);

/// Remove all objects and set the out_count the the number of removed objects.
/// @param out_count Pointer to retrieve the number of removed objects; optional: may be NULL.
OBX_C_API obx_err obx_box_remove_all(OBX_box* box, uint64_t* out_count);

/// Check whether there are any objects for this entity and updates the out_is_empty accordingly
OBX_C_API obx_err obx_box_is_empty(OBX_box* box, bool* out_is_empty);

/// Count the number of objects in the box, up to the given maximum.
/// You can pass limit=0 to count all objects without any limitation.
OBX_C_API obx_err obx_box_count(OBX_box* box, uint64_t limit, uint64_t* out_count);

/// Fetch IDs of all objects that link back to the given object (ID) using the given relation property (ID).
/// Note: This method refers to "property based relations" unlike the "stand-alone relations" (see obx_box_rel_*).
/// @param property_id the relation property, which must belong to the entity type represented by this box
/// @param id object ID; the object's type is the target of the relation property (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_C_API OBX_id_array* obx_box_get_backlink_ids(OBX_box* box, obx_schema_id property_id, obx_id id);

//----------------------------------------------
// Box + stand-alone relation; obx_box_rel_*
//----------------------------------------------

/// Insert a standalone relation entry between two objects.
/// @param relation_id must be a standalone relation ID with source entity belonging to this box
/// @param source_id identifies an object from this box
/// @param target_id identifies an object from the target box (as per the relation definition)
OBX_C_API obx_err obx_box_rel_put(OBX_box* box, obx_schema_id relation_id, obx_id source_id, obx_id target_id);

/// Remove a standalone relation entry between two objects.
/// See obx_box_rel_put() for parameters documentation.
OBX_C_API obx_err obx_box_rel_remove(OBX_box* box, obx_schema_id relation_id, obx_id source_id, obx_id target_id);

/// Fetch IDs of all objects in this Box related to the given object (typically from another Box).
/// Used for a stand-alone relation and its "regular" direction; this Box represents the target of the relation.
/// @param relation_id ID of a standalone relation, whose target type matches this Box
/// @param id object ID of the relation source type (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_C_API OBX_id_array* obx_box_rel_get_ids(OBX_box* box, obx_schema_id relation_id, obx_id id);

/// Fetch IDs of all objects in this Box related to the given object (typically from another Box).
/// Used for a stand-alone relation and its "backlink" direction; this Box represents the source of the relation.
/// @param relation_id ID of a standalone relation, whose source type matches this Box
/// @param id object ID of the relation target type (typically from another Box)
/// @returns resulting IDs representing objects in this Box, or NULL in case of an error
OBX_C_API OBX_id_array* obx_box_rel_get_backlink_ids(OBX_box* box, obx_schema_id relation_id, obx_id id);

//------------------------------------------------------------------
// Misc ops not tied to cursor/box
// - obx_expired_objects_*
//------------------------------------------------------------------

/// Callback for simple async functions that only deliver a obx_err status.
/// @param status The result status of the async operation
/// @param user_data The data initially passed to the async function call is passed back.
typedef void obx_status_callback(obx_err status, void* user_data);

/// Removes expired objects of one or all entity types.
/// @param entity_id Type of the objects to be remove; if zero, all types are included.
///        Hint: if you only have the entity type's name, use obx_store_entity_id() to get the ID.
/// @param out_removed_count If given (non-null), it will receive the count of removed objects.
/// @see OBXPropertyFlags_EXPIRATION_TIME to define a property for the expiration time.
OBX_C_API obx_err obx_expired_objects_remove(OBX_txn* txn, obx_schema_id entity_id, size_t* out_removed_count);

/// Asynchronously removes expired objects of one or all entity types.
/// @param entity_id Type of the objects to be remove; if zero, all types are included.
///        Hint: if you only have the entity type's name, use obx_store_entity_id() to get the ID.
/// @see OBXPropertyFlags_EXPIRATION_TIME to define a property for the expiration time.
OBX_C_API obx_err obx_expired_objects_remove_async(OBX_store* store, obx_schema_id entity_id,
                                                   obx_status_callback* callback, void* user_data);

//----------------------------------------------
// Time series
//----------------------------------------------

/// Time series: get the limits (min/max time values) over all objects
/// @param out_min_id pointer to receive an output (may be NULL)
/// @param out_min_value pointer to receive an output (may be NULL)
/// @param out_max_id pointer to receive an output (may be NULL)
/// @param out_max_value pointer to receive an output (may be NULL)
/// @returns OBX_NOT_FOUND if no objects are stored
OBX_C_API obx_err obx_box_ts_min_max(OBX_box* box, obx_id* out_min_id, int64_t* out_min_value, obx_id* out_max_id,
                                     int64_t* out_max_value);

/// Time series: get the limits (min/max time values) over objects within the given time range
/// @param out_min_id pointer to receive an output (may be NULL)
/// @param out_min_value pointer to receive an output (may be NULL)
/// @param out_max_id pointer to receive an output (may be NULL)
/// @param out_max_value pointer to receive an output (may be NULL)
/// @returns OBX_NOT_FOUND if no objects are stored in the given range
OBX_C_API obx_err obx_box_ts_min_max_range(OBX_box* box, int64_t range_begin, int64_t range_end, obx_id* out_min_id,
                                           int64_t* out_min_value, obx_id* out_max_id, int64_t* out_max_value);

//----------------------------------------------
// Async
//----------------------------------------------

struct OBX_async;  // doxygen (only) picks up the typedef struct below

/// Created by obx_box_async, used for async operations like obx_async_put.
typedef struct OBX_async OBX_async;

/// Note: DO NOT close this OBX_async; its lifetime is tied to the OBX_box instance.
OBX_C_API OBX_async* obx_async(OBX_box* box);

/// Put asynchronously with standard put semantics (insert or update).
OBX_C_API obx_err obx_async_put(OBX_async* async, obx_id id, const void* data, size_t size);

/// Put asynchronously using the given mode.
OBX_C_API obx_err obx_async_put5(OBX_async* async, obx_id id, const void* data, size_t size, OBXPutMode mode);

/// Put asynchronously with inserts semantics (won't put if object already exists).
OBX_C_API obx_err obx_async_insert(OBX_async* async, obx_id id, const void* data, size_t size);

/// Put asynchronously with update semantics (won't put if object is not yet present).
OBX_C_API obx_err obx_async_update(OBX_async* async, obx_id id, const void* data, size_t size);

/// Reserve an ID, which is returned immediately for future reference, and put asynchronously.
/// Note: of course, it can NOT be guaranteed that the entity will actually be put successfully in the DB.
/// @param data the given bytes are mutated to update the contained ID data.
/// @returns id of the new object, 0 on error
OBX_C_API obx_id obx_async_put_object(OBX_async* async, void* data, size_t size);

/// FB ID slot must be present in the given data; new entities must have an ID value of zero or OBX_ID_NEW
/// @param data writable data buffer, which may be updated for the ID
/// @returns id of the new object, 0 on error, e.g. the entity can't be put according to OBXPutMode
OBX_C_API obx_id obx_async_put_object4(OBX_async* async, void* data, size_t size, OBXPutMode mode);

/// Reserve an ID, which is returned immediately for future reference, and insert asynchronously.
/// Note: of course, it can NOT be guaranteed that the entity will actually be inserted successfully in the DB.
/// @param data the given bytes are mutated to update the contained ID data.
/// @returns id of the new object, 0 on error
OBX_C_API obx_id obx_async_insert_object(OBX_async* async, void* data, size_t size);

/// Remove asynchronously.
OBX_C_API obx_err obx_async_remove(OBX_async* async, obx_id id);

/// Create a custom OBX_async instance that has to be closed using obx_async_close().
/// Note: for standard tasks, prefer obx_box_async() giving you a shared instance that does not have to be closed.
OBX_C_API OBX_async* obx_async_create(OBX_box* box, uint64_t enqueue_timeout_millis);

/// Close a custom OBX_async instance created with obx_async_create().
/// @return OBX_ERROR_ILLEGAL_ARGUMENT if you pass the shared instance from obx_box_async()
OBX_C_API obx_err obx_async_close(OBX_async* async);

//----------------------------------------------
// Query Builder
//----------------------------------------------

struct OBX_query_builder;  // doxygen (only) picks up the typedef struct below

/// You use QueryBuilder to specify criteria and create a Query which actually executes the query and returns matching
/// objects.
typedef struct OBX_query_builder OBX_query_builder;

/// Not really an enum, but binary flags to use across languages
typedef enum {
    /// Reverse the order from ascending (default) to descending.
    OBXOrderFlags_DESCENDING = 1,

    /// Sort upper case letters (e.g. "Z") before lower case letters (e.g. "a").
    /// If not specified, the default is case insensitive for ASCII characters.
    OBXOrderFlags_CASE_SENSITIVE = 2,

    /// For scalars only: change the comparison to unsigned (default is signed).
    OBXOrderFlags_UNSIGNED = 4,

    /// null values will be put last.
    /// If not specified, by default null values will be put first.
    OBXOrderFlags_NULLS_LAST = 8,

    /// null values should be treated equal to zero (scalars only).
    OBXOrderFlags_NULLS_ZERO = 16,
} OBXOrderFlags;

/// Query Builder condition identifier
/// - returned by condition creating functions,
/// - used to combine conditions with any/all, thus building more complex conditions
typedef int obx_qb_cond;

/// Create a query builder which is used to collect conditions using the obx_qb_* functions.
/// Once all conditions are applied, use obx_query() to build a OBX_query that is used to actually retrieve data.
/// Use obx_qb_close() to close (free) the query builder.
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_query_builder* obx_query_builder(OBX_store* store, obx_schema_id entity_id);

/// Close the query builder; note that OBX_query objects outlive their builder and thus are not affected by this call.
/// @param builder may be NULL
OBX_C_API obx_err obx_qb_close(OBX_query_builder* builder);

/// @returns the entity type ID that was used to construct the query builder.
OBX_C_API obx_schema_id obx_qb_type_id(OBX_query_builder* builder);

/// To minimise the amount of error handling code required when building a query, the first error is stored in the query
/// and can be obtained here. All the obx_qb_XXX functions are null operations after the first query error has occurred.
OBX_C_API obx_err obx_qb_error_code(OBX_query_builder* builder);

/// To minimise the amount of error handling code required when building a query, the first error is stored in the query
/// and can be obtained here. All the obx_qb_XXX functions are null operations after the first query error has occurred.
OBX_C_API const char* obx_qb_error_message(OBX_query_builder* builder);

/// Add null check to the query
OBX_C_API obx_qb_cond obx_qb_null(OBX_query_builder* builder, obx_schema_id property_id);

/// Add not-null check to the query
OBX_C_API obx_qb_cond obx_qb_not_null(OBX_query_builder* builder, obx_schema_id property_id);

// String conditions ---------------------------

/// For string and flex properties, this looks for a value equal to the given one.
OBX_C_API obx_qb_cond obx_qb_equals_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                           bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_not_equals_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                               bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_contains_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                             bool case_sensitive);

/// For container types (e.g. string vector or flex), this looks for a element that equals the given string.
OBX_C_API obx_qb_cond obx_qb_contains_element_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                     const char* value, bool case_sensitive);

/// For flex properties that have a map as root value, this looks for matching key/value pair.
OBX_C_API obx_qb_cond obx_qb_contains_key_value_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                       const char* key, const char* value, bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_starts_with_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                const char* value, bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_ends_with_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                              bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_greater_than_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                 const char* value, bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_greater_or_equal_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                     const char* value, bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_less_than_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                              bool case_sensitive);

OBX_C_API obx_qb_cond obx_qb_less_or_equal_string(OBX_query_builder* builder, obx_schema_id property_id,
                                                  const char* value, bool case_sensitive);

/// Note that all string values are copied and thus do not need to be maintained by the calling code.
OBX_C_API obx_qb_cond obx_qb_in_strings(OBX_query_builder* builder, obx_schema_id property_id,
                                        const char* const values[], size_t count, bool case_sensitive);

/// For OBXPropertyType_StringVector - matches if at least one vector item equals the given value.
OBX_C_API obx_qb_cond obx_qb_any_equals_string(OBX_query_builder* builder, obx_schema_id property_id, const char* value,
                                               bool case_sensitive);

// Integral conditions -------------------------

OBX_C_API obx_qb_cond obx_qb_equals_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
OBX_C_API obx_qb_cond obx_qb_not_equals_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);

OBX_C_API obx_qb_cond obx_qb_greater_than_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
OBX_C_API obx_qb_cond obx_qb_greater_or_equal_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
OBX_C_API obx_qb_cond obx_qb_less_than_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
OBX_C_API obx_qb_cond obx_qb_less_or_equal_int(OBX_query_builder* builder, obx_schema_id property_id, int64_t value);
OBX_C_API obx_qb_cond obx_qb_between_2ints(OBX_query_builder* builder, obx_schema_id property_id, int64_t value_a,
                                           int64_t value_b);

/// Note that all values are copied and thus do not need to be maintained by the calling code.
OBX_C_API obx_qb_cond obx_qb_in_int64s(OBX_query_builder* builder, obx_schema_id property_id, const int64_t values[],
                                       size_t count);

/// Note that all values are copied and thus do not need to be maintained by the calling code.
OBX_C_API obx_qb_cond obx_qb_not_in_int64s(OBX_query_builder* builder, obx_schema_id property_id,
                                           const int64_t values[], size_t count);

/// Note that all values are copied and thus do not need to be maintained by the calling code.
OBX_C_API obx_qb_cond obx_qb_in_int32s(OBX_query_builder* builder, obx_schema_id property_id, const int32_t values[],
                                       size_t count);

/// Note that all values are copied and thus do not need to be maintained by the calling code.
OBX_C_API obx_qb_cond obx_qb_not_in_int32s(OBX_query_builder* builder, obx_schema_id property_id,
                                           const int32_t values[], size_t count);

// FP conditions -------------------------------
// Note: works for float and double properties

OBX_C_API obx_qb_cond obx_qb_greater_than_double(OBX_query_builder* builder, obx_schema_id property_id, double value);
OBX_C_API obx_qb_cond obx_qb_greater_or_equal_double(OBX_query_builder* builder, obx_schema_id property_id,
                                                     double value);
OBX_C_API obx_qb_cond obx_qb_less_than_double(OBX_query_builder* builder, obx_schema_id property_id, double value);
OBX_C_API obx_qb_cond obx_qb_less_or_equal_double(OBX_query_builder* builder, obx_schema_id property_id, double value);
OBX_C_API obx_qb_cond obx_qb_between_2doubles(OBX_query_builder* builder, obx_schema_id property_id, double value_a,
                                              double value_b);

// Bytes (blob) conditions ---------------------

OBX_C_API obx_qb_cond obx_qb_equals_bytes(OBX_query_builder* builder, obx_schema_id property_id, const void* value,
                                          size_t size);

OBX_C_API obx_qb_cond obx_qb_greater_than_bytes(OBX_query_builder* builder, obx_schema_id property_id,
                                                const void* value, size_t size);

OBX_C_API obx_qb_cond obx_qb_greater_or_equal_bytes(OBX_query_builder* builder, obx_schema_id property_id,
                                                    const void* value, size_t size);

OBX_C_API obx_qb_cond obx_qb_less_than_bytes(OBX_query_builder* builder, obx_schema_id property_id, const void* value,
                                             size_t size);

OBX_C_API obx_qb_cond obx_qb_less_or_equal_bytes(OBX_query_builder* builder, obx_schema_id property_id,
                                                 const void* value, size_t size);

// Other conditions ---------------------

/// An object matches, if it has a given number of related objects pointing to it.
/// At this point, there are a couple of limitations (later version may improve on that):
/// 1) 1:N relations only, 2) the complexity is O(n * (relationCount + 1)) and cannot be improved via indexes,
/// 3) The relation count cannot be set as an parameter.
/// @param relation_entity_id ID of the entity type the relation comes from.
/// @param relation_property_id ID of the property in the related entity type representing the relation.
/// @param relation_count Number of related object an object must have to match. May be 0 if objects shall be matched
///        that do not have related objects. (Typically low numbers are used for the count.)
OBX_C_API obx_qb_cond obx_qb_relation_count_property(OBX_query_builder* builder, obx_schema_id relation_entity_id,
                                                     obx_schema_id relation_property_id, int32_t relation_count);

// Other condition related functions ---------------------

/// Combine conditions[] to a new condition using operator AND (all).
OBX_C_API obx_qb_cond obx_qb_all(OBX_query_builder* builder, const obx_qb_cond conditions[], size_t count);

/// Combine conditions[] to a new condition using operator OR (any).
OBX_C_API obx_qb_cond obx_qb_any(OBX_query_builder* builder, const obx_qb_cond conditions[], size_t count);

/// Create an alias for the previous condition (the one added just before calling this function).
/// This is useful when you have a query with multiple conditions of the same property (e.g. height < 20 or height > 50)
/// and you want to use obx_query_param_* to change the values. Consider the following simplified example.
///
/// **Example:** Create a query with two aliased params and set their values later during query execution:
///
///          OBX_query_builder* qb = obx_query_builder(store, entity_id);
///          obx_qb_less_than_int(qb, height_prop_id, 0)
///          obx_qb_param_alias(qb, "height-lt")
///          obx_qb_greater_than_int(qb, height_prop_id, 0)
///          obx_qb_param_alias(qb, "height-gt")
///          OBX_query* query = obx_query(OBX_query_builder* qb);
///          ...
///          obx_query_param_alias_int(query, "height-lt", 20)
///          obx_query_param_alias_int(query, "height-gt", 50)
///          OBX_bytes_array* results = obx_query_find(query)
///          obx_query_param_alias_int(query, "height-lt", 100)
///          obx_query_param_alias_int(query, "height-gt", 500)
///          OBX_bytes_array* results2 = obx_query_find(query)
/// @param alias any non-empty string
OBX_C_API obx_err obx_qb_param_alias(OBX_query_builder* builder, const char* alias);

/// Configures an order of results in the query
OBX_C_API obx_err obx_qb_order(OBX_query_builder* builder, obx_schema_id property_id, OBXOrderFlags flags);

/// Create a link based on a property-relation (many-to-one)
OBX_C_API OBX_query_builder* obx_qb_link_property(OBX_query_builder* builder, obx_schema_id property_id);

/// Create a backlink based on a property-relation used in reverse (one-to-many)
OBX_C_API OBX_query_builder* obx_qb_backlink_property(OBX_query_builder* builder, obx_schema_id source_entity_id,
                                                      obx_schema_id source_property_id);

/// Create a link based on a standalone relation (many-to-many)
OBX_C_API OBX_query_builder* obx_qb_link_standalone(OBX_query_builder* builder, obx_schema_id relation_id);

/// Create a backlink based on a standalone relation (many-to-many, reverse direction)
OBX_C_API OBX_query_builder* obx_qb_backlink_standalone(OBX_query_builder* builder, obx_schema_id relation_id);

/// Link the (time series) entity type to another entity space using a time point or range defined in the given
/// linked entity type and properties.
/// Note: time series functionality must be available to use this.
/// @param linked_entity_id Entity type that defines a time point or range
/// @param begin_property_id Property of the linked entity defining a time point or the begin of a time range.
///        Must be a date type (e.g. PropertyType_Date or PropertyType_DateNano).
/// @param end_property_id Optional property of the linked entity defining the end of a time range.
///        Pass zero to only define a time point (begin_property_id).
///        Must be a date type (e.g. PropertyType_Date or PropertyType_DateNano).
OBX_C_API OBX_query_builder* obx_qb_link_time(OBX_query_builder* builder, obx_schema_id linked_entity_id,
                                              obx_schema_id begin_property_id, obx_schema_id end_property_id);

//----------------------------------------------
// Query
//----------------------------------------------

struct OBX_query;  // doxygen (only) picks up the typedef struct below

/// Query holds the information necessary to execute a database query. It's prepared by QueryBuilder and may be reused
/// any number of times. It also supports parametrization before executing, further improving the reusability.
/// Query is NOT thread safe and must only be used from a single thread at the same time. If you prefer to avoid locks,
/// you may want to create clonse using obx_query_clone().
typedef struct OBX_query OBX_query;

/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_query* obx_query(OBX_query_builder* builder);

/// Close the query and free resources.
OBX_C_API obx_err obx_query_close(OBX_query* query);

/// Create a clone of the given query such that it can be run on a separate thread
OBX_C_API OBX_query* obx_query_clone(OBX_query* query);

/// Configure an offset for this query - all methods that support offset will return/process objects starting at this
/// offset. Example use case: use together with limit to get a slice of the whole result, e.g. for "result paging".
/// Call with offset=0 to reset to the default behavior, i.e. starting from the first element.
OBX_C_API obx_err obx_query_offset(OBX_query* query, size_t offset);

/// Configure an offset and a limit for this query - all methods that support an offset/limit will return/process
/// objects starting at this offset and up to the given limit. Example use case: get a slice of the whole result, e.g.
/// for "result paging". Call with offset/limit=0 to reset to the default behavior, i.e. starting from the first element
/// without limit.
OBX_C_API obx_err obx_query_offset_limit(OBX_query* query, size_t offset, size_t limit);

/// Configure a limit for this query - all methods that support limit will return/process only the given number of
/// objects. Example use case: use together with offset to get a slice of the whole result, e.g. for "result paging".
/// Call with limit=0 to reset to the default behavior - zero limit means no limit applied.
OBX_C_API obx_err obx_query_limit(OBX_query* query, size_t limit);

/// Find entities matching the query. NOTE: the returned data is only valid as long the transaction is active!
OBX_C_API OBX_bytes_array* obx_query_find(OBX_query* query);

/// Find the first object matching the query.
/// @returns OBX_NOT_FOUND if no object matches.
/// The exposed data comes directly from the OS to allow zero-copy access, which limits the data lifetime:
/// @warning Currently ignores offset, taking the the first matching element.
/// @attention The exposed data is only valid as long as the (top) transaction is still active and no write
///            operation (e.g. put/remove) was executed. Accessing data after this is undefined behavior.
OBX_C_API obx_err obx_query_find_first(OBX_query* query, const void** data, size_t* size);

/// Find the only object matching the query.
/// @returns OBX_NOT_FOUND if no object matches, an error if there are multiple objects matching the query.
/// The exposed data comes directly from the OS to allow zero-copy access, which limits the data lifetime:
/// @warning Currently ignores offset and limit, considering all matching elements.
/// @attention The exposed data is only valid as long as the (top) transaction is still active and no write
///            operation (e.g. put/remove) was executed. Accessing data after this is undefined behavior.
OBX_C_API obx_err obx_query_find_unique(OBX_query* query, const void** data, size_t* size);

/// Walk over matching objects using the given data visitor
OBX_C_API obx_err obx_query_visit(OBX_query* query, obx_data_visitor* visitor, void* user_data);

/// Return the IDs of all matching objects
OBX_C_API OBX_id_array* obx_query_find_ids(OBX_query* query);

/// Return the number of matching objects
OBX_C_API obx_err obx_query_count(OBX_query* query, uint64_t* out_count);

/// Remove all matching objects from the database & return the number of deleted objects
OBX_C_API obx_err obx_query_remove(OBX_query* query, uint64_t* out_count);

/// The returned char* is valid until another call to describe() is made on the query or until the query is freed
OBX_C_API const char* obx_query_describe(OBX_query* query);

/// The returned char* is valid until another call to describe_params() is made on the query or until the query is freed
OBX_C_API const char* obx_query_describe_params(OBX_query* query);

//----------------------------------------------
// Query using Cursor (lower level API)
//----------------------------------------------
OBX_C_API obx_err obx_query_cursor_visit(OBX_query* query, OBX_cursor* cursor, obx_data_visitor* visitor,
                                         void* user_data);

/// Find entities matching the query; NOTE: the returned data is only valid as long the transaction is active!
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_bytes_array* obx_query_cursor_find(OBX_query* query, OBX_cursor* cursor);

OBX_C_API OBX_id_array* obx_query_cursor_find_ids(OBX_query* query, OBX_cursor* cursor);

OBX_C_API obx_err obx_query_cursor_count(OBX_query* query, OBX_cursor* cursor, uint64_t* out_count);

/// Remove (delete!) all matching objects.
OBX_C_API obx_err obx_query_cursor_remove(OBX_query* query, OBX_cursor* cursor, uint64_t* out_count);

//----------------------------------------------
// Query parameters - obx_query_param_{type}(s).
// Change previously set condition value in an existing query - this improves reusability of the query object.
//----------------------------------------------
OBX_C_API obx_err obx_query_param_string(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                         const char* value);
OBX_C_API obx_err obx_query_param_2strings(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                           const char* value, const char* value2);
OBX_C_API obx_err obx_query_param_strings(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                          const char* const values[], size_t count);
OBX_C_API obx_err obx_query_param_int(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                      int64_t value);
OBX_C_API obx_err obx_query_param_2ints(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                        int64_t value_a, int64_t value_b);
OBX_C_API obx_err obx_query_param_int64s(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                         const int64_t values[], size_t count);
OBX_C_API obx_err obx_query_param_int32s(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                         const int32_t values[], size_t count);
OBX_C_API obx_err obx_query_param_double(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                         double value);
OBX_C_API obx_err obx_query_param_2doubles(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                           double value_a, double value_b);
OBX_C_API obx_err obx_query_param_bytes(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id,
                                        const void* value, size_t size);

/// Gets the size of the property type used in a query condition.
/// A typical use case of this is to allow language bindings (e.g. Swift) use the right type (e.g. 32 bit ints) even
/// if the language has a bias towards another type (e.g. 64 bit ints).
/// @returns the size of the underlying property
/// @returns 0 if it does not have a fixed size (e.g. strings, vectors) or an error occurred
OBX_C_API size_t obx_query_param_get_type_size(OBX_query* query, obx_schema_id entity_id, obx_schema_id property_id);

//----------------------------------------------
// Query parameters with alias - obx_query_param_alias_{type}(s).
// Similar to obx_query_param_{type}, but when an alias was used for a parameter, see obx_qb_param_alias().
//----------------------------------------------
OBX_C_API obx_err obx_query_param_alias_string(OBX_query* query, const char* alias, const char* value);
OBX_C_API obx_err obx_query_param_alias_strings(OBX_query* query, const char* alias, const char* const values[],
                                                size_t count);
OBX_C_API obx_err obx_query_param_alias_int(OBX_query* query, const char* alias, int64_t value);
OBX_C_API obx_err obx_query_param_alias_2ints(OBX_query* query, const char* alias, int64_t value_a, int64_t value_b);
OBX_C_API obx_err obx_query_param_alias_int64s(OBX_query* query, const char* alias, const int64_t values[],
                                               size_t count);
OBX_C_API obx_err obx_query_param_alias_int32s(OBX_query* query, const char* alias, const int32_t values[],
                                               size_t count);
OBX_C_API obx_err obx_query_param_alias_double(OBX_query* query, const char* alias, double value);
OBX_C_API obx_err obx_query_param_alias_2doubles(OBX_query* query, const char* alias, double value_a, double value_b);
OBX_C_API obx_err obx_query_param_alias_bytes(OBX_query* query, const char* alias, const void* value, size_t size);

/// Gets the size of the property type used in a query condition.
/// A typical use case of this is to allow language bindings (e.g. Swift) use the right type (e.g. 32 bit ints) even
/// if the language has a bias towards another type (e.g. 64 bit ints).
/// @returns the size of the underlying property
/// @returns 0 if it does not have a fixed size (e.g. strings, vectors) or an error occurred
OBX_C_API size_t obx_query_param_alias_get_type_size(OBX_query* query, const char* alias);

//----------------------------------------------
// Property-Query - getting a single property instead of whole objects
//----------------------------------------------

struct OBX_query_prop;  // doxygen (only) picks up the typedef struct below

/// PropertyQuery - getting a single property instead of whole objects. Also provides aggregation over properties.
typedef struct OBX_query_prop OBX_query_prop;

/// Create a "property query" with results referring to single property (not complete objects).
/// Also provides aggregates like for example obx_query_prop_avg().
OBX_C_API OBX_query_prop* obx_query_prop(OBX_query* query, obx_schema_id property_id);

/// Close the property query and release resources.
OBX_C_API obx_err obx_query_prop_close(OBX_query_prop* query);

/// Configure the property query to work only on distinct values.
/// @note not all methods support distinct, those that don't will return an error
OBX_C_API obx_err obx_query_prop_distinct(OBX_query_prop* query, bool distinct);

/// Configure the property query to work only on distinct values.
/// This version is reserved for string properties and defines the case sensitivity for distinctness.
/// @note not all methods support distinct, those that don't will return an error
OBX_C_API obx_err obx_query_prop_distinct_case(OBX_query_prop* query, bool distinct, bool case_sensitive);

/// Count the number of non-NULL values of the given property across all objects matching the query
OBX_C_API obx_err obx_query_prop_count(OBX_query_prop* query, uint64_t* out_count);

/// Calculate an average value for the given numeric property across all objects matching the query.
/// @param query the query to run
/// @param out_average the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. a floating point NaN value will trigger the short cut as the average will be a NaN no matter
///                  what values will follow.
OBX_C_API obx_err obx_query_prop_avg(OBX_query_prop* query, double* out_average, int64_t* out_count);

/// Calculate an average value for the given numeric property across all objects matching the query.
/// @param query the query to run
/// @param out_average the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
/// @returns OBX_ERROR_NUMERIC_OVERFLOW if the result does not fit into an int64_t
OBX_C_API obx_err obx_query_prop_avg_int(OBX_query_prop* query, int64_t* out_average, int64_t* out_count);

/// Find the minimum value of the given floating-point property across all objects matching the query.
/// @param query the query to run
/// @param out_minimum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. if an index is used, it will be set to 0 or -1, instead of the actual count of objects.
OBX_C_API obx_err obx_query_prop_min(OBX_query_prop* query, double* out_minimum, int64_t* out_count);

/// Find the maximum value of the given floating-point property across all objects matching the query
/// @param query the query to run
/// @param out_maximum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. if an index is used, it will be set to 0 or -1, instead of the actual count of objects.
OBX_C_API obx_err obx_query_prop_max(OBX_query_prop* query, double* out_maximum, int64_t* out_count);

/// Calculate the sum of the given floating-point property across all objects matching the query.
/// @param query the query to run
/// @param out_sum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. a floating point NaN value will trigger the short cut as the average will be a NaN no matter
///                  what values will follow.
OBX_C_API obx_err obx_query_prop_sum(OBX_query_prop* query, double* out_sum, int64_t* out_count);

/// Find the minimum value of the given property across all objects matching the query.
/// @param query the query to run
/// @param out_minimum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. if an index is used, it will be set to 0 or -1, instead of the actual count of objects.
OBX_C_API obx_err obx_query_prop_min_int(OBX_query_prop* query, int64_t* out_minimum, int64_t* out_count);

/// Find the maximum value of the given property across all objects matching the query.
/// @param query the query to run
/// @param out_maximum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
///                  E.g. if an index is used, it will be set to 0 or -1, instead of the actual count of objects.
OBX_C_API obx_err obx_query_prop_max_int(OBX_query_prop* query, int64_t* out_maximum, int64_t* out_count);

/// Calculate the sum of the given property across all objects matching the query.
/// @param query the query to run
/// @param out_sum the result of the query
/// @param out_count (optional, may be NULL) number of objects contributing to the result (counted on the fly).
///                  A negative count indicates that the computation used a short cut and thus the count is incomplete.
/// @returns OBX_ERROR_NUMERIC_OVERFLOW if the result does not fit into an int64_t
OBX_C_API obx_err obx_query_prop_sum_int(OBX_query_prop* query, int64_t* out_sum, int64_t* out_count);

/// Return an array of strings stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified field are skipped
OBX_C_API OBX_string_array* obx_query_prop_find_strings(OBX_query_prop* query, const char* value_if_null);

/// Return an array of ints stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_int64_array* obx_query_prop_find_int64s(OBX_query_prop* query, const int64_t* value_if_null);

/// Return an array of ints stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_int32_array* obx_query_prop_find_int32s(OBX_query_prop* query, const int32_t* value_if_null);

/// Return an array of ints stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_int16_array* obx_query_prop_find_int16s(OBX_query_prop* query, const int16_t* value_if_null);

/// Return an array of ints stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_int8_array* obx_query_prop_find_int8s(OBX_query_prop* query, const int8_t* value_if_null);

/// Return an array of doubles stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_double_array* obx_query_prop_find_doubles(OBX_query_prop* query, const double* value_if_null);

/// Return an array of int stored as the given property across all objects matching the query.
/// @param value_if_null value that should be used in place of NULL values on object fields;
///     if value_if_null=NULL is given, objects with NULL values of the specified are skipped
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_float_array* obx_query_prop_find_floats(OBX_query_prop* query, const float* value_if_null);

struct OBX_observer;  // doxygen (only) picks up the typedef struct below

/// Observers are called back when data has changed in the database.
/// See obx_observe(), or obx_observe_single_type() to listen to a changes that affect a single entity type
typedef struct OBX_observer OBX_observer;

/// Callback for obx_observe()
/// @param user_data user data given to obx_observe()
/// @param type_ids array of object type IDs that had changes
/// @param type_ids_count number of IDs of type_ids
typedef void obx_observer(const obx_schema_id* type_ids, size_t type_ids_count, void* user_data);

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
/// @returns NULL if a illegal locking situation was detected, e.g. called from an observer itself or a
///          timeout/deadlock was detected (OBX_ERROR_ILLEGAL_STATE).
OBX_C_API OBX_observer* obx_observe(OBX_store* store, obx_observer* callback, void* user_data);

/// Create an observer (callback) to be notified about data changes for a given object type.
/// The callback is invoked right after a successful commit.
/// \note  If you intend to observe more than one type, it is more efficient to use obx_observe().
/// \par Threading note
/// The given callback is called on the thread issuing the commit for the data change, e.g. via obx_txn_success().
/// Future versions might change that to a background thread, so be careful with threading assumptions.
/// Also, it's a usually good idea to make the callback return quickly to let the calling thread continue.
/// \attention Currently, you can not call any data operations from inside the call back.
/// \attention More accurately, no transaction may be started. (This restriction may be removed in a later version.)
/// @param type_id ID of the object type to be observer.
/// @param user_data any value you want to be forwarded to the given observer callback (usually some context info).
/// @param callback pointer to be called when the observed data changes.
/// @returns NULL if a illegal locking situation was detected, e.g. called from an observer itself or a
///          timeout/deadlock was detected (OBX_ERROR_ILLEGAL_STATE).
OBX_C_API OBX_observer* obx_observe_single_type(OBX_store* store, obx_schema_id type_id,
                                                obx_observer_single_type* callback, void* user_data);

/// Free the memory used by the given observer and unsubscribe it from its box or query.
/// @returns OBX_ERROR_ILLEGAL_STATE if a illegal locking situation was detected, e.g. called from an observer itself
///          or a timeout/deadlock was detected. In that case, the caller must try to close again in a valid situation
///          not causing lock failures.
OBX_C_API obx_err obx_observer_close(OBX_observer* observer);

//----------------------------------------------
// Tree
//----------------------------------------------
typedef enum {
    /// If true, debug logs are always disabled for this tree regardless of the store's debug flags.
    OBXTreeOptionFlags_DebugLogsDisable = 1,

    /// If true, debug logs are always enabled for this tree regardless of the store's debug flags.
    OBXTreeOptionFlags_DebugLogsEnable = 2,

    /// By default, a path such as "a/b/c" can address a branch and a leaf at the same time.
    /// E.g. under the common parent path "a/b", a branch "c" and a "c" leaf may exist.
    /// To disable this, set this flag to true.
    /// This will enable an additional check when inserting new leafs and new branches for the existence of the other.
    OBXTreeOptionFlags_EnforceUniquePath = 4,

    /// In some scenarios, e.g. when using Sync, multiple node objects of the same type (e.g. branch or leaf) at the
    /// same path may exist temporarily. By enabling this flag, this is not considered an error situation. Instead, the
    /// first node is picked.
    OBXTreeOptionFlags_AllowNonUniqueNodes = 8,

    /// Nodes described in AllowNonUniqueNodes will be automatically detected to consolidate them (manually).
    OBXTreeOptionFlags_DetectNonUniqueNodes = 16,

    /// Nodes described in AllowNonUniqueNodes will be automatically consolidated to make them unique.
    /// This consolidation happens e.g. on put/remove operations.
    /// Using this value implies DetectNonUniqueNodes.
    OBXTreeOptionFlags_AutoConsolidateNonUniqueNodes = 32,

} OBXTreeOptionFlags;

struct OBX_tree_options;  // doxygen (only) picks up the typedef struct below

/// @brief Structural/behavioral options for a tree passed during tree creation.
typedef struct OBX_tree_options OBX_tree_options;

struct OBX_tree;  // doxygen (only) picks up the typedef struct below

/// @brief Top level tree API representing a tree structure/schema associated with a store.
///
/// Data is accessed via "tree cursors" that can be created using obx_tree_cursor().
typedef struct OBX_tree OBX_tree;

struct OBX_tree_cursor;  // doxygen (only) picks up the typedef struct below
typedef struct OBX_tree_cursor OBX_tree_cursor;

/// Creates a options object that is passed to tree creation via obx_tree().
OBX_C_API OBX_tree_options* obx_tree_options();

/// Free the tree options if they were not used (and "consumed") by obx_tree().
/// Note: Only free *unused* options, obx_tree() frees the options internally.
OBX_C_API void obx_tree_options_free(OBX_tree_options* options);

/// Adjusts the path delimiter character, which is by default "/".
OBX_C_API obx_err obx_tree_opt_path_delimiter(OBX_tree_options* options, char path_delimiter);

/// Sets the given OBXTreeOptionFlags at the tree options.
/// Combine multiple flags using bitwise OR.
OBX_C_API obx_err obx_tree_opt_flags(OBX_tree_options* options, uint32_t flags);

/// Creates a tree for the given store.
/// @param options can be null for default options (e.g. use '/' as a path delimiter).
///        The options are always freed when calling this function.
OBX_C_API OBX_tree* obx_tree(OBX_store* store, OBX_tree_options* options);
OBX_C_API void obx_tree_close(OBX_tree* tree);

/// To get/put data from tree, you need to create a tree cursor using this method.
/// @param txn may be null if the transaction will be set later via obx_tree_cursor_txn()
OBX_C_API OBX_tree_cursor* obx_tree_cursor(OBX_tree* tree, OBX_txn* txn);

/// Gets the number of currently tracked node conflicts (non-unique nodes at the same path).
/// This count gets resent when conflicts get consolidated.
/// Only tracked if OBXTreeOptionFlags_DetectNonUniqueNodes (or OBXTreeOptionFlags_AutoConsolidateNonUniqueNodes) is
/// set.
OBX_C_API size_t obx_tree_node_conflict_count(OBX_tree* tree);

/// \brief Closes the tree cursor, e.g. before a transaction is ending.
/// The cursor cannot be used afterwards (consider obx_tree_cursor_reset() if you want to keep using it).
OBX_C_API void obx_tree_cursor_close(OBX_tree_cursor* cursor);

/// \brief Sets or clears a transaction from the tree cursor.
///
/// A typical use case for this function is to cache the tree cursor for reusing it later.
/// Note: before closing a transaction, ensure to clear it here first (set to null).
OBX_C_API obx_err obx_tree_cursor_txn(OBX_tree_cursor* cursor, OBX_txn* txn);

/// \brief A "low-level" get operation to access a tree leaf using the raw FlatBuffer bytes stored in the database.
/// As usual, the data is only valid during the lifetime of the transaction and before the first write to the DB.
/// An advantage of using "raw" operations is that custom properties can be passed in the FlatBuffer.
/// @param data receiver of the data pointer (non-null pointer to a pointer), which will be pointing to FlatBuffers
///        bytes for the data leaf after this call.
/// @param metadata optional FlatBuffers receiver (nullable pointer to a pointer) for the meta leaf.
/// @returns general success/error codes and some tree-specific ones (OBX_ERROR_TREE_*)
OBX_C_API obx_err obx_tree_cursor_get_raw(OBX_tree_cursor* cursor, const char* path, const void** data, size_t* size,
                                          const void** metadata, size_t* metadata_size);

/// \brief A "low-level" put operation for a tree leaf using given raw FlatBuffer bytes.
/// Any non-existing branches and meta nodes are put on the fly if an optional meta-leaf FlatBuffers is given.
/// A typical usage pattern is to first try without the meta-leaf, and if it does not work, create the meta-leaf and
/// call again with the meta leaf. This approach takes into account that meta-leaves typically exist, and thus no
/// resources are wasted for meta-leaf creation if it's not required.
/// An advantage of using "raw" operations is that custom properties can be passed in the FlatBuffers.
/// @param leaf_data prepared FlatBuffers bytes for the data leaf (non-const as the data buffer will be mutated);
///        note: slots for IDs must be already be prepared (zero values)
/// @param type value type of the given leafObject: it has to be passed to verify against stored metadata
/// @param leaf_metadata optional FlatBuffers for the meta leaf; at minimum, "name" and "type" must be set.
///        Note: slots for IDs must be already be prepared (zero values).
///        Passing null indicates that the branches and meta nodes must already exist; otherwise
///        the operation will fail and OBX_NOT_FOUND will be returned.
/// @param leaf_put_mode For the data leaf only (the actual value tree node)
/// @returns OBX_NOT_FOUND if the path did not exist (and no metadata was given)
/// @returns OBX_NO_SUCCESS if the put was not successful according to the put mode
/// @returns also other standard result codes like OBX_SUCCESS or OBX_ERR_*
OBX_C_API obx_err obx_tree_cursor_put_raw(OBX_tree_cursor* cursor, const char* path, void* leaf_data,
                                          size_t leaf_data_size, OBXPropertyType type, obx_id* out_id,
                                          void* leaf_metadata, size_t leaf_metadata_size, OBXPutMode leaf_put_mode);

/// Explicitly consolidates tree node conflicts (non unique nodes sharing a common path).
/// See also obx_tree_async_consolidate_node_conflicts() for an asynchronous version.
OBX_C_API obx_err obx_tree_cursor_consolidate_node_conflicts(OBX_tree_cursor* cursor, size_t* out_consolidated_count);

/// Callback for obx_tree_async_put_raw().
/// \note If the given status is an error, you can use functions like obx_last_error_message() to gather more info
/// during this callback (error state is thread bound and the callback uses an internal thread).
/// @param status The result status of the async operation
/// @param id If the operation was successful, the ID of the leaf, which was put (otherwise zero).
/// @param user_data The data initially passed to the async function call is passed back.
typedef void obx_tree_async_put_callback(obx_err status, obx_id id, void* user_data);

/// Like obx_tree_cursor_put_raw(), but asynchronous.
/// @param callback Optional (may be null) function that is called with results once the async operation completes.
/// @param callback_user_data Any value you can supply, which is passed on to the callback (e.g. to identify user
///        specific context).
OBX_C_API obx_err obx_tree_async_put_raw(OBX_tree* tree, const char* path, void* leaf_data, size_t leaf_data_size,
                                         OBXPropertyType type, void* leaf_metadata, size_t leaf_metadata_size,
                                         OBXPutMode leaf_put_mode, obx_tree_async_put_callback* callback,
                                         void* callback_user_data);

/// Explicitly consolidates tree node conflicts (non unique nodes sharing a common path) asynchronously.
/// See also obx_tree_cursor_consolidate_node_conflicts() for a synchronous version.
OBX_C_API obx_err obx_tree_async_consolidate_node_conflicts(OBX_tree* tree);

//----------------------------------------------
// Weak store
//----------------------------------------------

struct OBX_weak_store;  // doxygen (only) picks up the typedef struct below

/// \brief A "weak store" is a weakly referenced store.
///
/// While you can "attach" to a store multiple time; each instance will keep the store open.
/// Alternatively for "secondary stores" (e.g. used in a background thread), a weakly referenced store does not keep
/// the store open by itself. Only strong references (e.g. the store created and attached stores) do.
/// A weak store needs to be "locked" to get a regular store (strong reference) in order to actually use it.
typedef struct OBX_weak_store OBX_weak_store;

/// Creates a weak reference to the given store.
OBX_C_API OBX_weak_store* obx_weak_store(OBX_store* store);

/// Frees a weak store reference.
OBX_C_API void obx_weak_store_free(OBX_weak_store* weak_store);

/// Tries to get a weak reference to the store with the given ID.
/// @returns null if no store was found for the given ID
OBX_C_API OBX_weak_store* obx_weak_store_by_id(uint64_t store_id);

/// Tries to create a regular store instance for the given weak reference.
/// @returns null if the store was already closed (all "strong references" are gone)
OBX_C_API OBX_store* obx_weak_store_lock(OBX_weak_store* weak_store);

//----------------------------------------------
// Utilities for bytes/ids/arrays
//----------------------------------------------
OBX_C_API void obx_bytes_free(OBX_bytes* bytes);

/// Allocate a bytes array struct of the given size, ready for the data to be pushed
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_bytes_array* obx_bytes_array(size_t count);

/// Set the given data as the index in the bytes array. The data is not copied, just referenced through the pointer
OBX_C_API obx_err obx_bytes_array_set(OBX_bytes_array* array, size_t index, const void* data, size_t size);

/// Free the bytes array struct
OBX_C_API void obx_bytes_array_free(OBX_bytes_array* array);

/// Create an ID array struct, copying the given IDs as the contents
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_id_array* obx_id_array(const obx_id* ids, size_t count);

/// Free the array struct
OBX_C_API void obx_id_array_free(OBX_id_array* array);

/// Free the array struct
OBX_C_API void obx_string_array_free(OBX_string_array* array);

/// Free the array struct
OBX_C_API void obx_int64_array_free(OBX_int64_array* array);

/// Free the array struct
OBX_C_API void obx_int32_array_free(OBX_int32_array* array);

/// Free the array struct
OBX_C_API void obx_int16_array_free(OBX_int16_array* array);

/// Free the array struct
OBX_C_API void obx_int8_array_free(OBX_int8_array* array);

/// Free the array struct
OBX_C_API void obx_double_array_free(OBX_double_array* array);

/// Free the array struct
OBX_C_API void obx_float_array_free(OBX_float_array* array);

/// Only for Apple platforms: set the prefix to use for mutexes based on POSIX semaphores.
/// You must supply the application group identifier for sand-boxed macOS apps here; see also:
/// https://developer.apple.com/library/archive/documentation/Security/Conceptual/AppSandboxDesignGuide/AppSandboxInDepth/AppSandboxInDepth.html#//apple_ref/doc/uid/TP40011183-CH3-SW24
/// @param prefix must be at most 20 characters long and end with a forward slash '/'.
OBX_C_API obx_err obx_posix_sem_prefix_set(const char* prefix);

//----------------------------------------------
// Admin web UI (e.g. object browser, DB & Sync stats, admin tasks)
//----------------------------------------------

struct OBX_admin_options;  // doxygen (only) picks up the typedef struct below

/// Options to configure the OBX_admin to be created
typedef struct OBX_admin_options OBX_admin_options;

/// Create a default set of admin options.
/// @returns NULL on failure, a default set of options on success
OBX_C_API OBX_admin_options* obx_admin_opt();

/// Configure admin with an existing, open OBX_store*.
/// @see also obx_admin_opt_store_path() as an alternative.
OBX_C_API obx_err obx_admin_opt_store(OBX_admin_options* opt, OBX_store* store);

/// Configure admin with an existing, open OBX_store*.
/// @see also obx_admin_opt_store() as an alternative.
OBX_C_API obx_err obx_admin_opt_store_path(OBX_admin_options* opt, const char* directory);

/// Set the address and port on which the underlying http-server should server the admin web UI.
/// Defaults to "http://127.0.0.1:8081"
/// @note: you can use for e.g. "http://127.0.0.1:0" for automatic free port assignment - see obx_admin_bound_port().
OBX_C_API obx_err obx_admin_opt_bind(OBX_admin_options* opt, const char* uri);

/// Configure the server to use SSL, with the given certificate.
/// @param cert_path - the file must be in PEM format, and it must have both private key and certificate (public key).
OBX_C_API obx_err obx_admin_opt_ssl_cert(OBX_admin_options* opt, const char* cert_path);

/// Sets the number of worker threads the http-server uses to serve requests. Default: 4
OBX_C_API obx_err obx_admin_opt_num_threads(OBX_admin_options* opt, size_t num_threads);

/// Disables authentication to make the web app accessible to anyone, e.g. to allow recovery if users have locked
/// themselves out. Otherwise, if users have been defined, authentication is enforced.
OBX_C_API obx_err obx_admin_opt_unsecured_no_authentication(OBX_admin_options* opt, bool value);

/// Disable user management in the database - this is usually only done during development, with local-only access.
OBX_C_API obx_err obx_admin_opt_user_management(OBX_admin_options* opt, bool value);

/// Logs request info, e.g. timing for serving a request
OBX_C_API obx_err obx_admin_opt_log_requests(OBX_admin_options* opt, bool value);

/// Free the options.
/// Note: Only free *unused* options, obx_admin() frees the options internally
OBX_C_API obx_err obx_admin_opt_free(OBX_admin_options* opt);

struct OBX_admin;  // doxygen (only) picks up the typedef struct below

/// The admin web UI is an optional feature consisting of an embedded http server, that comes with an object (data)
/// browser and an administration UI.
typedef struct OBX_admin OBX_admin;

/// Initialize the http-server with the given options.
/// Note: the given options are always freed by this function, including when an error occurs.
/// @param options required parameter holding the options (see obx_admin_opt_*())
/// @returns NULL if the operation failed, see functions like obx_last_error_code() to get error details
OBX_C_API OBX_admin* obx_admin(OBX_admin_options* options);

/// Returns a port this server listens on. This is especially useful if the port was assigned arbitrarily
/// (a "0" port was used in the URL given to obx_admin_opt_bind()).
OBX_C_API uint16_t obx_admin_port(OBX_admin* admin);

/// Stop the http-server and free all the resources.
/// @param admin may be NULL
OBX_C_API obx_err obx_admin_close(OBX_admin* admin);

// NOLINTEND(modernize-use-using)

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_H

/**@}*/  // end of doxygen group