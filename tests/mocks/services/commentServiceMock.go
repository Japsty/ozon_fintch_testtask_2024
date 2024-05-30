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

func NewCommentService(commentRepo model.CommentRepo, postRepo model.PostRepo) *CommentServiceMock {
	return &CommentServiceMock{CommentRepo: commentRepo, PostRepo: postRepo}
}

func (cs *CommentServiceMock) CommentPost(ctx context.Context, postID, commentText, parentCommentID string) ([]*model.Comment, error) {
	args := cs.Called(ctx, postID, commentText, parentCommentID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (cs *CommentServiceMock) GetCommentByParentID(ctx context.Context, parentID string) ([]*model.Comment, error) {
	args := cs.Called(ctx, parentID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (cs *CommentServiceMock) GetCommentsByPostID(ctx context.Context, postID string, limit int, offset int) ([]*model.Comment, error) {
	args := cs.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.Comment), args.Error(1)
}
