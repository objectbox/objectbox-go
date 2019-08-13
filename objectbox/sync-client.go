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
#include "objectbox.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// SyncClient provides automated data synchronization with other clients connected to the same server
type SyncClient struct {
	ob      *ObjectBox
	cClient *C.OBX_sync
	authSet bool
	started bool
	state   syncClientState

	// these are unregistered when closing
	callbackIds []cCallbackId
}

// NewSyncClient starts a creation of a new sync client.
// See other methods for configuration and then use Start() to begin synchronization.
func NewSyncClient(ob *ObjectBox, serverUri string) (err error, client *SyncClient) {
	if !C.obx_sync_available() {
		return errors.New("sync client is not available"), nil
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

	return err, client
}

type syncClientCode uint64

const (
	syncClientCodeOk                  syncClientCode = C.OBXSyncCode_OK
	syncClientCodeRequestRejected     syncClientCode = C.OBXSyncCode_REQ_REJECTED
	syncClientCodeCredentialsRejected syncClientCode = C.OBXSyncCode_CREDENTIALS_REJECTED
	syncClientCodeUnknown             syncClientCode = C.OBXSyncCode_UNKNOWN
	syncClientCodeAuthUnreachable     syncClientCode = C.OBXSyncCode_AUTH_UNREACHABLE
	syncClientCodeBadVersion          syncClientCode = C.OBXSyncCode_BAD_VERSION
	syncClientCodeClientIdTaken       syncClientCode = C.OBXSyncCode_CLIENT_ID_TAKEN
	syncClientCodeTxViolatedUnique    syncClientCode = C.OBXSyncCode_TX_VIOLATED_UNIQUE
)

func (client *SyncClient) registerCallbacks() error {
	// login
	if callbackId, err := cCallbackRegister(cVoidCallback(func() {
		client.state.Update(func(state *syncClientState) {
			state.loggedIn = true
			state.loginError = nil
		})
	})); err != nil {
		return err
	} else {
		client.callbackIds = append(client.callbackIds, callbackId)
		C.obx_sync_listener_login(client.cClient, (*C.OBX_sync_listener_login)(cVoidCallbackDispatchPtr), unsafe.Pointer(&callbackId))
	}

	// login failed
	if callbackId, err := cCallbackRegister(cVoidUint64Callback(func(code uint64) {
		client.state.Update(func(state *syncClientState) {
			state.loggedIn = false
			switch syncClientCode(code) {
			case syncClientCodeCredentialsRejected:
				state.loginError = errors.New("credentials rejected")
			case syncClientCodeAuthUnreachable:
				state.loginError = errors.New("authentication unreachable")
			default:
				state.loginError = fmt.Errorf("error code %v", code)
			}
		})
	})); err != nil {
		return err
	} else {
		client.callbackIds = append(client.callbackIds, callbackId)
		C.obx_sync_listener_login_failure(client.cClient, (*C.OBX_sync_listener_login_failure)(cVoidUint64CallbackDispatchPtr), unsafe.Pointer(&callbackId))
	}

	return nil
}

// Close stops synchronization and frees the resources.
func (client *SyncClient) Close() error {
	if client.cClient == nil {
		return nil
	}

	for _, cbId := range client.callbackIds {
		cCallbackUnregister(cbId)
	}

	return cCall(func() C.obx_err {
		defer func() { client.cClient = nil }()
		return C.obx_sync_close(client.cClient)
	})
}

// AuthSharedSecret configures the client to use shared-secret authentication
func (client *SyncClient) AuthSharedSecret(data []byte) error {
	client.authSet = true
	return cCall(func() C.obx_err {
		var dataPtr unsafe.Pointer = nil
		if len(data) > 0 {
			dataPtr = unsafe.Pointer(&data[0])
		}
		return C.obx_sync_credentials(client.cClient, C.OBXSyncCredentialsType_SHARED_SECRET, dataPtr, C.size_t(len(data)))
	})
}

type syncClientUpdatesMode uint

const (
	// SyncClientUpdatesManual configures the client to only get updates when triggered manually using RequestUpdates()
	SyncClientUpdatesManual syncClientUpdatesMode = C.OBXRequestUpdatesMode_MANUAL

	// SyncClientUpdatesAutomatic configures the client to get all updates automatically
	SyncClientUpdatesAutomatic syncClientUpdatesMode = C.OBXRequestUpdatesMode_AUTO

	// SyncClientUpdatesOnLogin configures the client to get all updates during log-in (initial and reconnects)
	SyncClientUpdatesOnLogin syncClientUpdatesMode = C.OBXRequestUpdatesMode_AUTO_NO_PUSHES
)

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

// UpdatesMode configures how/when the server will send the changes to us (the client).
// Can only be called before Start(). See SyncClientUpdatesManual, SyncClientUpdatesAutomatic, SyncClientUpdatesOnLogin.
func (client *SyncClient) UpdatesMode(mode syncClientUpdatesMode) error {
	return cCall(func() C.obx_err {
		return C.obx_sync_request_updates_mode(client.cClient, C.OBXRequestUpdatesMode(mode))
	})
}

// State returns the current state of the sync client
func (client *SyncClient) State() SyncClientState {
	return SyncClientState(C.obx_sync_state(client.cClient))
}

// Start initiates the connection to the server and begins the synchronization
func (client *SyncClient) Start() error {
	client.started = true
	// If no authentication was provided by the user, try if the server accepts clients without any credentials at all.
	// That's what the client code/setup implies. Maybe the c-api should do this automatically.
	if !client.authSet {
		if err := cCall(func() C.obx_err {
			return C.obx_sync_credentials(client.cClient, C.OBXSyncCredentialsType_UNCHECKED, nil, 0)
		}); err != nil {
			return err
		}
	}

	return cCall(func() C.obx_err {
		return C.obx_sync_start(client.cClient)
	})
}

// Stop stops the synchronization and close the connection to the server
func (client *SyncClient) Stop() error {
	return cCall(func() C.obx_err {
		return C.obx_sync_stop(client.cClient)
	})
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

// RequestUpdates can be used to manually synchronize incomming changes in case the client is running in "Manual" or
// "OnLogin" mode (i.e. it doesn't get the updates automatically). Additionally, it can be used to subscribe for future
// pushes (similar to the "Auto" mode).
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

// DoFullSync is useful for new clients to quickly bring the local database up-to-date in a single transaction, without
// transmitting the whole history.
func (client *SyncClient) DoFullSync() error {
	return cCall(func() C.obx_err {
		return C.obx_sync_full(client.cClient)
	})
}

type syncClientState struct {
	sync.Mutex
	loggedIn   bool
	loginError error
}

// Update changes the state under a mutex and signals the conditional variable
func (state *syncClientState) Update(fn func(*syncClientState)) {
	state.Lock()
	defer func() {
		state.Unlock()
	}()
	fn(state)
}
