package inmem

import (
	"Ozon_testtask/internal/models"
	"context"
	"sync"
	"time"
)

type CommentInMemoryRepository struct {
	data  map[string][]models.Comment
	mutex sync.RWMutex
}

func NewCommentInMemoryRepository() *CommentInMemoryRepository {
	return &CommentInMemoryRepository{
		data:  map[string][]models.Comment{},
		mutex: sync.RWMutex{},
	}
}

func (cr *CommentInMemoryRepository) CreateComment(_ context.Context, id, content, userID, postID, parentCommentID string) ([]models.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	newComment := models.Comment{
		ID:              id,
		Content:         content,
		AuthorID:        userID,
		PostID:          postID,
		ParentCommentID: parentCommentID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Replies:         []*models.Comment{},
	}

	cr.data[postID] = append(cr.data[postID], newComment)

	return cr.data[postID], nil
}

func (cr *CommentInMemoryRepository) GetCommentByParentID(_ context.Context, parentID string) ([]models.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	comments := []models.Comment{}
	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.ParentCommentID == parentID {
				comments = append(comments, comm)
			}
		}
	}

	return comments, nil
}

func (cr *CommentInMemoryRepository) GetCommentsByPostID(_ context.Context, postID string) ([]models.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	comments := []models.Comment{}
	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.PostID == postID {
				comments = append(comments, comm)
			}
		}
	}

	return comments, nil
}
