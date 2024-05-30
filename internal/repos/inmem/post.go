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

type PostInMemoryRepository struct {
	data  map[string]*model.Post
	mutex sync.RWMutex
}

func NewPostInMemoryRepository() *PostInMemoryRepository {
	return &PostInMemoryRepository{
		data:  map[string]*model.Post{},
		mutex: sync.RWMutex{},
	}
}

func (pr *PostInMemoryRepository) GetAllPosts(_ context.Context) ([]*model.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	result := []*model.Post{}

	for _, val := range pr.data {
		result = append(result, val)
	}

	return result, nil
}

func (pr *PostInMemoryRepository) CreatePost(_ context.Context, id, title, content, userID string, commentsAllowed bool) (model.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post := &model.Post{
		ID:              id,
		Title:           title,
		Content:         content,
		Comments:        []*model.Comment{},
		CommentsAllowed: commentsAllowed,
		CreatedAt:       time.Now().String(),
	}

	pr.data[id] = post

	return *post, nil
}

func (pr *PostInMemoryRepository) GetPostByPostID(_ context.Context, id string) (model.Post, error) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	v, ok := pr.data[id]
	if !ok {
		return model.Post{}, errNotFound
	}

	return *v, nil
}

func (pr *PostInMemoryRepository) UpdatePostCommentsStatus(ctx context.Context, id, userID string, commentsAllowed bool) (model.Post, error) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()

	post, err := pr.GetPostByPostID(ctx, id)
	if err != nil {
		return model.Post{}, err
	}
	post.CommentsAllowed = commentsAllowed
	pr.data[id] = &post

	return post, nil
}
