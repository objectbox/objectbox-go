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

package objectbox_test

import (
	"errors"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
	"testing"
	"time"
)

func TestSyncAuth(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	// actually starting a server for this test case is not necessary
	const serverURI = "ws://127.0.0.1:9999"

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.AuthSharedSecret(nil))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.AuthSharedSecret([]byte{1, 2, 3}))
		assert.NoErr(t, client.Close())
	}
}

func TestSyncUpdatesMode(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	// actually starting a server for this test case is not necessary
	const serverURI = "ws://127.0.0.1:9999"

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesManual))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesAutomatic))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesOnLogin))
		assert.NoErr(t, client.Close())
	}
}

func TestSyncState(t *testing.T) {
	var server = NewTestSyncServer(t)
	defer func() {
		if server != nil {
			server.Close()
		}
	}()

	var env = model.NewTestEnv(t)
	defer env.Close()

	var client = env.SyncClient(server.URI())
	assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesManual))

	assert.Eq(t, objectbox.SyncClientStateCreated, client.State())

	assert.NoErr(t, client.Start())
	assert.Eq(t, objectbox.SyncClientStateStarted, client.State())

	// we're trying both connected and logged in because we might miss the connected state if the login is too fast
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		var state = client.State()
		return state == objectbox.SyncClientStateConnected || state == objectbox.SyncClientStateLoggedIn, nil
	}))

	// now the client must be logged in
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		return client.State() == objectbox.SyncClientStateLoggedIn, nil
	}))

	// stop the server while the client is connected
	server.Close()
	server = nil // prevent a double Close() by defer

	// must be disconnected now
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		return client.State() == objectbox.SyncClientStateDisconnected, nil
	}))

	assert.NoErr(t, client.Stop())
	assert.Eq(t, objectbox.SyncClientStateStopped, client.State())

	assert.NoErr(t, client.Close())
}

func waitUntil(timeout time.Duration, fn func() (bool, error)) error {
	var endtime = time.After(timeout)
	tick := time.Tick(time.Millisecond)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-endtime:
			return errors.New("timed out while waiting for a condition to become true")
		// Got a tick, we should check on doSomething()
		case <-tick:
			if ok, err := fn(); err != nil {
				return err
			} else if ok {
				return nil
			}
			// fn() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}