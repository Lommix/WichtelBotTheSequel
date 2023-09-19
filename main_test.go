package main_test
// some integration tests to make my life easier
import (
	"database/sql"
	"fmt"
	"lommix/wichtelbot/server/store"
	"math/rand"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)
const TestPartySize = 100

// love sqlite in tests. its so easy
func OpenTestDb(t *testing.T) *sql.DB{
	db, err := sql.Open("sqlite3", "wichtel_test.db")
	if err != nil {
		t.Fatal(err)
	}

	err = store.SchemaDown(db)
	if err != nil {
		t.Fatal(err)
	}

	err = store.SchemaUp(db)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

// simple test to make sure everybody has a partner
func TestPlay(t *testing.T) {
	db := OpenTestDb(t)
	party,err := store.CreateParty(db)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < TestPartySize; i++ {
		_, err = store.CreateUser(
			db,
			party.Id,
			fmt.Sprint("test_name_",i),
			"test",
			"test_notice",
			store.DefaultUser,
		)
		if err != nil {
			t.Fatal(err)
		}
	}


	err = party.RollPartners(db, false)
	if err != nil {
		t.Fatal(err)
	}
	users, err := store.FindUsersByPartyId(db, party.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != TestPartySize {
		t.Fatal("less user in db than expected")
	}
	for _, user := range users {
		if user.PartnerId == 0 {
			t.Fatal("user has no partner")
		}

	}

	unique := make(map[int64]bool)
	for _, u := range users {
		if u.PartnerId == 0 {
			t.Fatal("no partner")
		}
		unique[u.PartnerId] = true
	}


	if len(unique) != TestPartySize {
		t.Fatal("not all users have partners")
	}

}

// simple test to make sure everybody still has a partner when using the blacklist
func TestBlacklistPlay(t *testing.T) {
	db := OpenTestDb(t)
	party,err := store.CreateParty(db)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < TestPartySize; i++ {
		_, err = store.CreateUser(
			db,
			party.Id,
			fmt.Sprint("test_name_",i),
			"test",
			"test_notice",
			store.DefaultUser,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	users, err := store.FindUsersByPartyId(db, party.Id)
	if err != nil {
		t.Fatal(err)
	}

	// test blacklist
	idStack := make([]int64, TestPartySize)
	for i := 0; i < TestPartySize; i++ {
		idStack[i] = int64(i + 1)
	}
	rand.Shuffle(len(idStack), func(i, j int) {
		idStack[i], idStack[j] = idStack[j], idStack[i]
	})
	for i, u := range users {
		u.ExcludeId = idStack[i]
		u.PartnerId = 0
		u.Update(db)
	}

	err = party.RollPartners(db, true)
	if err != nil {
		t.Fatal(err)
	}
	users, err = store.FindUsersByPartyId(db, party.Id)
	if err != nil {
		t.Fatal(err)
	}

	unique := make(map[int64]bool)
	for _, u := range users {
		if u.PartnerId == 0 {
			t.Fatal("no partner")
		}
		unique[u.PartnerId] = true
	}
	if len(unique) != TestPartySize {
		t.Fatal("blacklist failed, not all users have partners")
	}

}

// test if the blacklist fails correctly if impossible party
func TestBlacklistFail(t *testing.T) {
	db := OpenTestDb(t)
	party,err := store.CreateParty(db)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		_, err = store.CreateUser(
			db,
			party.Id,
			fmt.Sprint("test_name_",i),
			"test",
			"test_notice",
			store.DefaultUser,
		)
		if err != nil {
			t.Fatal(err)
		}

	}

	users, err := store.FindUsersByPartyId(db, party.Id)
	if err != nil {
		t.Fatal(err)
	}

	users[0].ExcludeId = users[2].Id
	users[1].ExcludeId = users[0].Id
	users[2].ExcludeId = users[0].Id

	for _, u := range users {
		u.Update(db)
	}

	err = party.RollPartners(db, true)
	if err == nil {
		t.Fatal("should have failed")
	}
}

// testing single sql query page load
func TestFastQuery(t *testing.T){
	db := OpenTestDb(t)
	party,err := store.CreateParty(db)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < TestPartySize; i++ {
		_, err = store.CreateUser(
			db,
			party.Id,
			fmt.Sprint("test_name_",i),
			"test",
			"test_notice",
			store.DefaultUser,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = party.RollPartners(db, false)
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.FindUserWithPartyFast(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(*user.Party.Users) < TestPartySize {
		t.Fatal("missing party members")
	}

	if user.PartnerId == 0 {
		t.Fatal("no partner id")
	}

	if user.Party == nil {
		t.Fatal("no party")
	}

	if user.Partner == nil {
		t.Fatal("no partner")
	}

	if user.Partner.Id != user.PartnerId{
		t.Fatal("partner id does not match partner")
	}
}
