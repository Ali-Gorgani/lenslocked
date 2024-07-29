package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		HashPassword(os.Args[2])
	case "compare":
		ComparePassword(os.Args[2], os.Args[3])
	default:
		fmt.Println("Invalid command")
	}
}

func HashPassword(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %s", err)
	}
	hashedPassword := string(hashedBytes)

	fmt.Println(hashedPassword)
}

func ComparePassword(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Passwords match")
}
