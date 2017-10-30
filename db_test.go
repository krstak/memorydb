package memorydb

import (
	"github.com/krstak/testify"
	"sync"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	db := Create()

	user := testUser{Age: 15}

	id := db.Add(&user, "users")

	res, err := db.FindById(id, "users")
	actUser := *(res.(*testUser))

	testify.Nil(t)(err)
	testify.Equal(t)(user, actUser)
}

func TestAddConcurrency(t *testing.T) {
	db := Create()
	var w sync.WaitGroup

	v := func() {
		user := testUser{Age: 10}
		db.Add(&user, "users")
		time.Sleep(1 * time.Millisecond)
		w.Done()
	}

	userNum := 1000
	w.Add(userNum)
	for i := 0; i < userNum; i++ {
		go v()
	}
	w.Wait()

	res := db.FindAll("users")

	testify.Equal(t)(userNum, len(res))

	// there must not be duplicated IDs
	mp := make(map[string]int)
	for i := 0; i < len(res); i++ {
		mp[(res[i]).(*testUser).Id]++
		testify.Equal(t)(1, mp[(res[i]).(*testUser).Id])
	}
}

func TestUpdate(t *testing.T) {
	db := Create()

	user := testUser{Age: 15, Name: "John"}
	id := db.Add(&user, "users")

	userUpdate := testUser{Age: 20, Name: "John"}
	db.Update(id, &userUpdate, "users")

	res, err := db.FindById(id, "users")
	actUser := *(res.(*testUser))

	testify.Nil(t)(err)
	testify.Equal(t)(20, actUser.Age)
	testify.Equal(t)("John", actUser.Name)
}

func TestRemove(t *testing.T) {
	db := Create()
	coll := "users"

	user1 := testUser{Age: 15, Name: "John"}
	id1 := db.Add(&user1, "users")

	user2 := testUser{Age: 17, Name: "Lucy"}
	id2 := db.Add(&user2, coll)

	db.Remove(id1, coll)

	res := db.FindAll(coll)
	actUser := *(res[0].(*testUser))

	testify.Equal(t)(1, len(res))
	testify.Equal(t)(id2, actUser.Id)
	testify.Equal(t)("Lucy", actUser.Name)
	testify.Equal(t)(17, actUser.Age)
}

func TestFindBy(t *testing.T) {
	db := Create()
	coll := "users"

	user1 := testUser{Age: 15, Name: "John"}
	db.Add(&user1, "users")

	user2 := testUser{Age: 17, Name: "Lucy"}
	db.Add(&user2, coll)

	res, err := db.FindBy("Name", "John", coll)
	actUser := *(res.(*testUser))

	testify.Nil(t)(err)
	testify.Equal(t)("John", actUser.Name)
	testify.Equal(t)(15, actUser.Age)
}

type testUser struct {
	Id   string
	Age  int
	Name string
}
