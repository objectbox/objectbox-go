/*
 * Copyright 2018-2025 ObjectBox Ltd. All rights reserved.
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

package objectbox

/*
#include <stdlib.h>
#include "objectbox-sync.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

// SyncClient is used to connect to an ObjectBox sync server.
type SyncClient struct {
	ob      *ObjectBox
	cClient *C.OBX_sync
	started bool

	// callback IDs for the "trampoline", see c-callbacks.go and how it's used in listeners
	cCallbacks [7]cCallbackId
}

// indexes into cCallbacks array
const (
	cCallbackIndexConnection = iota // 0
	cCallbackIndexDisconnection
	cCallbackIndexLogin
	cCallbackIndexLoginFailure
	cCallbackIndexCompletion
	cCallbackIndexChange
	cCallbackIndexServerTime // 6
)

// NewSyncClient creates a sync client associated with the given store and configures it with the given options.
// This does not initiate any connection attempts yet, call SyncClient.Start() to do so.
//
// Before SyncClient.Start(), you can still configure some aspects, e.g. SyncClient.SetRequestUpdatesMode().
func NewSyncClient(ob *ObjectBox, serverUri string, credentials *SyncCredentials) (*SyncClient, error) {
	if ob.syncClient != nil {
		return nil, errors.New("only one sync client can be active for a store, use ObjectBox.SyncClient() to access it")
	}

	var err error
	var client = &SyncClient{ob: ob}

	// close the sync client if some part of the initialization fails
	defer func() {
		if err != nil {
			if err2 := client.Close(); err2 != nil {
				err = fmt.Errorf("%s; %s", err, err2)
			}
			client = nil
		}
	}()

	err = cCallBool(func() bool {
		var cUri = C.CString(serverUri)
		defer C.free(unsafe.Pointer(cUri))
		client.cClient = C.obx_sync(ob.store, cUri)
		return client.cClient != nil
	})

	if err == nil {
		err = client.SetCredentials(credentials)
	}

	if err == nil {
		ob.syncClient = client
	}

	return client, err
}

// Close stops synchronization and frees the resources.
func (client *SyncClient) Close() error {
	if client.cClient == nil {
		return nil
	}

	for _, cbId := range client.cCallbacks {
		cCallbackUnregister(cbId)
	}

	client.ob.syncClient = nil

	return cCall(func() C.obx_err {
		defer func() { client.cClient = nil }()
		return C.obx_sync_close(client.cClient)
	})
}

// IsClosed returns true if this sync client is closed and can no longer be used.
func (client *SyncClient) IsClosed() bool {
	return client.cClient == nil
}

// SetCredentials configures authentication credentials, depending on your server config.
func (client *SyncClient) SetCredentials(credentials *SyncCredentials) error {
	return client.setOrAddCredentials(credentials, false, true)
}

// SetMultipleCredentials configures multiple authentication methods.
func (client *SyncClient) SetMultipleCredentials(multipleCredentials []*SyncCredentials) error {
	if len(multipleCredentials) == 0 {
		return errors.New("credentials array must not be empty")
	}

	for i := 0; i < len(multipleCredentials); i++ {
		var complete = i == len(multipleCredentials)-1
		err := client.setOrAddCredentials(multipleCredentials[i], true, complete)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *SyncClient) setOrAddCredentials(credentials *SyncCredentials, doAdd bool, complete bool) error {
	if credentials == nil {
		return errors.New("credentials must not be nil")
	}

	var err error
	if credentials.cType == C.OBXSyncCredentialsType_OBX_ADMIN_USER ||
		credentials.cType == C.OBXSyncCredentialsType_USER_PASSWORD {
		err = cCall(func() C.obx_err {
			var username = C.CString(credentials.username)
			defer C.free(unsafe.Pointer(username))
			var password = C.CString(credentials.password)
			defer C.free(unsafe.Pointer(password))
			if doAdd {
				return C.obx_sync_credentials_add_user_password(client.cClient, credentials.cType, username, password,
					C.bool(complete))
			} else {
				return C.obx_sync_credentials_user_password(client.cClient, credentials.cType, username, password)
			}
		})
	} else if credentials.cType == C.OBXSyncCredentialsType_JWT_ID ||
		credentials.cType == C.OBXSyncCredentialsType_JWT_ACCESS ||
		credentials.cType == C.OBXSyncCredentialsType_JWT_REFRESH ||
		credentials.cType == C.OBXSyncCredentialsType_JWT_CUSTOM {
		err = cCall(func() C.obx_err {
			var token = unsafe.Pointer(C.CString(credentials.dataString))
			defer C.free(token)
			var tokenLen = C.size_t(len(credentials.dataString))
			if doAdd {
				return C.obx_sync_credentials_add(client.cClient, credentials.cType, token, tokenLen, C.bool(complete))
			} else {
				return C.obx_sync_credentials(client.cClient, credentials.cType, token, tokenLen)
			}
		})
	} else {
		err = cCall(func() C.obx_err {
			var dataPtr unsafe.Pointer = nil
			if len(credentials.data) > 0 {
				dataPtr = unsafe.Pointer(&credentials.data[0])
			}
			var dataLen = C.size_t(len(credentials.data))
			if doAdd {
				return C.obx_sync_credentials_add(client.cClient, credentials.cType, dataPtr, dataLen, C.bool(complete))
			} else {
				return C.obx_sync_credentials(client.cClient, credentials.cType, dataPtr, dataLen)
			}
		})
	}
	return err
}

type syncRequestUpdatesMode uint

const (
	// SyncRequestUpdatesManual configures the client to only get updates when triggered manually using RequestUpdates()
	SyncRequestUpdatesManual syncRequestUpdatesMode = C.OBXRequestUpdatesMode_MANUAL

	// SyncRequestUpdatesAutomatic configures the client to get all updates automatically
	SyncRequestUpdatesAutomatic syncRequestUpdatesMode = C.OBXRequestUpdatesMode_AUTO

	// SyncRequestUpdatesAutoNoPushes configures the client to get all updates right after a log-in (initial and reconnects)
	SyncRequestUpdatesAutoNoPushes syncRequestUpdatesMode = C.OBXRequestUpdatesMode_AUTO_NO_PUSHES
)

// SetRequestUpdatesMode configures how/when the server will send the changes to us (the client). Can only be called
// before Start(). See SyncRequestUpdatesManual, SyncRequestUpdatesAutomatic, SyncRequestUpdatesAutoNoPushes.
func (client *SyncClient) SetRequestUpdatesMode(mode syncRequestUpdatesMode) error {
	return cCall(func() C.obx_err {
		return C.obx_sync_request_updates_mode(client.cClient, C.OBXRequestUpdatesMode(mode))
	})
}

type SyncClientState uint

const (
	// SyncClientStateCreated means the sync client has just been created
	SyncClientStateCreated SyncClientState = C.OBXSyncState_CREATED

	// SyncClientStateStarted means the sync client has been started (using start() method)
	SyncClientStateStarted SyncClientState = C.OBXSyncState_STARTED

	// SyncClientStateConnected means the sync client has a connection to the server (not logged in yet)
	SyncClientStateConnected SyncClientState = C.OBXSyncState_CONNECTED

	// SyncClientStateLoggedIn means the sync client has successfully logged in to the server
	SyncClientStateLoggedIn SyncClientState = C.OBXSyncState_LOGGED_IN

	// SyncClientStateDisconnected means the sync client has lost/closed the connection to the server
	SyncClientStateDisconnected SyncClientState = C.OBXSyncState_DISCONNECTED

	// SyncClientStateStopped means the sync client has stopped synchronization
	SyncClientStateStopped SyncClientState = C.OBXSyncState_STOPPED

	// SyncClientStateDead means the sync client is in an unrecoverable state
	SyncClientStateDead SyncClientState = C.OBXSyncState_DEAD
)

// State returns the current state of the sync client
func (client *SyncClient) State() SyncClientState {
	return SyncClientState(C.obx_sync_state(client.cClient))
}

// Start initiates the connection to the server and begins the synchronization
func (client *SyncClient) Start() error {
	client.started = true
	return cCall(func() C.obx_err {
		return C.obx_sync_start(client.cClient)
	})
}

// Stop stops the synchronization and closes the connection to the server Does nothing if it is already stopped.
func (client *SyncClient) Stop() error {
	var err = cCall(func() C.obx_err {
		return C.obx_sync_stop(client.cClient)
	})
	client.started = false
	return err
}

// WaitForLogin - waits for the sync client to get into the given state or until the given timeout is reached.
// For an asynchronous alternative, please check the listeners. Start() is called automatically if it hasn't been yet.
// Returns:
//
//	(true, nil) in case the login was successful;
//	(false, nil) in case of a time out;
//	(false, error) if an error occurred (such as wrong credentials)
func (client *SyncClient) WaitForLogin(timeout time.Duration) (successful bool, err error) {
	if !client.started {
		if err := client.Start(); err != nil {
			return false, err
		}
	}

	var timeoutMs = timeout.Nanoseconds() / 1000 / 1000 // .Milliseconds() only since Go 1.13+
	if timeoutMs < 0 {
		return false, fmt.Errorf("timeout must be >= 1 millisecond, %d given", timeoutMs)
	}

	var code = C.obx_sync_wait_for_logged_in_state(client.cClient, C.uint64_t(timeoutMs))
	switch code {
	case C.OBX_SUCCESS:
		return true, nil
	case C.OBX_TIMEOUT:
		return false, nil
	case C.OBX_NO_SUCCESS:
		// From native function documentation: a state was reached within the given timeout that is unlikely to result
		// in a successful login, e.g. "disconnected". Note: obx_last_error_code() is not set
		return false, errors.New("a state unlikely to lead to a successful login was encountered, " +
			"this usually indicates an issue with credentials (not set yet, wrong credentials)")
	default:
		return false, createError()
	}
}

// RequestUpdates can be used to manually synchronize incoming changes in case the client is running in "Manual" or
// "AutoNoPushes" mode (i.e. it doesn't get the updates automatically). Additionally, it can be used to subscribe for
// future pushes (similar to the "Auto" mode).
func (client *SyncClient) RequestUpdates(alsoSubscribe bool) error {
	return cCall(func() C.obx_err {
		return C.obx_sync_updates_request(client.cClient, C.bool(alsoSubscribe))
	})
}

// CancelUpdates can be used to unsubscribe from manually requested updates (see `RequestUpdates(true)`).
func (client *SyncClient) CancelUpdates() error {
	return cCall(func() C.obx_err {
		return C.obx_sync_updates_cancel(client.cClient)
	})
}

// SyncChange describes a single incoming data event received by the sync client
type SyncChange struct {
	EntityId TypeId
	Puts     []uint64
	Removals []uint64
}

type syncChangeListener func(changes []*SyncChange)
type syncConnectionListener func()
type syncDisconnectionListener func()
type syncLoginListener func()
type syncLoginFailureListener func(code SyncLoginFailure)
type syncCompletionListener func()
type syncTimeListener func(time.Time)

type SyncLoginFailure uint64 // TODO enumerate possible values

// SetConnectionListener sets or overrides a previously set listener for a "connection" event.
func (client *SyncClient) SetConnectionListener(callback syncConnectionListener) error {
	if callback == nil {
		C.obx_sync_listener_connect(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexConnection])
	} else {
		if cbId, err := cCallbackRegister(cVoidCallback(func() {
			callback()
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_connect(client.cClient, (*C.OBX_sync_listener_connect)(cVoidCallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexConnection, cbId)
		}
	}
	return nil
}

// SetDisconnectionListener sets or overrides a previously set listener for a "disconnection" event.
func (client *SyncClient) SetDisconnectionListener(callback syncDisconnectionListener) error {
	if callback == nil {
		C.obx_sync_listener_disconnect(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexDisconnection])
	} else {
		if cbId, err := cCallbackRegister(cVoidCallback(func() {
			callback()
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_disconnect(client.cClient, (*C.OBX_sync_listener_disconnect)(cVoidCallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexDisconnection, cbId)
		}
	}
	return nil
}

// SetLoginListener sets or overrides a previously set listener for a "login" event.
func (client *SyncClient) SetLoginListener(callback syncLoginListener) error {
	if callback == nil {
		C.obx_sync_listener_login(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexLogin])
	} else {
		if cbId, err := cCallbackRegister(cVoidCallback(func() {
			callback()
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_login(client.cClient, (*C.OBX_sync_listener_login)(cVoidCallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexLogin, cbId)
		}
	}
	return nil
}

// SetLoginFailureListener sets or overrides a previously set listener for a "login" event.
func (client *SyncClient) SetLoginFailureListener(callback syncLoginFailureListener) error {
	if callback == nil {
		C.obx_sync_listener_login_failure(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexLoginFailure])
	} else {
		if cbId, err := cCallbackRegister(cVoidUint64Callback(func(code uint64) {
			callback(SyncLoginFailure(code))
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_login_failure(client.cClient, (*C.OBX_sync_listener_login_failure)(cVoidUint64CallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexLoginFailure, cbId)
		}
	}
	return nil
}

// SetCompletionListener sets or overrides a previously set listener for a "login" event.
func (client *SyncClient) SetCompletionListener(callback syncCompletionListener) error {
	if callback == nil {
		C.obx_sync_listener_complete(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexCompletion])
	} else {
		if cbId, err := cCallbackRegister(cVoidCallback(func() {
			callback()
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_complete(client.cClient, (*C.OBX_sync_listener_complete)(cVoidCallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexCompletion, cbId)
		}
	}
	return nil
}

// SetServerTimeListener sets or overrides a previously set listener for a "login" event.
func (client *SyncClient) SetServerTimeListener(callback syncTimeListener) error {
	if callback == nil {
		C.obx_sync_listener_server_time(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexServerTime])
	} else {
		if cbId, err := cCallbackRegister(cVoidInt64Callback(func(timestampNs int64) {
			const nsInSec = 1000 * 1000 * 1000
			callback(time.Unix(timestampNs/nsInSec, timestampNs%nsInSec))
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_server_time(client.cClient, (*C.OBX_sync_listener_server_time)(cVoidInt64CallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexServerTime, cbId)
		}
	}
	return nil
}

// SetChangeListener sets or overrides a previously set listener for incoming changes notifications.
// SyncChange event is issued after a transaction is applied to the local database.
func (client *SyncClient) SetChangeListener(callback syncChangeListener) error {
	if callback == nil {
		C.obx_sync_listener_change(client.cClient, nil, nil)
		cCallbackUnregister(client.cCallbacks[cCallbackIndexChange])
	} else {
		if cbId, err := cCallbackRegister(cVoidConstVoidCallback(func(cChangeList unsafe.Pointer) {
			callback(cSyncChangeArrayToGo((*C.OBX_sync_change_array)(cChangeList)))
		})); err != nil {
			return err
		} else {
			C.obx_sync_listener_change(client.cClient, (*C.OBX_sync_listener_change)(cVoidConstVoidCallbackDispatchPtr), cbId.cPtr())
			client.swapCallbackId(cCallbackIndexChange, cbId)
		}
	}
	return nil
}

func (client *SyncClient) swapCallbackId(index uint, newId cCallbackId) {
	cCallbackUnregister(client.cCallbacks[index])
	client.cCallbacks[index] = newId
}

func cSyncChangeArrayToGo(cArray *C.OBX_sync_change_array) []*SyncChange {
	var size = uint(cArray.count)
	var changes = make([]*SyncChange, 0, size)
	if size > 0 {
		var cArrayStart = unsafe.Pointer(cArray.list)
		var cItemSize = unsafe.Sizeof(*cArray.list)
		for i := uint(0); i < size; i++ {
			var itemPtr = (*C.OBX_sync_change)(unsafe.Pointer(uintptr(cArrayStart) + uintptr(i)*cItemSize))

			var change = &SyncChange{
				EntityId: TypeId(itemPtr.entity_id),
			}

			if itemPtr.puts != nil {
				change.Puts = cIdsArrayToGo(itemPtr.puts)
			}

			if itemPtr.removals != nil {
				change.Removals = cIdsArrayToGo(itemPtr.removals)
			}

			changes = append(changes, change)
		}
	}
	return changes
}
