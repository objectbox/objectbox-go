// Code generated by ObjectBox; DO NOT EDIT.
// Learn more about defining entities and generating this file - visit https://golang.objectbox.io/entity-annotations

package perf

import (
	"errors"
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type entity_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var EntityBinding = entity_EntityInfo{
	Entity: objectbox.Entity{
		Id: 1,
	},
	Uid: 1737161401460991620,
}

// Entity_ contains type-based Property helpers to facilitate some common operations such as Queries.
var Entity_ = struct {
	ID      *objectbox.PropertyUint64
	Int32   *objectbox.PropertyInt32
	Int64   *objectbox.PropertyInt64
	String  *objectbox.PropertyString
	Float64 *objectbox.PropertyFloat64
}{
	ID: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &EntityBinding.Entity,
		},
	},
	Int32: &objectbox.PropertyInt32{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &EntityBinding.Entity,
		},
	},
	Int64: &objectbox.PropertyInt64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &EntityBinding.Entity,
		},
	},
	String: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     4,
			Entity: &EntityBinding.Entity,
		},
	},
	Float64: &objectbox.PropertyFloat64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     5,
			Entity: &EntityBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (entity_EntityInfo) GeneratorVersion() int {
	return 4
}

// AddToModel is called by ObjectBox during model build
func (entity_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("Entity", 1, 1737161401460991620)
	model.Property("ID", 6, 1, 7373286741377356014)
	model.PropertyFlags(1)
	model.Property("Int32", 5, 2, 4837914178321008766)
	model.Property("Int64", 6, 3, 3841825182616422591)
	model.Property("String", 9, 4, 6473251296493454829)
	model.Property("Float64", 8, 5, 8933082277725371577)
	model.EntityLastPropertyId(5, 8933082277725371577)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (entity_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*Entity).ID, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (entity_EntityInfo) SetId(object interface{}, id uint64) {
	object.(*Entity).ID = id
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (entity_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (entity_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*Entity)
	var offsetString = fbutils.CreateStringOffset(fbb, obj.String)

	// build the FlatBuffers object
	fbb.StartObject(5)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetInt32Slot(fbb, 1, obj.Int32)
	fbutils.SetInt64Slot(fbb, 2, obj.Int64)
	fbutils.SetUOffsetTSlot(fbb, 3, offsetString)
	fbutils.SetFloat64Slot(fbb, 4, obj.Float64)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (entity_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'Entity' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}
	var id = table.GetUint64Slot(4, 0)

	return &Entity{
		ID:      id,
		Int32:   fbutils.GetInt32Slot(table, 6),
		Int64:   fbutils.GetInt64Slot(table, 8),
		String:  fbutils.GetStringSlot(table, 10),
		Float64: fbutils.GetFloat64Slot(table, 12),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (entity_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*Entity, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (entity_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*Entity), nil)
	}
	return append(slice.([]*Entity), object.(*Entity))
}

// Box provides CRUD access to Entity objects
type EntityBox struct {
	*objectbox.Box
}

// BoxForEntity opens a box of Entity objects
func BoxForEntity(ob *objectbox.ObjectBox) *EntityBox {
	return &EntityBox{
		Box: ob.InternalBox(1),
	}
}

// Put synchronously inserts/updates a single object.
// In case the ID is not specified, it would be assigned automatically (auto-increment).
// When inserting, the Entity.ID property on the passed object will be assigned the new ID as well.
func (box *EntityBox) Put(object *Entity) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the Entity.Id property on the passed object will be assigned the new ID as well.
func (box *EntityBox) Insert(object *Entity) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *EntityBox) Update(object *Entity) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *EntityBox) PutAsync(object *Entity) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case IDs are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the Entity.ID property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the Entity.ID assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *EntityBox) PutMany(objects []*Entity) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *EntityBox) Get(id uint64) (*Entity, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*Entity), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *EntityBox) GetMany(ids ...uint64) ([]*Entity, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*Entity), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *EntityBox) GetManyExisting(ids ...uint64) ([]*Entity, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*Entity), nil
}

// GetAll reads all stored objects
func (box *EntityBox) GetAll() ([]*Entity, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*Entity), nil
}

// Remove deletes a single object
func (box *EntityBox) Remove(object *Entity) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *EntityBox) RemoveMany(objects ...*Entity) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.ID
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the Entity_ struct to create conditions.
// Keep the *EntityQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *EntityBox) Query(conditions ...objectbox.Condition) *EntityQuery {
	return &EntityQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the Entity_ struct to create conditions.
// Keep the *EntityQuery if you intend to execute the query multiple times.
func (box *EntityBox) QueryOrError(conditions ...objectbox.Condition) (*EntityQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &EntityQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See EntityAsyncBox for more information.
func (box *EntityBox) Async() *EntityAsyncBox {
	return &EntityAsyncBox{AsyncBox: box.Box.Async()}
}

// EntityAsyncBox provides asynchronous operations on Entity objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type EntityAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForEntity creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use EntityBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForEntity(ob *objectbox.ObjectBox, timeoutMs uint64) *EntityAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 1, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 1: %s" + err.Error())
	}
	return &EntityAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the Id property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *EntityAsyncBox) Put(object *Entity) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The Id property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *EntityAsyncBox) Insert(object *Entity) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *EntityAsyncBox) Update(object *Entity) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *EntityAsyncBox) Remove(object *Entity) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all Entity which ID is either 42 or 47:
// 		box.Query(Entity_.ID.In(42, 47)).Find()
type EntityQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *EntityQuery) Find() ([]*Entity, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*Entity), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *EntityQuery) Offset(offset uint64) *EntityQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *EntityQuery) Limit(limit uint64) *EntityQuery {
	query.Query.Limit(limit)
	return query
}
