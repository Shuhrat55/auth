package app

import (
	"github.com/Shuhrat55/auth/internal/delivery/gin"
	"github.com/Shuhrat55/auth/internal/repository"
	"github.com/Shuhrat55/auth/internal/usecase"
	"github.com/Shuhrat55/auth/pkg/database"
	"github.com/Shuhrat55/auth/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os"
)

func Run() {
	db, err := database.NewSQLiteConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.SetupRouter(userUseCase)

	port := os.Getenv("AUTH_HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		logger.Logger.Fatal("Ошибка запуска сервера",
			zap.Error(err),
			zap.String("app", "database"))
	}
	logger.Logger.Info("Микросервис auth стартует",
		zap.String("port", port))
}
