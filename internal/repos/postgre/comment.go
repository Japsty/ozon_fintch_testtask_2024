package postgre

import (
	"Ozon_testtask/internal/models"
	"Ozon_testtask/internal/repos/postgre/querries"
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
	_, err := cr.db.Exec(ctx, querries.CreateComment, id, content, userID, postID, parentCommentID)
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, querries.GetCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}

	var comments []models.DbComment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.Author,
			&comment.PostID,
			&comment.ParentCommentID,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentMemoryRepository) GetCommentByParentID(ctx context.Context, parentID string) ([]models.Comment, error) {
	rows, err := cr.db.Query(ctx, querries.GetCommentsByParentID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.Author,
			&comment.PostID,
			&comment.ParentCommentID,
			&comment.CreatedAt,
			&comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentMemoryRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]models.Comment, error) {
	rows, err := cr.db.Query(ctx, querries.GetCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentCommentID,
			&comment.CreatedAt,
			&comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
