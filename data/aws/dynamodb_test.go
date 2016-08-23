package aws

import (
	"bytes"
	"testing"

	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
)

const (
	endpoint = "http://localhost:8000"
	region   = "us-west-2"
)

func newTestDynamoDB() *DynamoDB {
	return NewDynamoDB(region, endpoint)
}

func compareUsersForEquality(t *testing.T, u1 *models.User, u2 *models.User) {
	if u1.ID != u2.ID ||
		u1.Registered != u2.Registered ||
		u1.Email != u2.Email ||
		u1.PhoneNumber != u2.PhoneNumber ||
		u1.GCMRegID != u2.GCMRegID ||
		u1.APNSDeviceToken != u2.APNSDeviceToken ||
		u1.Name != u2.Name ||
		!bytes.Equal(u1.Picture, u2.Picture) ||
		u1.JWTToken != u2.JWTToken {
		t.Fatal("user fields are invalid")
	}
}

func TestListTables(t *testing.T) {
	db := newTestDynamoDB()
	tbl, err := db.listTables()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tbl)
}

func TestSeed(t *testing.T) {
	db := newTestDynamoDB()
	err := db.Seed(true)
	if err != nil {
		t.Fatal(err)
	}

	tbl, err := db.listTables()
	if err != nil {
		t.Fatal(err)
	}

	if len(tbl) < 1 {
		t.Fatal("tables not created")
	}
}

func TestGetByID(t *testing.T) {
	db := newTestDynamoDB()

	for _, user := range data.SeedUsers {
		u, ok := db.GetByID(user.ID)
		if !ok {
			t.Fatal("coulnd't get user")
		}

		compareUsersForEquality(t, u, &user)
	}
}

func TestGetByMail(t *testing.T) {
	db := newTestDynamoDB()

	for _, user := range data.SeedUsers {
		u, ok := db.GetByEmail(user.ID)
		if !ok {
			t.Fatal("coulnd't get user")
		}

		compareUsersForEquality(t, u, &user)
	}
}

func TestSaveUser(t *testing.T) {
	db := newTestDynamoDB()

	// create a user
	u := models.User{
		Email:    "test@user",
		Name:     "Test User",
		JWTToken: "345565",
	}

	if err := db.SaveUser(&u); err != nil {
		t.Fatal("cannot create user")
	} else if u.ID == "" {
		t.Fatal("user was not assigned a unique ID")
	}

	// update the user
	u.Email = "test2@user"
	if err := db.SaveUser(&u); err != nil {
		t.Fatal("cannot create user")
	}

	ur, ok := db.GetByID(u.ID)
	if !ok {
		t.Fatal("coulnd't get user")
	}

	compareUsersForEquality(t, ur, &u)
}