# In-Memory Database

A small in-memory database. It is useful for testing and test environment. Should not be used in production.

## Install

```go
go get github.com/krstak/memorydb
```

## API

```go
// Add adds a new item in the collection.
// Returns an error if occurs.
Add(collection string, item interface{}) error

// FindAll finds all items in the collection.
// Returns an error if occurs.
FindAll(collection string, st interface{}) error

// FindById finds an item in the collection by given id. 
// It returns true if the item is found, otherwise false.
// Returns an error if occurs.
FindById(collection string, idval interface{}, st interface{}) (bool, error)

// FindBy finds an item in the collection by given field and value.
// Returns an error if occurs.
FindBy(collection string, field string, value interface{}, st interface{}) error

// Remove removes an item with a given id from the collection
// Returns an error if occurs.
Remove(collection string, idval interface{}) error
```

## Usage

### Create database

```go
db := memorydb.New()
```

### Entity to persist

In order to save a struct in the db, it should have a required identifier. Default identifier is `Id`.
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
db := memorydb.New()
db.Identifier("ID")
```

### Write to database

```go
memBook := book{Id: 12, isbn: "1234567", name: "In-Memory DB", author: "Marko Krstic"}
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