package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/db"
	"github.com/lewiscasewell/hotel-reservation/types"
)

func makeTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: "Lewis",
		LastName:  "Casewell",
		Email:     "lewis@test.com",
		Password:  "password",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "lewis@test.com",
		Password: "password",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	var authResponse AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Errorf("expected token to be present")
	}

	// Set encrppted password to empty string so we can compare the rest of the user
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		t.Errorf("expected user to be present")
	}

}

func TestAuthenticateFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "lewis@test.com",
		Password: "password-not-corrrect",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", resp.StatusCode)
	}

	var genResp genericResponse

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Errorf("expected gen response type to be error but got %s", genResp.Type)
	}

	if genResp.Msg != "Invalid credentials" {
		t.Errorf("expected gen response msg to be Invalid credentials but got %s", genResp.Msg)
	}
}
