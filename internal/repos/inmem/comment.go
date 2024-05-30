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
		ID:        id,
		Content:   content,
		AuthorID:  userID,
		PostID:    postID,
		ParentID:  parentCommentID,
		CreatedAt: time.Now().String(),
		Replies:   []*model.Comment{},
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
			if comm.ParentID == parentID {
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

func (cr *CommentInMemoryRepository) GetCommentsByPostIDPaginated(_ context.Context, postID string, limit, offset int) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	comments := []*model.Comment{}
	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.PostID == postID {
				if offset > 0 {
					offset--
					continue
				} else if limit > 0 {
					comments = append(comments, comm)
					limit--
				} else {
					break
				}
			}
		}
	}

	return comments, nil
}
