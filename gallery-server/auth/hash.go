package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 12)
	return string(bytes), err
}

func ValidateHash(pass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err == nil
}
