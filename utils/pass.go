package utils

import "golang.org/x/crypto/bcrypt"

//HashPassword -
func HashPassword(pass string) string {
	converted, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(converted)
}

//CompareHashedPassword -
func CompareHashedPassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err == nil {
		return true
	}
	return false
}
