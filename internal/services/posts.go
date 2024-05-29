package services

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")

type PostService struct {
	UserRepo    models.UserRepo
	PostRepo    models.PostRepo
	CommentRepo models.CommentRepo
}

func NewPostService(userRepo models.UserRepo, postRepo models.PostRepo, commentRepo models.CommentRepo) *PostService {
	return &PostService{UserRepo: userRepo, PostRepo: postRepo, CommentRepo: commentRepo}
}

func (ps *PostService) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := ps.PostRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (ps *PostService) AddPost(ctx context.Context, Title, Text string, status bool) (models.Post, error) {
	newPostID := uuid.NewString()
	post, err := ps.PostRepo.CreatePost(ctx, newPostID, Title, Text, status)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (ps *PostService) GetPostByPostID(ctx context.Context, postID string) (models.Post, error) {
	posts, err := ps.PostRepo.GetPostByPostID(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}
	return posts, nil
}

func (ps *PostService) UpdatePostCommentsStatus(ctx context.Context, postID string, status bool) (models.Post, error) {
	post, err := ps.PostRepo.UpdatePostCommentsStatus(ctx, postID, status)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}
