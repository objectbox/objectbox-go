package generator

import (
	"html/template"
	"os"
	"strings"

	// TODO check whether we can use this dependency
	// used to include template in the compiled binary
	"github.com/gobuffalo/packr"
)

func Process(filePath string) (err error) {
	var f *file
	if f, err = parseFile(filePath); err != nil {
		return err
	}

	var binding *Binding
	if binding, err = newBinding(); err != nil {
		return err
	}

	if err = binding.loadAstFile(f); err != nil {
		return err
	}

	if err = generateBindingFile(binding); err != nil {
		return err
	}

	return nil
}

func generateBindingFile(binding *Binding) (err error) {
	// load the template for the binding
	box := packr.NewBox("./templates")

	var tplText string
	if tplText, err = box.MustString("binding.tmpl"); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"StringTitle": strings.Title,
	}

	tpl := template.Must(template.New("binding").Funcs(funcMap).Parse(tplText))
	if err = tpl.Execute(os.Stdout, binding); err != nil { // TODO replace os.Stdout with file writer
		return err
	}

	return nil
}
