package handler

import (
	"context"
	"loading_time/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository  *repository.Repository
	MinioClient *minio.Client
}

func NewHandler(r *repository.Repository) *Handler {
	// Инициализация MinIO
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
		return &Handler{
			Repository:  r,
			MinioClient: nil,
		}
	}

	logrus.Info("MinIO connected successfully")

	// проверка подключения к бакету loading-time-img
	bucketName := "loading-time-img"
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		logrus.Errorf("Bucket check error: %v", err)
	} else if exists {
		logrus.Infof("Bucket %s is available", bucketName)
	} else {
		logrus.Warnf("Bucket %s does not exist", bucketName)
	}

	return &Handler{
		Repository:  r,
		MinioClient: client,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/ships", h.GetShips)
	router.GET("/ship/:id", h.GetShip)
	router.GET("/request/:id", h.GetRequest)

	// POST маршруты
	router.POST("/request/add/:ship_id", h.AddShipToRequest)
	router.POST("/request/delete/:id", h.DeleteRequest)

	router.POST("/request/:id/remove/:ship_id", h.RemoveShipFromRequest) //для удаления корабля из заявки
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
