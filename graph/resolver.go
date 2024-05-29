package graph

import (
	"Ozon_testtask/internal/models"
	"go.uber.org/zap"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const TimeoutTime = 500 * time.Millisecond

type Resolver struct {
	PostService    models.PostService
	CommentService models.CommentService
	Logger         *zap.SugaredLogger
}
