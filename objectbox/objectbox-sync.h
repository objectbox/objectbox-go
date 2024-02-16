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
 * @defgroup c-sync ObjectBox Sync C API
 * @{
 */

// Single header file for the ObjectBox C Sync API.
// Check https://objectbox.io/sync/ for ObjectBox Sync.
//
// Naming conventions
// ------------------
// * methods: obx_sync_thing_action()
// * structs: OBX_sync_thing {}
//

#ifndef OBJECTBOX_SYNC_H
#define OBJECTBOX_SYNC_H

#include "objectbox.h"

#if defined(static_assert) || defined(__cplusplus)
static_assert(OBX_VERSION_MAJOR == 0 && OBX_VERSION_MINOR == 21 && OBX_VERSION_PATCH == 0,  // NOLINT
              "Versions of objectbox.h and objectbox-sync.h files do not match, please update");
#endif

#ifdef __cplusplus
extern "C" {
#endif

// NOLINTBEGIN(modernize-use-using)

//----------------------------------------------
// Sync (client)
//----------------------------------------------

struct OBX_sync;
typedef struct OBX_sync OBX_sync;

/// Specifies user-side credential types as well as server-side authenticator types.
/// Some credentail types do not make sense as authenticators such as OBXSyncCredentialsType_USER_PASSWORD which
/// specifies a generic client-side credential type.
typedef enum {
    OBXSyncCredentialsType_NONE = 1,
    OBXSyncCredentialsType_SHARED_SECRET = 2,
    OBXSyncCredentialsType_GOOGLE_AUTH = 3,
    OBXSyncCredentialsType_SHARED_SECRET_SIPPED = 4,
    OBXSyncCredentialsType_OBX_ADMIN_USER = 5,
    OBXSyncCredentialsType_USER_PASSWORD = 6,
} OBXSyncCredentialsType;

// TODO sync prefix
typedef enum {
    /// no updates by default, obx_sync_updates_request() must be called manually
    OBXRequestUpdatesMode_MANUAL = 0,

    /// same as calling obx_sync_updates_request(sync, TRUE)
    /// default mode unless overridden by obx_sync_request_updates_mode
    OBXRequestUpdatesMode_AUTO = 1,

    /// same as calling obx_sync_updates_request(sync, FALSE)
    OBXRequestUpdatesMode_AUTO_NO_PUSHES = 2
} OBXRequestUpdatesMode;

typedef enum {
    OBXSyncState_CREATED = 1,
    OBXSyncState_STARTED = 2,
    OBXSyncState_CONNECTED = 3,
    OBXSyncState_LOGGED_IN = 4,
    OBXSyncState_DISCONNECTED = 5,
    OBXSyncState_STOPPED = 6,
    OBXSyncState_DEAD = 7
} OBXSyncState;

typedef enum {
    OBXSyncCode_OK = 20,
    OBXSyncCode_REQ_REJECTED = 40,
    OBXSyncCode_CREDENTIALS_REJECTED = 43,
    OBXSyncCode_UNKNOWN = 50,
    OBXSyncCode_AUTH_UNREACHABLE = 53,
    OBXSyncCode_BAD_VERSION = 55,
    OBXSyncCode_CLIENT_ID_TAKEN = 61,
    OBXSyncCode_TX_VIOLATED_UNIQUE = 71,
} OBXSyncCode;

/// Sync-level error reporting codes, passed via obx_sync_listener_error().
typedef enum {
    /// Sync client received rejection of transaction writes due to missing permissions.
    /// Until reconnecting with new credentials client will run in receive-only mode.
    OBXSyncError_REJECT_TX_NO_PERMISSION = 1
} OBXSyncError;

typedef enum {
    OBXSyncObjectType_FlatBuffers = 1,
    OBXSyncObjectType_String = 2,
    OBXSyncObjectType_Raw = 3,
} OBXSyncObjectType;

typedef struct OBX_sync_change {
    obx_schema_id entity_id;
    const OBX_id_array* puts;
    const OBX_id_array* removals;
} OBX_sync_change;

typedef struct OBX_sync_change_array {
    const OBX_sync_change* list;
    size_t count;
} OBX_sync_change_array;

/// A single data object contained in a OBX_sync_msg_objects message.
typedef struct OBX_sync_object {
    OBXSyncObjectType type;
    uint64_t id;       ///< optional value that the application can use identify the object (may be zero)
    const void* data;  ///< Pointer to object data, which is to be interpreted according to its type
    size_t size;       ///< Size of the object data (including the trailing \0 in case of OBXSyncObjectType_String)
} OBX_sync_object;

/// Incubating message that carries multiple data "objects" (e.g. FlatBuffers, strings, raw bytes).
/// Interpretation is up to the application. Does not involve any persistence or delivery guarantees at the moment.
typedef struct OBX_sync_msg_objects {
    const void* topic;
    size_t topic_size;  ///< topic is usually a string, but could also be binary (up to the application)
    const OBX_sync_object* objects;
    size_t count;
} OBX_sync_msg_objects;

/// An outgoing sync objects-message.
struct OBX_sync_msg_objects_builder;
typedef struct OBX_sync_msg_objects_builder OBX_sync_msg_objects_builder;

/// Called when connection is established
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_connect(void* arg);

/// Called when connection is closed/lost
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_disconnect(void* arg);

/// Called on successful login
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_login(void* arg);

/// Called on a login failure
/// @param arg is a pass-through argument passed to the called API
/// @param code error code indicating why the login failed
typedef void OBX_sync_listener_login_failure(void* arg, OBXSyncCode code);

/// Called when synchronization is complete
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_complete(void* arg);

/// Callend when sync-level errors occur
/// @param arg is a pass-through argument passed to the called API
/// @param error error code indicating sync-level error events
typedef void OBX_sync_listener_error(void* arg, OBXSyncError error);

/// Called with fine grained sync changes (IDs of put and removed entities)
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_change(void* arg, const OBX_sync_change_array* changes);

/// Called when a server time information is received on the client.
/// @param arg is a pass-through argument passed to the called API
/// @param timestamp_ns is timestamp in nanoseconds since Unix epoch
typedef void OBX_sync_listener_server_time(void* arg, int64_t timestamp_ns);

typedef void OBX_sync_listener_msg_objects(void* arg, const OBX_sync_msg_objects* msg_objects);

/// Creates a sync client associated with the given store and sync server URL.
/// This does not initiate any connection attempts yet: call obx_sync_start() to do so.
/// Before obx_sync_start(), you must configure credentials via obx_sync_credentials.
/// By default a sync client automatically receives updates from the server once login succeeded.
/// To configure this differently, call obx_sync_request_updates_mode() with the wanted mode.
OBX_C_API OBX_sync* obx_sync(OBX_store* store, const char* server_url);

/// Creates a sync client associated with the given store and a list of sync server URL.
/// For details, see obx_sync()
OBX_C_API OBX_sync* obx_sync_urls(OBX_store* store, const char* server_urls[], size_t server_urls_count);

/// Stops and closes (deletes) the sync client, freeing its resources.
OBX_C_API obx_err obx_sync_close(OBX_sync* sync);

/// Sets credentials to authenticate the client with the server.
/// See OBXSyncCredentialsType for available options.
/// The accepted OBXSyncCredentials type depends on your sync-server configuration.
/// @param data may be NULL in combination with OBXSyncCredentialsType_NONE
OBX_C_API obx_err obx_sync_credentials(OBX_sync* sync, OBXSyncCredentialsType type, const void* data, size_t size);

/// Set username/password credentials to authenticate the client with the server.
/// This is suitable for OBXSyncCredentialsType_OBX_ADMIN_USER and OBXSyncCredentialsType_USER_PASSWORD credential
/// types. Use obx_sync_credentials() for other credential types.
/// @param type should be OBXSyncCredentialsType_OBX_ADMIN_USER or OBXSyncCredentialsType_USER_PASSWORD
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if credential type does not support username/password authentication.
OBX_C_API obx_err obx_sync_credentials_user_password(OBX_sync* sync, OBXSyncCredentialsType type, const char* username,
                                                     const char* password);

/// Configures the maximum number of outgoing TX messages that can be sent without an ACK from the server.
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if value is not in the range 1-20
OBX_C_API obx_err obx_sync_max_messages_in_flight(OBX_sync* sync, int value);

/// Triggers a reconnection attempt immediately.
/// By default, an increasing backoff interval is used for reconnection attempts.
/// But sometimes the user of this API has additional knowledge and can initiate a reconnection attempt sooner.
OBX_C_API obx_err obx_sync_trigger_reconnect(OBX_sync* sync);

/// Sets the interval in which the client sends "heartbeat" messages to the server, keeping the connection alive.
/// To detect disconnects early on the client side, you can also use heartbeats with a smaller interval.
/// Use with caution, setting a low value (i.e. sending heartbeat very often) may cause an excessive network usage
/// as well as high server load (with many clients).
/// @param interval_ms interval in milliseconds; the default is 25 minutes (1 500 000 milliseconds),
///        which is also the allowed maximum.
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if value is not in the allowed range, e.g. larger than the maximum (1 500 000).
OBX_C_API obx_err obx_sync_heartbeat_interval(OBX_sync* sync, uint64_t interval_ms);

/// Triggers the heartbeat sending immediately. This lets you check the network connection at any time.
/// @see obx_sync_heartbeat_interval()
OBX_C_API obx_err obx_sync_send_heartbeat(OBX_sync* sync);

/// Switches operation mode that's initialized after successful login
/// Must be called before obx_sync_start() (returns OBX_ERROR_ILLEGAL_STATE if it was already started)
OBX_C_API obx_err obx_sync_request_updates_mode(OBX_sync* sync, OBXRequestUpdatesMode mode);

// TODO add along with the flag values and obx_sync_server_flags()
// Can be called before or after obx_sync_start()
// OBX_C_API obx_err obx_sync_flags(OBX_sync* sync, uint32_t flags);

/// Once the sync client is configured, you can "start" it to initiate synchronization.
/// This method triggers communication in the background and will return immediately.
/// If the synchronization destination is reachable, this background thread will connect to the server,
/// log in (authenticate) and, depending on "update request mode", start syncing data.
/// If the device, network or server is currently offline, connection attempts will be retried later using
/// increasing backoff intervals.
OBX_C_API obx_err obx_sync_start(OBX_sync* sync);

/// Stops this sync client and thus stopping the synchronization. Does nothing if it is already stopped.
OBX_C_API obx_err obx_sync_stop(OBX_sync* sync);

/// Gets the current state of the sync client (0 on error, e.g. given sync was NULL)
OBX_C_API OBXSyncState obx_sync_state(OBX_sync* sync);

/// Waits for the sync client to get into the given state or until the given timeout is reached.
/// For an asynchronous alternative, please check the listeners.
/// @param timeout_millis Must be greater than 0
/// @returns OBX_SUCCESS if the LOGGED_IN state has been reached within the given timeout
/// @returns OBX_TIMEOUT if the given timeout was reached before a relevant state change was detected.
///          Note: obx_last_error_code() is not set.
/// @returns OBX_NO_SUCCESS if a state was reached within the given timeout that is unlikely to result in a
///          successful login, e.g. "disconnected". Note: obx_last_error_code() is not set.
OBX_C_API obx_err obx_sync_wait_for_logged_in_state(OBX_sync* sync, uint64_t timeout_millis);

/// Request updates from the server since we last synced our database.
/// @param subscribe_for_pushes keep sending us future updates as they come in.
/// This should only be called in "logged in" state and there are no delivery guarantees given.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
OBX_C_API obx_err obx_sync_updates_request(OBX_sync* sync, bool subscribe_for_pushes);

/// Cancel updates from the server (once received, the server stops sending updates).
/// The counterpart to obx_sync_updates_request().
/// This should only be called in "logged in" state and there are no delivery guarantees given.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
OBX_C_API obx_err obx_sync_updates_cancel(OBX_sync* sync);

/// Count the number of messages in the outgoing queue, i.e. those waiting to be sent to the server.
/// @param limit pass 0 to count all messages without any limitation or a lower number that's enough for your app logic.
/// @note This calls uses a (read) transaction internally: 1) it's not just a "cheap" return of a single number.
///       While this will still be fast, avoid calling this function excessively.
///       2) the result follows transaction view semantics, thus it may not always match the actual value.
OBX_C_API obx_err obx_sync_outgoing_message_count(OBX_sync* sync, uint64_t limit, uint64_t* out_count);

/// Experimental. This API is likely to be replaced/removed in a future version.
/// Quickly bring our database up-to-date in a single transaction, without transmitting all the history.
/// Good for initial sync of new clients.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
OBX_C_API obx_err obx_sync_full(OBX_sync* sync);

/// Estimates the current server timestamp based on the last known server time and local steady clock.
/// @param out_timestamp_ns - unix timestamp in nanoseconds - may be set to zero if the last server's time is unknown.
OBX_C_API obx_err obx_sync_time_server(OBX_sync* sync, int64_t* out_timestamp_ns);

/// Returns the estimated difference between the server time and the local timestamp.
/// Equivalent to calculating obx_sync_time_server() - "current time" (nanos since epoch).
/// @param out_diff_ns time difference in nanoseconds; e.g. positive if server time is ahead of local time.
///                    Set to 0 if there has not been a server contact yet and thus the server's time is unknown.
OBX_C_API obx_err obx_sync_time_server_diff(OBX_sync* sync, int64_t* out_diff_ns);

/// Gets the protocol version this client uses.
OBX_C_API uint32_t obx_sync_protocol_version();

/// Gets the protocol version of the server after a connection is established (or attempted), zero otherwise.
OBX_C_API uint32_t obx_sync_protocol_version_server(OBX_sync* sync);

/// Start here to prepare an 'objects message'.
/// Use obx_sync_msg_objects_builder_add() to set at least one object and then call obx_sync_send_msg_objects() or
/// obx_sync_server_send_msg_objects() to initiate the sending process.
/// @param topic optional, application-specific message qualifier (may be NULL), usually a string but can also be binary
OBX_C_API OBX_sync_msg_objects_builder* obx_sync_msg_objects_builder(const void* topic, size_t topic_size);

/// Adds an object to the given message (builder). There must be at least one message before sending.
/// @param id an optional (pass 0 if you don't need it) value that the application can use identify the object
OBX_C_API obx_err obx_sync_msg_objects_builder_add(OBX_sync_msg_objects_builder* message, OBXSyncObjectType type,
                                                   const void* data, size_t size, uint64_t id);

/// Free the given message if you end up not sending it. Sending frees it already so never call this after obx_*_send().
OBX_C_API obx_err obx_sync_msg_objects_builder_discard(OBX_sync_msg_objects_builder* message);

/// Sends the given 'objects message' from the client to the currently connected server.
/// @param message the prepared outgoing message; it will be freed along with any associated resources during this call
///        (regardless of the call's success/failure outcome).
/// @returns OBX_SUCCESS if the message was scheduled to be sent (no guarantees for actual sending/transmission given).
/// @returns OBX_NO_SUCCESS if the message was not sent (no error will be set).
/// @returns error code if an unexpected error occurred.
OBX_C_API obx_err obx_sync_send_msg_objects(OBX_sync* sync, OBX_sync_msg_objects_builder* message);

/// Set or overwrite a previously set 'connect' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_connect(OBX_sync* sync, OBX_sync_listener_connect* listener, void* listener_arg);

/// Set or overwrite a previously set 'disconnect' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_disconnect(OBX_sync* sync, OBX_sync_listener_disconnect* listener, void* listener_arg);

/// Set or overwrite a previously set 'login' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_login(OBX_sync* sync, OBX_sync_listener_login* listener, void* listener_arg);

/// Set or overwrite a previously set 'login failure' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_login_failure(OBX_sync* sync, OBX_sync_listener_login_failure* listener,
                                               void* listener_arg);

/// Set or overwrite a previously set 'complete' listener - notifies when the latest sync has finished.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_complete(OBX_sync* sync, OBX_sync_listener_complete* listener, void* listener_arg);

/// Set or overwrite a previously set 'change' listener - provides information about incoming changes.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_change(OBX_sync* sync, OBX_sync_listener_change* listener, void* listener_arg);

/// Set or overwrite a previously set 'serverTime' listener - provides current time updates from the sync-server.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_server_time(OBX_sync* sync, OBX_sync_listener_server_time* listener,
                                             void* listener_arg);

