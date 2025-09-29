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

	// для подсчета кораблей в заявке
	requestCount := 0
	request, err := h.Repository.GetOrCreateUserDraft(1)
	if err == nil {
		// ДОБАВЬТЕ ОТЛАДОЧНУЮ ИНФОРМАЦИЮ
		logrus.Infof("Найдена заявка ID=%d, количество кораблей в заявке: %d", request.ID, len(request.Ships))
		for i, shipInRequest := range request.Ships {
			logrus.Infof("Корабль %d: %s, количество: %d", i, shipInRequest.Ship.Name, shipInRequest.Count)
			requestCount += shipInRequest.Count
		}
	} else {
		logrus.Errorf("Ошибка получения заявки: %v", err)
	}

	logrus.Infof("Итоговый счетчик для отображения: %d", requestCount)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"ships":         ships,
		"search":        searchQuery,
		"request_count": requestCount,
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
