package services

import (
	"Ozon_testtask/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type PostService struct {
	PostRepo    model.PostRepo
	CommentRepo model.CommentRepo
}

func NewPostService(postRepo model.PostRepo, commentRepo model.CommentRepo) *PostService {
	return &PostService{PostRepo: postRepo, CommentRepo: commentRepo}
}

func (ps *PostService) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts, err := ps.PostRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		comments, err := ps.CommentRepo.GetCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
	}

	return posts, nil
}

func (ps *PostService) AddPost(ctx context.Context, title, text string, status bool) (model.Post, error) {
	newPostID := uuid.NewString()

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return model.Post{}, errors.New("unauthorized")
	}

	post, err := ps.PostRepo.CreatePost(ctx, newPostID, title, text, userID, status)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (ps *PostService) GetPostByPostID(ctx context.Context, postID string) (model.Post, error) {
	post, err := ps.PostRepo.GetPostByPostID(ctx, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, errNotFound
		}

		return model.Post{}, err
	}

	return post, nil
}

func (ps *PostService) UpdatePostCommentsStatus(ctx context.Context, postID string, status bool) (model.Post, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return model.Post{}, errors.New("unauthorized")
	}

	post, err := ps.PostRepo.UpdatePostStatus(ctx, postID, userID, status)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}
