package pkg

import "fmt"

type CredentialError struct{}

func (e *CredentialError) Error() string {
	return "wrong username or password"
}

type UnknownServerError struct{}

func (e *UnknownServerError) Error() string {
	return "unknown error"
}

type DoesntExistError struct {
	ResourceName string
}

func (e *DoesntExistError) Error() string {
	eMsg := fmt.Sprintf("You either lack permission or the %s may not exist", e.ResourceName)
	return eMsg
}
