package handler

import (
	"net/http"
	"time"

	"github.com/adexcell/go-tutorial/internal/domain"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service domain.NotificationService
}

type NotificationRequest struct {
	Message string    `json:"message" binding:"required"`
	SendAt  time.Time `json:"send_at" binding:"required"`
}

func NewNotificationHandler(service domain.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Schedule(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, ok := val.(int64)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}


	nR := &NotificationRequest{}

	if err := c.ShouldBindJSON(&nR); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n := &domain.Notification{
		UserID: userID,
		Message: nR.Message,
		SendAt: nR.SendAt,
	}

	if err := h.service.Schedule(c, n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
