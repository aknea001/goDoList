package pkg

type User struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type Task struct {
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"desc"`
}

type LoginRes struct {
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

type CredentialError struct{}

func (e *CredentialError) Error() string {
	return "wrong username or password"
}

type UnknownServerError struct{}

func (e *UnknownServerError) Error() string {
	return "unknown error"
}
