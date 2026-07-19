package main

import (
	"fmt"

	"github.com/shubhkr72/helix/internal/password"
)

func main() {

	hash, err := password.HashPassword("password123")
	if err != nil {
		panic(err)
	}

	fmt.Println("Hash:", hash)

	err = password.VerifyPassword("password123", hash)
	fmt.Println("Correct password:", err == nil)

	err = password.VerifyPassword("wrongpassword", hash)
	fmt.Println("Wrong password:", err == nil)
}