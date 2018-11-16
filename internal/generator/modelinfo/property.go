package modelinfo

import "fmt"

type Property struct {
	Id      IdUid  `json:"id"`
	Name    string `json:"name"`
	IndexId *IdUid `json:"indexId,omitempty"`
}

func CreateProperty(id id, uid uid) *Property {
	return &Property{
		Id: CreateIdUid(id, uid),
	}
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (property *Property) Validate() error {
	if err := property.Id.Validate(); err != nil {
		return err
	}

	if property.IndexId != nil {
		if err := property.IndexId.Validate(); err != nil {
			return fmt.Errorf("indexId: %s", err)
		}
	}

	if len(property.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	return nil
}
