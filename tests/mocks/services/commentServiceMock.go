package mocks

import (
	"Ozon_testtask/internal/model"
	"context"
	"github.com/stretchr/testify/mock"
)

type CommentServiceMock struct {
	CommentRepo model.CommentRepo
	PostRepo    model.PostRepo
	mock.Mock
}

func (cs *CommentServiceMock) CommentPost(ctx context.Context, id, text, parentID string) ([]*model.Comment, error) {
	args := cs.Called(ctx, id, text, parentID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (cs *CommentServiceMock) GetCommentsByPostID(ctx context.Context, id string, l, o int) ([]*model.Comment, error) {
	args := cs.Called(ctx, id, l, o)
	return args.Get(0).([]*model.Comment), args.Error(1)
}