/// Set or overwrite a previously set 'objects message' listener to receive application specific data objects.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_msg_objects(OBX_sync* sync, OBX_sync_listener_msg_objects* listener,
                                             void* listener_arg);

/// Set or overwrite a previously set 'error' listener - provides information about occurred sync-level errors.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API void obx_sync_listener_error(OBX_sync* sync, OBX_sync_listener_error* error, void* listener_arg);

//----------------------------------------------
// Sync Stats
//----------------------------------------------

/// Stats counter type IDs as passed to obx_sync_stats_u64().
typedef enum {
    /// Total number of connects (u64)
    OBXSyncStats_connects = 1,

    /// Total number of succesful logins (u64)
    OBXSyncStats_logins = 2,

    /// Total number of messages received (u64)
    OBXSyncStats_messagesReceived = 3,

    /// Total number of messages sent (u64)
    OBXSyncStats_messagesSent = 4,

    /// Total number of errors during message sending (u64)
    OBXSyncStats_messageSendFailures = 5,

    /// Total number of bytes received via messages.
    /// Note: this is measured on the application level and thus may not match e.g. the network level. (u64)
    OBXSyncStats_messageBytesReceived = 6,

    /// Total number of bytes sent via messages.
    /// Note: this is measured on the application level and thus may not match e.g. the network level.
    /// E.g. messages may be still enqueued so at least the timing will differ. (u64)
    OBXSyncStats_messageBytesSent = 7,

} OBXSyncStats;

