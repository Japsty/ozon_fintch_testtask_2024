package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"database/sql"
	"time"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (cr *CommentRepository) CreateComment(ctx context.Context, id, text, uID, pID, pcID string) ([]*model.Comment, error) {
	_, err := cr.db.ExecContext(ctx, querries.CreateComment, id, text, uID, pID, pcID)
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.QueryContext(ctx, querries.GetCommentsByPostID, pID)
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
			&comment.ParentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}

		comment.CreatedAt = createdAtTime.String()

		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentRepository) getRepliesForComment(ctx context.Context, comment *model.Comment) error {
	rows, err := cr.db.QueryContext(ctx, querries.GetCommentsByParentID, comment.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	replies := []*model.Comment{}

	for rows.Next() {
		reply := &model.Comment{}

		var createdAtTime time.Time

		if err = rows.Scan(
			&reply.ID,
			&reply.Content,
			&reply.AuthorID,
			&reply.PostID,
			&reply.ParentID,
			&createdAtTime,
		); err != nil {
			return err
		}

		reply.CreatedAt = createdAtTime.String()
		replies = append(replies, reply)
	}

	for _, reply := range replies {
		comment.Replies = append(comment.Replies, reply)

		if err = cr.getRepliesForComment(ctx, reply); err != nil {
			return err
		}
	}

	return nil
}

func (cr *CommentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	rows, err := cr.db.QueryContext(ctx, querries.GetCommentsByPostID, postID)
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
			&comment.ParentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}

		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	for _, comment := range comments {
		if err := cr.getRepliesForComment(ctx, comment); err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func (cr *CommentRepository) GetCommentsByPostIDPaginated(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	rows, err := cr.db.QueryContext(ctx, querries.GetCommentsByPostIDPaginated, postID, limit, offset)
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
			&comment.ParentID,
			&createdAtTime,
		); err != nil {
			return nil, err
		}

		comment.CreatedAt = createdAtTime.String()
		comments = append(comments, comment)
	}

	for _, comment := range comments {
		if err := cr.getRepliesForComment(ctx, comment); err != nil {
			return nil, err
		}
	}

	return comments, nil
}
