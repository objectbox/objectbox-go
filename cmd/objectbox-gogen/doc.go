/*
Generates objectbox related code for ObjectBox entities (Go structs)

It can be used by adding `//go:generate objectbox-gogen` comment inside a .go file
containing the struct that you want to persist and executing `go generate` in the module


Alternatively, you can run the command manually:

	objectbox-gogen [flags]


The flags are

  -byValue
        getters should return a struct value (a copy) instead of a struct pointer
  -persist string
        path to the model information persistence file
  -source string
        path to the source file containing structs to process



To learn more about different configuration and annotations for entities, see docs at https://golang.objectbox.io/
*/
package main
