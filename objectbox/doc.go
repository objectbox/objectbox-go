/*
Package objectbox provides a super-fast, light-weight object persistence framework.

You can define your entity as a standard .go struct, with a comment signalling to generate ObjectBox code

	//go:generate objectbox-gogen

	type Person struct {
	   Id        uint64 `id`
	   FirstName string
	   LastName  string
	}


Now, just init ObjectBox using the generated code (don't forget to errors in your real code, they are discarded here to keep the example concise)

	ob, _ := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	defer ob.Close()

	box := BoxForPerson(ob)
	defer box.Close()

	// Create
	id, _ := box.Put(&Person{
	   FirstName: "Joe",
	   LastName:  "Green",
	})

	// Read
	person, _ := box.Get(id)

	// Update
	person.LastName = "Black"
	box.Put(person)

	// Delete
	box.Remove(person)


To learn more, see https://golang.objectbox.io/
*/
package objectbox
