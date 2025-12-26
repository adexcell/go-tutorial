package handler

import (
	"errors"
	"net/http"

	"github.com/adexcell/go-tutorial/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	u := &RegisterRequest{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(c, u.Email, u.Password)
	if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		c.JSON(http.StatusConflict, nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully created",
	})
}
