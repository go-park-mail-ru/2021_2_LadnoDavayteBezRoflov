package hasher

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen = 8
	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

func makeSalt() (salt []byte, err error) {
	salt = make([]byte, saltLen)
	_, err = rand.Read(salt)
	return
}

func saltPassword(salt []byte, password string) []byte {
	hashedPassword := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	return append(salt, hashedPassword...)
}

func HashPassword(password string) ([]byte, error) {
	salt, err := makeSalt()
	if err != nil {
		return nil, err
	}
	return saltPassword(salt, password), nil
}

func IsPasswordsEqual(plainPassword string, hashedPassword []byte) bool {
	var salt []byte
	salt = append(salt, hashedPassword[0:saltLen]...)
	return bytes.Equal(hashedPassword, saltPassword(salt, plainPassword))
}
