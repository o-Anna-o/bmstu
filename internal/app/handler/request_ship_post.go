package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// добавляет корабль в заявку
func (h *Handler) AddShipToRequestShip(c *gin.Context) {
	shipIDStr := c.Param("ship_id")
	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	request_ship, err := h.Repository.GetOrCreateUserDraft(1)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	err = h.Repository.AddShipToRequestShip(request_ship.ID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Корабль %d добавлен в заявку %d через ORM", shipID, request_ship.ID)
	c.Redirect(http.StatusFound, fmt.Sprintf("/request_ship/%d", request_ship.ID))
}

// логическое удаление заявки
func (h *Handler) DeleteRequestShip(c *gin.Context) {
	request_shipIDStr := c.Param("id")
	request_shipID, err := strconv.Atoi(request_shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteRequestShipSQL(request_shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Заявка %d удалена через SQL UPDATE (статус изменен на 'удалён')", request_shipID)
	c.Redirect(http.StatusFound, "/ships")
}

// удаляет корабль из заявки
func (h *Handler) RemoveShipFromRequestShip(c *gin.Context) {
	request_shipIDStr := c.Param("id")
	shipIDStr := c.Param("ship_id")

	request_shipID, err := strconv.Atoi(request_shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.RemoveShipFromRequestShip(request_shipID, shipID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/request_ship/"+request_shipIDStr)
}
