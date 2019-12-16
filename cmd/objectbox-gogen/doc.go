/*
Generates objectbox related code for ObjectBox entities (Go structs)

It can be used by adding `//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen` comment inside a .go file
containing the struct that you want to persist and executing `go generate` in the module


Alternatively, you can run the command manually:

	objectbox-gogen [flags] [path-pattern]
		to generate the binding code

or

	objectbox-gogen clean [path-pattern]
		to remove the generated files instead of creating them - this removes *.obx.go and objectbox-model.go but keeps objectbox-model.json

path-pattern:
  * a path or a valid path pattern as accepted by the go tool (e.g. ./...)
  * if not given, the generator expects GOFILE environment variable to be set

Available flags:
  -byValue
        getters should return a struct value (a copy) instead of a struct pointer
  -help
        print this help
  -persist string
        path to the model information persistence file
  -source string
        @deprecated, equivalent to passing the given source file path as as the path-pattern argument
  -version
        print the generator version info



To learn more about different configuration and annotations for entities, see docs at https://golang.objectbox.io/
*/
package main
