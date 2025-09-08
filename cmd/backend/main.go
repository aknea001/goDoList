package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
}

func main() {
	var username string
	var hash string

	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Password: ")
	fmt.Scan(&hash)

	newUser := []User{
		{Username: username, Hash: hash},
	}

	byteValue, err := json.Marshal(newUser[0])

	if err != nil {
		panic(err)
	}

	err = os.WriteFile("users.json", byteValue, 0644)

	if err != nil {
		panic(err)
	}
}
