package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAwtToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("Should generate awt token", func(t *testing.T) {

		result, err := GenerateJwtToken("arandom@email.com", false)

		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestVerifyJwtToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("Should verify correctly a correct token", func(t *testing.T) {
		token, _ := GenerateJwtToken("arandom@email.com", false)

		result, err := VerifyJwtToken(token)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.EqualValues(t, "arandom@email.com", result.Email)
	})

	t.Run("Should return invalid key when not a valid key", func(t *testing.T) {
		token, _ := GenerateJwtToken("arandom@email.com", false)

		_, err := VerifyJwtToken(token + "askguab")

		assert.Containsf(t, err.Error(), "ERROR invalid token", err.Error())
	})
}
