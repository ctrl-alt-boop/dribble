package memory

import "strings"

// pathStore is a simple store keyed by paths
type pathStore struct {
	store map[string]Item
}

func newPathStore() *pathStore {
	return &pathStore{
		store: make(map[string]Item),
	}
}

func (c *pathStore) Pop(path ...string) (Item, bool) {
	key := strings.Join(path, "/")
	value, ok := c.store[key]
	if ok {
		delete(c.store, key)
	}
	return value, ok
}

func (c *pathStore) Del(path ...string) {
	key := strings.Join(path, "/")
	delete(c.store, key)
}

func (c *pathStore) Add(value Item, path ...string) {
	key := strings.Join(path, "/")
	c.store[key] = value
}

func (c *pathStore) Get(path ...string) (Item, bool) {
	key := strings.Join(path, "/")
	value, ok := c.store[key]
	return value, ok
}

func (c *pathStore) Has(path ...string) bool {
	key := strings.Join(path, "/")
	_, ok := c.store[key]
	return ok
}

func (c *pathStore) Clear() {
	c.store = make(map[string]Item)
}

func (c *pathStore) Size() int {
	return len(c.store)
}
