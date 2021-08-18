package server

type SignupInfo struct {
	LoginName string `json:"login_name"`
	ShowName  string `json:"show_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginInfo struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type UpdateImage struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}
