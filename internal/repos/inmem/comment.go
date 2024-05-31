package inmem

import (
	"Ozon_testtask/internal/model"
	"context"
	"sync"
	"time"
)

type CommentRepository struct {
	data  map[string][]*model.Comment
	mutex sync.RWMutex
}

func NewCommentInMemoryRepository() *CommentRepository {
	return &CommentRepository{
		data:  map[string][]*model.Comment{},
		mutex: sync.RWMutex{},
	}
}

func (cr *CommentRepository) CreateComment(_ context.Context, id, text, uID, pID, pcID string) ([]*model.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	newComment := &model.Comment{
		ID:        id,
		Content:   text,
		AuthorID:  uID,
		PostID:    pID,
		ParentID:  &pcID,
		CreatedAt: time.Now().String(),
		Replies:   []*model.Comment{},
	}

	cr.data[pID] = append(cr.data[pID], newComment)

	return cr.data[pID], nil
}

func (cr *CommentRepository) getRepliesForComment(comment *model.Comment) error {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	var replies []*model.Comment

	for _, comms := range cr.data {
		for _, comm := range comms {
			if *comm.ParentID == comment.ID {
				replies = append(replies, comm)
			}
		}
	}

	for _, reply := range replies {
		comment.Replies = append(comment.Replies, reply)

		if err := cr.getRepliesForComment(reply); err != nil {
			return err
		}
	}

	return nil
}

func (cr *CommentRepository) GetCommentsByPostID(_ context.Context, pID string) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	var comments []*model.Comment

	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.PostID == pID {
				comments = append(comments, comm)
			}
		}
	}

	for _, comment := range comments {
		if err := cr.getRepliesForComment(comment); err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func (cr *CommentRepository) GetCommentsByPostIDPaginated(_ context.Context, pID string, l, o int) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	var comments []*model.Comment

	for _, comms := range cr.data {
		for _, comm := range comms {
			if comm.PostID == pID {
				switch {
				case o > 0:
					o--
					continue
				case l > 0:
					comments = append(comments, comm)
					l--
				default:
					break
				}
			}
		}
	}

	for _, comment := range comments {
		if err := cr.getRepliesForComment(comment); err != nil {
			return nil, err
		}
	}

	return comments, nil
}
