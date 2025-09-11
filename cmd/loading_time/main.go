package main

// go run cmd/loading_time/main.go

import (
	"context"
	"loading_time/internal/api"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

// функция инициализации MinIO
func InitMinIO() {
	endpoint := "localhost:9000"
	accessKey := "minio_login_001"
	secretKey := "minio_login_001"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("MinIO connection error: %v", err)
	}

	minioClient = client
	log.Println("MinIO connected successfully")

	bucketName := "bmstu-photos"
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Printf("Bucket check error: %v", err)
		return
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Bucket creation error: %v", err)
		} else {
			log.Printf("Bucket %s created", bucketName)
		}
	}
}

func main() {
	log.Println("Application started!")
	api.StartServer()
	log.Println("Application terminated!")
}
