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

package objectbox

/*
#include <stdlib.h>
#include "objectbox-sync.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// SyncClient is used to connect to an ObjectBox sync server.
type SyncClient struct {
	ob      *ObjectBox
	cClient *C.OBX_sync
	started bool
	state   syncClientInternalState

	// these are unregistered when closing
	cCallbackIds []cCallbackId
}

// NewSyncClient creates a sync client associated with the given store and configures it with the given options.
// This does not initiate any connection attempts yet, call SyncClient.Start() to do so.
//
// Before SyncClient.Start(), you can still configure some aspects, e.g. SyncClient.SetRequestUpdatesMode().
func NewSyncClient(ob *ObjectBox, serverUri string, credentials *SyncCredentials) (err error, client *SyncClient) {
	if ob.syncClient != nil {
		return errors.New("only one sync client can be active for a store"), nil
	}

	client = &SyncClient{ob: ob}

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
		err = client.registerCallbacks()
	}

	if err == nil {
		err = client.SetCredentials(credentials)
	}

	if err == nil {
		ob.syncClient = client
	}

	return err, client
}

func (client *SyncClient) registerCallbacks() error {
	// login
	if callbackId, err := cCallbackRegister(cVoidCallback(func() {
		client.state.Update(func(state *syncClientInternalState) {
			state.loggedIn = true
			state.loginError = nil
		})
	})); err != nil {
		return err
	} else {
		client.cCallbackIds = append(client.cCallbackIds, callbackId)
		C.obx_sync_listener_login(client.cClient, (*C.OBX_sync_listener_login)(cVoidCallbackDispatchPtr), unsafe.Pointer(&callbackId))
	}

	// login failed
	if callbackId, err := cCallbackRegister(cVoidUint64Callback(func(code uint64) {
		client.state.Update(func(state *syncClientInternalState) {
			state.loggedIn = false
			switch code {
			case C.OBXSyncCode_CREDENTIALS_REJECTED:
				state.loginError = errors.New("credentials rejected")
			case C.OBXSyncCode_AUTH_UNREACHABLE:
				state.loginError = errors.New("authentication unreachable")
			default:
				state.loginError = fmt.Errorf("error code %v", code)
			}
		})
	})); err != nil {
		return err
	} else {
		client.cCallbackIds = append(client.cCallbackIds, callbackId)
		C.obx_sync_listener_login_failure(client.cClient, (*C.OBX_sync_listener_login_failure)(cVoidUint64CallbackDispatchPtr), unsafe.Pointer(&callbackId))
	}

	return nil
}

// Close stops synchronization and frees the resources.
func (client *SyncClient) Close() error {
	if client.cClient == nil {
		return nil
	}

	for _, cbId := range client.cCallbackIds {
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
	if credentials == nil {
		return errors.New("credentials must not be nil")
	}

	return cCall(func() C.obx_err {
		var dataPtr unsafe.Pointer = nil
		if len(credentials.data) > 0 {
			dataPtr = unsafe.Pointer(&credentials.data[0])
		}
		return C.obx_sync_credentials(client.cClient, credentials.cType, dataPtr, C.size_t(len(credentials.data)))
	})
}

type syncRequestUpdatesMode uint

const (
	// SyncRequestUpdatesManual configures the client to only get updates when triggered manually using RequestUpdates()
	SyncRequestUpdatesManual syncRequestUpdatesMode = C.OBXRequestUpdatesMode_MANUAL

	// SyncRequestUpdatesAutomatic configures the client to get all updates automatically
	SyncRequestUpdatesAutomatic syncRequestUpdatesMode = C.OBXRequestUpdatesMode_AUTO

	// SyncRequestUpdatesAutoNoPushes configures the client to get all updates during log-in (initial and reconnects)
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

// WaitForLogin initiates the connection to the server and waits for a login response (either a success or a failure)
// Returns:
// 		(true, nil) in case of a time out;
// 		(false, nil) in case the login was successful;
// 		(false, error) if an error occurred (such as wrong credentials)
func (client *SyncClient) WaitForLogin(timeout time.Duration) (timedOut bool, err error) {
	if !client.started {
		if err := client.Start(); err != nil {
			return false, err
		}
	}

	return waitUntil(timeout, time.Millisecond, func() (result bool, err error) {
		client.state.Lock()
		result = client.state.loggedIn
		err = client.state.loginError
		client.state.Unlock()
		return
	})
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

// SetChangeListener attaches a callback to receive incoming changes notifications.
// SyncChange event is issued after a transaction is applied to the local database.
func (client *SyncClient) SetChangeListener(callback syncChangeListener) error {
	if client.started {
		return errors.New("cannot attach an SetChangeListener listener - already started")
	}

	if callbackId, err := cCallbackRegister(cVoidConstVoidCallback(func(cChangeList unsafe.Pointer) {
		callback(cSyncChangeArrayToGo((*C.OBX_sync_change_array)(cChangeList)))
	})); err != nil {
		return err
	} else {
		client.cCallbackIds = append(client.cCallbackIds, callbackId)
		C.obx_sync_listener_change(client.cClient, (*C.OBX_sync_listener_change)(cVoidConstVoidCallbackDispatchPtr), unsafe.Pointer(&callbackId))
	}

	return nil
}

type syncClientInternalState struct {
	sync.Mutex
	loggedIn   bool
	loginError error
}

// Update changes the state under a mutex and signals the conditional variable
func (state *syncClientInternalState) Update(fn func(*syncClientInternalState)) {
	state.Lock()
	defer func() {
		state.Unlock()
	}()
	fn(state)
}
