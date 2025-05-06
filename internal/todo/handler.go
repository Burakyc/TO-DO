package todo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// =================== REQUEST STRUCT'ları ===================
type CreateTodoRequest struct {
	Name string `json:"name" binding:"required"`
}

// =================== TO-DO LIST HANDLERLARI ===================
func GetTodos(context *gin.Context) {
	username := context.GetString("username")
	role := context.GetString("role")

	var updatedTodoLists []ToDoList

	for _, todolist := range TodoLists {
		if todolist.DeletedAt != nil {
			continue
		}

		if role != "admin" && todolist.CreatedBy != username {
			continue
		}

		totalItems := 0
		doneItems := 0
		for _, item := range TodoItems {
			if item.ListID == todolist.ID && item.DeletedAt == nil {
				totalItems++
				if item.Done {
					doneItems++
				}
			}
		}

		todolist.CompletionRate = 0
		if totalItems > 0 {
			todolist.CompletionRate = (float32(doneItems) / float32(totalItems)) * 100
		}

		updatedTodoLists = append(updatedTodoLists, todolist)
	}

	context.JSON(http.StatusOK, updatedTodoLists)
}

func CreateTodo(context *gin.Context) {
	var req CreateTodoRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "geçersiz veri"})
		return
	}

	role := context.GetString("role")
	if role == "admin" {
		context.JSON(http.StatusForbidden, gin.H{"error": "Admin kullanıcıları liste oluşturamaz"})
		return
	}

	username := context.GetString("username")

	newTodo := ToDoList{
		ID:             uuid.New().String(),
		Name:           req.Name,
		CreatedBy:      username,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CompletionRate: 0,
	}

	TodoLists = append(TodoLists, newTodo)
	context.JSON(http.StatusCreated, newTodo)
}

func UpdateTodo(context *gin.Context) {
	id := context.Param("id")
	var req CreateTodoRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	found := false
	for i, todolist := range TodoLists {
		if todolist.ID == id && todolist.DeletedAt == nil {
			TodoLists[i].Name = req.Name
			TodoLists[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "To-Do Listesi bulunamadı"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "To-Do Listesi güncellendi"})
}

func DeleteTodo(context *gin.Context) {
	id := context.Param("id")

	found := false
	for i, todolist := range TodoLists {
		if todolist.ID == id && todolist.DeletedAt == nil {
			now := time.Now()
			TodoLists[i].DeletedAt = &now
			found = true
			break
		}
	}

	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "To-Do Listesi bulunamadı veya zaten silinmiş"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "To-Do Listesi silindi"})
}
