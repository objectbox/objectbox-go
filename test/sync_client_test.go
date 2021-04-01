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

package objectbox_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

func skipTestIfSyncNotAvailable(t *testing.T) {
	if !objectbox.SyncIsAvailable() {
		t.Skip("Sync is not available in the currently loaded ObjectBox native library")
	}
}

func TestSyncBasics(t *testing.T) {
	skipTestIfSyncNotAvailable(t)

	var env = model.NewTestEnv(t)
	defer env.Close()

	// actually starting a server for this test case is not necessary
	const serverURI = "ws://127.0.0.1:9999"

	client, err := env.ObjectBox.SyncClient()
	assert.Err(t, err)
	assert.True(t, client == nil)

	client = env.SyncClient(serverURI)
	assert.True(t, !client.IsClosed())

	clientFromStore, err := env.ObjectBox.SyncClient()
	assert.NoErr(t, err)
	assert.Eq(t, clientFromStore, client)

	assert.NoErr(t, client.Close())
	assert.True(t, client.IsClosed())

	client, err = env.ObjectBox.SyncClient()
	assert.Err(t, err)
	assert.True(t, client == nil)
}

func TestSyncAuth(t *testing.T) {
	skipTestIfSyncNotAvailable(t)

	var env = model.NewTestEnv(t)
	defer env.Close()

	// actually starting a server for this test case is not necessary
	const serverURI = "ws://127.0.0.1:9999"

	var client = env.SyncClient(serverURI)
	assert.Err(t, client.SetCredentials(nil))
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsSharedSecret([]byte{1, 2, 3})))
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsGoogleAuth([]byte{4, 5, 6})))
	assert.NoErr(t, client.Close())
}

func TestSyncUpdatesMode(t *testing.T) {
	skipTestIfSyncNotAvailable(t)

	var env = model.NewTestEnv(t)
	defer env.Close()

	// actually starting a server for this test case is not necessary
	const serverURI = "ws://127.0.0.1:9999"

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.SetRequestUpdatesMode(objectbox.SyncRequestUpdatesManual))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.SetRequestUpdatesMode(objectbox.SyncRequestUpdatesAutomatic))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient(serverURI)
		assert.NoErr(t, client.SetRequestUpdatesMode(objectbox.SyncRequestUpdatesAutoNoPushes))
		assert.NoErr(t, client.Close())
	}
}

