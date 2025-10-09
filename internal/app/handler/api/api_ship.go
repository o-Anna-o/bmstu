package api

import (
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShipHandler struct {
	Repository interface {
		GetShips() ([]ds.Ship, error)
		GetShipsByName(name string) ([]ds.Ship, error)
		GetShip(id int) (ds.Ship, error)
	}
}

// GetShipsAPI - GET /api/ships - список кораблей с фильтрацией
func (h *ShipHandler) GetShipsAPI(c *gin.Context) {
	nameFilter := c.Query("name")

	var ships []ds.Ship
	var err error

	if nameFilter == "" {
		ships, err = h.Repository.GetShips()
	} else {
		ships, err = h.Repository.GetShipsByName(nameFilter)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":      "error",
			"description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ships,
	})
}

// GetShipAPI - GET /api/ships/:id - один корабль
func (h *ShipHandler) GetShipAPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "Invalid ship ID",
		})
		return
	}

	ship, err := h.Repository.GetShip(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":      "error",
			"description": "Ship not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ship,
	})
}