/// Get u64 value for sync statistics.
/// @param counter_type the counter value to be read.
/// @param out_count receives the counter value.
/// @return OBX_SUCCESS if the counter has been successfully retrieved.
/// @return OBX_ERROR_ILLEGAL_ARGUMENT if counter_type is undefined.
OBX_C_API obx_err obx_sync_stats_u64(OBX_sync* sync, OBXSyncStats counter_type, uint64_t* out_count);

struct OBX_sync_server;
typedef struct OBX_sync_server OBX_sync_server;

/// Prepares an ObjectBox Sync Server to run within your application (embedded server) at the given URI.
/// Note that you need a special sync edition, which includes the server components. Check https://objectbox.io/sync/.
/// This call opens a store with the given options (also see obx_store_open()).
/// The server's store is tied to the server itself and is closed when the server is closed.
/// Before actually starting the server via obx_sync_server_start(), you can configure:
/// - accepted credentials via obx_sync_server_credentials() (always required)
/// - SSL certificate info via obx_sync_server_certificate_path() (required if you use wss)
/// \note The model given via store_options is also used to verify the compatibility of the models presented by clients.
///       E.g. a client with an incompatible model will be rejected during login.
/// @param store_options Options for the server's store.
///        It is freed automatically (same as with obx_store_open()) - don't use or free it afterwards.
/// @param url The URL (following the pattern "protocol://IP:port") the server should listen on.
///        Supported \b protocols are "ws" (WebSockets) and "wss" (secure WebSockets).
///        To use the latter ("wss"), you must also call obx_sync_server_certificate_path().
///        To bind to all available \b interfaces, including those that are available from the "outside", use 0.0.0.0 as
///        the IP. On the other hand, "127.0.0.1" is typically (may be OS dependent) only available on the same device.
///        If you do not require a fixed \b port, use 0 (zero) as a port to tell the server to pick an arbitrary port
///        that is available. The port can be queried via obx_sync_server_port() once the server was started.
///        \b Examples: "ws://0.0.0.0:9999" could be used during development (WS no certificate config needed),
///        while in a production system, you may want to use WSS and a specific IP for security reasons.
/// @see obx_store_open()
/// @returns NULL if server could not be created (e.g. the store could not be opened, bad URL, etc.)
OBX_C_API OBX_sync_server* obx_sync_server(OBX_store_options* store_options, const char* url);

