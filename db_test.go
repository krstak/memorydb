package memorydb

import (
	"sync"
	"testing"
	"time"

	"github.com/krstak/testify"
)

func TestAddAndFindById(t *testing.T) {
	db := Create()
	db.Identifier("ID")

	// ******** str *********
	user := testUserStr{Age: 15}
	id := db.Add(&user, "users")
	res, err := db.FindById(id, "users")
	actUser := *(res.(*testUserStr))

	testify.Nil(t)(err)
	testify.Equal(t)(user, actUser)

	// ******** int *********
	userInt := testUserInt{Age: 15}
	idInt := db.Add(&userInt, "users")
	resInt, err := db.FindById(idInt, "users")
	actUserInt := *(resInt.(*testUserInt))

	testify.Nil(t)(err)
	testify.Equal(t)(userInt, actUserInt)

	// ******** uint *********
	userUint := testUserUint{Age: 15}
	idUint := db.Add(&userUint, "users")
	resUint, err := db.FindById(idUint, "users")
	actUserUint := *(resUint.(*testUserUint))

	testify.Nil(t)(err)
	testify.Equal(t)(userUint, actUserUint)
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
	db.Identifier("ID")

	// ******* int *******

	user := testUserInt{Age: 15, Name: "John"}
	id := db.Add(&user, "users")
	userUpdate := testUserInt{Age: 20, Name: "John"}
	db.Update(id, &userUpdate, "users")

	res, err := db.FindById(id, "users")
	actUser := *(res.(*testUserInt))
	testify.Nil(t)(err)
	testify.Equal(t)(20, actUser.Age)
	testify.Equal(t)("John", actUser.Name)

	// ******* int16 *******

	user16 := testUserInt16{Age: 15, Name: "John"}
	id16 := db.Add(&user16, "users")
	userUpdate16 := testUserInt16{Age: 20, Name: "John"}
	db.Update(id16, &userUpdate16, "users")

	res16, err := db.FindById(id16, "users")
	actUser16 := *(res16.(*testUserInt16))
	testify.Nil(t)(err)
	testify.Equal(t)(20, actUser16.Age)
	testify.Equal(t)("John", actUser16.Name)

	// ******* uint32 *******

	user32 := testUserUint32{Age: 15, Name: "John"}
	id32 := db.Add(&user32, "users")
	userUpdate32 := testUserUint32{Age: 20, Name: "John"}
	db.Update(id32, &userUpdate32, "users")

	res32, err := db.FindById(id32, "users")
	actUser32 := *(res32.(*testUserUint32))
	testify.Nil(t)(err)
	testify.Equal(t)(20, actUser32.Age)
	testify.Equal(t)("John", actUser32.Name)

	// ******* string *******

	userStr := testUserStr{Age: 15, Name: "John"}
	idStr := db.Add(&userStr, "users")
	userUpdateStr := testUserStr{Age: 20, Name: "John"}
	db.Update(idStr, &userUpdateStr, "users")

	resStr, err := db.FindById(idStr, "users")
	actUserStr := *(resStr.(*testUserStr))
	testify.Nil(t)(err)
	testify.Equal(t)(20, actUserStr.Age)
	testify.Equal(t)("John", actUserStr.Name)
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

type testUserUint struct {
	ID   string
	Age  int
	Name string
}

type testUserStr struct {
	ID   string
	Age  int
	Name string
}

type testUserInt struct {
	ID   int
	Age  int
	Name string
}

type testUserInt8 struct {
	ID   int8
	Age  int
	Name string
}

type testUserInt16 struct {
	ID   int16
	Age  int
	Name string
}

type testUserUint32 struct {
	ID   uint32
	Age  int
	Name string
}
