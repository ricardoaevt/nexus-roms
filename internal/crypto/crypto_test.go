package crypto_test

import (
	"romsRename/internal/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	key, _ := crypto.GetKey()
	plaintext := "my secret message"

	t.Run("Encrypt and decrypt", func(t *testing.T) {
		encoded, err := crypto.Encrypt(plaintext, key)
		assert.NoError(t, err)
		assert.NotEqual(t, plaintext, encoded)

		decoded, err := crypto.Decrypt(encoded, key)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decoded)
	})

	t.Run("Decrypt invalid text", func(t *testing.T) {
		_, err := crypto.Decrypt("invalid!!", key)
		assert.Error(t, err)
	})
}
