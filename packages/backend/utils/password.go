package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengenkripsi password string plain text menjadi Bcrypt hash text
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash mencocokkan password plain dengan hash yang ada di database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
