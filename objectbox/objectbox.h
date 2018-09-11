// Single header file for the ObjectBox C API
//
// Naming conventions
// ------------------
// * methods: ob_thing_action()
// * structs: OB_thing {}
// * error codes: OB_ERROR_REASON
// * enums: ?
//
#ifndef OBJECTBOX_H
#define OBJECTBOX_H

#include <stdint.h>
#include <stdio.h>

//----------------------------------------------
// Return codes
//----------------------------------------------

/// Value returned when no error occurred (0)
#define OB_SUCCESS 0

/// Returned by e.g. get operations if nothing was found for a specific ID.
/// This is NOT an error condition, and thus no last error info is set.
#define OB_NOT_FOUND 404

// General errors
#define OB_ERROR_ILLEGAL_STATE 10001
#define OB_ERROR_ILLEGAL_ARGUMENT 10002
#define OB_ERROR_ALLOCATION 10003
#define OB_ERROR_NO_ERROR_INFO 10097
#define OB_ERROR_GENERAL 10098
#define OB_ERROR_UNKNOWN 10099

// Storage errors (often have a secondary error code)
#define OB_ERROR_DB_FULL 10101
#define OB_ERROR_MAX_READERS_EXCEEDED 10102
#define OB_ERROR_STORE_MUST_SHUTDOWN 10103
#define OB_ERROR_STORAGE_GENERAL 10199

// Data errors
#define OB_ERROR_UNIQUE_VIOLATED 10201
#define OB_ERROR_NON_UNIQUE_RESULT 10202
#define OB_ERROR_CONSTRAINT_VIOLATED 10299

// STD errors
#define OB_ERROR_STD_ILLEGAL_ARGUMENT 10301
#define OB_ERROR_STD_OUT_OF_RANGE 10302
#define OB_ERROR_STD_LENGTH 10303
#define OB_ERROR_STD_BAD_ALLOC 10304
#define OB_ERROR_STD_RANGE 10305
#define OB_ERROR_STD_OVERFLOW 10306
#define OB_ERROR_STD_OTHER 10399

// Inconsistencies detected
#define OB_ERROR_SCHEMA 10501
#define OB_ERROR_FILE_CORRUPT 10502

//----------------------------------------------
// Error info
//----------------------------------------------

int ob_last_error_code();

const char* ob_last_error_message();

int ob_last_error_secondary();

void ob_last_error_clear();

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
    /// Unused yet, used by References for 1) back-references and 2) to clear references to deleted objects (required for ID reuse)
            PropertyFlags_INDEX_PARTIAL_SKIP_ZERO = 512,
    /// Virtual properties may not have a dedicated field in their entity class, e.g. target IDs of to-one relations
            PropertyFlags_VIRTUAL = 1024,
    /// Index uses a 32 bit hash instead of the value
    /// (32 bits is shorter on disk, runs well on 32 bit systems, and should be OK even with a few collisions)

    PropertyFlags_INDEX_HASH = 2048,
    /// Index uses a 64 bit hash instead of the value
    /// (recommended mostly for 64 bit machines with values longer >200 bytes; small values are faster with a 32 bit hash)
            PropertyFlags_INDEX_HASH64 = 4096

} OBPropertyFlags;

struct OB_model;
typedef struct OB_model OB_model;

OB_model* ob_model_create();

/// Only call when not calling ob_store_open (which will destroy it internally)
int ob_model_destroy(OB_model* model);

void ob_model_last_entity_id(OB_model*, uint32_t id, uint64_t uid);

void ob_model_last_index_id(OB_model* model, uint32_t id, uint64_t uid);

void ob_model_last_relation_id(OB_model* model, uint32_t id, uint64_t uid);

int ob_model_entity(OB_model* model, const char* name, uint32_t id, uint64_t uid);

int ob_model_entity_last_property_id(OB_model* model, uint32_t id, uint64_t uid);

int ob_model_property(OB_model* model, const char* name, OBPropertyType type, uint32_t id, uint64_t uid);

int ob_model_property_flags(OB_model* model, OBPropertyFlags flags);

int ob_model_property_index_id(OB_model* model, uint32_t id, uint64_t uid);

