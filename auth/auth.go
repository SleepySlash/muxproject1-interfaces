package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(username string) (string, error) {
    log.Println("Creating token for user:", username)
	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 12).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
    secretKey := os.Getenv("SECRET_KEY")
    log.Println("Verifying token:", tokenString)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        log.Println("Token parsing error:", err)
        return err
    }

    if !token.Valid {
        log.Println("Token is invalid")
        return fmt.Errorf("invalid token")
    }

    log.Println("Token is valid")
    return nil
}

