package inmem

import (
	"Ozon_testtask/internal/model"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	errNotFound = errors.New("not found")
)

type PostRepository struct {
	data  map[string]*model.Post
	mutex sync.RWMutex
}

func NewPostInMemoryRepository() *PostRepository {
	return &PostRepository{
		data:  map[string]*model.Post{},
		mutex: sync.RWMutex{},
	}
}

func (pr *PostRepository) GetAllPosts(_ context.Context) ([]*model.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := []*model.Post{}

	for _, val := range pr.data {
		result = append(result, val)
	}

	return result, nil
}

func (pr *PostRepository) CreatePost(_ context.Context, id, title, text, uID string, cStatus bool) (model.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post := &model.Post{
		ID:              id,
		Title:           title,
		Content:         text,
		UserID:          uID,
		Comments:        []*model.Comment{},
		CommentsAllowed: cStatus,
		CreatedAt:       time.Now().String(),
	}

	pr.data[id] = post

	return *post, nil
}

func (pr *PostRepository) GetPostByPostID(_ context.Context, id string) (model.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	v, ok := pr.data[id]
	if !ok {
		return model.Post{}, errNotFound
	}

	return *v, nil
}

func (pr *PostRepository) UpdatePostStatus(ctx context.Context, id, _ string, status bool) (model.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post, err := pr.GetPostByPostID(ctx, id)
	if err != nil {
		return model.Post{}, err
	}

	post.CommentsAllowed = status

	pr.data[id] = &post

	return post, nil
}
