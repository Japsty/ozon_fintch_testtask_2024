package services

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	errNotFound           = errors.New("not found")
	errForbidden          = errors.New("forbidden")
	errCommentsNotAllowed = errors.New("comments not allowed")
)

type CommentService struct {
	CommentRepo models.CommentRepo
	PostRepo    models.PostRepo
}

func NewCommentService(commentRepo models.CommentRepo, postRepo models.PostRepo) *CommentService {
	return &CommentService{CommentRepo: commentRepo, PostRepo: postRepo}
}

func (cs *CommentService) CommentPost(ctx context.Context, userID string, postID string, commentText string) ([]models.Comment, error) {
	commentID := uuid.NewString()
	post, err := cs.PostRepo.GetPostByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if !post.CommentsAllowed {
		return nil, errCommentsNotAllowed
	}

	_, err = cs.CommentRepo.CreateComment(ctx, commentID, userID, commentText, postID, "")
	if err != nil {
		return nil, err
	}

	updatedComments, err := cs.CommentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return updatedComments, nil
}

func (cs *CommentService) GetCommentByParentID(ctx context.Context, parentID string) ([]models.Comment, error) {
	comments, err := cs.CommentRepo.GetCommentByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (cs *CommentService) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	comments, err := cs.CommentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
