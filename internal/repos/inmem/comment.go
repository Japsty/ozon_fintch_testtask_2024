package inmem

import (
	"Ozon_testtask/internal/model"
	"context"
	"sync"
	"time"
)

type CommentInMemoryRepository struct {
	data  map[string][]*model.Comment
	mutex sync.RWMutex
}

func NewCommentInMemoryRepository() *CommentInMemoryRepository {
	return &CommentInMemoryRepository{
		data:  map[string][]*model.Comment{},
		mutex: sync.RWMutex{},
	}
}

func (cr *CommentInMemoryRepository) CreateComment(_ context.Context, id, content, userID, postID, parentCommentID string) ([]*model.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	newComment := &model.Comment{
		ID:              id,
		Content:         content,
		Author:          userID,
		PostID:          postID,
		ParentCommentID: &parentCommentID,
		CreatedAt:       time.Now().String(),
		Replies:         []*model.Comment{},
	}

	cr.data[postID] = append(cr.data[postID], newComment)

	return cr.data[postID], nil
}

func (cr *CommentInMemoryRepository) GetCommentByParentID(_ context.Context, parentID string) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	comments := []*model.Comment{}
	for _, comms := range cr.data {
		for _, comm := range comms {
			if *comm.ParentCommentID == parentID {
				comments = append(comments, comm)
			}
		}
	}

	return comments, nil
}

func (cr *CommentInMemoryRepository) GetCommentsByPostID(_ context.Context, postID string) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	comments := []*model.Comment{}
	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.PostID == postID {
				comments = append(comments, comm)
			}
		}
	}

	return comments, nil
}