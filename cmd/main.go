package main

import (
	"github/Burakyc/TO-DO/internal/auth"
	"github/Burakyc/TO-DO/internal/todo"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/login", auth.Login)
	protected := r.Group("/api")
	protected.Use(auth.JWTAuthMiddleware())

	protected.GET("/todos", todo.GetTodos)
	protected.POST("/todos", todo.CreateTodo)
	protected.PUT("/todos/:id", todo.UpdateTodo)
	protected.DELETE("/todos/:id", todo.DeleteTodo)

	protected.POST("/todos/:id/items", todo.CreateTodoItem)
	protected.PUT("/items/:id", todo.UpdateTodoItem)
	protected.DELETE("/items/:id", todo.DeleteTodoItem)
	protected.GET("/todos/:id/items", todo.GetTodoItems)

	r.Run(":5000")
}
