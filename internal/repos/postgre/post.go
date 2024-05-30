package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostMemoryRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostMemoryRepository {
	return &PostMemoryRepository{db: db}
}
func (pr *PostMemoryRepository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := pr.db.QueryContext(ctx, querries.GetAllPosts)
	if err != nil {
		return nil, err
	}

	var posts []*model.Post

	for rows.Next() {
		post := &model.Post{}
		var createdAtTime time.Time
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.CommentsAllowed,
			&createdAtTime,
		); err != nil {
			return nil, err
		}
		post.CreatedAt = createdAtTime.String()
		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostMemoryRepository) CreatePost(ctx context.Context, id, title, content, userID string, commentsAllowed bool) (model.Post, error) {
	var post model.Post
	var createdAtTime time.Time

	err := pr.db.QueryRowContext(ctx, querries.CreatePost, id, title, content, userID, commentsAllowed).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CommentsAllowed,
		&createdAtTime,
	)
	if err != nil {
		return model.Post{}, err
	}
	post.CreatedAt = createdAtTime.String()

	return post, nil
}

func (pr *PostMemoryRepository) GetPostByPostID(ctx context.Context, id string) (model.Post, error) {
	var post model.Post
	var createdAtTime time.Time
	err := pr.db.QueryRowContext(ctx, querries.GetPostByID, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CommentsAllowed,
		&createdAtTime,
	)
	if err != nil {
		return model.Post{}, err
	}
	post.CreatedAt = createdAtTime.String()

	return post, nil
}

func (pr *PostMemoryRepository) UpdatePostCommentsStatus(ctx context.Context, id, userID string, commentsAllowed bool) (model.Post, error) {
	var foundPost model.Post
	var createdAtTime time.Time

	err := pr.db.QueryRowContext(ctx, querries.GetPostByID, id).Scan(
		&foundPost.ID,
		&foundPost.Title,
		&foundPost.Content,
		&foundPost.UserID,
		&foundPost.CommentsAllowed,
		&createdAtTime,
	)
	if err != nil {
		return model.Post{}, nil
	}

	if foundPost.UserID != userID {
		return model.Post{}, errors.New("userID ")
	}

	var post model.Post
	err = pr.db.QueryRowContext(ctx, querries.UpdatePostCommentsStatus, id, commentsAllowed).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.CommentsAllowed,
		&createdAtTime,
	)
	if err != nil {
		return model.Post{}, nil
	}
	post.CreatedAt = createdAtTime.String()

	return post, nil
}
