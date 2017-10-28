# In-Memory Database
Small in-memory database. It is useful for testing and test environment. Should not be used in production.

## Install
```
go get -u github.com/krstak/memorydb
```

## API

```go
Add      func(item interface{}, collection string) string
FindAll  func(collection string) ([]interface{}, error)
FindById func(id string, collection string) (interface{}, error)
FindBy   func(field, value, collection string) (interface{}, error)
Update   func(id string, item interface{}, collection string) error
Remove   func(id, collection string) error
```

## Usage

##### Create database
```go
db := memorydb.Create()
```

##### Entity to persist
In order to save some struct in db, it should have the required field `Id` type `string`. Everything else is optional.
```go
type book struct {
	Id     string
	isbn   string
	name   string
	author string
}
```

##### Write to database
Important: Reference should be passed in.
```go
memBook := book{isbn: "1234567", name: "In-Memory DB", author: "Marko Krstic"}
id := db.Add(&memBook, "books")
```

##### Find all
Important: slice of interface{} is returned
```go
books, err := db.FindAll("books")
bk := books[0].(*book)
```

##### Find by Id
Important: interface{} is returned
```go
bk, err := db.FindById(id, "books")
b := bk.(*book)
```

##### Find by custom field
Important: interface{} is returned
```go
bk, err := db.FindBy("isbn", "1234567", "books")
b := bk.(*book)
```

##### Update
Important: Reference should be passed in.
```go
bk := book{isbn: "222", name: "In-Memory DB", author: "Marko Krstic"}
err := db.Update(id, &bk, "books")
```

##### Remove
```go
err := db.Remove(id, "books")
```