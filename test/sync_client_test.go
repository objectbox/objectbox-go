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
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsJwtId("{}")))
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsJwtAccess("{}")))
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsJwtRefresh("{}")))
	assert.NoErr(t, client.SetCredentials(objectbox.SyncCredentialsJwtCustom("{}")))
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

func populateSyncedBox(t *testing.T, box *model.TestEntitySyncedBox, count uint) {
	for i := uint(0); i < count; i++ {
		_, err := box.Put(&model.TestEntitySynced{Name: "foo"})
		assert.NoErr(t, err)
	}
}

func TestSyncDataAutomatic(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var a = env.NamedClient("a")
	a.Start()

	var b = env.NamedClient("b")
	b.Start()

	var aBox = model.BoxForTestEntitySynced(a.env.ObjectBox)
	isEmpty, err := aBox.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	var bBox = model.BoxForTestEntitySynced(b.env.ObjectBox)
	isEmpty, err = bBox.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	// insert into one box
	var count uint = 10
	populateSyncedBox(t, aBox, count)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := bBox.Count()
		return bCount == uint64(count), err
	}))

	var assertEqualBoxes = func(boxA, boxB *model.TestEntitySyncedBox) {
		itemsA, err := aBox.GetAll()
		assert.NoErr(t, err)

		itemsB, err := bBox.GetAll()
		assert.NoErr(t, err)

		assert.Eq(t, count, uint(len(itemsA)))
		assert.Eq(t, count, uint(len(itemsB)))
		assert.Eq(t, itemsA, itemsB)
	}
	assertEqualBoxes(aBox, bBox)

	// remove from one of the boxes
	removed, err := bBox.RemoveIds(1, 3, 6)
	assert.NoErr(t, err)
	assert.True(t, 3 == removed)
	count = count - uint(removed)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := aBox.Count()
		return bCount == uint64(count), err
	}))

	assertEqualBoxes(aBox, bBox)
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

	var aBox = model.BoxForTestEntitySynced(a.env.ObjectBox)
	isEmpty, err := aBox.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	var bBox = model.BoxForTestEntitySynced(b.env.ObjectBox)
	isEmpty, err = bBox.IsEmpty()
	assert.NoErr(t, err)
	assert.True(t, isEmpty)

	// insert into one box
	var count uint = 10
	populateSyncedBox(t, aBox, count)

	// this will time out because we haven't manually initiated an update
	assert.True(t, strings.Contains(waitUntil(500*time.Millisecond, func() (bool, error) {
		bCount, err := bBox.Count()
		return bCount == uint64(count), err
	}).Error(), "timeout"))

	// manually trigger the data synchronization
	assert.NoErr(t, b.sync.RequestUpdates(false))

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := bBox.Count()
		return bCount == uint64(count), err
	}))

	assert.NoErr(t, aBox.RemoveAll())
	count = 0

	// this will time out because we haven't subscribed for all further updates
	assert.True(t, strings.Contains(waitUntil(500*time.Millisecond, func() (bool, error) {
		bCount, err := bBox.Count()
		return bCount == uint64(count), err
	}).Error(), "timeout"))

	// subscribe for further updates
	assert.NoErr(t, b.sync.RequestUpdates(true))

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := bBox.Count()
		return bCount == uint64(count), err
	}))

	count = 10
	populateSyncedBox(t, aBox, count)

	// wait for the data to be synced to the other box
	assert.NoErr(t, waitUntil(time.Second, func() (bool, error) {
		bCount, err := bBox.Count()
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

func TestSyncJwtLoginValidToken(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10)

	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

	// Valid JWT token
	var jwtToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJzeW5jLXNlcnZlciIsImlzcyI6Im9iamVjdGJveC1hdXRoIn0.YZSt5XIp7KLSIEtYegEGInea2IvyZajEOWEXcH8p0kYTvhU07LFcxbPWxnNeBtQSjkGp0U0XQUQkCaRjRbNDiHKHCtQHOsUtLefAfQc-WENzSSrGqbb7YKw7FHgsGCQX7FRblcdw3ExU9w8NBgt0xQaDqnwBYfltfu6bmJG5QabGnljcmLGB3Q5EcppxBgWZdLzhmVRiqkiIsCp8kBtELz3Lk8a2LIJP80khJWdls1zIK_NR0XtV6Dbbac1fFN0v5F2VN61VjL9HXZWm68zf2ueW_jobN8IBcJkOAfefgsQu_1e5B0iVAxyRki6F99V1B8Ci_5wbTXRs4bob1Nsl2Q"
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
	assert.NoErr(t, env.SyncClient().SetCredentials(objectbox.SyncCredentialsJwtId(jwtToken)))
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "success", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
}

func TestSyncJwtLoginInvalidToken(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10)

	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

	// Invalid JWT token (expired)
	var expiredJwtToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJzeW5jLXNlcnZlciIsImlzcyI6Im9iamVjdGJveC1hdXRoIiwiZXhwIjoxNzM4MjE1NjAwLCJpYXQiOjE3MzgyMTc0MDN9.3auqtgaSEqpFqXhuCyoDM-LbfTOIEGGF6X0AjCcykJ2Nv1WN6LaVbuMDjMf-tKSLyeqFkzQbIckP4FvLHh7wQJ6rafDiT4H2pb6xhouU1QH3szK2S_7VDl_4BhxRbW5pEUt9086HXaVFHEZVS0417pxomlPHxrc1n4Z_A4QxZM5_xh5xcHV8PiGgXWb6_2basjBj5z6POTrazRs67IOQ-ob6ROIsOUGu3om6b8i0h_QSMmeJbujfr2EZqhYWTKijeyidbjRWZ97NFxtGRYN_jPOvy-T3gANXs2a32Er8XvgZTjr_-O8tl_1fHPo2kDE-UCNdwUfBQFhTokDUdJ81bg"
	assert.NoErr(t, env.SyncClient().SetCredentials(objectbox.SyncCredentialsJwtId(expiredJwtToken)))
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "failure 43", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
}

// func TestSyncMultipleCredentialsSuccess(t *testing.T) {
// 	var env = NewSyncTestEnv(t)
// 	defer env.Close()

// 	var messages = make(chan string, 10)

// 	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
// 	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

// 	// Valid JWT token (valid)
// 	var jwtToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJzeW5jLXNlcnZlciIsImlzcyI6Im9iamVjdGJveC1hdXRoIn0.YZSt5XIp7KLSIEtYegEGInea2IvyZajEOWEXcH8p0kYTvhU07LFcxbPWxnNeBtQSjkGp0U0XQUQkCaRjRbNDiHKHCtQHOsUtLefAfQc-WENzSSrGqbb7YKw7FHgsGCQX7FRblcdw3ExU9w8NBgt0xQaDqnwBYfltfu6bmJG5QabGnljcmLGB3Q5EcppxBgWZdLzhmVRiqkiIsCp8kBtELz3Lk8a2LIJP80khJWdls1zIK_NR0XtV6Dbbac1fFN0v5F2VN61VjL9HXZWm68zf2ueW_jobN8IBcJkOAfefgsQu_1e5B0iVAxyRki6F99V1B8Ci_5wbTXRs4bob1Nsl2Q"
// 	// Shared secret (valid)
// 	var invalidSharedSecret = "shared-secret"
// 	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
// 	assert.NoErr(t, env.SyncClient().SetMultipleCredentials([]*objectbox.SyncCredentials{
// 		objectbox.SyncCredentialsJwtId([]byte(jwtToken)),
// 		objectbox.SyncCredentialsSharedSecret([]byte(invalidSharedSecret)),
// 	}))
// 	assert.NoErr(t, env.SyncClient().Start())
// 	assert.StringChannelExpect(t, "success", messages, env.defaultTimeout)
// 	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
// }

func TestSyncMultipleCredentialsFailure(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10)

	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

	// Valid JWT token (valid)
	var jwtToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJzeW5jLXNlcnZlciIsImlzcyI6Im9iamVjdGJveC1hdXRoIn0.YZSt5XIp7KLSIEtYegEGInea2IvyZajEOWEXcH8p0kYTvhU07LFcxbPWxnNeBtQSjkGp0U0XQUQkCaRjRbNDiHKHCtQHOsUtLefAfQc-WENzSSrGqbb7YKw7FHgsGCQX7FRblcdw3ExU9w8NBgt0xQaDqnwBYfltfu6bmJG5QabGnljcmLGB3Q5EcppxBgWZdLzhmVRiqkiIsCp8kBtELz3Lk8a2LIJP80khJWdls1zIK_NR0XtV6Dbbac1fFN0v5F2VN61VjL9HXZWm68zf2ueW_jobN8IBcJkOAfefgsQu_1e5B0iVAxyRki6F99V1B8Ci_5wbTXRs4bob1Nsl2Q"
	// Shared secret (invalid)
	var invalidSharedSecret = "invalid-shared-secret"
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
	assert.NoErr(t, env.SyncClient().SetMultipleCredentials([]*objectbox.SyncCredentials{
		objectbox.SyncCredentialsJwtId(jwtToken),
		objectbox.SyncCredentialsSharedSecret([]byte(invalidSharedSecret)),
	}))
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "failure 43", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
}

