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

func TestAddAndFindById(t *testing.T) {
	db := New()
	user1 := newuser(34, "John")
	user2 := newuser(22, "Merry")
	user3 := newuser(67, "Jo")

	db.Add("users", user1)
	db.Add("users", user2)
	db.Add("users", user3)

	user := testuser{}
	db.FindById("users", user2.Id, &user)

	testify.Equal(t)(user2, user)
}

func TestAddAndFindByCustomId(t *testing.T) {
	db := New()
	db.Identifier("ID")
	user1 := newuser(34, "John")
	user2 := newuser(22, "Merry")
	user3 := newuser(67, "Jo")

	db.Add("users", user1)
	db.Add("users", user2)
	db.Add("users", user3)

	user := testuser{}
	db.FindById("users", user2.ID, &user)

	testify.Equal(t)(user2, user)
}

func TestFindBy(t *testing.T) {
	db := New()
	db.Identifier("ID")
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

	db.Remove("users", user2.Id)

	var users = []testuser{}
	db.FindAll("users", &users)

	testify.Equal(t)(2, len(users))
	testify.Equal(t)(user1.Id, users[0].Id)
	testify.Equal(t)(user3.Id, users[1].Id)
}

type testuser struct {
	Id   uuid.UUID
	ID   uuid.UUID
	Age  int
	Name string
}

func newuser(age int, name string) testuser {
	return testuser{
		Id:   uuid.New(),
		ID:   uuid.New(),
		Age:  age,
		Name: name,
	}
}
