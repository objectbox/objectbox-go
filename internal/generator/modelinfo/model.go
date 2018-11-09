package modelinfo

import (
	"fmt"
	"os"
)

type id = uint32
type uid = uint64

type ModelInfo struct {
	Comment      []string
	Entities     []*Entity
	LastEntityId IdUid
	LastIndexId  IdUid
	//ModelVersion        int
	//Version             int
	RetiredEntityUids   []uid
	RetiredIndexUids    []uid
	RetiredPropertyUids []uid

	file *os.File // file handle, locked while the model is open
}

func createModelInfo() *ModelInfo {
	return &ModelInfo{
		Comment: []string{
			"KEEP THIS FILE! Check it into a version control system (VCS) like git.",
			"ObjectBox manages crucial IDs for your object model. See docs for details.",
			"If you have VCS merge conflicts, you must resolve them according to ObjectBox docs.",
		},
		Entities:            make([]*Entity, 0),
		RetiredEntityUids:   make([]uid, 0),
		RetiredIndexUids:    make([]uid, 0),
		RetiredPropertyUids: make([]uid, 0),
	}
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (model *ModelInfo) Validate() (err error) {
	if model.Entities == nil {
		return fmt.Errorf("entities are not defined or not an array")
	}

	for _, entity := range model.Entities {
		err = entity.Validate()
		if err != nil {
			return fmt.Errorf("entity %s %s is invalid: %s", entity.Name, string(entity.Id), err)
		}
	}

	if len(model.Entities) > 0 {
		if err = model.LastEntityId.Validate(); err != nil {
			return fmt.Errorf("lastEntityId: %s", err)
		}

		var lastId = model.LastEntityId.getIdSafe()
		var lastUid = model.LastEntityId.getUidSafe()

		var found = false
		for _, entity := range model.Entities {
			if lastId == entity.Id.getIdSafe() {
				if lastUid != entity.Id.getUidSafe() {
					return fmt.Errorf("lastEntityId %s doesn't match entity %s %s",
						model.LastEntityId, entity.Name, entity.Id)
				}
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("lastEntityId %s doesn't match any entity", model.LastEntityId)
		}
	}

	if len(model.LastIndexId) > 0 {
		if err = model.LastIndexId.Validate(); err != nil {
			return fmt.Errorf("lastEntityId: %s", err)
		}
	}

	if model.RetiredEntityUids == nil {
		return fmt.Errorf("retiredEntityUids are not defined or not an array")
	}

	if model.RetiredIndexUids == nil {
		return fmt.Errorf("retiredIndexUids are not defined or not an array")
	}

	if model.RetiredPropertyUids == nil {
		return fmt.Errorf("retiredPropertyUids are not defined or not an array")
	}

	return nil
}

func (model *ModelInfo) FindEntityByUid(uid uid) (*Entity, error) {
	for _, entity := range model.Entities {
		entityUid, _ := entity.Id.GetUid()
		if entityUid == uid {
			return entity, nil
		}
	}

	return nil, fmt.Errorf("entity with Uid %d not found", uid)
}

func (model *ModelInfo) FindEntityByName(name string) (*Entity, error) {
	for _, entity := range model.Entities {
		if entity.Name == name {
			return entity, nil
		}
	}

	return nil, fmt.Errorf("entity with Name %s not found", name)
}

func (model *ModelInfo) CreateEntity() (*Entity, error) {
	var id id = 1
	if len(model.Entities) > 0 {
		id = model.LastEntityId.getIdSafe() + 1
	}

	// generate a unique UID
	uniqueUid, err := generateUid(func(uid uid) bool {
		item, err := model.FindEntityByUid(uid)
		return item == nil && err != nil
	})

	if err != nil {
		return nil, err
	}

	var entity = CreateEntity(id, uniqueUid)

	model.Entities = append(model.Entities, entity)
	model.LastEntityId = entity.Id

	return entity, nil
}
