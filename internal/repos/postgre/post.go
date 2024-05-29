package postgre

import (
	"Ozon_testtask/internal/model"
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

func (pr *PostMemoryRepository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := pr.db.Query(ctx, querries.GetAllPosts)
	if err != nil {
		return nil, err
	}

	var posts []*model.Post

	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CommentsAllowed,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostMemoryRepository) CreatePost(ctx context.Context, id, title, content string, commentsAllowed bool) (model.Post, error) {
	var post model.Post
	err := pr.db.QueryRow(ctx, querries.CreatePost, id, title, content, commentsAllowed).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CommentsAllowed,
	)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (pr *PostMemoryRepository) GetPostByPostID(ctx context.Context, id string) (model.Post, error) {
	var post model.Post
	err := pr.db.QueryRow(ctx, querries.GetPostByID, id).Scan(&post.ID, &post.Title, &post.Content, &post.CommentsAllowed)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (pr *PostMemoryRepository) UpdatePostCommentsStatus(ctx context.Context, id string, commentsAllowed bool) (model.Post, error) {
	var post model.Post
	err := pr.db.QueryRow(ctx, querries.UpdatePostCommentsStatus, id, commentsAllowed).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CommentsAllowed,
	)
	if err != nil {
		return model.Post{}, nil
	}

	return post, nil
}

func (pr *PostMemoryRepository) UpdatePostComments(ctx context.Context, id string, comments []*model.Comment) (model.Post, error) {
	var post model.Post
	err := pr.db.QueryRow(ctx, querries.UpdatePostCommentsStatus, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CommentsAllowed,
	)
	if err != nil {
		return model.Post{}, nil
	}

	return post, nil
}
