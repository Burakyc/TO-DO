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
type CreateTodoItemRequest struct {
	Content string `json:"content" binding:"required"`
}
type UpdateTodoItemRequest struct {
	Content string `json:"content,omitempty"`
	Done    bool   `json:"done"`
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

// ================== TO-DO ITEM HANDLERLARI ==================

func CreateTodoItem(context *gin.Context) {
	listID := context.Param("id")

	var req CreateTodoItemRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	// Liste var mı kontrol et
	found := false
	for _, todolist := range TodoLists {
		if todolist.ID == listID && todolist.DeletedAt == nil {
			found = true
			break
		}
	}
	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "To-Do Listesi bulunamadı"})
		return
	}

	newItem := TodoItem{
		ID:        uuid.New().String(),
		ListID:    listID,
		Content:   req.Content,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	TodoItems = append(TodoItems, newItem)

	context.JSON(http.StatusCreated, newItem)
}
func UpdateTodoItem(context *gin.Context) {
	id := context.Param("id")

	var req UpdateTodoItemRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	found := false
	for i, item := range TodoItems {
		if item.ID == id && item.DeletedAt == nil {
			if req.Content != "" {
				TodoItems[i].Content = req.Content
			}
			TodoItems[i].Done = req.Done
			TodoItems[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "Adım bulunamadı"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Adım güncellendi"})
}
func DeleteTodoItem(context *gin.Context) {
	id := context.Param("id")

	found := false
	for i, item := range TodoItems {
		if item.ID == id && item.DeletedAt == nil {
			now := time.Now()
			TodoItems[i].DeletedAt = &now
			found = true
			break
		}
	}

	if !found {
		context.JSON(http.StatusNotFound, gin.H{"error": "Adım bulunamadı veya zaten silinmiş"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Adım silindi"})
}
func GetTodoItems(context *gin.Context) {
	listID := context.Param("id")

	var items []TodoItem
	for _, item := range TodoItems {
		if item.ListID == listID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}

	context.JSON(http.StatusOK, items)
}
