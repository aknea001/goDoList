package pkg

type User struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type CredentialError struct {
	Message string
}

func (e *CredentialError) Error() string {
	return e.Message
}
