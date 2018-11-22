// Package generator provides tools to generate ObjectBox entity bindings between GO structs & ObjectBox schema
package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"

	"github.com/objectbox/objectbox-go/internal/generator/templates"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

func BindingFileName(sourceFile string) string {
	var extension = path.Ext(sourceFile)
	return sourceFile[0:len(sourceFile)-len(extension)] + "-binding" + extension
}

// Process is the main API method of the package
// it takes source file & model-information file paths and generates bindings (as a sibling file to the source file)
func Process(sourceFile, modelInfoFile string) (err error) {
	var err2 error

	var f *file
	if f, err = parseFile(sourceFile); err != nil {
		return fmt.Errorf("can't parse GO file %s: %s", sourceFile, err)
	}

	var modelInfo *modelinfo.ModelInfo
	if modelInfo, err = modelinfo.LoadOrCreateModel(modelInfoFile); err != nil {
		return fmt.Errorf("can't init ModelInfo: %s", err)
	} else {
		defer modelInfo.Close()
	}

	if err = modelInfo.Validate(); err != nil {
		return fmt.Errorf("invalid ModelInfo loaded: %s", err)
	}

	var binding *Binding
	if binding, err = newBinding(); err != nil {
		return fmt.Errorf("can't init Binding: %s", err)
	}

	if err = binding.createFromAst(f); err != nil {
		return fmt.Errorf("can't prepare bindings for %s: %s", sourceFile, err)
	}

	if err = mergeBindingWithModelInfo(binding, modelInfo); err != nil {
		return fmt.Errorf("can't merge binding model information: %s", err)
	}

	var bindingSource []byte
	if bindingSource, err = generateBinding(binding); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	var bindingFile = BindingFileName(sourceFile)
	if formattedSource, err := format.Source(bindingSource); err != nil {
		// we just store error but still writ the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated binding file %s: %s", bindingFile, err)
	} else {
		bindingSource = formattedSource
	}

	if err = writeBindingFile(sourceFile, bindingFile, bindingSource); err != nil {
		return fmt.Errorf("can't write binding file %s: %s", sourceFile, err)
	} else if err2 != nil {
		// now when the binding has been written (for debugging purposes), we can return the error
		return err2
	}

	if err = modelInfo.Write(); err != nil {
		return fmt.Errorf("can't write model-info file %s: %s", modelInfoFile, err)
	}

	return nil
}

func generateBinding(binding *Binding) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err = templates.BindingTemplate.Execute(writer, binding); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func writeBindingFile(sourceFile string, bindingFile string, data []byte) (err error) {
	var perm os.FileMode
	// copy permissions either from the existing bindings file or from the source file
	if info, _ := os.Stat(bindingFile); info != nil {
		perm = info.Mode()
	} else if info, err := os.Stat(sourceFile); info != nil {
		perm = info.Mode()
	} else {
		return err
	}

	return ioutil.WriteFile(bindingFile, data, perm)
}
