package memorydb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

type M interface {
	// Add adds a new item in the collection.
	// Returns an error if occurs.
	Add(collection string, item interface{}) error

	// FindAll finds all items in the collection and maps them into the items.
	// Returns an error if occurs.
	FindAll(collection string, items interface{}) error

	// FindBy finds all items in the collection by given field and value and maps them into the items.
	// Returns an error if occurs.
	FindBy(collection string, field string, value interface{}, items interface{}) error

	// Remove removes an item with by given field and value.
	// Returns an error if occurs.
	Remove(collection string, field string, value interface{}) error
}

type manager struct {
	db  map[string]collectionitem
	syn sync.RWMutex
}

func New() M {
	return &manager{
		db: make(map[string]collectionitem),
	}
}

type collectionitem struct {
	t     reflect.Type
	items [][]byte
}

func (m *manager) Add(collection string, item interface{}) error {
	m.syn.Lock()
	defer m.syn.Unlock()
	collectionItem, ok := m.db[collection]
	if !ok {
		collectionItem = collectionitem{
			t:     reflect.TypeOf(item), //todo: check if it's a pointer. It might be an error
			items: make([][]byte, 0),
		}
	}
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	collectionItem.items = append(collectionItem.items, b)
	m.db[collection] = collectionItem
	return nil
}

func (m *manager) FindAll(collection string, st interface{}) error {
	m.syn.RLock()
	defer m.syn.RUnlock()
	return m.find(collection, "", "", st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		return reflect.Append(list, s.Elem())
	})
}

func (m *manager) FindBy(collection string, field string, value interface{}, st interface{}) error {
	m.syn.RLock()
	defer m.syn.RUnlock()
	return m.find(collection, field, value, st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		val := valueOf(s.Interface(), field)
		if val == value {
			list = reflect.Append(list, s.Elem())
		}
		return list
	})
}

func (m *manager) find(collection string, field string, value interface{}, st interface{}, append func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value) error {
	collectionItem, ok := m.db[collection]

	if !ok {
		return nil
	}

	vp := reflect.ValueOf(st).Elem()

	for _, b := range collectionItem.items {
		stitem := reflect.New(vp.Type().Elem())
		err := json.Unmarshal(b, stitem.Interface())
		if err != nil {
			return err
		}

		vp = append(field, value, stitem, vp)
	}

	reflect.ValueOf(st).Elem().Set(vp)
	return nil
}

func (m *manager) Remove(collection string, field string, value interface{}) error {
	m.syn.Lock()
	defer m.syn.Unlock()
	collectionItem, ok := m.db[collection]

	if !ok {
		return nil
	}

	indexToRemove := -1
	for i, b := range collectionItem.items {
		st := reflect.New(collectionItem.t)
		err := json.Unmarshal(b, st.Interface())
		if err != nil {
			return err
		}

		val := valueOf(st.Interface(), field)
		if val == value {
			indexToRemove = i
			break
		}
	}

	if indexToRemove >= 0 {
		collectionItem.items = append(collectionItem.items[:indexToRemove], collectionItem.items[indexToRemove+1:]...)
		m.db[collection] = collectionItem
	}
	return nil
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
