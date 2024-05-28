package postgre

import (
	"Ozon_testtask/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PostRepo interface {
	GetAllPosts() ([]Post, error)
	CreatePost(string, string, string, bool) (Post, error)
	GetPostByPostID(string) (Post, error)
	UpdatePostComments(string, bool) (Post, error)
	DeletePostByID(string) error
}

type PostMemoryRepository struct {
	db *pgxpool.Pool
}

func NewPostMemoryRepository(db *pgxpool.Pool) *PostMemoryRepository {
	return &PostMemoryRepository{db: db}
}

func (g *PostMemoryRepository) ExistionCheck(ctx context.Context, goodId, projectId int) (bool, error) {
	var exists bool
	err := g.db.QueryRow(ctx, querries.CheckRecord, goodId, projectId).Scan(&exists)
	if err != nil {
		log.Printf("UpdateGood QueryRow CheckRecord Error: %v", err)
		return false, err
	}
	if !exists {
		log.Printf("Record doesn't exists")
		return exists, nil
	}
	return exists, nil
}

func (pr *PostMemoryRepository) GetAllPosts() ([]models.Post, error) {
	pr.db.
}

func (pr *PostMemoryRepository) CreatePost(id, title, content string, commentsAllowed bool) (models.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post := models.Post{
		ID:              id,
		Title:           title,
		Content:         content,
		Comments:        []*models.Comment{},
		CommentsAllowed: commentsAllowed,
	}

	pr.data[id] = post

	return post, nil
}

func (pr *PostMemoryRepository) GetPostByPostID(id string) (models.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	v, ok := pr.data[id]
	if !ok {
		return models.Post{}, errNotFound
	}

	return v, nil
}

func (pr *PostMemoryRepository) UpdatePostComments(id string, commentsAllowed bool) (models.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post, err := pr.GetPostByPostID(id)
	if err != nil {
		return models.Post{}, err
	}
	post.CommentsAllowed = commentsAllowed
	pr.data[id] = post

	return post, nil
}

func (pr *PostMemoryRepository) DeletePostByID(id string) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, ok := pr.data[id]; !ok {
		return errNotFound
	}

	delete(pr.data, id)

	return nil
}
