package graph

import (
	"Ozon_testtask/internal/model"
	"go.uber.org/zap"
	"sync"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const TimeoutTime = 500 * time.Millisecond

type Resolver struct {
	subscribers    map[string]chan *model.Comment
	mutex          sync.Mutex
	PostService    model.PostService
	CommentService model.CommentService
	Logger         *zap.SugaredLogger
}

func NewResolver(ps model.PostService, cs model.CommentService, logger *zap.SugaredLogger) *Resolver {
	return &Resolver{
		subscribers:    make(map[string]chan *model.Comment),
		mutex:          sync.Mutex{},
		PostService:    ps,
		CommentService: cs,
		Logger:         logger,
	}
}
