package api

import (
	"loading_time/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repository interface {
		DB() *gorm.DB
	}
}

// RegisterUserAPI - POST /api/users/register - регистрация
func (h *UserHandler) RegisterUserAPI(c *gin.Context) {
	var user ds.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	db := h.Repository.DB()

	err := db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   user,
	})
}

// GetUserProfileAPI - GET /api/users/profile - поля пользователя
func (h *UserHandler) GetUserProfileAPI(c *gin.Context) {
	const fixedUserID = 1

	var user ds.User

	db := h.Repository.DB()

	err := db.First(&user, fixedUserID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}
