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
	"bytes"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

// testSyncServer wraps a sync-server binary and executes it.
// The binary must be present in the test folder in order for sync tests to work.
type testSyncServer struct {
	t   *testing.T
	err error
	cmd *exec.Cmd
	env *model.TestEnv
}

func NewTestSyncServer(t *testing.T) *testSyncServer {
	// sync-server executable must be located in the `test` directory (current directory when executing the tests)
	const executable = "sync-server"

	// check if the executable exists
	if _, err := os.Stat(executable); os.IsNotExist(err) {
		cwd, err := os.Getwd()
		assert.NoErr(t, err)

		t.Skipf("%s executable not found in %v", executable, cwd)
	} else {
		assert.NoErr(t, err)
	}

	var server = &testSyncServer{
		t: t,
	}

	// prepare a database directory
	server.env = model.NewTestEnv(server.t)
	server.env.ObjectBox.Close() // close the database so that the server can open it

	server.cmd = exec.Command("./"+executable,
		"--unsecure-no-authentication",
		"--db-directory", server.env.Directory,
		"--bind", server.URI(),
	)
	server.cmd.Stdout = &bytes.Buffer{}
	server.cmd.Stderr = &bytes.Buffer{}

	// start the server
	assert.NoErr(t, server.cmd.Start())

	// wait for the server to start listening for connections
	conn, err := net.DialTimeout("tcp", "golang.org:80", 5 * time.Second)
	assert.NoErr(t, err)
	assert.NoErr(t, conn.Close())

	return server
}

func (server *testSyncServer) Close() {
	defer server.env.Close()

	// wait for the server to finish
	assert.NotNil(server.t, server.cmd.Process)
	assert.NoErr(server.t, server.cmd.Process.Signal(os.Interrupt))
	var err = server.cmd.Wait()

	// print the output
	server.t.Log("sync-server output: \n" + server.cmd.Stdout.(*bytes.Buffer).String())

	if server.cmd.Stderr.(*bytes.Buffer).Len() > 0 {
		server.t.Log("sync-server errors: \n" + server.cmd.Stderr.(*bytes.Buffer).String())
		server.t.Fail()
	}

	assert.NoErr(server.t, err)
}

func (server *testSyncServer) URI() string {
	return "ws://127.0.0.1:9999"
}

func TestSyncServer(t *testing.T) {
	var server = NewTestSyncServer(t)
	server.Close()
}
