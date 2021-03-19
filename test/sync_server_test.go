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
	"bytes"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model"
)

// find a free (available) port to bind to
func findFreeTCPPort(t *testing.T) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoErr(t, err)

	listener, err := net.ListenTCP("tcp", addr)
	assert.NoErr(t, err)
	var port = listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

// testSyncServer wraps a sync-server binary and executes it.
// The binary must be present in the test folder or a known executable (available in $PATH).
type testSyncServer struct {
	t    *testing.T
	err  error
	cmd  *exec.Cmd
	env  *model.TestEnv
	port int
}

func NewTestSyncServer(t *testing.T) *testSyncServer {
	var execPath = findSyncServerExecutable(t)

	var server = &testSyncServer{
		t:    t,
		port: findFreeTCPPort(t),
	}

	// will be executed in case of error in this function
	var cleanup = server.Close
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()

	// prepare a database directory
	server.env = model.NewTestEnv(server.t)
	server.env.ObjectBox.Close() // close the database so that the server can open it

	server.cmd = exec.Command(execPath,
		"--unsecured-no-authentication",
		"--db-directory="+server.env.Directory,
		"--bind="+server.URI(),
		"--browser-bind=127.0.0.1:"+strconv.FormatInt(int64(findFreeTCPPort(t)), 10),
	)
	server.cmd.Stdout = &bytes.Buffer{}
	server.cmd.Stderr = &bytes.Buffer{}

	// start the server
	assert.NoErr(t, server.cmd.Start())

	// wait for the server to start listening for connections
	uri, err := url.Parse(server.URI())
	assert.NoErr(t, err)
	assert.NoErr(t, waitUntil(5*time.Second, func() (b bool, e error) {
		conn, err := net.DialTimeout("tcp", uri.Hostname()+":"+uri.Port(), 5*time.Second)

		// if connection was successful, stop waiting (return true)
		if err == nil {
			return true, conn.Close()
		}

		// if the connection was refused, try again next time
		if strings.Contains(err.Error(), "connection refused") {
			return false, nil
		}

		// fail immediately on other errors
		return false, err
	}))

	cleanup = nil // no error, don't close the server
	return server
}

func findSyncServerExecutable(t *testing.T) string {
	// sync-server executable must be located in the `test` directory (CWD when executing tests) or available in $PATH
	const executable = "sync-server"
	var path = "./" + executable

	// check if the executable exists in the CWD - that one will have preference
	if _, err := os.Stat(executable); err != nil {
		if !os.IsNotExist(err) {
			assert.NoErr(t, err)
		}

		// if not found in CWD, try to look up in $PATH
		path, err = exec.LookPath(executable)
		if err != nil {
			cwd, err := os.Getwd()
			assert.NoErr(t, err)

			t.Skipf("%s executable not found in %v and is not available in $PATH either", executable, cwd)
		}
	}
	return path
}

func (server *testSyncServer) Close() {
	if server.env == nil {
		return
	}

	defer func() {
		server.env.Close()
		server.env = nil
		server.cmd = nil
	}()

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
	return "ws://127.0.0.1:" + strconv.FormatInt(int64(server.port), 10)
}

func TestSyncServer(t *testing.T) {
	var server = NewTestSyncServer(t)
	server.Close()
}
