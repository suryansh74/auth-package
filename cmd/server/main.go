package main

import (
	"fmt"
	"log"

	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/dto"
	"github.com/suryansh74/auth-package/internal/services"
)

func main() {
	// 1. Init DB
	authDB := db.NewAuth()
	defer authDB.Close()

	// 2. Init Authenticator service
	authService := services.NewAuthenticator(authDB)

	// ========== TEST REGISTER ==========
	registerReq := dto.UserRegisterRequest{
		Name:     "chit",
		Email:    "chit@example.com",
		Password: "123456",
	}

	fmt.Println("➡️ Calling Register...")

	registerRes, err := authService.Register(registerReq)
	if err != nil {
		log.Println("❌ Register Error:", err)
	} else {
		fmt.Println("✅ Register Response:", registerRes)
	}

	// ========== TEST LOGIN ==========
	loginReq := dto.UserLoginRequest{
		Email:    "chit@example.com",
		Password: "123456",
	}

	fmt.Println("➡️ Calling Login...")

	loginRes, err := authService.Login(loginReq)
	if err != nil {
		log.Println("❌ Login Error:", err)
	} else {
		fmt.Println("✅ Login Response:", loginRes)
	}

	fmt.Println("🎉 Done!")
}
