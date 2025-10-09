package api

import (
	"loading_time/internal/app/ds"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestShipHandler struct {
	Repository interface {
		GetOrCreateUserDraft(userID int) (ds.RequestShip, error)
		DB() *gorm.DB
	}
}

// GetRequestsBasketAPI - GET /api/requests/basket - иконка корзины
func (h *RequestShipHandler) GetRequestsBasketAPI(c *gin.Context) {
	const fixedUserID = 1

	requestShip, err := h.Repository.GetOrCreateUserDraft(fixedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	totalShipsCount := 0
	for _, ship := range requestShip.Ships {
		totalShipsCount += ship.ShipsCount
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"request_ship_id": requestShip.RequestShipID,
			"ships_count":     totalShipsCount,
		},
	})
}

// GetRequestsAPI - GET /api/requests - список заявок с фильтрацией
func (h *RequestShipHandler) GetRequestsAPI(c *gin.Context) {
	statusFilter := c.Query("status")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	var requests []ds.RequestShip

	db := h.Repository.DB()

	query := db.Where("status != ? AND status != ?", "удалён", "черновик")

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	if dateFrom != "" {
		if from, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("formation_date >= ?", from)
		}
	}

	if dateTo != "" {
		if to, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("formation_date <= ?", to)
		}
	}

	err := query.Find(&requests).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   requests,
	})
}
