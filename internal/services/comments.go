package services

import "Ozon_testtask/internal/models"

var (
	errNotFound  = errors.New("not found")
	errForbidden = errors.New("forbidden")
)

type CommentService struct {
	CommentRepo models.CommentRepo
	PostRepo    models.PostRepo
}

func NewCommentService(commentRepo models.CommentRepo, postRepo models.PostRepo) *CommentService {
	return &CommentService{CommentRepo: commentRepo, PostRepo: postRepo}
}

func (cs *CommentService) CommentPost(ctx context.Context, author models.Author, postID string, commentText string) (models.Post, error) {
	err := cs.CommentRepo.CreateComment(ctx, author, postID, commentText)
	if err != nil {
		return models.Post{}, err
	}

	updatedComments, err := cs.CommentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}

	updatedPost, err := cs.PostRepo.UpdateComment(ctx, postID, updatedComments)
	if err != nil {
		return models.Post{}, err
	}

	return updatedPost, nil
}

func (cs *CommentService) DeleteComment(ctx context.Context, userID, postID, commentID string) (models.Post, error) {
	comment, err := cs.CommentRepo.GetCommentByIDs(ctx, postID, commentID)
	if err != nil {
		return models.Post{}, err
	}

	if comment.Author.UserID != userID {
		return models.Post{}, errForbidden
	}

	err = cs.CommentRepo.DeleteComment(ctx, postID, commentID)
	if err != nil {
		return models.Post{}, err
	}

	comments, err := cs.CommentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil && errors.Is(err, errNotFound) {
		return models.Post{}, errNotFound
	} else if err != nil {
		return models.Post{}, err
	}

	updatedPost, err := cs.PostRepo.UpdateComment(ctx, postID, comments)
	if err != nil {
		return models.Post{}, err
	}

	return updatedPost, nil
}
