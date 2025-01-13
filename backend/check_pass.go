package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Replace these with your database values
	plainPassword := "Louder!23" // The user's input
	hashedPassword := "$2a$12$nGUphdnT0s2xUCSjJbzb/esIiNbPhQtgOt5j16VOu5BJcu12GYcc6" // The hash from DB

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		fmt.Println("Password mismatch:", err)
	} else {
		fmt.Println("Password match!")
	}
}
