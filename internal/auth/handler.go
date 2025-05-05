package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(context *gin.Context) {
	req, err := bindLoginRequest(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	user := authenticateUser(req.Username, req.Password)
	if user == nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı adı veya şifre"})
		return
	}

	token, err := GenerateJWT(user.Username, user.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Token üretilmedi"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

func bindLoginRequest(context *gin.Context) (*LoginRequest, error) {
	var req LoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func authenticateUser(username, password string) *User {
	for _, u := range users {
		if u.Username == username && u.Password == password {
			return &u
		}
	}
	return nil
}
