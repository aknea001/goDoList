package backend

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aknea001/goDoList/pkg"
)

func FirstJson(byteValue []byte, file *os.File) error {
	jsonList := fmt.Appendf(nil, "[\n%s\n]", byteValue)

	_, err := file.Write(jsonList)
	if err != nil {
		return err
	}

	return nil
}

func AppendJson(byteValue []byte, file *os.File) error {
	formattedData := fmt.Appendf(nil, "\n,\n%s\n]", byteValue)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = file.Truncate(fileInfo.Size() - 2)
	if err != nil {
		return err
	}

	_, err = file.Write(formattedData)
	if err != nil {
		return err
	}

	return nil
}

func RegisterJson(username string, passwd string) error {
	// hash passwd !!
	newUser := pkg.User{
		Username: username, Passwd: passwd,
	}

	byteValue, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	jsonFileName := "user.json"
	jsonFile, err := os.OpenFile(
		jsonFileName,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	fileInfo, err := jsonFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() <= 0 {
		err = FirstJson(byteValue, jsonFile)
		if err != nil {
			return err
		}
		return nil
	}

	err = AppendJson(byteValue, jsonFile)
	if err != nil {
		return err
	}

	return nil
}

func LoginJson(username string, passwd string) error {
	jsonFileName := "user.json"

	jsonFile, err := os.OpenFile(
		jsonFileName,
		os.O_RDONLY,
		0,
	)
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		return &pkg.CredentialError{}
	}

	defer jsonFile.Close()

	jsonScanner := bufio.NewScanner(jsonFile)
	for jsonScanner.Scan() {
		currentBytes := jsonScanner.Bytes()

		if len(currentBytes) <= 1 {
			continue
		}

		var currentUser pkg.User

		err := json.Unmarshal(currentBytes, &currentUser)
		if err != nil {
			return err
		}

		if currentUser.Username != username {
			continue
		}

		if currentUser.Passwd != passwd {
			return &pkg.CredentialError{}
		}

		return nil
	}

	return &pkg.CredentialError{}
}

func GetTaskJson(user string) ([]pkg.Task, error) {
	jsonFileName := "task.json"

	jsonFile, err := os.OpenFile(
		jsonFileName,
		os.O_RDONLY,
		0,
	)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	tasks := make([]pkg.Task, 0)

	jsonScanner := bufio.NewScanner(jsonFile)
	for jsonScanner.Scan() {
		currentBytes := jsonScanner.Bytes()

		if len(currentBytes) <= 1 {
			continue
		}

		var currentTask pkg.Task

		err := json.Unmarshal(currentBytes, &currentTask)
		if err != nil {
			return nil, err
		}

		if currentTask.Owner != user {
			continue
		}

		tasks = append(tasks, currentTask)
	}

	return tasks, nil
}

func NewTaskJson(task pkg.Task) error {
	jsonFileName := "task.json"

	byteValue, err := json.Marshal(task)
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(
		jsonFileName,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	fileInfo, err := jsonFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() <= 0 {
		err = FirstJson(byteValue, jsonFile)
		if err != nil {
			return err
		}
		return nil
	}

	err = AppendJson(byteValue, jsonFile)
	if err != nil {
		return err
	}
	return nil
}
