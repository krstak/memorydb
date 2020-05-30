# In-Memory Database

A small in-memory database. It is useful for testing and test environment. Should not be used in production.

## Install

```go
go get github.com/krstak/memorydb
```

## API

```go
// Add adds a new item into the given collection.
// Returns an error if occurs.
Add(collection string, item interface{}) error

// FindAll finds all items in the collection and maps them into the items.
// Returns an error if occurs.
FindAll(collection string, st interface{}) error

// FindBy finds all items in the collection by given field and value and maps them into the items.
// Returns an error if occurs.
FindBy(collection string, field string, value interface{}, st interface{}) error

// Remove removes an item with a given field and value.
// Returns an error if occurs.
Remove(collection string, idval interface{}) error
```

## Usage

### Entity to persist

Any struct is a valid object to be persisted.

```go
type book struct {
	ID     uuid.UUID
	Isbn   int
	Name   string
	Author string
}
```

### Create database

```go
db := memorydb.New()
```

### Write to database

```go
bk := book{ID: uuid.New(), Isbn: 123456, Name: "In-Memory DB", Author: "Marko Krstic"}
err := db.Add("books", bk)
```

### Find all

```go
var books []book
err := db.FindAll("books", &books)
```

### Find by

```go
var books []book
err := db.FindBy("books", "Isbn", 123456, &books)
```

### Remove

```go
err := db.Remove("books", "Isbn", 123456)
```