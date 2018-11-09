package modelinfo

import "fmt"

type Relation struct {
	Id   IdUid
	Name string
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (relation *Relation) Validate() error {
	if err := relation.Id.Validate(); err != nil {
		return err
	}

	if len(relation.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	return nil
}
