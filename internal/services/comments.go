package services

import (
	"Ozon_testtask/internal/model"
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
	CommentRepo model.CommentRepo
	PostRepo    model.PostRepo
}

func NewCommentService(commentRepo model.CommentRepo, postRepo model.PostRepo) *CommentService {
	return &CommentService{CommentRepo: commentRepo, PostRepo: postRepo}
}

func (cs *CommentService) CommentPost(ctx context.Context, postID, commentText, parentCommentID string) ([]*model.Comment, error) {
	commentID := uuid.NewString()
	post, err := cs.PostRepo.GetPostByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if !post.CommentsAllowed {
		return nil, errCommentsNotAllowed
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	_, err = cs.CommentRepo.CreateComment(ctx, commentID, userID, commentText, postID, parentCommentID)
	if err != nil {
		return nil, err
	}

	updatedComments, err := cs.CommentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return updatedComments, nil
}

func (cs *CommentService) GetCommentByParentID(ctx context.Context, parentID string) ([]*model.Comment, error) {
	comments, err := cs.CommentRepo.GetCommentByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (cs *CommentService) GetCommentsByPostID(ctx context.Context, postID string, limit int, offset int) ([]*model.Comment, error) {
	comments, err := cs.CommentRepo.GetCommentsByPostIDPaginated(ctx, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