//----------------------------------------------
// Store
//----------------------------------------------

struct OB_store;
typedef struct OB_store OB_store;

struct OB_store_options {
    /// Use NULL for default value ("objectbox")
    char* directory;

    /// Use zero for default value
    uint64_t maxDbSizeInKByte;

    /// Use zero for default value
    unsigned int fileMode;

    /// Use zero for default value
    unsigned int maxReaders;
};

typedef struct OB_store_options OB_store_options;

enum DebugFlags {
    DebugFlags_LOG_TRANSACTIONS_READ = 1,
    DebugFlags_LOG_TRANSACTIONS_WRITE = 2,
    DebugFlags_LOG_QUERIES = 4,
    DebugFlags_LOG_QUERY_PARAMETERS = 8,
    DebugFlags_LOG_ASYNC_QUEUE = 16,
};

struct OB_bytes {
    void* data;
    size_t size;
};
typedef struct OB_bytes OB_bytes;

struct OB_bytes_array {
    OB_bytes* bytes;
    size_t size;
};
typedef struct OB_bytes_array OB_bytes_array;

struct OB_table_array {
    void* tables;
    size_t size;
};
typedef struct OB_table_array OB_table_array;


OB_store* ob_store_open_bytes(const void* modelBytes, size_t modelSize, const OB_store_options* options);

/// Note: the model is destroyed by calling this method
OB_store* ob_store_open(OB_model* model, const OB_store_options* options);

int ob_store_await_async_completion(OB_store* store);

int ob_store_debug_flags(OB_store* store, uint32_t debugFlags);

int ob_store_close(OB_store* store);

//----------------------------------------------
// Transaction
//----------------------------------------------

struct OB_txn;
typedef struct OB_txn OB_txn;

OB_txn* ob_txn_begin(OB_store* store);

OB_txn* ob_txn_begin_read(OB_store* store);

int ob_txn_destroy(OB_txn* txn);

int ob_txn_abort(OB_txn* txn);

int ob_txn_commit(OB_txn* txn);

//----------------------------------------------
// Cursor
//----------------------------------------------

struct OB_cursor;
typedef struct OB_cursor OB_cursor;

OB_cursor* ob_cursor_create(OB_txn* txn, uint32_t schemaEntityId);

OB_cursor* ob_cursor_create2(OB_txn* txn, const char* schemaEntityName);

int ob_cursor_destroy(OB_cursor* cursor);

uint64_t ob_cursor_id_for_put(OB_cursor* cursor, uint64_t idOrZero);

int ob_cursor_put(OB_cursor* cursor, uint64_t entityId, const void* data, size_t size, int checkForPreviousValueFlag);

int ob_cursor_get(OB_cursor* cursor, uint64_t entityId, void** data, size_t* size);

int ob_cursor_first(OB_cursor* cursor, void** data, size_t* size);

int ob_cursor_next(OB_cursor* cursor, void** data, size_t* size);

int ob_cursor_remove(OB_cursor* cursor, uint64_t entityId);

int ob_cursor_remove_all(OB_cursor* cursor);

int ob_cursor_count(OB_cursor* cursor, uint64_t* outCount);

//----------------------------------------------
// Box
//----------------------------------------------

struct OB_box;
typedef struct OB_box OB_box;

OB_box* ob_box_create(OB_store* store, uint32_t schemaEntityId);

int ob_box_destroy(OB_box* box);

uint64_t ob_box_id_for_put(OB_box* box, uint64_t idOrZero);

int ob_box_put_async(OB_box* box, uint64_t entityId, const void* data, size_t size, int checkForPreviousValueFlag);

//----------------------------------------------
// Query
//----------------------------------------------

OB_table_array* ob_simple_query_string(OB_cursor* cursor, uint32_t propertyId, const char* value, uint32_t valueSize);

OB_bytes_array* ob_query_by_string(OB_cursor* cursorStruct, uint32_t propertyId, const char* value);

void ob_bytes_destroy(OB_bytes* bytes);

void ob_bytes_array_destroy(OB_bytes_array* bytesArray);

void ob_table_array_destroy(OB_table_array* tableArray);

#endif //OBJECTBOX_H