/// Stops and closes (deletes) the sync server, freeing its resources.
/// This includes the store associated with the server; it gets closed and must not be used anymore after this call.
OBX_C_API obx_err obx_sync_server_close(OBX_sync_server* server);

/// Gets the store this server uses. This is owned by the server and must NOT be closed manually with obx_store_close().
OBX_C_API OBX_store* obx_sync_server_store(OBX_sync_server* server);

/// Sets SSL certificate for the server to use. Use before obx_sync_server_start().
OBX_C_API obx_err obx_sync_server_certificate_path(OBX_sync_server* server, const char* certificate_path);

/// Sets credentials for the server to accept. Use before obx_sync_server_start().
/// @param data may be NULL in combination with OBXSyncCredentialsType_NONE
OBX_C_API obx_err obx_sync_server_credentials(OBX_sync_server* server, OBXSyncCredentialsType type, const void* data,
                                              size_t size);

/// Enables authenticator for server. Can be called multiple times. Use before obx_sync_server_start().
/// Use obx_sync_server_credentials() for authenticators which requires additional credentials data (i.e. Google Auth
/// and shared secrets authenticators).
/// @param type should be one of the available authentications, it should not be OBXSyncCredentialsType_USER_PASSWORD.
OBX_C_API obx_err obx_sync_server_enable_auth(OBX_sync_server* server, OBXSyncCredentialsType type);

/// Sets the number of worker threads. Use before obx_sync_server_start().
/// @param thread_count The default is "0" which is hardware dependent, e.g. a multiple of CPU "cores".
OBX_C_API obx_err obx_sync_server_worker_threads(OBX_sync_server* server, int thread_count);

