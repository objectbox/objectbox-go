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
