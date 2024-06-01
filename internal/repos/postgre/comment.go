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

func (cr *CommentRepository) CreateComment(ctx context.Context, id, text, uID, pID, pcID string) error {
	txOptions := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	tx, err := cr.db.BeginTx(ctx, &txOptions)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	_, err = tx.ExecContext(ctx, querries.CreateComment, id, text, uID, pID, pcID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (cr *CommentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	txOptions := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	tx, err := cr.db.BeginTx(ctx, &txOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	rows, err := tx.QueryContext(ctx, querries.GetAllCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentsMap := make(map[string]*model.Comment)
	roots := []*model.Comment{}

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
		commentsMap[comment.ID] = comment
	}

	for _, comment := range commentsMap {
		if comment.ParentID == nil {
			roots = append(roots, comment)
		} else {
			parent := commentsMap[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return roots, nil
}

func (cr *CommentRepository) GetCommentsByPostIDPaginated(ctx context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	txOptions := sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}

	tx, err := cr.db.BeginTx(ctx, &txOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	rows, err := tx.QueryContext(ctx, querries.GetAllCommentsByPostID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentsMap := make(map[string]*model.Comment)
	roots := []*model.Comment{}

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
		commentsMap[comment.ID] = comment
	}

	for _, comment := range commentsMap {
		if comment.ParentID == nil {
			if offset > 0 {
				offset--
			} else if limit > 0 {
				roots = append(roots, comment)
				limit--
			}
		} else {
			parent := commentsMap[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return roots, nil
}