/// Sets a maximum size for sync history entries to limit storage: old entries are removed to stay below this limit.
/// Deleting older history entries may require clients to do a full sync if they have not contacted the server for
/// a certain time.
/// @param max_in_kb Once this maximum size is reached, old history entries are deleted (default 0: no limit).
/// @param target_in_kb If this value is non-zero, the deletion of old history entries is extended until reaching this
///                     target (lower than the maximum) allowing deletion "batching", which may be more efficient.
///                     If zero, the deletion stops already stops when reaching the max size (or lower).
OBX_C_API obx_err obx_sync_server_history_max_size_in_kb(OBX_sync_server* server, uint64_t max_in_kb,
                                                         uint64_t target_in_kb);

/// Set or overwrite a previously set 'change' listener - provides information about incoming changes.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API obx_err obx_sync_server_listener_change(OBX_sync_server* server, OBX_sync_listener_change* listener,
                                                  void* listener_arg);

/// Set or overwrite a previously set 'objects message' listener to receive application specific data objects.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
OBX_C_API obx_err obx_sync_server_listener_msg_objects(OBX_sync_server* server, OBX_sync_listener_msg_objects* listener,
                                                       void* listener_arg);

/// After the sync server is fully configured (e.g. credentials), this will actually start the server.
/// Once this call returns, the server is ready to accept client connections. Also, port and URL will be available.
OBX_C_API obx_err obx_sync_server_start(OBX_sync_server* server);

/// Stops this sync server. Does nothing if it is already stopped.
OBX_C_API obx_err obx_sync_server_stop(OBX_sync_server* server);

/// Whether the server is up and running.
OBX_C_API bool obx_sync_server_running(OBX_sync_server* server);

/// Returns a URL this server is listening on, including the bound port (see obx_sync_server_port().
/// The returned char* is valid until another call to obx_sync_server_url() or the server is closed.
OBX_C_API const char* obx_sync_server_url(OBX_sync_server* server);

/// Returns a port this server listens on. This is especially useful if the port was assigned arbitrarily
/// (a "0" port was used in the URL given to obx_sync_server()).
OBX_C_API uint16_t obx_sync_server_port(OBX_sync_server* server);

/// Returns the number of clients connected to this server.
OBX_C_API uint64_t obx_sync_server_connections(OBX_sync_server* server);

//----------------------------------------------
// Sync Server Stats
//----------------------------------------------

/// Stats counter type IDs as passed to obx_sync_server_stats_u64() (for u64 values) and obx_sync_server_stats_f64()
/// (for double (f64) values).
typedef enum {
    /// Total number of client connections established (u64)
    OBXSyncServerStats_connects = 1,

    /// Total number of messages received from clients (u64)
    OBXSyncServerStats_messagesReceived = 2,

    /// Total number of messages sent to clients (u64)
    OBXSyncServerStats_messagesSent = 3,

    /// Total number of bytes received from clients via messages. (u64)
    /// Note: this is measured on the application level and thus may not match e.g. the network level.
    OBXSyncServerStats_messageBytesReceived = 4,

    /// Total number of bytes sent to clients via messages. (u64)
    /// Note: this is measured on the application level and thus may not match e.g. the network level.
    /// E.g. messages may be still enqueued so at least the timing will differ.
    OBXSyncServerStats_messageBytesSent = 5,

    /// Full syncs performed (u64)
    OBXSyncServerStats_fullSyncs = 6,

    /// Processing was aborted due to clients disconnected (u64)
    OBXSyncServerStats_disconnectAborts = 7,

    /// Total number of client transactions applied (u64)
    OBXSyncServerStats_clientTxsApplied = 8,

    /// Total size in bytes of client transactions applied (u64)
    OBXSyncServerStats_clientTxBytesApplied = 9,

    /// Total size in number of operations of transactions applied (u64)
    OBXSyncServerStats_clientTxOpsApplied = 10,

    /// Total number of local (server initiated) transactions applied (u64)
    OBXSyncServerStats_localTxsApplied = 11,

    /// AsyncQ committed TXs (u64)
    OBXSyncServerStats_asyncDbCommits = 12,

    /// Total number of skipped transactions duplicates (have been already applied before) (u64)
    OBXSyncServerStats_skippedTxDups = 13,

    /// Total number of login successes (u64)
    OBXSyncServerStats_loginSuccesses = 14,

    /// Total number of login failures (u64)
    OBXSyncServerStats_loginFailures = 15,

    /// Total number of login failures due to bad user credentials (u64)
    OBXSyncServerStats_loginFailuresUserBadCredentials = 16,

    /// Total number of login failures due to authenticator not available (u64)
    OBXSyncServerStats_loginFailuresAuthUnavailable = 17,

    /// Total number of login failures due to user has no permissions (u64)
    OBXSyncServerStats_loginFailuresUserNoPermission = 18,

    /// Total number of errors during message sending (u64)
    OBXSyncServerStats_messageSendFailures = 19,

    /// Total number of protocol errors; e.g. offending clients (u64)
    OBXSyncServerStats_errorsProtocol = 20,

    /// Total number of errors in message handlers (u64)
    OBXSyncServerStats_errorsInHandlers = 21,

    /// Total number of times a client has been disconnected due to heart failure (u64)
    OBXSyncServerStats_heartbeatFailures = 22,

    /// Total number of received client heartbeats (u64)
    OBXSyncServerStats_heartbeatsReceived = 23,

    /// Total APPLY_TX messages HistoryPusher has sent out (u64)
    OBXSyncServerStats_historicUpdateTxsSent = 24,

    /// Total APPLY_TX messages newDataPusher has sent out (u64)
    OBXSyncServerStats_newUpdateTxsSent = 25,

    /// Total number of messages received from clients (u64)
    OBXSyncServerStats_forwardedMessagesReceived = 26,

    /// Total number of messages sent to clients (u64)
    OBXSyncServerStats_forwardedMessagesSent = 27,

    /// Total number of global-to-local cache hits (u64)
    OBXSyncServerStats_cacheGlobalToLocalHits = 28,

    /// Total number of global-to-local cache misses (u64)
    OBXSyncServerStats_cacheGlobalToLocalMisses = 29,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cacheGlobalToLocalSize = 30,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cachePeerToLocalHits = 31,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cachePeerToLocalMisses = 32,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cacheLocalToPeerHits = 33,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cacheLocalToPeerMisses = 34,

    /// Internal dev stat for ID Map caching  (u64)
    OBXSyncServerStats_cachePeerSize = 35,

    /// Current cluster peer state (0 = unknown, 1 = leader, 2 = follower, 3 = candidate) (u64)
    OBXSyncServerStats_clusterPeerState = 36,

    /// Number of transactions between the current Tx and the oldest Tx currently ACKed on any client (current)
    /// (f64)
    OBXSyncServerStats_clientTxsBehind = 37,

    /// Number of transactions between the current Tx and the oldest Tx currently ACKed on any client (minimum)
    /// (u64)
    OBXSyncServerStats_clientTxsBehind_min = 38,

    /// Number of transactions between the current Tx and the oldest Tx currently ACKed on any client (maximum)
    /// (u64)
    OBXSyncServerStats_clientTxsBehind_max = 39,

    /// Number of connected clients (current) (f64)
    OBXSyncServerStats_connectedClients = 40,

    /// Number of connected clients (minimum) (u64)
    OBXSyncServerStats_connectedClients_min = 41,

    /// Number of connected clients (maximum) (u64)
    OBXSyncServerStats_connectedClients_max = 42,

    /// Length of the queue for regular Tasks (current) (f64)
    OBXSyncServerStats_queueLength = 43,

    /// Length of the queue for regular Tasks (minimum) (u64)
    OBXSyncServerStats_queueLength_min = 44,

    /// Length of the queue for regular Tasks (maximum) (u64)
    OBXSyncServerStats_queueLength_max = 45,

    /// Length of the async queue (current) (f64)
    OBXSyncServerStats_queueLengthAsync = 46,

    /// Length of the async queue (minimum) (u64)
    OBXSyncServerStats_queueLengthAsync_min = 47,

    /// Length of the async queue (maximum) (u64)
    OBXSyncServerStats_queueLengthAsync_max = 48,

    /// Sequence number of TX log history (current) (f64)
    OBXSyncServerStats_txHistorySequence = 49,

    /// Sequence number of TX log history (minimum) (u64)
    OBXSyncServerStats_txHistorySequence_min = 50,

    /// Sequence number of TX log history (maximum) (u64)
    OBXSyncServerStats_txHistorySequence_max = 51,

} OBXSyncServerStats;

