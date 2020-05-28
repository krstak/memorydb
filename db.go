package memorydb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

var ErrItemNotFound = errors.New("item not found")

type Manager struct {
	identifier *string

	// Add adds new item in the collection. It return new id as interface{} value
	Add func(item interface{}, collection string) interface{}

	// FindAll finds all items in the collection. It return slice of interface{}
	FindAll func(collection string) []interface{}

	// FindById finds an item in the collection by given id. It returns item as interface{} and error if item can't be found
	FindById func(id interface{}, collection string) (interface{}, error)

	// FindBy finds an item in the collection by given field and value. It returns item as interface{} and error if item can't be found
	FindBy func(field, value, collection string) (interface{}, error)

	// Update updates an item with given id in the collection. It returns error if item can't be updated
	Update func(id interface{}, item interface{}, collection string) error

	// Remove removes an item from the collection. It returns error if item can't be removed
	Remove func(id interface{}, collection string) error
}

// Identifier sets identifier for an entity
// Example: Identifier("ID")
// Default identifier is "Id"
func (m *Manager) Identifier(key string) {
	*m.identifier = key
}

type col struct {
	lastId int64
	items  []interface{}
}

type collectionitem struct {
	t     reflect.Type
	items [][]byte
}

func New() manager {
	return manager{
		db:         make(map[string]collectionitem),
		identifier: "Id",
	}

	// var syn sync.RWMutex
	// ident := "Id"
	// return Manager{
	// 	identifier: &ident,
	// 	Add:        add(&ident, db, &syn),
	// 	FindAll:    getAll(db, &syn),
	// 	FindById:   findById(&ident, db, &syn),
	// 	FindBy:     findBy(db, &syn),
	// 	Remove:     remove(&ident, db, &syn),
	// 	Update:     update(&ident, db, &syn),
	// }
}

type manager struct {
	db         map[string]collectionitem
	syn        sync.RWMutex
	identifier string
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
			t:     reflect.TypeOf(item), //todo: check if it's a pointer
			items: make([][]byte, 0),
		}
	}
	b, _ := json.Marshal(item)

	collectionItem.items = append(collectionItem.items, b)
	m.db[collection] = collectionItem
}

func (m *manager) FindAll(collection string, st interface{}) {
	m.find(collection, "", "", st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		return reflect.Append(list, s.Elem())
	})
}

func (m *manager) FindBy(collection string, field string, value interface{}, st interface{}) {
	m.find(collection, field, value, st, func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value {
		val := valueOf(s.Interface(), field)
		if val == value {
			list = reflect.Append(list, s.Elem())
		}
		return list
	})
}

func (m *manager) find(collection string, field string, value interface{}, st interface{}, append func(field string, value interface{}, s reflect.Value, list reflect.Value) reflect.Value) {
	m.syn.RLock()
	collectionItem, ok := m.db[collection]
	m.syn.RUnlock()

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
	collectionItem, ok := m.db[collection]
	m.syn.RUnlock()

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

func getAll(db map[string]col, syn *sync.RWMutex) func(string) []interface{} {
	return func(collection string) []interface{} {
		syn.RLock()
		defer syn.RUnlock()

		c, ok := db[collection]
		if !ok {
			return make([]interface{}, 0)
		}
		return c.items
	}
}

func add(identifierName *string, db map[string]col, syn *sync.RWMutex) func(interface{}, string) interface{} {
	return func(item interface{}, collection string) interface{} {
		t := checkType(item, *identifierName)

		v := reflect.ValueOf(item).Elem().FieldByName(*identifierName)

		syn.Lock()
		defer syn.Unlock()

		lastId := db[collection].lastId
		newId := lastId + 1
		var id interface{}

		switch t.Kind() {
		case reflect.Int:
			id = int(newId)
		case reflect.Int8:
			id = int8(newId)
		case reflect.Int16:
			id = int16(newId)
		case reflect.Int32:
			id = int32(newId)
		case reflect.Int64:
			id = int64(newId)
		case reflect.Uint:
			id = uint(newId)
		case reflect.Uint8:
			id = uint8(newId)
		case reflect.Uint16:
			id = uint16(newId)
		case reflect.Uint32:
			id = uint32(newId)
		case reflect.Uint64:
			id = uint64(newId)
		case reflect.String:
			strId := strconv.FormatInt(newId, 10)
			id = strId
		default:
			panic(fmt.Sprintf("type %s is not supported as identifier" + t.Kind().String()))
		}

		idValue := reflect.ValueOf(id)
		v.Set(idValue)

		items := append(db[collection].items, item)
		c := col{lastId: newId, items: items}
		db[collection] = c

		return id
	}
}

func remove(identifierName *string, db map[string]col, syn *sync.RWMutex) func(interface{}, string) error {
	return func(id interface{}, collection string) error {
		_, index, err := find(db, *identifierName, id, collection, syn)
		if err != nil {
			return err
		}

		syn.Lock()
		defer syn.Unlock()

		items := append(db[collection].items[:index], db[collection].items[index+1:]...)
		c := col{lastId: db[collection].lastId, items: items}
		db[collection] = c

		return nil
	}
}

func update(identifierName *string, db map[string]col, syn *sync.RWMutex) func(interface{}, interface{}, string) error {
	return func(id interface{}, item interface{}, collection string) error {
		checkType(item, *identifierName)

		_, index, err := find(db, *identifierName, id, collection, syn)
		if err != nil {
			return err
		}

		syn.Lock()
		defer syn.Unlock()
		idValue := reflect.ValueOf(id)
		v := reflect.ValueOf(item).Elem().FieldByName(*identifierName)
		v.Set(idValue)

		db[collection].items[index] = item
		return nil
	}
}

func findById(identifierName *string, db map[string]col, syn *sync.RWMutex) func(interface{}, string) (interface{}, error) {
	return func(id interface{}, collection string) (interface{}, error) {
		item, _, err := find(db, *identifierName, id, collection, syn)
		return item, err
	}
}

func findBy(db map[string]col, syn *sync.RWMutex) func(string, string, string) (interface{}, error) {
	return func(field, value, collection string) (interface{}, error) {
		item, _, err := find(db, field, value, collection, syn)
		return item, err
	}
}

func find(db map[string]col, field string, value interface{}, collection string, syn *sync.RWMutex) (interface{}, int, error) {
	items := getAll(db, syn)(collection)

	for i, v := range items {
		val := valueOf(v, field)
		if val == value {
			return v, i, nil
		}
	}

	return nil, 0, ErrItemNotFound
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
