package handler

import (
	"fmt"
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/request_ship", func(ctx *gin.Context) {
		request_ship, err := h.Repository.GetOrCreateUserDraft(1)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusInternalServerError, "request_ship.html", gin.H{
				"request_ship": ds.RequestShip{},
				"error":        "Не удалось создать черновик",
			})
			return
		}
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/request_ship/%d", request_ship.ID))
	})

	// Основной маршрут для конкретной заявки
	router.GET("/request_ship/:id", h.GetRequestShip)
}

func (h *Handler) GetRequestShip(ctx *gin.Context) {
	idStr := ctx.Param("id")

	request_shipID, err := strconv.Atoi(idStr)
	if err != nil || request_shipID <= 0 {
		logrus.Errorf("Неверный ID заявки: %v", idStr)
		ctx.HTML(http.StatusBadRequest, "request_ship.html", gin.H{
			"request_ship": ds.RequestShip{},
			"error":        "Некорректный ID заявки",
		})
		return
	}

	request_ship, err := h.Repository.GetRequestShipExcludingDeleted(request_shipID)
	if err != nil {
		logrus.Errorf("Заявка не найдена или удалена: %v", err)
		ctx.HTML(http.StatusNotFound, "request_ship.html", gin.H{
			"request_ship": ds.RequestShip{},
			"error":        "Заявка не найдена",
		})
		return
	}

	ctx.HTML(http.StatusOK, "request_ship.html", gin.H{
		"request_ship": request_ship,
	})
}
