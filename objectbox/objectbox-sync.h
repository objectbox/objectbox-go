/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
static_assert(OBX_VERSION_MAJOR == 0 && OBX_VERSION_MINOR == 15 && OBX_VERSION_PATCH == 0,
              "Versions of objectbox.h and objectbox-sync.h files do not match, please update");
#endif

#ifdef __cplusplus
extern "C" {
#endif

//----------------------------------------------
// Sync (client)
//----------------------------------------------

struct OBX_sync;
typedef struct OBX_sync OBX_sync;

typedef enum {
    OBXSyncCredentialsType_NONE = 1,
    OBXSyncCredentialsType_SHARED_SECRET = 2,
    OBXSyncCredentialsType_GOOGLE_AUTH = 3,
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

/// Called with fine grained sync changes (IDs of put and removed entities)
/// @param arg is a pass-through argument passed to the called API
typedef void OBX_sync_listener_change(void* arg, const OBX_sync_change_array* changes);

/// Called when a server time information is received on the client.
/// @param arg is a pass-through argument passed to the called API
/// @param timestamp_ns is timestamp in nanoseconds since Unix epoch
typedef void OBX_sync_listener_server_time(void* arg, int64_t timestamp_ns);

typedef void OBX_sync_listener_msg_objects(void* arg, const OBX_sync_msg_objects* msg_objects);

/// Creates a sync client associated with the given store and sync server URI.
/// This does not initiate any connection attempts yet: call obx_sync_start() to do so.
/// Before obx_sync_start(), you must configure credentials via obx_sync_credentials.
/// By default a sync client automatically receives updates from the server once login succeeded.
/// To configure this differently, call obx_sync_request_updates_mode() with the wanted mode.
OBX_C_API OBX_sync* obx_sync(OBX_store* store, const char* server_uri);

/// Stops and closes (deletes) the sync client, freeing its resources.
OBX_C_API obx_err obx_sync_close(OBX_sync* sync);

/// Sets credentials to authenticate the client with the server.
/// See OBXSyncCredentialsType for available options.
/// The accepted OBXSyncCredentials type depends on your sync-server configuration.
/// @param data may be NULL in combination with OBXSyncCredentialsType_NONE
OBX_C_API obx_err obx_sync_credentials(OBX_sync* sync, OBXSyncCredentialsType type, const void* data, size_t size);

/// Configures the maximum number of outgoing TX messages that can be sent without an ACK from the server.
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if value is not in the range 1-20
OBX_C_API obx_err obx_sync_max_messages_in_flight(OBX_sync* sync, int value);

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
/// @param uri The URI (following the pattern protocol:://IP:port) the server should listen on.
///        Supported \b protocols are "ws" (WebSockets) and "wss" (secure WebSockets).
///        To use the latter ("wss"), you must also call obx_sync_server_certificate_path().
///        To bind to all available \b interfaces, including those that are available from the "outside", use 0.0.0.0 as
///        the IP. On the other hand, "127.0.0.1" is typically (may be OS dependent) only available on the same device.
///        If you do not require a fixed \b port, use 0 (zero) as a port to tell the server to pick an arbitrary port
///        that is available. The port can be queried via obx_sync_server_port() once the server was started.
///        \b Examples: "ws://0.0.0.0:9999" could be used during development (no certificate config needed),
///        while in a production system, you may want to use wss and a specific IP for security reasons.
/// @see obx_store_open()
/// @returns NULL if server could not be created (e.g. the store could not be opened, bad uri, etc.)
OBX_C_API OBX_sync_server* obx_sync_server(OBX_store_options* store_options, const char* uri);

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
/// (a "0" port was used in the URI given to obx_sync_server()).
OBX_C_API uint16_t obx_sync_server_port(OBX_sync_server* server);

/// Returns the number of clients connected to this server.
OBX_C_API uint64_t obx_sync_server_connections(OBX_sync_server* server);

/// Get server runtime statistics.
/// The returned char* is valid until another call to obx_sync_server_stats_string() or the server is closed.
OBX_C_API const char* obx_sync_server_stats_string(OBX_sync_server* server, bool include_zero_values);

// TODO admin UI ("browser")

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_SYNC_H

/**@}*/  // end of doxygen group