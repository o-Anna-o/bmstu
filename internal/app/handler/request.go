package handler

import (
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetRequest(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заявки из урла (то есть из /request/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше

	if idStr == "" {
		ctx.HTML(http.StatusOK, "request.html", gin.H{
			"request": ds.Request{
				ID:    0,
				Ships: []ds.ShipInRequest{},
			},
		})
		return
	}

	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	request, err := h.Repository.GetRequest(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}
