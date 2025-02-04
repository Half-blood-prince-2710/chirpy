package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hashed_password,err:=bcrypt.GenerateFromPassword([]byte(password),12)
	if err!=nil{
		slog.Error("err",err)
		return "",err
	}
	return string(hashed_password), nil
}

func CheckPasswordHash(password, hash string) error {
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err!=nil{
		slog.Error("err",err)
		return err
	}
	return nil
}