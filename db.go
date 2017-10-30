package memorydb

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
)

const (
	identifier = "Id"
)

var ErrItemNotFound = errors.New("item not found")

type Manager struct {
	Add      func(item interface{}, collection string) string
	FindAll  func(collection string) []interface{}
	FindById func(id string, collection string) (interface{}, error)
	FindBy   func(field, value, collection string) (interface{}, error)
	Update   func(id string, item interface{}, collection string) error
	Remove   func(id, collection string) error
}

type col struct {
	lastId int64
	items  []interface{}
}

func Create() Manager {
	db := make(map[string]col)
	var syn sync.RWMutex
	return Manager{
		Add:      add(db, syn),
		FindAll:  getAll(db),
		FindById: findById(db),
		FindBy:   findBy(db),
		Remove:   remove(db, syn),
		Update:   update(db, syn),
	}
}

func getAll(db map[string]col) func(string) []interface{} {
	return func(collection string) []interface{} {
		c, ok := db[collection]
		if !ok {
			return make([]interface{}, 0)
		}
		return c.items
	}
}

func add(db map[string]col, syn sync.RWMutex) func(interface{}, string) string {
	return func(item interface{}, collection string) string {
		syn.Lock()
		defer syn.Unlock()
		checkIdField(item)

		lastId := db[collection].lastId
		newId := lastId + 1
		strId := strconv.FormatInt(newId, 10)

		v := reflect.ValueOf(item).Elem().FieldByName(identifier)
		v.SetString(strId)

		items := append(db[collection].items, item)
		c := col{lastId: newId, items: items}
		db[collection] = c
		return strId
	}
}

func remove(db map[string]col, syn sync.RWMutex) func(string, string) error {
	return func(id, collection string) error {
		syn.Lock()
		defer syn.Unlock()
		_, index, err := find(db, identifier, id, collection)
		if err != nil {
			return err
		}

		items := append(db[collection].items[:index], db[collection].items[index+1:]...)
		c := col{lastId: db[collection].lastId, items: items}
		db[collection] = c

		return nil
	}
}

func update(db map[string]col, syn sync.RWMutex) func(string, interface{}, string) error {
	return func(id string, item interface{}, collection string) error {
		syn.Lock()
		defer syn.Unlock()
		_, index, err := find(db, identifier, id, collection)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(item).Elem().FieldByName(identifier)
		v.SetString(id)

		db[collection].items[index] = item
		return nil
	}
}

func findById(db map[string]col) func(string, string) (interface{}, error) {
	return func(id, collection string) (interface{}, error) {
		item, _, err := find(db, identifier, id, collection)
		return item, err
	}
}

func findBy(db map[string]col) func(string, string, string) (interface{}, error) {
	return func(field, value, collection string) (interface{}, error) {
		item, _, err := find(db, field, value, collection)
		return item, err
	}
}

func find(db map[string]col, field, value, collection string) (interface{}, int, error) {
	items := getAll(db)(collection)

	for i, v := range items {
		if reflect.ValueOf(v).Elem().FieldByName(field).String() == value {
			return v, i, nil
		}
	}

	return nil, 0, ErrItemNotFound
}

func checkIdField(item interface{}) {
	elem := reflect.TypeOf(item).Elem()
	val, ok := elem.FieldByName(identifier)
	if !ok || val.Type.String() != "string" {
		panic(elem.Name() + " must have field `" + identifier + "` type string")
	}
}
