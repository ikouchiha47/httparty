package collections

import (
	"sync"
)

type OrderedMap struct {
	mu sync.Mutex

	m map[interface{}]interface{}
	k []interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{m: make(map[interface{}]interface{}), k: []interface{}{}}
}

func (om *OrderedMap) Get(value interface{}) (interface{}, bool) {
	om.mu.Lock()
	defer om.mu.Unlock()

	value, ok := om.m[value]
	return value, ok
}

func (om *OrderedMap) Contains(value interface{}) bool {
	om.mu.Lock()
	defer om.mu.Unlock()

	_, ok := om.m[value]
	return ok
}

func (om *OrderedMap) Add(key, value interface{}) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, ok := om.m[key]; !ok {
		om.k = append(om.k, key)
	}

	om.m[key] = value
}

func (om *OrderedMap) Remove(key interface{}) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, ok := om.m[key]; !ok {
		return
	}

	delete(om.m, key)

	for i := range om.k {
		if om.k[i] == key {
			om.k = append(om.k[:i], om.k[i+1:]...)
			return
		}
	}
}

func (om *OrderedMap) Values() []interface{} {
	result := make([]interface{}, len(om.k))

	for i := range om.k {
		result[i] = om.m[i]
	}

	return result
}

//func (om *OrderedMap) M() map[interface{}]interface{} {
//	return om.m
//}