func TestSyncState(t *testing.T) {
	skipTestIfSyncNotAvailable(t)

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

type syncTestEnv struct {
	t              *testing.T
	server         *testSyncServer
	clients        map[string]*testSyncClient
	defaultTimeout time.Duration
}

func NewSyncTestEnv(t *testing.T) *syncTestEnv {
	skipTestIfSyncNotAvailable(t)

	return &syncTestEnv{
		t:              t,
		server:         NewTestSyncServer(t),
		clients:        map[string]*testSyncClient{},
		defaultTimeout: 500 * time.Millisecond,
	}
}

func (env *syncTestEnv) NamedClient(name string) *testSyncClient {
	if env.clients[name] == nil {
		env.clients[name] = NewTestSyncClient(env.t, env.server.URI(), name)
	}
	return env.clients[name]
}

func (env *syncTestEnv) Client() *testSyncClient {
	return env.NamedClient("")
}

func (env *syncTestEnv) SyncClient() *objectbox.SyncClient {
	return env.Client().sync
}

func (env *syncTestEnv) Close() {
	if env.server != nil {
		for _, client := range env.clients {
			client.Close()
		}
		env.clients = nil
		env.server.Close()
		env.server = nil
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
	if client.sync != nil {
		assert.NoErr(client.t, client.sync.Close())
		client.env.Close()
		client.sync = nil
	}
}

func (client *testSyncClient) Start() {
	assert.NoErr(client.t, client.sync.Start())

	assert.NoErr(client.t, waitUntil(time.Second, func() (bool, error) {
		return client.sync.State() == objectbox.SyncClientStateLoggedIn, nil
	}))
}

func TestSyncDataAutomatic(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var a = env.NamedClient("a")
	a.Start()

	var b = env.NamedClient("b")
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
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var a = env.NamedClient("a")
	assert.NoErr(t, a.sync.SetRequestUpdatesMode(objectbox.SyncRequestUpdatesManual))
	a.Start()

	var b = env.NamedClient("b")
	assert.NoErr(t, b.sync.SetRequestUpdatesMode(objectbox.SyncRequestUpdatesManual))
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
	var env = NewSyncTestEnv(t)
	defer env.Close()

	// success
	var a = env.NamedClient("a")
	successful, err := a.sync.WaitForLogin(time.Second)
	assert.NoErr(t, err)
	assert.True(t, successful)

	// failure
	var b = env.NamedClient("b")
	assert.NoErr(t, b.sync.SetCredentials(objectbox.SyncCredentialsSharedSecret([]byte{1})))
	successful, err = b.sync.WaitForLogin(time.Second)
	if !strings.Contains(err.Error(), "credentials") {
		assert.Failf(t, "Error was expected to contain 'credentials': %v", err)
	}
	assert.True(t, !successful)

	// time out
	var c = env.NamedClient("c")
	env.server.Close()
	successful, err = c.sync.WaitForLogin(time.Millisecond)
	assert.NoErr(t, err)
	assert.True(t, !successful)
}

func TestSyncConnectionListener(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10) // buffer up to 10 values

	assert.NoErr(t, env.SyncClient().SetConnectionListener(func() { messages <- "connected" }))
	assert.NoErr(t, env.SyncClient().SetDisconnectionListener(func() { messages <- "disconnected" }))

	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "connected", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no more messages
	env.server.Close()
	assert.StringChannelExpect(t, "disconnected", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no more messages
}

func TestSyncConnectionListenerReset(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10) // buffer up to 10 values

	assert.NoErr(t, env.SyncClient().SetConnectionListener(func() { messages <- "connected" }))
	assert.NoErr(t, env.SyncClient().SetDisconnectionListener(func() { messages <- "disconnected" }))

	// reset should remove the listener
	assert.NoErr(t, env.SyncClient().SetConnectionListener(nil))
	assert.NoErr(t, env.SyncClient().SetDisconnectionListener(nil))

	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no messages
	env.server.Close()
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no messages
}

func TestSyncLoginListener(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10) // buffer up to 10 values

	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
	assert.NoErr(t, env.SyncClient().SetCredentials(objectbox.SyncCredentialsSharedSecret([]byte("invalid secret"))))
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "failure 43", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no more messages
	assert.NoErr(t, env.SyncClient().SetCredentials(objectbox.SyncCredentialsNone()))
	assert.StringChannelExpect(t, "success", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10) // no more messages
}

func TestSyncServerTimeListener(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan time.Time, 10) // buffer up to 10 values

	var before = time.Now()
	assert.NoErr(t, env.SyncClient().SetServerTimeListener(func(serverTime time.Time) { messages <- serverTime }))
	assert.NoErr(t, env.SyncClient().Start())

	select {
	case received := <-messages:
		var after = time.Now().Add(time.Second)
		t.Logf("Server time received %v\n", received)
		t.Logf("Checking if its between %v and %v\n", before, after)
		assert.True(t, received.Nanosecond() >= before.Nanosecond())
		assert.True(t, received.Nanosecond() <= after.Nanosecond())
	case <-time.After(env.defaultTimeout):
		assert.Failf(t, "Waiting for a server-time listener timed out after %v", env.defaultTimeout)
	}
}

func TestSyncChangeListener(t *testing.T) {
	skipTestIfSyncNotAvailable(t)

	var env = NewSyncTestEnv(t)

	var a = env.NamedClient("a")
	a.Start()

	var b = env.NamedClient("b")

	// make a couple of buffered channels
	var putIDs = make(chan uint64, 100)
	var removedIDs = make(chan uint64, 100)
	var messages = make(chan string, 10)

	assert.NoErr(t, b.sync.SetCompletionListener(func() { messages <- "sync-completed" }))
	assert.NoErr(t, b.sync.SetChangeListener(func(changes []*objectbox.SyncChange) {
		t.Logf("received %d changes", len(changes))
		for i, change := range changes {
			t.Logf("change %d: %v", i, change)

			// only count the main entity, not relations
			if change.EntityId == model.EntityBinding.Id {
				for j := 0; j < len(change.Puts); j++ {
					putIDs <- change.Puts[j]
				}
				for j := 0; j < len(change.Removals); j++ {
					removedIDs <- change.Removals[j]
				}
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

	assert.StringChannelExpect(t, "sync-completed", messages, env.defaultTimeout)

	// check expected IDs
	for id := uint(1); id <= count; id++ {
		select {
		case received := <-putIDs:
			assert.Eq(t, uint64(id), received)
		default:
			assert.Failf(t, "Didn't receive the expected id %v in putIDs", id)
		}
	}
}
