package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/MXkodo/cash-server/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	adminToken  string
}

func NewAuthHandler(authService *service.AuthService, adminToken string) *AuthHandler {
	return &AuthHandler{authService, adminToken}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
		Login string `json:"login" binding:"required"`
		Pswd  string `json:"pswd" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные параметры"})
		return
	}

	if req.Token != h.adminToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав доступа"})
		return
	}

	if err := validatePassword(req.Pswd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(req.Login, req.Pswd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нежданчик"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req struct {
		Login string `json:"login" binding:"required"`
		Pswd  string `json:"pswd" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные параметры"})
		return
	}

	token, err := h.authService.Authenticate(req.Login, req.Pswd)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": gin.H{"token": token}})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.Param("token")

	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Нежданчик"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": gin.H{token: true}})
}

func validatePassword(password string) error {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit   = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if !hasMinLen || !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("пароль должен содержать минимум 8 символов, включая заглавные и строчные буквы, цифры и специальные символы")
	}
	return nil
}
