package mocks

import (
	"Ozon_testtask/internal/model"
	"context"
	"github.com/stretchr/testify/mock"
)

type PostRepoMock struct {
	mock.Mock
}

func (pr *PostRepoMock) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	args := pr.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (pr *PostRepoMock) CreatePost(ctx context.Context, id, title, text, uID string, status bool) (model.Post, error) {
	args := pr.Called(ctx, id, title, text, uID, status)
	return args.Get(0).(model.Post), args.Error(1)
}

func (pr *PostRepoMock) GetPostByPostID(ctx context.Context, id string) (model.Post, error) {
	args := pr.Called(ctx, id)
	return args.Get(0).(model.Post), args.Error(1)
}

func (pr *PostRepoMock) UpdatePostStatus(ctx context.Context, id, uID string, status bool) (model.Post, error) {
	args := pr.Called(ctx, id, uID, status)
	return args.Get(0).(model.Post), args.Error(1)
}
