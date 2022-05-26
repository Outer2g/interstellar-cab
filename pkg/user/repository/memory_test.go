package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestUserDatabase() *database {
	email := "existing@email.com"
	return &database{map[string]User{"existing@email.com": aUserObject(email)}}
}

func TestGetUser(t *testing.T) {
	repository := newTestUserDatabase()

	t.Run("Should return Nil when user does not exists", func(t *testing.T) {
		user := repository.GetUser("anonexistent@email.com")

		assert.Nil(t, user)
	})

	t.Run("Should return user when it exists", func(t *testing.T) {
		expectedUser := aUserObject("existing@email.com")

		user := repository.GetUser("existing@email.com")

		assert.EqualValues(t, expectedUser, *user)
	})
}

func TestAddUser(t *testing.T) {

	t.Run("Should return false if the user did not exists", func(t *testing.T) {
		repository := newTestUserDatabase()
		result := repository.AddUser("anonexistent@email.com", "123", false)

		assert.False(t, result)
	})

	t.Run("Should return true when the user exists", func(t *testing.T) {
		repository := newTestUserDatabase()
		result := repository.AddUser("existing@email.com", "123", false)

		assert.True(t, result)
	})
}

func aUserObject(email string) User {
	return User{email, "123", false}
}
