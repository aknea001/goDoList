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

type GetTasksRes struct {
	Msg   string `json:"msg"`
	Tasks []Task `json:"tasks"`
}
