package util

import "strings"

type Cacheable interface {
	Value() any
	Name() string
}

type CachableString string

func (s CachableString) Value() any {
	return string(s)
}

func (s CachableString) Name() string {
	return string(s)
}

type SimpleCacheable struct {
	name  string
	value any
}

func (s SimpleCacheable) Value() any {
	return s.value
}

func (s SimpleCacheable) Name() string {
	return s.name
}

func NewCacheable(name string, value any) Cacheable {
	return SimpleCacheable{
		name:  name,
		value: value,
	}
}

type Cache struct {
	cache       map[string]Cacheable
	currentPath []string
}

func NewCache() Cache {
	return Cache{
		cache: make(map[string]Cacheable),
	}
}

func (c *Cache) Add(value Cacheable) {
	c.currentPath = append(c.currentPath, value.Name())
	key := strings.Join(c.currentPath, "/")
	c.cache[key] = value
}

func (c *Cache) Get() Cacheable {
	key := strings.Join(c.currentPath, "/")
	return c.cache[key]
}

func (c *Cache) Forward(s string) {
	c.currentPath = append(c.currentPath, s)
}

func (c *Cache) Back() string {
	if len(c.currentPath) > 1 {
		c.currentPath = c.currentPath[:len(c.currentPath)-1]
		return strings.Join(c.currentPath, "/")
	}
	return ""
}

func (c *Cache) Pop() Cacheable {
	key := strings.Join(c.currentPath, "/")
	value := c.cache[key]
	delete(c.cache, key)
	c.currentPath = c.currentPath[:len(c.currentPath)-1]
	return value
}

func (c *Cache) Del(path ...string) {
	key := strings.Join(path, "/")
	delete(c.cache, key)
}

func (c *Cache) AddAt(value Cacheable, path ...string) {
	key := strings.Join(path, "/")
	c.cache[key] = value
}

func (c *Cache) GetAt(path ...string) Cacheable {
	key := strings.Join(path, "/")
	return c.cache[key]
}

func (c *Cache) Has(path ...string) bool {
	key := strings.Join(path, "/")
	_, ok := c.cache[key]
	return ok
}

func (c *Cache) Clear() {
	c.cache = make(map[string]Cacheable)
}

func (c *Cache) Size() int {
	return len(c.cache)
}