/// Get u64 value for sync server statistics.
/// @param counter_type the counter value to be read (make sure to choose a uint64_t (u64) metric value type).
/// @param out_count receives the counter value.
/// @return OBX_SUCCESS if the counter has been successfully retrieved.
/// @return OBX_ERROR_ILLEGAL_ARGUMENT if counter_type is undefined (this also happens if the wrong type is requested)
/// @return OBX_ERROR_ILLEGAL_STATE if the server is not started.
OBX_C_API obx_err obx_sync_server_stats_u64(OBX_sync_server* server, OBXSyncServerStats counter_type,
                                            uint64_t* out_value);

/// Get double value for sync server statistics.
/// @param counter_type the counter value to be read (make sure to use a double (f64) metric value type).
/// @param out_count receives the counter value.
/// @return OBX_SUCCESS if the counter has been successfully retrieved.
/// @return OBX_ERROR_ILLEGAL_ARGUMENT if counter_type is undefined (this also happens if the wrong type is requested)
/// @return OBX_ERROR_ILLEGAL_STATE if the server is not started.
OBX_C_API obx_err obx_sync_server_stats_f64(OBX_sync_server* server, OBXSyncServerStats counter_type,
                                            double* out_value);

/// Get server runtime statistics.
/// The returned char* is valid until another call to obx_sync_server_stats_string() or the server is closed.
OBX_C_API const char* obx_sync_server_stats_string(OBX_sync_server* server, bool include_zero_values);

/// Broadcast the given 'objects message' from the server to all currently connected (and logged-in) clients.
/// @param message the prepared outgoing message; it will be freed along with any associated resources during this call
///        (regardless of the call's success/failure outcome).
/// @returns OBX_SUCCESS if the message was scheduled to be sent (no guarantees for actual sending/transmission given).
/// @returns OBX_NO_SUCCESS if the message was not sent (no error will be set).
/// @returns error code if an unexpected error occurred.
OBX_C_API obx_err obx_sync_server_send_msg_objects(OBX_sync_server* server, OBX_sync_msg_objects_builder* message);

/// Configure admin with a sync server, attaching the store and enabling custom sync-server functionality in the UI.
/// This is a replacement for obx_admin_opt_store() and obx_admin_opt_store_path() - don't set them for the server.
/// After configuring, this acts as obx_admin() - see for more details.
/// You must use obx_admin_close() to stop & free resources after you're done; obx_sync_server_stop() doesn't do that.
/// @param options configuration set up with obx_admin_opt_*. You can pass NULL to use the default options.
OBX_C_API OBX_admin* obx_sync_server_admin(OBX_sync_server* server, OBX_admin_options* options);

//---------------------------------------------------------------------------
// Custom messaging server
//---------------------------------------------------------------------------

