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
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
	"testing"
)

func TestSyncAuth(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.AuthSharedSecret(nil))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.AuthSharedSecret([]byte{1, 2, 3}))
		assert.NoErr(t, client.Close())
	}
}

func TestSyncUpdatesMode(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesManual))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesAutomatic))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesOnLogin))
		assert.NoErr(t, client.Close())
	}
}

func TestSyncState(t *testing.T) {
	var env = model.NewTestEnv(t)
	defer env.Close()

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesManual))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesAutomatic))
		assert.NoErr(t, client.Close())
	}

	{
		var client = env.SyncClient()
		assert.NoErr(t, client.UpdatesMode(objectbox.SyncClientUpdatesOnLogin))
		assert.NoErr(t, client.Close())
	}
}
