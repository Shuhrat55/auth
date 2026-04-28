package app

import (
	_ "github.com/lib/pq"
	"github.com/Shuhrat55/auth/internal/delivery/gin"
	"github.com/Shuhrat55/auth/internal/repository"
	"github.com/Shuhrat55/auth/internal/usecase"
	"github.com/Shuhrat55/auth/pkg/database"
	"github.com/Shuhrat55/auth/pkg/logger"
	"go.uber.org/zap"
	"log"
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

	if err := router.Run(":8080"); err != nil {
		logger.Logger.Fatal("Ошибка запуска сервера на порту :8080",
			zap.Error(err),
			zap.String("app", "database"))
	}
	logger.Logger.Info("Микросервис стартует на порту :8080")
}
