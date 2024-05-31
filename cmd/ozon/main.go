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
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
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

	var pr model.PostRepo

	var cr model.CommentRepo

	storageType := os.Getenv("STORAGE")
	switch storageType {
	case "postgres":
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

		pr = postgre.NewPostRepository(postgresConnect)
		cr = postgre.NewCommentRepository(postgresConnect)
	case "inmemory":
		pr = inmem.NewPostInMemoryRepository()
		cr = inmem.NewCommentInMemoryRepository()
	default:
		logger.Fatal("Unknown storage type: ", storageType)
		return
	}

	ps := services.NewPostService(pr, cr)
	cs := services.NewCommentService(cr, pr)

	resolver := graph.NewResolver(ps, cs, logger)

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return middleware.AccessLog(logger, next)
	})
	r.Use(func(next http.Handler) http.Handler {
		return middleware.Auth(logger, next)
	})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Infof("Starting client on port: %s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatal("Server failed to start", err)
	}
}
