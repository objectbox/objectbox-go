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

// Package generator provides tools to generate ObjectBox entity bindings between GO structs & ObjectBox schema
package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
	"github.com/objectbox/objectbox-go/internal/generator/templates"
)

const Version = 3

func BindingFile(sourceFile string) string {
	var extension = filepath.Ext(sourceFile)
	return sourceFile[0:len(sourceFile)-len(extension)] + ".obx" + extension
}

func ModelInfoFile(dir string) string {
	return filepath.Join(dir, "objectbox-model.json")
}

func ModelFile(modelInfoFile string) string {
	var extension = filepath.Ext(modelInfoFile)
	return modelInfoFile[0:len(modelInfoFile)-len(extension)] + ".go"
}

func isGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.go" || strings.HasSuffix(name, ".obx.go")
}

// Process is the main API method of the package
// it takes source file & model-information file paths and generates bindings (as a sibling file to the source file)
func Process(sourceFile string, options Options) error {
	var err error

	// if no random generator is provided, we create and seed a new one
	if options.Rand == nil {
		options.Rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	}

	if len(options.ModelInfoFile) == 0 {
		options.ModelInfoFile = ModelInfoFile(filepath.Dir(sourceFile))
	}

	var modelInfo *modelinfo.ModelInfo
	if modelInfo, err = modelinfo.LoadOrCreateModel(options.ModelInfoFile); err != nil {
		return fmt.Errorf("can't init ModelInfo: %s", err)
	} else {
		modelInfo.Rand = options.Rand
		defer modelInfo.Close()
	}

	if err = modelInfo.Validate(); err != nil {
		return fmt.Errorf("invalid ModelInfo loaded: %s", err)
	}

	// if the model is valid, upgrade it to the latest version
	modelInfo.MinimumParserVersion = modelinfo.ModelVersion
	modelInfo.ModelVersion = modelinfo.ModelVersion

	if err = createBinding(sourceFile, modelInfo, options); err != nil {
		return err
	}

	if err = createModel(options.ModelInfoFile, modelInfo); err != nil {
		return err
	}

	return nil
}

func createBinding(sourceFile string, modelInfo *modelinfo.ModelInfo, options Options) error {
	var err, err2 error
	var binding *Binding
	var f *file

	if f, err = parseFile(sourceFile); err != nil {
		return fmt.Errorf("can't parse GO file %s: %s", sourceFile, err)
	}

	if binding, err = newBinding(); err != nil {
		return fmt.Errorf("can't init Binding: %s", err)
	}

	if err = binding.createFromAst(f); err != nil {
		return fmt.Errorf("can't prepare bindings for %s: %s", sourceFile, err)
	}

	if err = mergeBindingWithModelInfo(binding, modelInfo); err != nil {
		return fmt.Errorf("can't merge binding model information: %s", err)
	}

	if err = modelInfo.CheckRelationCycles(); err != nil {
		return err
	}

	var bindingSource []byte
	if bindingSource, err = generateBindingFile(binding, options); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	var bindingFile = BindingFile(sourceFile)
	if formattedSource, err := format.Source(bindingSource); err != nil {
		// we just store error but still writ the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated binding file %s: %s", bindingFile, err)
	} else {
		bindingSource = formattedSource
	}

	if err = writeFile(bindingFile, bindingSource, sourceFile); err != nil {
		return fmt.Errorf("can't write binding file %s: %s", sourceFile, err)
	} else if err2 != nil {
		// now when the binding has been written (for debugging purposes), we can return the error
		return err2
	}

	return nil
}

func generateBindingFile(binding *Binding, options Options) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Binding          *Binding
		GeneratorVersion int
		Options          Options
	}{binding, Version, options}

	if err = templates.BindingTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

// writes data to targetFile, while using permissions either from the targetFile or permSource
func writeFile(file string, data []byte, permSource string) error {
	var perm os.FileMode
	// copy permissions either from the existing file or from the source file
	if info, _ := os.Stat(file); info != nil {
		perm = info.Mode()
	} else if info, err := os.Stat(permSource); info != nil {
		perm = info.Mode()
	} else {
		return err
	}

	return ioutil.WriteFile(file, data, perm)
}

func createModel(modelInfoFile string, modelInfo *modelinfo.ModelInfo) error {
	var err, err2 error

	if err = modelInfo.Write(); err != nil {
		return fmt.Errorf("can't write model-info file %s: %s", modelInfoFile, err)
	}

	var modelFile = ModelFile(modelInfoFile)
	var modelSource []byte

	if modelSource, err = generateModelFile(modelInfo); err != nil {
		return fmt.Errorf("can't generate model file %s: %s", modelFile, err)
	}

	if formattedSource, err := format.Source(modelSource); err != nil {
		// we just store error but still writ the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated model file %s: %s", modelFile, err)
	} else {
		modelSource = formattedSource
	}

	if err = writeFile(modelFile, modelSource, modelInfoFile); err != nil {
		return fmt.Errorf("can't write model file %s: %s", modelFile, err)
	} else if err2 != nil {
		// now when the model has been written (for debugging purposes), we can return the error
		return err2
	}

	return nil
}

func generateModelFile(model *modelinfo.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Model            *modelinfo.ModelInfo
		GeneratorVersion int
	}{model, Version}

	if err = templates.ModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

const recursionSuffix = "/..."

// Clean removes generated files in the given path.
// Removes *.obx.go and objectbox-model.go but keeps objectbox-model.json
func Clean(path string) error {
	var recursive bool

	// if it's a pattern
	if strings.HasSuffix(path, recursionSuffix) {
		recursive = true
		path = path[0:len(path)-len(recursionSuffix)] + "/*"
	} else {
		// if it's a directory
		if finfo, err := os.Stat(path); err == nil && finfo.IsDir() {
			path = path + "/*"
		}
	}

	matches, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, subpath := range matches {
		finfo, err := os.Stat(subpath)
		if err != nil {
			return err
		}

		if recursive && finfo.Mode().IsDir() {
			err = Clean(subpath + recursionSuffix)
		} else if finfo.Mode().IsRegular() && isGeneratedFile(subpath) {
			fmt.Printf("Removing %s\n", subpath)
			err = os.Remove(subpath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
