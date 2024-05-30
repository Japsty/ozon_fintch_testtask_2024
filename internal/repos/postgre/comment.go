package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type CommentMemoryRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentMemoryRepository {
	return &CommentMemoryRepository{db: db}
}

func (cr *CommentMemoryRepository) CreateComment(ctx context.Context, id, content, userID, postID, parentCommentID string) ([]*model.Comment, error) {
	_, err := cr.db.Exec(ctx, querries.CreateComment, id, content, userID, postID, parentCommentID)
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, querries.GetCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}

	comments := []*model.Comment{}
	for rows.Next() {
		comment := &model.Comment{}
		var createdAtTime time.Time
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentCommentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentMemoryRepository) GetCommentByParentID(ctx context.Context, parentID string) ([]*model.Comment, error) {
	rows, err := cr.db.Query(ctx, querries.GetCommentsByParentID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}
	for rows.Next() {
		comment := &model.Comment{}
		var createdAtTime time.Time
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentCommentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentMemoryRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	rows, err := cr.db.Query(ctx, querries.GetCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}
	for rows.Next() {
		comment := &model.Comment{}
		var createdAtTime time.Time
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentCommentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentMemoryRepository) GetCommentsByPostIDPaginated(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	rows, err := cr.db.Query(ctx, querries.GetCommentsByPostIDPaginated, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}
	for rows.Next() {
		comment := &model.Comment{}
		var createdAtTime time.Time
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentCommentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}
		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	return comments, nil
}
