package postgre

import (
	"Ozon_testtask/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentMemoryRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentMemoryRepository {
	return &CommentMemoryRepository{db: db}
}

func (cr *CommentMemoryRepository) CreateComment(id, content, userID, postID, parentCommentID string) ([]models.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	newComment := models.Comment{
		ID:              id,
		Content:         content,
		AuthorID:        userID,
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Replies:         []*models.Comment{},
	}

	cr.data[postID] = append(cr.data[postID], newComment)

	return cr.data[postID], nil
}

func (cr *CommentMemoryRepository) DeleteComment(userID, postID, commentID string) ([]models.Comment, error) {
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
