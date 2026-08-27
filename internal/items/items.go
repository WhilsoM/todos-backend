package items

type TodoItem struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Success bool   `json:"success"`
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
