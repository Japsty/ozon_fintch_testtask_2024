package postgre

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/postgre/querries"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
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

func (pr *PostMemoryRepository) CreatePost(ctx context.Context, id, title, content string, commentsAllowed bool) (model.Post, error) {
	var post model.Post
	var createdAtTime time.Time

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return model.Post{}, errors.New("unauthorized")
	}

	err := pr.db.QueryRow(ctx, querries.CreatePost, id, title, content, userID, commentsAllowed).Scan(
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
	err := pr.db.QueryRow(ctx, querries.GetPostByID, id).Scan(
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

func (pr *PostMemoryRepository) UpdatePostCommentsStatus(ctx context.Context, id string, commentsAllowed bool) (model.Post, error) {
	var foundPost model.Post
	var createdAtTime time.Time

	err := pr.db.QueryRow(ctx, querries.GetPostByID, id).Scan(
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

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return model.Post{}, errors.New("unauthorized")
	}

	if foundPost.UserID != userID {
		return model.Post{}, errors.New("userID ")
	}

	var post model.Post
	err = pr.db.QueryRow(ctx, querries.UpdatePostCommentsStatus, id, commentsAllowed).Scan(
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

func (pr *PostMemoryRepository) UpdatePostComments(ctx context.Context, id string, comments []*model.Comment) (model.Post, error) {
	var post model.Post
	var createdAtTime time.Time

	err := pr.db.QueryRow(ctx, querries.UpdatePostCommentsStatus, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CommentsAllowed,
		&createdAtTime,
	)
	if err != nil {
		return model.Post{}, nil
	}
	post.CreatedAt = createdAtTime.String()
	post.Comments = comments

	return post, nil
}
