package handler

type TodoItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type CreateTodoResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type DeleteTodoByIDResponse struct {
	OK bool `json:"ok"`
}

type UpdateTodoByIDRequest struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UpdateTodoByIDResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
