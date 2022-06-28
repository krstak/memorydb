package memorydb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/krstak/testify"
)

func TestAddAndFindAll(t *testing.T) {
	db := New()
	user1 := newuser(34, "John")
	user2 := newuser(22, "Merry")

	db.Add("users", user1)
	db.Add("users", user2)

	var users = []testuser{}
	db.FindAll("users", &users)

	testify.Equal(t)(2, len(users))
	testify.Equal(t)(user1.Id, users[0].Id)
	testify.Equal(t)(user1.Age, users[0].Age)
	testify.Equal(t)(user1.Name, users[0].Name)
	testify.Equal(t)(user2.Id, users[1].Id)
	testify.Equal(t)(user2.Age, users[1].Age)
	testify.Equal(t)(user2.Name, users[1].Name)
}

func TestFindBy(t *testing.T) {
	db := New()
	user1 := newuser(34, "John")
	user2 := newuser(22, "Merry")
	user3 := newuser(67, "Jo")
	user4 := newuser(22, "Sofia")

	db.Add("users", user1)
	db.Add("users", user2)
	db.Add("users", user3)
	db.Add("users", user4)

	users := []testuser{}
	db.FindBy("users", "Age", 22, &users)

	testify.Equal(t)(2, len(users))
	testify.Equal(t)(user2.Id, users[0].Id)
	testify.Equal(t)(user4.Id, users[1].Id)
}

func TestRemove(t *testing.T) {
	db := New()
	user1 := newuser(34, "John")
	user2 := newuser(22, "Merry")
	user3 := newuser(67, "Jo")

	db.Add("users", user1)
	db.Add("users", user2)
	db.Add("users", user3)

	db.Remove("users", "Id", user2.Id)

	var users = []testuser{}
	db.FindAll("users", &users)

	testify.Equal(t)(2, len(users))
	testify.Equal(t)(user1.Id, users[0].Id)
	testify.Equal(t)(user3.Id, users[1].Id)
}

func TestUpdate(t *testing.T) {
	db := New()
	user1 := newuser(34, "John")
	user2 := newuser(22, "John")
	user3 := newuser(67, "Jo")
	user4 := newuser(21, "Sofia")

	db.Add("users", user1)
	db.Add("users", user2)
	db.Add("users", user3)
	db.Add("users", user4)

	user2.Name = "Johnny"
	err := db.Update("users", user2, []Fileds{{Key: "Name", Value: "John"}, {Key: "Age", Value: 22}})
	testify.Nil(t)(err)

	users := []testuser{}
	db.FindBy("users", "Age", 22, &users)

	testify.Equal(t)(1, len(users))
	testify.Equal(t)(user2.Id, users[0].Id)
	testify.Equal(t)(user2.Name, "Johnny")

	users = []testuser{}
	db.FindBy("users", "Age", 34, &users)

	testify.Equal(t)(1, len(users))
	testify.Equal(t)(user1.Id, users[0].Id)
	testify.Equal(t)(user1.Name, "John")
}

type testuser struct {
	Id   uuid.UUID
	Age  int
	Name string
}

func newuser(age int, name string) testuser {
	return testuser{
		Id:   uuid.New(),
		Age:  age,
		Name: name,
	}
}
