package mocks

import (
	"Ozon_testtask/internal/model"
	"context"
	"github.com/stretchr/testify/mock"
)

type CommentRepoMock struct {
	mock.Mock
}

func (m *CommentRepoMock) CreateComment(ctx context.Context, id, text, userID, postID, parentID string) (*model.Comment, error) {
	args := m.Called(ctx, id, text, userID, postID, parentID)
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *CommentRepoMock) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *CommentRepoMock) GetCommentsByPostIDPaginated(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.Comment), args.Error(1)
}
