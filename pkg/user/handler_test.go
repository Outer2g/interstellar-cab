package user

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRepository struct {
	user   *User
	exists bool
}

func (repo *testRepository) GetUser(email string) *User {
	return repo.user
}

func (repo *testRepository) AddUser(email string, passwordhash string, isVip bool) bool {
	return repo.exists
}

func newTestUserController(user *User, exists bool) *AuthController {
	impl := testRepository{user, exists}
	return &AuthController{&impl}
}
func Test(t *testing.T) {
	t.Run("Should register new user", func(t *testing.T) {
		service := newTestUserController(nil, false)
		req := httptest.NewRequest("POST", "/signup", aUserInJson())
		recorder := httptest.NewRecorder()

		service.HandleSignupUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 200, result.StatusCode)
		assert.NotEmpty(t, body)

	})

	t.Run("Should not return ok when user already exists", func(t *testing.T) {
		service := newTestUserController(nil, true)
		req := httptest.NewRequest("POST", "/signup", aUserInJson())
		recorder := httptest.NewRecorder()

		service.HandleSignupUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 406, result.StatusCode)
		assert.Empty(t, body)

	})

	t.Run("Should not return ok when invalid mail", func(t *testing.T) {
		service := newTestUserController(nil, false)
		req := httptest.NewRequest("POST", "/signup", aUserInJsonWithBrokenMail())
		recorder := httptest.NewRecorder()

		service.HandleSignupUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 406, result.StatusCode)
		assert.Empty(t, body)

	})
}

func TestHandleLoginUser(t *testing.T) {
	t.Run("Should login with existent user", func(t *testing.T) {
		service := newTestUserController(aUser("existing@email.com"), false)
		req := httptest.NewRequest("POST", "/login", aUserInJson())
		recorder := httptest.NewRecorder()

		service.HandleLoginUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 200, result.StatusCode)
		assert.NotEmpty(t, body)
	})

	t.Run("Should not return ok when user does not exist", func(t *testing.T) {
		service := newTestUserController(nil, false)
		req := httptest.NewRequest("POST", "/login", aUserInJson())
		recorder := httptest.NewRecorder()

		service.HandleLoginUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 404, result.StatusCode)
		assert.Empty(t, body)

	})

	t.Run("Should not return ok when invalid mail", func(t *testing.T) {
		service := newTestUserController(nil, false)
		req := httptest.NewRequest("POST", "/signup", aUserInJsonWithBrokenMail())
		recorder := httptest.NewRecorder()

		service.HandleLoginUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 406, result.StatusCode)
		assert.Empty(t, body)
	})

	t.Run("Should not return ok when passwords does not match", func(t *testing.T) {
		service := newTestUserController(aUserWithPassword("existing@email.com", "anotherpassword"), false)
		req := httptest.NewRequest("POST", "/login", aUserInJson())
		recorder := httptest.NewRecorder()

		service.HandleLoginUser(recorder, req)

		result := recorder.Result()
		body, _ := ioutil.ReadAll(result.Body)

		assert.EqualValues(t, 404, result.StatusCode)
		assert.Empty(t, body)

	})
}

func aUserInJson() *bytes.Buffer {
	var jsonData = []byte(`{
		"email": "arandom@email.com",
		"passwordHash": "123",
		"vip": false
	}`)
	return bytes.NewBuffer(jsonData)
}

func aUserInJsonWithBrokenMail() *bytes.Buffer {
	var jsonData = []byte(`{
		"email": "notanemail",
		"passwordHash": "123",
		"vip": false
	}`)
	return bytes.NewBuffer(jsonData)
}

func aUserWithPassword(email, password string) *User {
	return &User{email, password, false}
}

func aUser(email string) *User {
	return aUserWithPassword(email, "123")
}
