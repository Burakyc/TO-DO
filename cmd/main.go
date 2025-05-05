package main

import (
	"github/Burakyc/TO-DO/internal/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/login", auth.Login)
	protected := r.Group("/api")
	protected.Use(auth.JWTAuthMiddleware())

	r.Run(":5000")
}
