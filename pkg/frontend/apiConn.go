package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aknea001/goDoList/pkg"
)

type APIconn struct {
	BaseURL string
	Header  http.Header
	Token   string
}

func (api APIconn) Register(username string, passwd string) error {
	fullUrl := fmt.Sprintf("%s/users", api.BaseURL)

	newUser := pkg.User{
		Username: username, Passwd: passwd,
	}

	byteValue, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(byteValue)

	res, err := http.Post(fullUrl, "application/json", body)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("%d: unknown error", res.StatusCode)
	}

	return nil
}

func (api APIconn) Login(username string, passwd string) error {
	fullUrl := fmt.Sprintf("%s/auth/login", api.BaseURL)

	newUser := pkg.User{
		Username: username, Passwd: passwd,
	}

	byteValue, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	body := bytes.NewReader(byteValue)

	res, err := http.Post(fullUrl, "application/json", body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 401 {
		return &pkg.CredentialError{}
	} else if res.StatusCode != 200 {
		return fmt.Errorf("%d: unknown error", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var newLoginRes pkg.LoginRes

	err = json.Unmarshal(data, &newLoginRes)
	if err != nil {
		return err
	}

	api.Token = newLoginRes.Token

	bearerToken := fmt.Sprintf("Bearer %s", api.Token)
	api.Header.Set("Authorization", bearerToken)

	return nil
}

func (api APIconn) GetTasks() ([]pkg.Task, error) {
	fullUrl := fmt.Sprintf("%s/tasks", api.BaseURL)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header = api.Header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newGetTasksRes pkg.GetTasksRes

	err = json.Unmarshal(data, &newGetTasksRes)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%d: %s", res.StatusCode, newGetTasksRes.Msg)
	}

	return newGetTasksRes.Tasks, nil
}

func (api APIconn) NewTask(newTask pkg.Task) error {
	fullUrl := fmt.Sprintf("%s/tasks", api.BaseURL)

	byteValue, err := json.Marshal(newTask)
	if err != nil {
		return err
	}

	body := bytes.NewReader(byteValue)

	req, err := http.NewRequest("POST", fullUrl, body)
	if err != nil {
		return err
	}

	req.Header = api.Header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var defaultRes pkg.DefaultRes

	err = json.Unmarshal(data, &defaultRes)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%d: %s", res.StatusCode, defaultRes.Msg)
	}

	return nil
}

func (api APIconn) FinishTask(task pkg.Task) error {
	fullUrl := fmt.Sprintf("%s/tasks/%s", api.BaseURL, url.PathEscape(task.Title))

	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return err
	}

	req.Header = api.Header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	var defaultRes pkg.DefaultRes
	err = json.Unmarshal(data, &defaultRes)
	if err != nil {
		return err
	}

	if res.StatusCode == 403 {
		return &pkg.DoesntExistError{
			ResourceName: "task",
		}
	} else if res.StatusCode != 200 {
		return fmt.Errorf("%d: %s", res.StatusCode, defaultRes.Msg)
	}

	return nil
}

func NewAPIconn(baseURL string) APIconn {
	var newAPIconn APIconn

	newAPIconn.BaseURL = baseURL
	newAPIconn.Header = http.Header{}

	return newAPIconn
}
