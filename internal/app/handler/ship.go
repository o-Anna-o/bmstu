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

	searchQuery := ctx.Query("search") // получаем значение из нашего поля
	if searchQuery == "" {             // если поле поиска пусто, то просто получаем из репозитория все записи
		ships, err = h.Repository.GetShips()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ships, err = h.Repository.GetShipsByName(searchQuery) // в ином случае ищем корабль по названию
		if err != nil {
			logrus.Error(err)
		}
	}

	// для подсчета кораблей в заявке
	request, err := h.Repository.GetRequest(1)
	requestCount := 0
	if err == nil {
		for _, shipInRequest := range request.Ships {
			requestCount += shipInRequest.Count
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"ships":         ships,
		"search":        searchQuery,  // передаем введенный запрос обратно на страницу
		"request_count": requestCount, // количество кораблей в заявке
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
