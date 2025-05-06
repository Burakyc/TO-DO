package todo

import "time"

type ToDoList struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	CompletionRate float32    `json:"completion_rate"`
}

type TodoItem struct {
	ID        string     `json:"id"`
	ListID    string     `json:"list_id"`
	Content   string     `json:"content"`
	Done      bool       `json:"done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
