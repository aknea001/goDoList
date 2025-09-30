package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aknea001/goDoList/pkg"
)

type APIconn struct {
	BaseURL string
	Header  http.Header
	Token   string
}

func (api APIconn) Register(username string, passwd string) error {
	fullUrl := fmt.Sprintf("%s/register", api.BaseURL)

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
		return fmt.Errorf("%s: unknown error", strconv.Itoa(res.StatusCode))
	}

	return nil
}

func (api APIconn) Login(username string, passwd string) error {
	fullUrl := fmt.Sprintf("%s/login", api.BaseURL)

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
		return fmt.Errorf("%s: unknown error", strconv.Itoa(res.StatusCode))
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

func NewAPIconn(baseURL string) APIconn {
	var newAPIconn APIconn

	newAPIconn.BaseURL = baseURL
	newAPIconn.Header = http.Header{}

	return newAPIconn
}
