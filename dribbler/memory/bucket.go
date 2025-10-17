package memory

import "reflect"

type bucketStore struct {
	buckets map[string]any
}

func newBucketStore() *bucketStore {
	return &bucketStore{
		buckets: make(map[string]any),
	}
}

func Add[T any](name string, table *Bucket[T]) {
	Store.bucketStore.buckets[name] = table
}

func Get[T any](name string) (*Bucket[T], bool) {
	table, ok := Store.bucketStore.buckets[name]
	if !ok {
		return nil, false
	}
	return table.(*Bucket[T]), true
}

func GetFor[T any]() (*Bucket[T], bool) {
	name := reflect.TypeFor[T]().Name()
	return Get[T](name)
}

type Bucket[T any] struct {
	Name    string
	data    map[uint64]T
	nextKey uint64
}

func NewTable[T any](name string) *Bucket[T] {
	return &Bucket[T]{
		Name:    name,
		data:    make(map[uint64]T),
		nextKey: 1,
	}
}

func (t *Bucket[T]) Add(value T) uint64 {
	key := t.nextKey
	t.data[key] = value
	t.nextKey++
	return key
}

func (t *Bucket[T]) Get(key uint64) (T, bool) {
	value, ok := t.data[key]
	return value, ok
}

func (t *Bucket[T]) Update(key uint64, value T) bool {
	if _, ok := t.data[key]; ok {
		t.data[key] = value
		return true
	}
	return false
}

func (t *Bucket[T]) Delete(key uint64) bool {
	if _, ok := t.data[key]; ok {
		delete(t.data, key)
		return true
	}
	return false
}

func (t *Bucket[T]) GetAll() []T {
	values := make([]T, 0, len(t.data))
	for _, value := range t.data {
		values = append(values, value)
	}
	return values
}
