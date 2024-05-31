package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}
func (pr *PostRepository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
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

func (pr *PostRepository) CreatePost(ctx context.Context, id, title, text, uID string, status bool) (model.Post, error) {
	var post model.Post

	var createdAtTime time.Time

	err := pr.db.QueryRowContext(ctx, querries.CreatePost, id, title, text, uID, status).Scan(
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

func (pr *PostRepository) GetPostByPostID(ctx context.Context, id string) (model.Post, error) {
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

func (pr *PostRepository) UpdatePostStatus(ctx context.Context, id, uID string, status bool) (model.Post, error) {
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
		return model.Post{}, err
	}

	if foundPost.UserID != uID {
		return model.Post{}, errors.New("not author")
	}

	var post model.Post
	err = pr.db.QueryRowContext(ctx, querries.UpdatePostCommentsStatus, id, status).Scan(
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
