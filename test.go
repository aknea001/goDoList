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
	newUser := []User{
		{Username: "testUser", Hash: "amazning_password"},
	}

	byteValue, err := json.Marshal(newUser[0])
	if err != nil {
		panic(err)
	}

	dataToAdd := []byte(fmt.Sprintf(",%s]", byteValue))

	filename := "users.json"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	err = file.Truncate(info.Size() - 1)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(dataToAdd)
	if err != nil {
		panic(err)
	}
}
