package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type User struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
}

func FirstUser(byteValue []byte, file *os.File) error {
	jsonList := fmt.Appendf(nil, "[%s]", byteValue)

	_, err := file.Write(jsonList)
	if err != nil {
		return err
	}

	return nil
}

func AppendUser(byteValue []byte, file *os.File) error {
	formattedData := fmt.Appendf(nil, ",%s]", byteValue)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = file.Truncate(fileInfo.Size() - 1)
	if err != nil {
		return err
	}

	_, err = file.Write(formattedData)
	if err != nil {
		return err
	}

	return nil
}

func AddJsonUser(username string, passwd string) {
	// hash passwd !!
	newUser := []User{
		{Username: username, Hash: passwd},
	}

	byteValue, err := json.Marshal(newUser[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(byteValue)

	jsonFileName := "users.json"
	jsonFile, err := os.OpenFile(
		jsonFileName,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	fileInfo, err := jsonFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.Size() <= 0 {
		err = FirstUser(byteValue, jsonFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = AppendUser(byteValue, jsonFile)
	if err != nil {
		log.Fatal(err)
	}
}
