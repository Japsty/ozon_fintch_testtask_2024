package postgre

import (
	"Ozon_testtask/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentMemoryRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentMemoryRepository {
	return &CommentMemoryRepository{db: db}
}

func (cr *CommentMemoryRepository) CreateComment(ctx context.Context, id, content, userID, postID, parentCommentID string) ([]models.Comment, error) {
	return nil, nil
}

func (cr *CommentMemoryRepository) DeleteCommentByID(ctx context.Context, userID, postID, commentID string) ([]models.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	postComments := cr.data[postID]
	for idx, val := range postComments {
		if val.ID == commentID {
			if val.AuthorID != userID {
				return []models.Comment{}, errNotFound
			}
			postComments = append(postComments[:idx], postComments[idx+1:]...)
			cr.data[postID] = postComments
			return cr.data[postID], nil
		}
	}

	return []models.Comment{}, errNotFound
}
