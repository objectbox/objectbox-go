package modelinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func LoadOrCreateModel(path string) (model *ModelInfo, err error) {
	if fileExists(path) {
		return loadModelFromJsonFile(path)
	} else {
		return createModelJsonFile(path)
	}
}

// Close and unlock model
func (model *ModelInfo) Close() error {
	return model.file.Close()
}

// Write current model data to file
func (model *ModelInfo) Write() error {
	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	if err = model.file.Truncate(0); err != nil {
		return err
	}

	if _, err := model.file.WriteAt(data, 0); err != nil {
		return err
	}

	if err = model.file.Sync(); err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func loadModelFromJsonFile(path string) (model *ModelInfo, err error) {
	model = &ModelInfo{}

	if model.file, err = os.OpenFile(path, os.O_RDWR, 0); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(io.Reader(model.file))

	if err == nil {
		err = json.Unmarshal(data, model)
	}

	if err != nil {
		defer model.Close()
		return nil, fmt.Errorf("can't read file %s: %s", path, err)
	} else {
		return model, nil
	}
}

func createModelJsonFile(path string) (model *ModelInfo, err error) {
	model = createModelInfo()

	// create a file handle so to have an exclusive access
	if model.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		return nil, err
	}

	return model, nil
}