/// Callback to create a custom messaging server.
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param server_id the ID that was assigned to the custom server instance
/// @param config_user_data user provided data set at registration of the server
/// @returns server user data, which will be passed on to the subsequent callbacks (OBX_custom_msg_server_func_*)
/// @returns null to indicate an error that the server could not be created
typedef void* OBX_custom_msg_server_func_create(uint64_t server_id, const char* url, const char* cert_path,
                                                void* config_user_data);

/// Callback to start a custom server.
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param server_user_data User supplied data returned by the function that created the server
/// @param out_port When starting, the custom server can optionally supply a port by writing to the given pointer.
///        The port value is arbitrary and, for now, is only used for debug logs.
/// @returns OBX_SUCCESS if the server was successfully started
/// @returns Any other fitting error code (OBX_ERROR_*) if the server could be started
typedef obx_err OBX_custom_msg_server_func_start(void* server_user_data, uint64_t* out_port);

/// Callback to stop and close the custom server (e.g. further messages delivery will be rejected).
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// This includes the store associated with the server; it gets closed and must not be used anymore after this call.
/// @param server_user_data User supplied data returned by the function that created the server
typedef void OBX_custom_msg_server_func_stop(void* server_user_data);

/// Callback to shut the custom server down, freeing its resources (the custom server is not used after this point).
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param server_user_data User supplied data returned by the function that created the server
typedef void OBX_custom_msg_server_func_shutdown(void* server_user_data);

/// Callback to enqueue a message for sending.
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param bytes lazy bytes storing the message
/// @param server_user_data User supplied data returned by the function that created the server
typedef bool OBX_custom_msg_server_func_client_connection_send_async(OBX_bytes_lazy* bytes, void* server_user_data,
                                                                     void* connection_user_data);

/// Callback to close the sync client connection to the custom server.
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param server_user_data User supplied data returned by the function that created the server
typedef void OBX_custom_msg_server_func_client_connection_close(void* server_user_data, void* connection_user_data);

/// Callback to shutdown and free all resources associated with the sync client connection to the custom server.
/// Note that the custom server may already have been shutdown at this point (e.g. no server user data is supplied).
/// Must be provided to implement a custom server. See notes on OBX_custom_msg_server_functions for more details.
/// @param server_user_data User supplied data returned by the function that created the server
typedef void OBX_custom_msg_server_func_client_connection_shutdown(void* connection_user_data);

/// Struct of the custom server function callbacks. In order to implement the custom server, you must provide
/// custom methods for each of the members of this struct. This is then passed to obx_custom_msg_server_register()
/// to register the custom server.
typedef struct OBX_custom_msg_server_functions {
    /// Must be initialized with sizeof(OBX_custom_msg_server_functions) to "version" the struct.
    /// This allows the library (whi) to detect older or newer versions and react properly.
    size_t version;

    OBX_custom_msg_server_func_create* func_create;
    OBX_custom_msg_server_func_start* func_start;
    OBX_custom_msg_server_func_stop* func_stop;
    OBX_custom_msg_server_func_shutdown* func_shutdown;

    OBX_custom_msg_server_func_client_connection_send_async* func_conn_send_async;
    OBX_custom_msg_server_func_client_connection_close* func_conn_close;
    OBX_custom_msg_server_func_client_connection_shutdown* func_conn_shutdown;
} OBX_custom_msg_server_functions;

/// Must be called to register a protocol for a custom messaging server. Call before starting a server.
/// @param protocol the communication protocol to use, e.g. "tcp"
/// @param functions the custom server function callbacks
/// @param config_user_data user provided data set at registration of custom server
/// @returns OBX_SUCCESS if the operation was successful
/// @returns Any other fitting error code (OBX_ERROR_*) if the protocol could not be registered
OBX_C_API obx_err obx_custom_msg_server_register(const char* protocol, OBX_custom_msg_server_functions* functions,
                                                 void* config_user_data);

/// Must be called from the custom server when a new client connection becomes available.
/// @param server_id the ID that was assigned to the custom server instance
/// @param user_data user provided data set at registration of custom server
/// @returns a client connection ID that must be passed on to obx_custom_msg_server_receive_message_from_client().
/// @returns 0 in case the operation encountered an exceptional issue
OBX_C_API uint64_t obx_custom_msg_server_add_client_connection(uint64_t server_id, void* user_data);

/// Must be called from the custom server when a client connection becomes inactive (e.g. closed) and can be removed.
/// @param server_id the ID that was assigned to the custom server instance
/// @param client_connection_id the ID that was assigned to the custom client connection
/// @returns OBX_SUCCESS if the operation was successful
/// @returns OBX_NO_SUCCESS if no active server or active connection was found matching the given IDs
/// @returns OBX_ERROR_* in case the operation encountered an exceptional issue
OBX_C_API obx_err obx_custom_msg_server_remove_client_connection(uint64_t server_id, uint64_t client_connection_id);

/// Must be called from the custom server when a message is received from a client connection.
/// @param server_id the ID that was assigned to the custom server instance
/// @param client_connection_id the ID that was assigned to the custom client connection
/// @param message_data the message data in bytes
/// @returns OBX_SUCCESS if the operation was successful
/// @returns OBX_NO_SUCCESS if no active server or active connection was found matching the given IDs
/// @returns OBX_ERROR_* in case the operation encountered an exceptional issue
OBX_C_API obx_err obx_custom_msg_server_receive_message_from_client(uint64_t server_id, uint64_t client_connection_id,
                                                                    const void* message_data, size_t message_size);

//---------------------------------------------------------------------------
//  Custom messaging client
//---------------------------------------------------------------------------

