package handler

import (
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetShips(ctx *gin.Context) {
	var ships []ds.Ship
	var err error

	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		ships, err = h.Repository.GetShips()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ships, err = h.Repository.GetShipsByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// для подсчета кораблей в заявке - используем ORM
	request_shipCount := 0
	var request_shipID int // для ID заявки

	request_ship, err := h.Repository.GetOrCreateUserDraft(1)
	if err == nil {
		// Проверяем что заявка не удалена (GetOrCreateUserDraft уже исключает удаленные по статусу "черновик")
		logrus.Infof("Найдена заявка ID=%d, количество кораблей в заявке: %d", request_ship.ID, len(request_ship.Ships))
		for i, shipInRequestShip := range request_ship.Ships {
			logrus.Infof("Корабль %d: %s, количество: %d", i, shipInRequestShip.Ship.Name, shipInRequestShip.Count)
			request_shipCount += shipInRequestShip.Count
		}
		// СОХРАНЯЕМ ID ЗАЯВКИ ДЛЯ ШАБЛОНА
		request_shipID = request_ship.ID
	} else {
		logrus.Errorf("Ошибка получения заявки: %v", err)
		// Если ошибка, всё равно создаем/получаем заявку для получения ID
		request_ship, err = h.Repository.GetOrCreateUserDraft(1)
		if err == nil {
			request_shipID = request_ship.ID
		}
	}

	logrus.Infof("Итоговый счетчик для отображения: %d", request_shipCount)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"ships":              ships,
		"search":             searchQuery,
		"request_ship_count": request_shipCount,
		"request_ship_id":    request_shipID,
	})
}

func (h *Handler) GetShip(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id корабля из урла (то есть из /ship/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	ship, err := h.Repository.GetShip(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "ship.html", gin.H{
		"ship": ship,
	})
}
