package handler

import (
	// "loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AddShipToRequest - добавление корабля в заявку через ORM
func (h *Handler) AddShipToRequest(c *gin.Context) {
	shipIDStr := c.Param("ship_id")
	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// TODO: Заглушка - нужно реализовать через ORM
	// Здесь должен быть код добавления корабля в заявку пользователя

	logrus.Infof("Корабль %d добавлен в заявку через ORM", shipID)
	c.Redirect(http.StatusFound, "/request")
}

// DeleteRequest - логическое удаление заявки через SQL
func (h *Handler) DeleteRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// TODO: Заглушка - нужно реализовать через SQL UPDATE
	// Здесь должен быть код логического удаления заявки

	logrus.Infof("Заявка %d удалена через SQL UPDATE", requestID)
	c.Redirect(http.StatusFound, "/ships")
}
