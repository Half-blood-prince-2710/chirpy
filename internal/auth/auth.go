package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hashed_password,err:=bcrypt.GenerateFromPassword([]byte(password),12)
	if err!=nil{
		slog.Error("Error hashing password: ","err",err)
		return "",err
	}
	return string(hashed_password), nil
}

func CheckPasswordHash(password, hash string) error {
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err!=nil{
		slog.Error("Error Comparing Password","err",err)
		return err
	}
	return nil
}


func MakeJWT(userID uuid.UUID, tokenSecret string) (string,error) {
	// fmt.Print("makejwt: userid ",userID,"\n\n")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,&jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject: userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err!=nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken,nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID,error){

	token , err:=jwt.ParseWithClaims(tokenString,&jwt.RegisteredClaims{},func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret),nil
	})
	if err!=nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w",err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	claims , ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil , fmt.Errorf("invalid claims format")
	}
	// fmt.Print("\n\n validatejwt:  claims subject: ",claims.Subject,"\n\n")
	userID,err := uuid.Parse(claims.Subject)
	if err!=nil {
		return uuid.Nil , fmt.Errorf("invalid user ID in token: %w",err)
	}

	return userID , nil
}

func GetBearerToken(header http.Header) (string,error) {
	str:=header.Get("Authorization")
	if str== "" {
		return "",fmt.Errorf("no token")
	}
	token := strings.TrimSpace(strings.TrimPrefix(str,"Bearer"))
	if token == "" {
		return "", fmt.Errorf("token is empty after trimming")
	}

	return token ,nil

}


func MakeRefreshToken() (string,error){
	randomBytes :=  make([]byte,32)
	_,err:=rand.Read(randomBytes)
	if err!=nil {
		return "",fmt.Errorf("error: failed to random data, - %w ",err)
	}
	refreshToken := hex.EncodeToString(randomBytes)
	return refreshToken,nil
}

 func GetAPIKey(headers http.Header) (string, error) {
	str:= headers.Get("Authorization")
	if str== "" {
		return "",fmt.Errorf("no apiKey")
	}
	apiKey := strings.TrimSpace(strings.TrimPrefix(str,"ApiKey"))
	if apiKey == "" {
		return "", fmt.Errorf("apikey is empty after trimming")
	}
	return apiKey,nil
 }