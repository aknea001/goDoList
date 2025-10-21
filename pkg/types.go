package pkg

type User struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type Task struct {
	Owner       string `json:"owner"`
	Title       string `json:"title" binding:"required,min=2,max=20"`
	Description string `json:"desc"`
}

type LoginRes struct {
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

type GetTasksRes struct {
	Msg   string `json:"msg"`
	Tasks []Task `json:"tasks"`
}

type DefaultRes struct {
	Msg string `json:"msg"`
}
