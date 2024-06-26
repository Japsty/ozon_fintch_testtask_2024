package services

import (
	"Ozon_testtask/internal/model"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	errNotFound           = errors.New("not found")
	errCommentsNotAllowed = errors.New("comments not allowed")
)

type CommentService struct {
	CommentRepo model.CommentRepo
	PostRepo    model.PostRepo
}

func NewCommentService(commentRepo model.CommentRepo, postRepo model.PostRepo) *CommentService {
	return &CommentService{CommentRepo: commentRepo, PostRepo: postRepo}
}

func (cs *CommentService) CommentPost(ctx context.Context, id, text, parID string) ([]*model.Comment, *model.Comment, error) {
	commentID := uuid.NewString()
	post, err := cs.PostRepo.GetPostByPostID(ctx, id)

	if err != nil {
		return nil, nil, err
	}

	if !post.CommentsAllowed {
		return nil, nil, errCommentsNotAllowed
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, nil, errors.New("unauthorized")
	}

	comment, err := cs.CommentRepo.CreateComment(ctx, commentID, text, userID, id, parID)
	if err != nil {
		return nil, nil, err
	}

	updatedComments, err := cs.CommentRepo.GetCommentsByPostID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return updatedComments, comment, nil
}

func (cs *CommentService) GetCommentsByPostID(ctx context.Context, id string, l, o int) ([]*model.Comment, error) {
	if l < 0 || o < 0 {
		return nil, errors.New("limit or Offset must be > 0")
	}

	comments, err := cs.CommentRepo.GetCommentsByPostIDPaginated(ctx, id, l, o)

	if err != nil {
		return nil, err
	}

	return comments, nil
}
