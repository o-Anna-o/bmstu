package api

import (
	"loading_time/internal/app/config"
	"loading_time/internal/app/dsn"
	"loading_time/internal/app/handler"
	"loading_time/internal/app/repository"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	// Загрузка конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Создание DSN строки
	dsnString := dsn.FromEnv(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	log.Printf("Подключение к БД: %s@%s:%d/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Инициализация репозитория
	repo, err := repository.New(dsnString)
	if err != nil {
		logrus.Fatal("Ошибка инициализации репозитория:", err)
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	// Регистрация статики и шаблонов
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// Регистрация маршрутов
	registerRoutes(r, handler)

	log.Printf(" Сервер запущен на http://localhost:%d", cfg.ServicePort)
	r.Run(":" + strconv.Itoa(cfg.ServicePort))
}

func registerRoutes(r *gin.Engine, handler *handler.Handler) {
	r.GET("/ships", handler.GetShips)
	r.GET("/ship/:id", handler.GetShip)
	r.GET("/request", handler.GetRequest)
	r.GET("/request/:id", handler.GetRequest)
	r.GET("/request/add/:id", handler.AddToRequest)
	r.GET("/request/:id/remove/:ship_id", handler.RemoveShipFromRequest)
	r.POST("/request/:id/delete", handler.DeleteRequest)
}
