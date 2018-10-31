package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// TODO check whether we can use this dependency
	// used to include template in the compiled binary
	"github.com/gobuffalo/packr"
)

func Process(sourceFile string) (err error) {
	var f *file
	if f, err = parseFile(sourceFile); err != nil {
		return fmt.Errorf("can't parse GO file %s: %s", sourceFile, err)
	}

	var binding *Binding
	if binding, err = newBinding(); err != nil {
		return err
	}

	if err = binding.loadAstFile(f); err != nil {
		return fmt.Errorf("can't prepare bindings for %s: %s", sourceFile, err)
	}

	var bindingSource []byte
	if bindingSource, err = generateBinding(binding); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	var extension = path.Ext(sourceFile)
	var bindingFile = sourceFile[0:len(sourceFile)-len(extension)] + "binding" + extension

	if err = writeBindingFile(sourceFile, bindingFile, bindingSource); err != nil {
		return fmt.Errorf("can't write binding file %s: %s", sourceFile, err)
	}

	return nil
}

func generateBinding(binding *Binding) (data []byte, err error) {
	// load the template for the binding
	box := packr.NewBox("./templates")

	var tplText string
	if tplText, err = box.MustString("binding.tmpl"); err != nil {
		return nil, fmt.Errorf("can't load template from the binary distribution: %s", err)
	}

	funcMap := template.FuncMap{
		"StringTitle": strings.Title,
	}

	// prepare the template
	tpl := template.Must(template.New("binding").Funcs(funcMap).Parse(tplText))

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err = tpl.Execute(writer, binding); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	// format the bindings source (same as gofmt)
	if data, err = format.Source(b.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to format generated binding file: %s", err)
	}

	return data, nil
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
