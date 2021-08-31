package server

// SignupInfo user for router bind json
type SignupInfo struct {
	LoginName string `json:"login_name"`
	ShowName  string `json:"show_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// LoginInfo user for router bind json
type LoginInfo struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

// UpdateImage user for router bind json
type UpdateImage struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}
