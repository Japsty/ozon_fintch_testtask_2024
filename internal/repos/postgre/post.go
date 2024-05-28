package postgre

import (
	"Ozon_testtask/internal/models"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostMemoryRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostMemoryRepository {
	return &PostMemoryRepository{db: db}
}

func (pr *PostMemoryRepository) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	rows, err := pr.db.Query(ctx, querries.GetAllPosts)
	if err != nil {
		return nil, err
	}

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Comments,
			&post.CommentsAllowed,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostMemoryRepository) CreatePost(ctx context.Context, id, title, content string, commentsAllowed bool) (models.Post, error) {
	var post models.Post
	err := pr.db.QueryRow(ctx, querries.CreatePost, id, title, content, commentsAllowed).Scan(&post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (pr *PostMemoryRepository) GetPostByPostID(ctx context.Context, id string) (models.Post, error) {
	var post models.Post
	err := pr.db.QueryRow(ctx, querries.GetPostByID, id).Scan(&post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (pr *PostMemoryRepository) UpdatePostComments(ctx context.Context, id string, commentsAllowed bool) (models.Post, error) {
	var post models.Post
	err := pr.db.QueryRow(ctx, querries.UpdatePostComments, id, commentsAllowed)
	if err != nil {
		return models.Post{}, nil
	}

	return post, nil
}

func (pr *PostMemoryRepository) DeletePostByID(ctx context.Context, id string) error {
	//err := pr.db.Exec(ctx,querries.Dele)
	return nil
}
