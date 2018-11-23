package modelinfo

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type id = uint32
type uid = uint64

type ModelInfo struct {
	Comment      []string  `json:"comment"`
	Entities     []*Entity `json:"entities"`
	LastEntityId IdUid     `json:"lastEntityId"`
	LastIndexId  IdUid     `json:"lastIndexId"`
	//ModelVersion        int
	//Version             int
	RetiredEntityUids   []uid `json:"retiredEntityUids"`
	RetiredIndexUids    []uid `json:"retiredIndexUids"`
	RetiredPropertyUids []uid `json:"retiredPropertyUids"`

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
		if entity.model == nil {
			entity.model = model
		} else if entity.model != model {
			return fmt.Errorf("entity %s %s has incorrect parent model reference", entity.Name, entity.Id)
		}

		err = entity.Validate()
		if err != nil {
			return fmt.Errorf("entity %s %s is invalid: %s", entity.Name, entity.Id, err)
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
			} else if lastId < entity.Id.getIdSafe() {
				return fmt.Errorf("lastEntityId %s is lower than entity %s %s",
					model.LastEntityId, entity.Name, entity.Id)
			}
		}

		if !found && !searchSliceUid(model.RetiredEntityUids, lastUid) {
			return fmt.Errorf("lastEntityId %s doesn't match any entity", model.LastEntityId)
		}
	}

	if len(model.LastIndexId) > 0 {
		if err = model.LastIndexId.Validate(); err != nil {
			return fmt.Errorf("lastIndexId: %s", err)
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
		if strings.ToLower(entity.Name) == strings.ToLower(name) {
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

	uniqueUid, err := model.generateUid()

	if err != nil {
		return nil, err
	}

	var entity = CreateEntity(model, id, uniqueUid)

	model.Entities = append(model.Entities, entity)
	model.LastEntityId = entity.Id

	return entity, nil
}

func (model *ModelInfo) generateUid() (result uid, err error) {
	result = 0

	for i := 0; i < 1000; i++ {
		t := uid(rand.Int63())
		if !model.containsUid(t) {
			result = t
			break
		}
	}

	if result == 0 {
		err = fmt.Errorf("internal error = could not generate a unique UID")
	}

	return result, err
}

func (model *ModelInfo) createIndex() (IdUid, error) {
	var id id = 1
	if len(model.LastIndexId) > 0 {
		id = model.LastIndexId.getIdSafe() + 1
	}

	uniqueUid, err := model.generateUid()

	if err != nil {
		return "", err
	}

	model.LastIndexId = CreateIdUid(id, uniqueUid)
	return model.LastIndexId, nil
}

// recursively checks whether given UID is present in the model
func (model *ModelInfo) containsUid(searched uid) bool {
	if model.LastEntityId.getUidSafe() == searched {
		return true
	}

	if model.LastIndexId.getUidSafe() == searched {
		return true
	}

	if searchSliceUid(model.RetiredEntityUids, searched) {
		return true
	}

	if searchSliceUid(model.RetiredIndexUids, searched) {
		return true
	}

	if searchSliceUid(model.RetiredPropertyUids, searched) {
		return true
	}

	for _, entity := range model.Entities {
		if entity.containsUid(searched) {
			return true
		}
	}

	return false
}

// the passed slices are not too large so let's just do linear search
func searchSliceUid(slice []uid, searched uid) bool {
	for _, i := range slice {
		if i == searched {
			return true
		}
	}

	return false
}
