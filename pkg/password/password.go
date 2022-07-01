package password

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password []byte) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("HashPassword error: %s", err.Error())
		return nil
	}
	return bytes
}

func ComparePassword(password, hashedPassword []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return false
	}
	return true
}
