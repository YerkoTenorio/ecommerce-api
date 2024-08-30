package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

// CheckPasswordHash compara un hash de contraseña con su posible contraseña en texto plano

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}
