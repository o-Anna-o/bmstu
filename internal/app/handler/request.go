package handler

import (
	"loading_time/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var requestID int
	var err error

	if idStr == "" {
		// Если ID не указан в URL (/request), используем заявку пользователя
		request, err := h.Repository.GetOrCreateUserDraft(1)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusOK, "request.html", gin.H{
				"request": ds.Request{},
			})
			return
		}
		requestID = request.ID
	} else {
		// Если ID указан в URL (/request/1)
		requestID, err = strconv.Atoi(idStr)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusOK, "request.html", gin.H{
				"request": ds.Request{},
			})
			return
		}
	}

	request, err := h.Repository.GetRequest(requestID)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusOK, "request.html", gin.H{
			"request": ds.Request{},
		})
		return
	}

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}
