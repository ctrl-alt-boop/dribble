// Package memory contain framework for in-memory database
package memory

import (
	"reflect"

	"github.com/google/uuid"
)

// Store is an in-memory store for simple caching
var Store = &store{
	pathStore:   newPathStore(),
	bucketStore: newBucketStore(),
	treeStore:   newTree(),
}

// not sure which I really want to use
type store struct {
	pathStore   *pathStore
	bucketStore *bucketStore
	treeStore   *treeStore
}

// Item is a memory store item
type Item struct {
	uuid  uuid.UUID
	name  string
	value any
}

// Value .
func (i Item) Value() any {
	return i.value
}

// Name .
func (i Item) Name() string {
	return i.name
}

var CombiStore = newCombiStore()

type combiStore struct {
	store       map[uuid.UUID]*Item
	typeToUUIDs map[reflect.Type][]uuid.UUID
}

func newCombiStore() *combiStore {
	return &combiStore{
		store:       make(map[uuid.UUID]*Item),
		typeToUUIDs: make(map[reflect.Type][]uuid.UUID),
	}
}

func (c *combiStore) Add(item *Item) uuid.UUID {
	item.uuid = uuid.New()
	c.store[item.uuid] = item
	itemType := reflect.TypeOf(item.value)
	c.typeToUUIDs[itemType] = append(c.typeToUUIDs[itemType], item.uuid)
	return item.uuid
}

func (c *combiStore) Get(uuid uuid.UUID) (*Item, bool) {
	item, ok := c.store[uuid]
	return item, ok
}

func (c *combiStore) GetByType(itemType reflect.Type) ([]*Item, bool) {
	uuids, ok := c.typeToUUIDs[itemType]
	if !ok {
		return nil, false
	}
	items := make([]*Item, len(uuids))
	for i, uuid := range uuids {
		items[i] = c.store[uuid]
	}
	return items, true
}
