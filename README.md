# In-Memory Database

Small in-memory database. It is useful for testing and test environment. Should not be used in production.

## Install

```go
get -u github.com/krstak/memorydb
```

## API

```go
// Add adds new item in the collection.
// It return new id as interface{} value
Add func(item interface{}, collection string) interface{}

// FindAll finds all items in the collection.
// It return slice of interface{}
FindAll func(collection string) []interface{}

// FindById finds an item in the collection by given id. 
// It returns item as interface{} and error if item can't be found
FindById func(id interface{}, collection string) (interface{}, error)

// FindBy finds an item in the collection by given field and value.
// It returns item as interface{} and error if item can't be found
FindBy func(field, value, collection string) (interface{}, error)

// Update updates an item with given id in the collection.
// It returns error if item can't be updated
Update func(id interface{}, item interface{}, collection string) error

// Remove removes an item from the collection.
// It returns error if item can't be removed
Remove func(id interface{}, collection string) error
```

## Usage

### Create database

```go
db := memorydb.Create()
```

### Entity to persist

In order to save some struct in db, it should have a required identifier. Default identifier is `Id`. 
Possible identifier types are: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `string`
Everything else is optional.

```go
type book struct {
	Id     string
	isbn   string
	name   string
	author string
}
```

### Custom identifier

```go
db := memorydb.Create()
db.Identifier("ID")
```

### Write to database

```go
memBook := book{isbn: "1234567", name: "In-Memory DB", author: "Marko Krstic"}
id := db.Add(&memBook, "books")
```

### Find all

```go
books := db.FindAll("books")
bk := books[0].(*book)
```

### Find by Id

```go
bk, err := db.FindById(id, "books")
b := bk.(*book)
```

### Find by custom field

```go
bk, err := db.FindBy("isbn", "1234567", "books")
b := bk.(*book)
```

### Update

```go
bk := book{isbn: "222", name: "In-Memory DB", author: "Marko Krstic"}
err := db.Update(id, &bk, "books")
```

### Remove

```go
err := db.Remove(id, "books")
```