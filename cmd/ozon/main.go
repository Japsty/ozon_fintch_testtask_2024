package main

import (
	graph2 "Ozon_testtask/internal/graph"
	"Ozon_testtask/internal/middleware"
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/inmem"
	"Ozon_testtask/internal/repos/postgre"
	"Ozon_testtask/internal/services"
	"Ozon_testtask/pkg/storage/connect"
	"Ozon_testtask/pkg/storage/migrate"
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
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
		logger.Infof("inmemory enabled")
		pr = inmem.NewPostInMemoryRepository()
		cr = inmem.NewCommentInMemoryRepository()
	default:
		logger.Fatal("Unknown storage type: ", storageType)
		return
	}

	ps := services.NewPostService(pr, cr)
	cs := services.NewCommentService(cr, pr)

	resolver := graph2.NewResolver(ps, cs, logger)

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return middleware.Auth(logger, next)
	})

	srv := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: resolver}))

	srv.AroundOperations(
		func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			oc := graphql.GetOperationContext(ctx)
			logger.Infof(
				"GraphQL request. Operation Name: %v, Operation Query: %s",
				oc.Operation.Name,
				oc.RawQuery,
			)
			return next(ctx)
		},
	)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger.Infof("Starting on port: %s", port)

	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			logger.Fatal("Server failed to start", err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	logger.Infof("Shutting down...")
	os.Exit(0)
}
