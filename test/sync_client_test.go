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
	"strings"
	"testing"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
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
	defer client.Close()

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
	var endTime = time.After(timeout)
	tick := time.Tick(time.Millisecond)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-endTime:
			return errors.New("timeout while waiting for a condition to become true")
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

type testSyncClient struct {
	t    *testing.T
	env  *model.TestEnv
	sync *objectbox.SyncClient
}

func NewTestSyncClient(t *testing.T, serverURI, name string) *testSyncClient {
	var client = &testSyncClient{t: t}

	client.env = model.NewTestEnv(t)

	client.sync = client.env.SyncClient(serverURI)
	return client
}

func (client *testSyncClient) Close() {
	assert.NoErr(client.t, client.sync.Close())
	client.env.Close()
}

func (client *testSyncClient) Start() {
	assert.NoErr(client.t, client.sync.Start())

	assert.NoErr(client.t, waitUntil(time.Second, func() (bool, error) {
		return client.sync.State() == objectbox.SyncClientStateLoggedIn, nil
	}))
}

func TestSyncDataAutomatic(t *testing.T) {
	var server = NewTestSyncServer(t)
	defer server.Close()

	var a = NewTestSyncClient(t, server.URI(), "a")
	defer a.Close()
	a.Start()

	var b = NewTestSyncClient(t, server.URI(), "b")
	defer b.Close()
	b.Start()

	isEmpty, err := a.env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	isEmpty, err = b.env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	// insert into one box
	var count uint = 10
	a.env.Populate(count)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}))

	var assertEqualBoxes = func(boxA, boxB *model.EntityBox) {
		itemsA, err := a.env.Box.GetAll()
		assert.NoErr(t, err)

		itemsB, err := b.env.Box.GetAll()
		assert.NoErr(t, err)

		assert.Eq(t, count, uint(len(itemsA)))
		assert.Eq(t, count, uint(len(itemsB)))
		assert.Eq(t, itemsA, itemsB)
	}
	assertEqualBoxes(a.env.Box, b.env.Box)

	// remove from one of the boxes
	removed, err := b.env.Box.RemoveIds(1, 3, 6)
	assert.NoErr(t, err)
	assert.True(t, 3 == removed)
	count = count - uint(removed)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := a.env.Box.Count()
		return bCount == uint64(count), err
	}))

	assertEqualBoxes(a.env.Box, b.env.Box)
}

func TestSyncDataManual(t *testing.T) {
	var server = NewTestSyncServer(t)
	defer server.Close()

	var a = NewTestSyncClient(t, server.URI(), "a")
	defer a.Close()
	assert.NoErr(t, a.sync.UpdatesMode(objectbox.SyncClientUpdatesManual))
	a.Start()

	var b = NewTestSyncClient(t, server.URI(), "b")
	defer b.Close()
	assert.NoErr(t, b.sync.UpdatesMode(objectbox.SyncClientUpdatesManual))
	b.Start()

	isEmpty, err := a.env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	isEmpty, err = b.env.Box.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	// insert into one box
	var count uint = 10
	a.env.Populate(count)

	// this will time out because we haven't manually initiated an update
	assert.True(t, strings.Contains(waitUntil(500*time.Millisecond, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}).Error(), "timeout"))

	// manually trigger the data synchronization
	assert.NoErr(t, b.sync.RequestUpdates(false))

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}))

	assert.NoErr(t, a.env.Box.RemoveAll())
	count = 0

	// this will time out because we haven't subscribed for all further updates
	assert.True(t, strings.Contains(waitUntil(500*time.Millisecond, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}).Error(), "timeout"))

	// subscribe for further updates
	assert.NoErr(t, b.sync.RequestUpdates(true))

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}))

	count = 10
	a.env.Populate(count)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := b.env.Box.Count()
		return bCount == uint64(count), err
	}))
}

func TestSyncWaitForLogin(t *testing.T) {
	var server = NewTestSyncServer(t)
	defer server.Close()

	// success
	var a = NewTestSyncClient(t, server.URI(), "a")
	defer a.Close()
	timedOut, err := a.sync.WaitForLogin(time.Second)
	assert.NoErr(t, err)
	assert.True(t, !timedOut)

	// failure
	var b = NewTestSyncClient(t, server.URI(), "b")
	defer b.Close()
	assert.NoErr(t, b.sync.AuthSharedSecret([]byte{1}))
	timedOut, err = b.sync.WaitForLogin(time.Second)
	assert.True(t, strings.Contains(err.Error(), "credentials rejected"))
	assert.True(t, !timedOut)

	// time out
	var c = NewTestSyncClient(t, server.URI(), "b")
	defer c.Close()
	timedOut, err = c.sync.WaitForLogin(time.Nanosecond)
	assert.NoErr(t, err)
	assert.True(t, timedOut)
}

func TestSyncOnChange(t *testing.T) {
	var server = NewTestSyncServer(t)
	defer server.Close()

	var a = NewTestSyncClient(t, server.URI(), "a")
	defer a.Close()
	a.Start()

	var putIDs = make([]uint64, 0)
	var removedIDs = make([]uint64, 0)

	var b = NewTestSyncClient(t, server.URI(), "b")
	defer b.Close()
	assert.NoErr(t, b.sync.OnChange(func(changes []*objectbox.SyncChangeNotification) {
		t.Logf("received %d changes", len(changes))
		for i, change := range changes {
			t.Logf("change %d: %v", i, change)

			// only count the main entity, not relations
			if change.EntityId == model.EntityBinding.Id {
				putIDs = append(putIDs, change.PutIds...)
				removedIDs = append(removedIDs, change.RemovedIds...)
			}
		}
	}))
	b.Start()

	// insert on one client
	var count uint = 100
	a.env.Populate(count)

	// wait for the data to be received by another client - its onChange() listener is called
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		return count == uint(len(putIDs)), nil
	}))

	assert.Eq(t, 0, len(removedIDs))

	var expectedIds []uint64
	for id := uint(1); id <= count; id++ {
		expectedIds = append(expectedIds, uint64(id))
	}
	assert.Eq(t, expectedIds, putIDs)
}
