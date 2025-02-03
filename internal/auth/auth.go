package auth

import "golang.org/x/crypto/bcrypt"


func HashPassword(password string) (string, error) {
	hashed_password,err:=bcrypt.GenerateFromPassword([]byte(password),12)
}