func TestSyncUserPasswordLogin(t *testing.T) {
	var env = NewSyncTestEnv(t)
	defer env.Close()

	var messages = make(chan string, 10)

	assert.NoErr(t, env.SyncClient().SetLoginListener(func() { messages <- "success" }))
	assert.NoErr(t, env.SyncClient().SetLoginFailureListener(func(code objectbox.SyncLoginFailure) { messages <- "failure " + strconv.FormatUint(uint64(code), 10) }))

	// Valid JWT token
	var user = ""
	var password = ""
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
	assert.NoErr(t, env.SyncClient().SetCredentials(objectbox.SyncCredentialsUsernamePassword(user, password)))
	assert.NoErr(t, env.SyncClient().Start())
	assert.StringChannelExpect(t, "success", messages, env.defaultTimeout)
	assert.StringChannelMustTimeout(t, messages, env.defaultTimeout/10)
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

			assert.Eq(t, model.TestEntitySyncedBinding.Id, change.EntityId)

			for j := 0; j < len(change.Puts); j++ {
				putIDs <- change.Puts[j]
			}
			for j := 0; j < len(change.Removals); j++ {
				removedIDs <- change.Removals[j]
			}

		}
	}))
	b.Start()

	// insert on one client
	var count uint = 100
	var aBox = model.BoxForTestEntitySynced(a.env.ObjectBox)
	populateSyncedBox(t, aBox, count)

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
