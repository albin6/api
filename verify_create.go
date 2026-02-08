package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	baseURL := "http://localhost:8080"
	client := &http.Client{}

	// 1. Login
	fmt.Println("1. Logging in...")
	loginBody := []byte(`{"email":"test-one@gmail.com", "password":"Albin@123"}`)
	// Note: We need a valid user. "test-one" was used in debug. "test-two" was the one suggested in previous turn. 
	// The DB has "test-one" with "password" tag fixed? 
    // Wait, I fixed the code, but I didn't update the DB record for "test-one" which has empty password?
    // User said "test-one" has empty password hash in DB.
    // I need to use a user that *works*.
    // I can try to Create a new user first (Signup).
    
    fmt.Println("1. Signing up a new test user...")
    signupBody := []byte(`{"name":"Creator User","email":"creator@test.com","phone":"1234567890","password":"password123","role":"HEAD"}`)
    resp, err := client.Post(baseURL+"/auth/signup", "application/json", bytes.NewBuffer(signupBody))
    if err != nil {
        fmt.Println("Signup request failed:", err)
    } else {
        fmt.Println("Signup status:", resp.Status)
        resp.Body.Close()
    }

	fmt.Println("2. Logging in with new user...")
    loginBody = []byte(`{"email":"creator@test.com", "password":"password123"}`)
	resp, err = client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Login failed:", resp.Status)
        body, _ := io.ReadAll(resp.Body)
        fmt.Println(string(body))
		return
	}

	var loginResp map[string]string
	json.NewDecoder(resp.Body).Decode(&loginResp)
	token := loginResp["access_token"]
	fmt.Println("Login success! Token obtained.")

	// 2. Create Student (Success)
	fmt.Println("\n3. Creating Student (Valid)...")
	studentBody := []byte(`{"full_name":"New Student","email":"new.student@test.com","phone":"1122334455","program_status":true}`)
	req, _ := http.NewRequest("POST", baseURL+"/api/students", bytes.NewBuffer(studentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
	fmt.Println("Create Status:", resp.Status) // Should be 201 Created
    fmt.Println("Response:", string(body))

    // 3. Create Student (Duplicate)
    fmt.Println("\n4. Creating Student (Duplicate)...")
    req, _ = http.NewRequest("POST", baseURL+"/api/students", bytes.NewBuffer(studentBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err = client.Do(req)
    defer resp.Body.Close()
    body, _ = io.ReadAll(resp.Body)
    fmt.Println("Create Status:", resp.Status) // Should be 400 or 500
    fmt.Println("Response:", string(body))

    // 4. Create Student (Invalid)
    fmt.Println("\n5. Creating Student (Invalid)...")
    invalidBody := []byte(`{"full_name":"Incomplete"}`)
    req, _ = http.NewRequest("POST", baseURL+"/api/students", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err = client.Do(req)
    defer resp.Body.Close()
    body, _ = io.ReadAll(resp.Body)
    fmt.Println("Create Status:", resp.Status) // Should be 400
    fmt.Println("Response:", string(body))
}
