package handler

import (
	"loading_time/internal/app/handler/api"
	"loading_time/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
	// Добавляем API handlers
	ShipAPIHandler        *api.ShipHandler
	RequestShipAPIHandler *api.RequestShipHandler
	UserAPIHandler        *api.UserHandler
}

func NewHandler(r *repository.Repository) *Handler {
	endpoint := "localhost:9000"
	accessKey := "minio_login_001"
	secretKey := "minio_login_001"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Errorf("MinIO connection error: %v", err)
		client = nil
	} else {
		logrus.Info("MinIO connected successfully")
	}

	return &Handler{
		Repository:  r,
		MinioClient: client,
		// Инициализируем API handlers
		ShipAPIHandler:        &api.ShipHandler{Repository: r},
		RequestShipAPIHandler: &api.RequestShipHandler{Repository: r},
		UserAPIHandler:        &api.UserHandler{Repository: r},
	}
}

// RegisterHandler регистрирует маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// HTML методы
	router.GET("/ships", h.GetShips)
	router.GET("/ship/:id", h.GetShip)
	router.GET("/request_ship", h.CreateOrRedirectRequestShip)
	router.GET("/request_ship/:id", h.GetRequestShip)

	router.POST("/request_ship/add/:ship_id", h.AddShipToRequestShip)
	router.POST("/request_ship/delete/:id", h.DeleteRequestShip)
	router.POST("/request_ship/:id/remove/:ship_id", h.RemoveShipFromRequestShip)

	// API маршруты через API handlers
	apiGroup := router.Group("/api")
	{
		// Домен услуги (контейнеровозы)
		apiGroup.GET("/ships", h.ShipAPIHandler.GetShipsAPI)
		apiGroup.GET("/ships/:id", h.ShipAPIHandler.GetShipAPI)

		// Домен заявки
		apiGroup.GET("/requests/basket", h.RequestShipAPIHandler.GetRequestsBasketAPI)
		apiGroup.GET("/requests", h.RequestShipAPIHandler.GetRequestsAPI)

		// Домен пользователь
		apiGroup.POST("/users/register", h.UserAPIHandler.RegisterUserAPI)
		apiGroup.GET("/users/profile", h.UserAPIHandler.GetUserProfileAPI)
	}
}

// RegisterStatic регистрирует статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

// errorHandler
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
