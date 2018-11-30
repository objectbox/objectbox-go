/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
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

package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/objectbox/objectbox-go/objectbox"
)

type TestEnv struct {
	ObjectBox *objectbox.ObjectBox
	Box       *EntityBox

	t      *testing.T
	dbName string
}

func removeDb(name string) {
	os.Remove(filepath.Join(name, "data.mdb"))
	os.Remove(filepath.Join(name, "lock.mdb"))
}

func NewTestEnv(t *testing.T) *TestEnv {
	var dbName = "testdata"

	removeDb(dbName)

	ob, err := objectbox.NewBuilder().Directory(dbName).Model(ObjectBoxModel()).Build()
	if err != nil {
		t.Fatal(err)
	}
	return &TestEnv{
		ObjectBox: ob,
		Box:       BoxForEntity(ob),
		dbName:    dbName,
		t:         t,
	}
}

func (env *TestEnv) Close() {
	err := env.Box.Close()
	env.ObjectBox.Close()

	if err != nil {
		env.t.Error(err)
	}

	removeDb(env.dbName)
}
