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

#ifdef __cplusplus
extern "C" {
#endif

//----------------------------------------------
// Sync (client)
//----------------------------------------------

/// Before calling any of the other sync APIs, ensure that those are actually available.
/// If the application is linked a non-sync ObjectBox library, this allows you to fail gracefully.
/// @return true if this library comes with the sync API
/// @deprecated use obx_has_feature(OBXFeature_Sync)
bool obx_sync_available();

struct OBX_sync;
typedef struct OBX_sync OBX_sync;

typedef enum {
    OBXSyncCredentialsType_NONE = 0,
    OBXSyncCredentialsType_SHARED_SECRET = 1,
    OBXSyncCredentialsType_GOOGLE_AUTH = 2,
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

typedef struct OBX_sync_change {
    obx_schema_id entity_id;
    const OBX_id_array* puts;
    const OBX_id_array* removals;
} OBX_sync_change;

typedef struct OBX_sync_change_array {
    const OBX_sync_change* list;
    size_t count;
} OBX_sync_change_array;

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

/// Creates a sync client associated with the given store and sync server URI.
/// This does not initiate any connection attempts yet: call obx_sync_start() to do so.
/// Before obx_sync_start(), you must configure credentials via obx_sync_credentials.
/// By default a sync client automatically receives updates from the server once login succeeded.
/// To configure this differently, call obx_sync_request_updates_mode() with the wanted mode.
OBX_sync* obx_sync(OBX_store* store, const char* server_uri);

/// Stops and closes (deletes) the sync client freeing its resources.
obx_err obx_sync_close(OBX_sync* sync);

/// Sets credentials to authenticate the client with the server.
/// See OBXSyncCredentialsType for available options.
/// The accepted OBXSyncCredentials type depends on your sync-server configuration.
/// @param data may be NULL, i.e. in combination with OBXSyncCredentialsType_NONE
obx_err obx_sync_credentials(OBX_sync* sync, OBXSyncCredentialsType type, const void* data, size_t size);

/// Configures the maximum number of outgoing TX messages that can be sent without an ACK from the server.
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if value is not in the range 1-20
obx_err obx_sync_max_messages_in_flight(OBX_sync* sync, int value);

/// Sets the interval in which the client sends "heartbeat" messages to the server, keeping the connection alive.
/// To detect disconnects early on the client side, you can also use heartbeats with a smaller interval.
/// Use with caution, setting a low value (i.e. sending heartbeat very often) may cause an excessive network usage
/// as well as high server load (with many clients).
/// @param interval_ms interval in milliseconds; the default is 25 minutes (1 500 000 milliseconds),
///        which is also the allowed maximum.
/// @returns OBX_ERROR_ILLEGAL_ARGUMENT if value is not in the allowed range, e.g. larger than the maximum (1 500 000).
obx_err obx_sync_heartbeat_interval(OBX_sync* sync, uint64_t interval_ms);

/// Triggers the heartbeat sending immediately.
/// @see obx_sync_heartbeat_interval()
obx_err obx_sync_send_heartbeat(OBX_sync* sync);

/// Switches operation mode that's initialized after successful login
/// Must be called before obx_sync_start (returns OBX_ERROR_ILLEGAL_STATE if it was already started)
obx_err obx_sync_request_updates_mode(OBX_sync* sync, OBXRequestUpdatesMode mode);

/// Once the sync client is configured, you can "start" it to initiate synchronization.
/// This method triggers communication in the background and will return immediately.
/// If the synchronization destination is reachable, this background thread will connect to the server,
/// log in (authenticate) and, depending on "update request mode", start syncing data.
/// If the device, network or server is currently offline, connection attempts will be retried later using
/// increasing backoff intervals.
obx_err obx_sync_start(OBX_sync* sync);

/// Stops this sync client and thus stopping the synchronization. Does nothing if it is already stopped.
obx_err obx_sync_stop(OBX_sync* sync);

/// Gets the current state of the sync client (0 on error, e.g. given sync was NULL)
OBXSyncState obx_sync_state(OBX_sync* sync);

/// Waits for the sync client to get into the given state or until the given timeout is reached.
/// For an asynchronous alternative, please check the listeners.
/// @param timeout_millis Must be greater than 0
/// @returns OBX_SUCCESS if the LOGGED_IN state has been reached within the given timeout
/// @returns OBX_TIMEOUT if the given timeout was reached before a relevant state change was detected.
///          Note: obx_last_error_code() is not set.
/// @returns OBX_NO_SUCCESS if a state was reached within the given timeout that is unlikely to result in a
///          successful login, e.g. "disconnected". Note: obx_last_error_code() is not set.
obx_err obx_sync_wait_for_logged_in_state(OBX_sync* sync, uint64_t timeout_millis);

/// Request updates from the server since we last synced our database.
/// @param subscribe_for_pushes keep sending us future updates as they come in.
/// This should only be called in "logged in" state and there are no delivery guarantees given.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
obx_err obx_sync_updates_request(OBX_sync* sync, bool subscribe_for_pushes);

/// Cancel updates from the server (once received, the server stops sending updates).
/// The counterpart to obx_sync_updates_request().
/// This should only be called in "logged in" state and there are no delivery guarantees given.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
obx_err obx_sync_updates_cancel(OBX_sync* sync);

/// Count the number of messages in the outgoing queue, i.e. those waiting to be sent to the server.
/// @param limit pass 0 to count all messages without any limitation or a lower number that's enough for your app logic.
/// @note This calls uses a (read) transaction internally: 1) it's not just a "cheap" return of a single number.
///       While this will still be fast, avoid calling this function excessively.
///       2) the result follows transaction view semantics, thus it may not always match the actual value.
obx_err obx_sync_outgoing_message_count(OBX_sync* sync, uint64_t limit, uint64_t* out_count);

/// Experimental. This API is likely to be replaced/removed in a future version.
/// Quickly bring our database up-to-date in a single transaction, without transmitting all the history.
/// Good for initial sync of new clients.
/// @returns OBX_SUCCESS if the request was likely sent (e.g. the sync client is in "logged in" state)
/// @returns OBX_NO_SUCCESS if the request was not sent (and will not be sent in the future).
///          Note: obx_last_error_code() is not set.
obx_err obx_sync_full(OBX_sync* sync);

/// Estimates the current server timestamp based on the last known server time and local steady clock.
/// @param out_timestamp_ns - unix timestamp in nanoseconds - may be set to zero if the last server's time is unknown.
obx_err obx_sync_server_time(OBX_sync* sync, int64_t* out_timestamp_ns);

/// Returns the estimated difference between the server time and the local timestamp.
/// Equivalent to calculating obx_sync_server_time() - "current time" (nanos since epoch).
/// @param out_diff_ns time difference in nanoseconds; e.g. positive if server time is ahead of local time.
///                    Set to 0 if there has not been a server contact yet and thus the server's time is unknown.
obx_err obx_sync_server_time_diff(OBX_sync* sync, int64_t* out_diff_ns);

/// Set or overwrite a previously set 'connect' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_connect(OBX_sync* sync, OBX_sync_listener_connect* listener, void* listener_arg);

/// Set or overwrite a previously set 'disconnect' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_disconnect(OBX_sync* sync, OBX_sync_listener_disconnect* listener, void* listener_arg);

/// Set or overwrite a previously set 'login' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_login(OBX_sync* sync, OBX_sync_listener_login* listener, void* listener_arg);

/// Set or overwrite a previously set 'login failure' listener.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_login_failure(OBX_sync* sync, OBX_sync_listener_login_failure* listener, void* listener_arg);

/// Set or overwrite a previously set 'complete' listener - notifies when the latest sync has finished.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_complete(OBX_sync* sync, OBX_sync_listener_complete* listener, void* listener_arg);

/// Set or overwrite a previously set 'change' listener - provides information about incoming changes.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_change(OBX_sync* sync, OBX_sync_listener_change* listener, void* listener_arg);

/// Set or overwrite a previously set 'serverTime' listener - provides current time updates from the sync-server.
/// @param listener set NULL to reset
/// @param listener_arg is a pass-through argument passed to the listener
void obx_sync_listener_server_time(OBX_sync* sync, OBX_sync_listener_server_time* listener, void* listener_arg);

#ifdef __cplusplus
}
#endif

#endif  // OBJECTBOX_SYNC_H

/**@}*/  // end of doxygen group