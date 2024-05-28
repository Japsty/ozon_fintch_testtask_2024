package main

import (
	"Ozon_testtask/pkg/storage/connect"
	"Ozon_testtask/pkg/storage/migrate"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return
	}
	defer func(zapLogger *zap.Logger) {
		err = zapLogger.Sync()
		if err != nil {
			return
		}
	}(zapLogger)
	logger := zapLogger.Sugar()

	postgresConnect, err := connect.NewPostgresConnection()
	if err != nil {
		logger.Error("Connecting to SQL database error: ", err)
		return
	}
	defer postgresConnect.Close()

	repo := repos.New(db)

	//Поднимаем миграции
	err = migrate.UpMigration(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
	}

	router := handlers.NewGoodsHandler(repo, redis_repo, natsClient)

	// err = router.Run("localhost:8080") - если на локальной машине
	log.Println("Starting client on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server dropped")
	}
}
