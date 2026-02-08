package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := ""
	hash := "$2a$14$zZiAx9nSAOOQp8hQ240vROS7.N57TW8fGYlji6W7jt/2N.zjU2f0u"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Password mismatch: %v\n", err)
	} else {
		fmt.Println("Password match!")
	}
}
