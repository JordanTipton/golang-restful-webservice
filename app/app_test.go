package app_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/jordantipton/golang-restful-webservice/app"
)

type userResult struct {
	ID   int
	Name string
}

type createUserReq struct {
	Name string
}

func TestCreateGetUser(t *testing.T) {
	// Setup
	port := "8080"
	addr := "http://localhost:" + port
	a := app.App{}
	a.Initialize(os.Getenv("DSN"))
	go a.Run(":" + port)

	userName := "TestUser" + strconv.Itoa(rand.Intn(999999))
	user := createUserReq{Name: userName}
	jsonUser, _ := json.Marshal(user)
	byteBuffer := bytes.NewBuffer(jsonUser)

	// Create
	resp, err := http.Post(addr+"/users", "application/json", byteBuffer)
	if err != nil {
		t.Errorf("Create response err, expected: nil, got: %s", err.Error())
	}
	if resp.StatusCode != 201 {
		t.Errorf("Create response StatusCode, expected: %d, got: %d", 201, resp.StatusCode)
	}
	defer resp.Body.Close()
	var created userResult
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&created); err != nil {
		t.Errorf("Failed to decode created body. Error: %s", err.Error())
	}

	// Get
	resp, err = http.Get(addr + "/users/" + strconv.Itoa(created.ID))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	// Assert
	if err != nil {
		t.Errorf("Get response err, expeted: nil, got: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Get response StatusCode, expected: %d, got: %d", 200, resp.StatusCode)
	}
	var result userResult
	dec = json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		t.Errorf("Failed to decode user body from get. Error: %s", err.Error())
	}
	if result.Name != userName {
		t.Errorf("Get response user name, expected: %s, got: %s", userName, result.Name)
	}
}
