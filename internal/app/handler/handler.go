package handler

import (
	"context"
	"fmt"
	"loading_time/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

	// Создаем bucket
	bucketName := "loading_time_img"
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		logrus.Errorf("Bucket check error: %v", err)
	} else if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Errorf("Bucket creation error: %v", err)
		} else {
			logrus.Infof("Bucket %s created", bucketName)
		}
	}

	return &Handler{
		Repository:  r,
		MinioClient: client,
	}
}

func (h *Handler) GetShips(ctx *gin.Context) {
	var ships []repository.Ship
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		ships, err = h.Repository.GetShips()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ships, err = h.Repository.GetShipsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":  time.Now().Format("15:04:05"),
		"ships": ships,
		"query": searchQuery,
	})
}

func (h *Handler) GetShip(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
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

// НОВЫЙ метод для загрузки файлов в MinIO
func (h *Handler) UploadFile(ctx *gin.Context) {
	// Проверяем, что MinIO клиент инициализирован
	if h.MinioClient == nil {
		ctx.String(http.StatusInternalServerError, "Storage service unavailable")
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request: "+err.Error())
		return
	}
	defer file.Close()

	bucketName := "loading_time_img"
	objectName := header.Filename

	// Загружаем файл в MinIO
	_, err = h.MinioClient.PutObject(context.Background(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})

	if err != nil {
		logrus.Errorf("Upload error: %v", err)
		ctx.String(http.StatusInternalServerError, "Upload failed: "+err.Error())
		return
	}

	// Генерируем URL для просмотра изображения
	fileURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, objectName)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Message":  "File uploaded successfully to cloud storage!",
		"FileURL":  fileURL,
		"FileName": objectName,
	})
}

func getMinIOImageURL(fileName string) string {
	if fileName == "" {
		return "" // или URL заглушки
	}
	return fmt.Sprintf("http://localhost:9000/bmstu-photos/%s", fileName)
}
