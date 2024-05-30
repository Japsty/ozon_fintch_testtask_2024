package mocks

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/services"
	"context"
	"github.com/stretchr/testify/mock"
)

type PostServiceMock struct {
	sessions    services.SessionService
	PostRepo    model.PostRepo
	CommentRepo model.CommentRepo
	mock.Mock
}

func NewPostServiceMock(sessions services.SessionService, postRepo model.PostRepo, commentRepo model.CommentRepo) *PostServiceMock {
	return &PostServiceMock{sessions: sessions, PostRepo: postRepo, CommentRepo: commentRepo}
}

func (ps *PostServiceMock) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	args := ps.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (ps *PostServiceMock) AddPost(ctx context.Context, title, text string, status bool) (model.Post, error) {
	args := ps.Called(ctx, title, text, status)
	return args.Get(0).(model.Post), args.Error(1)
}

func (ps *PostServiceMock) GetPostByPostID(ctx context.Context, postID string) (model.Post, error) {
	args := ps.Called(ctx, postID)
	return args.Get(0).(model.Post), args.Error(1)
}

func (ps *PostServiceMock) UpdatePostCommentsStatus(ctx context.Context, postID string, status bool) (model.Post, error) {
	args := ps.Called(ctx, postID, status)
	return args.Get(0).(model.Post), args.Error(1)
}
