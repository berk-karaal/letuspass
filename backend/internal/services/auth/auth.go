package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const hashingCost = bcrypt.DefaultCost

// HashPassword returns hash value of given rawPassword.
func HashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), hashingCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword returns true if hashValue is the hashed version of rawPassword.
func ComparePassword(hashValue, rawPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(rawPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
