package main

import (
	"Ozon_testtask/graph"
	"Ozon_testtask/internal/middleware"
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/inmem"
	"Ozon_testtask/internal/repos/postgre"
	"Ozon_testtask/internal/services"
	"Ozon_testtask/pkg/storage/connect"
	"Ozon_testtask/pkg/storage/migrate"
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file: ", err)
	//}

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

	var ur model.UserRepo
	var pr model.PostRepo
	var cr model.CommentRepo

	storageType := os.Getenv("STORAGE")
	if storageType == "postgres" {
		logger.Infof("postgres enabled")
		postgresConnect, err := connect.NewPostgresConnection()
		if err != nil {
			logger.Error("Connecting to SQL database error: ", err)
			return
		}
		defer postgresConnect.Close()

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = migrate.UpMigration(ctxWithTimeout, postgresConnect)
		if err != nil {
			logger.Fatal("Failed to up migration: ", err)
		}

		ur = postgre.NewUserRepository(postgresConnect)
		pr = postgre.NewPostRepository(postgresConnect)
		cr = postgre.NewCommentRepository(postgresConnect)
	} else if storageType == "inmemory" {
		ur = inmem.NewUserRepository()
		pr = inmem.NewPostInMemoryRepository()
		cr = inmem.NewCommentInMemoryRepository()
	} else {
		logger.Fatal("Unknown storage type: ", storageType)
		return
	}

	ps := services.NewPostService(ur, pr, cr)
	cs := services.NewCommentService(cr, pr)

	resolver := graph.NewResolver(ps, cs, logger)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	//srv.Use(middleware.AccessLog(logger))

	//http.Handle("/graphql", middleware.AccessLog(sugar, middleware.Auth(sugar, srv)))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.AccessLog(logger, srv))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// err = router.Run("localhost:8080") - если на локальной машине
	logger.Infof("Starting client on port: %s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		logger.Fatal("Server failed to start", err)
	}
}
