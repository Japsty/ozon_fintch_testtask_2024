package services

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"sort"
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

	sort.Slice(posts, func(i, j int) bool { return posts[i].Score > posts[j].Score })

	return posts, nil
}

func (ps *PostService) AddPost(ctx context.Context, postAuthor models.Author, postInfo models.PostRequest) (models.Post, error) {
	newPostID := uuid.NewString()
	post, err := ps.PostRepo.CreatePost(ctx, postAuthor, newPostID, postInfo.PostType, postInfo.Category, postInfo.Title, postInfo.Text, postInfo.URL)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (ps *PostService) GetPostsCategorized(ctx context.Context, category string) ([]models.Post, error) {
	posts, err := ps.PostRepo.GetPostsCategorized(ctx, category)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (ps *PostService) GetPostWithComment(ctx context.Context, postID string) (models.Post, error) {
	foundPost, err := ps.PostRepo.PostUpdateViews(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}
	return foundPost, nil
}

func (ps *PostService) UpdatePostVote(ctx context.Context, userID, postID string, postVote int) (models.Post, error) {
	post, err := ps.PostRepo.UpdatePostVote(ctx, userID, postID, postVote)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (ps *PostService) DeletePost(ctx context.Context, userID, postID string) (models.DeleteMessage, error) {
	post, err := ps.PostRepo.GetPostByID(ctx, postID)
	if err != nil {
		return models.DeleteMessage{Message: "error"}, err
	}

	if post.Author.UserID != userID {
		return models.DeleteMessage{Message: "forbidden"}, ErrForbidden
	}

	err = ps.PostRepo.DeletePostByID(ctx, postID)
	if err != nil {
		return models.DeleteMessage{}, err
	}
	response := models.DeleteMessage{Message: "success"}

	return response, nil
}

func (ps *PostService) GetAllUsersPosts(ctx context.Context, userLogin string) ([]models.Post, error) {
	response, err := ps.PostRepo.GetAllUsersPosts(ctx, userLogin)
	if err != nil {
		return nil, err
	}
	return response, nil
}
