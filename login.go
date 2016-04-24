package auth

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func MakeToken() (token string, err error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["role"] = "user"
	token.Claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
	return tokenString, err
}

func Validate(token string) (result bool) {
	token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		//return myLookupKey(token.Header["kid"]), nil
		return nil, nil
	})

	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
}
