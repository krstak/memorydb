package memorydb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var ErrItemNotFound = errors.New("item not found")

type manager struct {
	db         map[string]collectionitem
	syn        sync.RWMutex
	identifier string
}

func New() manager {
	return manager{
		db:         make(map[string]collectionitem),
		identifier: "Id",
	}
}

type collectionitem struct {
	t     reflect.Type
	items [][]byte
}

// Identifier sets identifier for an entity
// Example: Identifier("ID")
// Default identifier is "Id"
func (m *manager) Identifier(key string) {
	m.identifier = key
}

func (m *manager) Add(collection string, item interface{}) {
	m.syn.Lock()
	defer m.syn.Unlock()
	collectionItem, ok := m.db[collection]
	if !ok {
		collectionItem = collectionitem{
			t:     reflect.TypeOf(item), //todo: check if it's a pointer. It might be an error
			items: make([][]byte, 0),
		}
	}
	b, _ := json.Marshal(item)

	collectionItem.items = append(collectionItem.items, b)
	m.db[collection] = collectionItem
}

func (m *manager) FindAll(collection string, st interface{}) {
	m.syn.RLock()
	defer m.syn.RUnlock()
	m.find(collection, "", "", st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		return reflect.Append(list, s.Elem())
	})
}

func (m *manager) FindBy(collection string, field string, value interface{}, st interface{}) {
	m.syn.RLock()
	defer m.syn.RUnlock()
	m.find(collection, field, value, st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		val := valueOf(s.Interface(), field)
		if val == value {
			list = reflect.Append(list, s.Elem())
		}
		return list
	})
}

func (m *manager) find(collection string, field string, value interface{}, st interface{}, append func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value) {
	collectionItem, ok := m.db[collection]

	if !ok {
		return
	}

	vp := reflect.ValueOf(st).Elem()

	for _, b := range collectionItem.items {
		stitem := reflect.New(vp.Type().Elem())
		json.Unmarshal(b, stitem.Interface())

		vp = append(field, value, stitem, vp)
	}

	reflect.ValueOf(st).Elem().Set(vp)
}

func (m *manager) FindById(collection string, idval interface{}, st interface{}) bool {
	m.syn.RLock()
	defer m.syn.RUnlock()
	collectionItem, ok := m.db[collection]

	if !ok {
		return false
	}

	found := false
	for _, b := range collectionItem.items {
		json.Unmarshal(b, st)

		val := valueOf(st, m.identifier)
		if val == idval {
			found = true
			break
		}
	}
	return found
}

func (m *manager) Remove(collection string, idval interface{}) {
	m.syn.Lock()
	defer m.syn.Unlock()
	collectionItem, ok := m.db[collection]

	if !ok {
		return
	}

	indexToRemove := -1
	for i, b := range collectionItem.items {
		st := reflect.New(collectionItem.t)
		err := json.Unmarshal(b, st.Interface())
		fmt.Println(err)

		val := valueOf(st.Interface(), m.identifier)
		if val == idval {
			indexToRemove = i
			break
		}
	}

	if indexToRemove >= 0 {
		collectionItem.items = append(collectionItem.items[:indexToRemove], collectionItem.items[indexToRemove+1:]...)
		m.db[collection] = collectionItem
	}
}

func valueOf(item interface{}, field string) interface{} {
	checkType(item, field)
	v := reflect.ValueOf(item).Elem().FieldByName(field)
	return v.Interface()
}

func checkType(item interface{}, field string) reflect.Type {
	t, ok := reflect.TypeOf(item).Elem().FieldByName(field)
	if !ok {
		panic(fmt.Sprintf("field '%s' doesn't exist", field))
	}

	return t.Type
}