/// Callback to create a custom messaging client.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_id the ID that was assigned to the client instance
/// @param config_user_data user provided data set at registration of the client
/// @returns client user data, which will be passed on to the subsequent callbacks (OBX_custom_msg_client_func_*)
/// @returns null to indicate an error that the client could not be created
typedef void* OBX_custom_msg_client_func_create(uint64_t client_id, const char* url, const char* cert_path,
                                                void* config_user_data);

/// Callback to start the client.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
/// @returns OBX_SUCCESS if the client was successfully started
/// @returns Any other fitting error code (OBX_ERROR_*) if the client could be started
typedef obx_err OBX_custom_msg_client_func_start(void* client_user_data);

/// Callback to stop and close the client (e.g. further messages delivery will be rejected).
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef void OBX_custom_msg_client_func_stop(void* client_user_data);

/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef void OBX_custom_msg_client_func_join(void* client_user_data);

/// Callback that tells the client it shall start trying to connect.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef bool OBX_custom_msg_client_func_connect(void* client_user_data);

/// Callback that tells the client it shall disconnect.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef void OBX_custom_msg_client_func_disconnect(bool clear_outgoing_messages, void* client_user_data);

/// Callback to shut the custom client down, freeing its resources.
/// The custom client is not used after this point.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef void OBX_custom_msg_client_func_shutdown(void* client_user_data);

/// Callback to enqueue a message for sending.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param bytes lazy bytes storing the message
/// @param client_user_data user supplied data returned by the function that created the client
typedef bool OBX_custom_msg_client_func_send_async(OBX_bytes_lazy* bytes, void* client_user_data);

/// Callback to clear all outgoing messages.
/// Must be provided to implement a custom client. See notes on OBX_custom_msg_client_functions for more details.
/// @param client_user_data user supplied data returned by the function that created the client
typedef void OBX_custom_msg_client_func_clear_outgoing_messages(void* client_user_data);

/// Struct of the custom client function callbacks. In order to implement the custom client, you must provide
/// custom methods for each of the members of this struct. This is then passed to obx_custom_msg_client_register()
/// to register the custom client.
typedef struct OBX_custom_msg_client_functions {
    /// Must be initialized with sizeof(OBX_custom_msg_client_functions) to "version" the struct.
    /// This allows the library to detect older or newer versions and react properly.
    size_t version;

    OBX_custom_msg_client_func_create* func_create;
    OBX_custom_msg_client_func_start* func_start;
    OBX_custom_msg_client_func_connect* func_connect;
    OBX_custom_msg_client_func_disconnect* func_disconnect;
    OBX_custom_msg_client_func_stop* func_stop;
    OBX_custom_msg_client_func_join* func_join;
    OBX_custom_msg_client_func_shutdown* func_shutdown;
    OBX_custom_msg_client_func_send_async* func_send_async;
    OBX_custom_msg_client_func_clear_outgoing_messages* func_clear_outgoing_messages;
} OBX_custom_msg_client_functions;

/// States of custom msg client that must be forwarded to obx_custom_msg_client_set_state().
typedef enum {
    OBXCustomMsgClientState_Connecting = 1,
    OBXCustomMsgClientState_Connected = 2,
    OBXCustomMsgClientState_Disconnected = 3,
} OBXCustomMsgClientState;

/// Must be called to register a protocol for your custom messaging client. Call before starting a client.
/// @param protocol the communication protocol to use, e.g. "tcp"
/// @returns OBX_SUCCESS if the operation was successful
/// @returns Any other fitting error code (OBX_ERROR_*) if the protocol could not be registered
OBX_C_API obx_err obx_custom_msg_client_register(const char* protocol, OBX_custom_msg_client_functions* functions,
                                                 void* config_user_data);

/// The custom msg client must call this whenever a message is received from the server.
/// @param client_id the ID that was assigned to the client instance
/// @param message_data the message data in bytes
/// @returns OBX_SUCCESS if the given message could be forwarded
/// @returns OBX_NO_SUCCESS if no active client or active connection was found matching the given ID
/// @returns OBX_ERROR_* in case the operation encountered an exceptional issue
OBX_C_API obx_err obx_custom_msg_client_receive_message_from_server(uint64_t client_id, const void* message_data,
                                                                    size_t message_size);

/// The custom msg client must call this whenever the state (according to given enum values) changes.
/// @param client_id the ID that was assigned to the client instance
/// @param state the state to transition the custom client to
/// @returns OBX_SUCCESS if the client was in a state that allowed the transition to the given state.
/// @returns OBX_NO_SUCCESS if no active client or active connection was found matching the given ID.
/// @returns OBX_NO_SUCCESS if no state transition was possible from the current to the given state (e.g. an internal
///          "closed" state was reached).
/// @returns OBX_ERROR_* in case the operation encountered an exceptional issue
OBX_C_API obx_err obx_custom_msg_client_set_state(uint64_t client_id, OBXCustomMsgClientState state);

/// The custom msg client may call this if it has knowledge when a reconnection attempt makes sense,
/// for example, when the network becomes available.
/// @param client_id the ID that was assigned to the client instance
/// @returns OBX_SUCCESS if a reconnect was actually triggered.
/// @returns OBX_NO_SUCCESS if no reconnect was triggered.
/// @returns OBX_ERROR_* in case the operation encountered an exceptional issue
OBX_C_API obx_err obx_custom_msg_client_trigger_reconnect(uint64_t client_id);

// NOLINTEND(modernize-use-using)

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_SYNC_H

/**@}*/  // end of doxygen group