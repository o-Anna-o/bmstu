package handler

import (
	// "loading_time/internal/app/ds"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) AddShipToRequest(c *gin.Context) {
	shipIDStr := c.Param("ship_id")
	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Получаем или создаем черновую заявку для пользователя
	request, err := h.Repository.GetOrCreateUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Добавляем корабль в заявку через ORM
	err = h.Repository.AddShipToRequest(request.ID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Корабль %d добавлен в заявку %d через ORM", shipID, request.ID)

	c.Redirect(http.StatusFound, fmt.Sprintf("/request/%d", request.ID))
}

// DeleteRequest - логическое удаление заявки через SQL
func (h *Handler) DeleteRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Логическое удаление заявки через SQL UPDATE (требование лабораторной)
	err = h.Repository.DeleteRequestSQL(requestID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Заявка %d удалена через SQL UPDATE (статус изменен на 'удалён')", requestID)
	c.Redirect(http.StatusFound, "/ships")
}

// RemoveShipFromRequest - удаление корабля из заявки
func (h *Handler) RemoveShipFromRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	shipIDStr := c.Param("ship_id")

	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.RemoveShipFromRequest(requestID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/request/"+requestIDStr)
}
