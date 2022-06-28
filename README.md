# In-Memory Database

A small in-memory database. It is useful for testing and test environment. Should not be used in production.

## Install

```go
go get github.com/krstak/memorydb/v3
```

## API

```go
// Add adds a new item into the given collection.
// Returns an error if occurs.
Add(collection string, item interface{}) error

// Add updates an existing item into the given collection.
// Returns an error if occurs.
Update(collection string, item interface{}, fields []Fileds) error

// FindAll finds all items in the collection and maps them into the items.
// Returns an error if occurs.
FindAll(collection string, items interface{}) error

// FindBy finds all items in the collection by given field and value and maps them into the items.
// Returns an error if occurs.
FindBy(collection string, field string, value interface{}, items interface{}) error

// Remove removes an item with a given field and value.
// Returns an error if occurs.
Remove(collection string, field string, value interface{}) error
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

### Update

```go
bk := book{...}
err := db.Update("books", bk, []Fileds{{Key: "ID", Value: 123456}})
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