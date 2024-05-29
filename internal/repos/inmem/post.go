package inmem

import (
	"Ozon_testtask/internal/models"
	"context"
	"errors"
	"sync"
)

var (
	errNotFound = errors.New("not found")
)

type PostInMemoryRepository struct {
	data  map[string]models.Post
	mutex sync.RWMutex
}

func NewPostInMemoryRepository() *PostInMemoryRepository {
	return &PostInMemoryRepository{
		data:  map[string]models.Post{},
		mutex: sync.RWMutex{},
	}
}

func (pr *PostInMemoryRepository) GetAllPosts(_ context.Context) ([]models.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := []models.Post{}

	for _, val := range pr.data {
		result = append(result, val)
	}

	return result, nil
}

func (pr *PostInMemoryRepository) CreatePost(_ context.Context, id, title, content string, commentsAllowed bool) (models.Post, error) {
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

func (pr *PostInMemoryRepository) GetPostByPostID(_ context.Context, id string) (models.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	v, ok := pr.data[id]
	if !ok {
		return models.Post{}, errNotFound
	}

	return v, nil
}

func (pr *PostInMemoryRepository) UpdatePostCommentsStatus(ctx context.Context, id string, commentsAllowed bool) (models.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post, err := pr.GetPostByPostID(ctx, id)
	if err != nil {
		return models.Post{}, err
	}
	post.CommentsAllowed = commentsAllowed
	pr.data[id] = post

	return post, nil
}

func (pr *PostInMemoryRepository) DeletePostByID(_ context.Context, id string) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	if _, ok := pr.data[id]; !ok {
		return errNotFound
	}

	delete(pr.data, id)

	return nil
}
