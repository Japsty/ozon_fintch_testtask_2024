package inmem

import (
	"Ozon_testtask/internal/model"
	"context"
	"errors"
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

func (cr *CommentRepository) CreateComment(ctx context.Context, id, text, uID, pID, pcID string) (*model.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	newComment := &model.Comment{
		ID:        id,
		Content:   text,
		AuthorID:  uID,
		PostID:    pID,
		ParentID:  &pcID,
		CreatedAt: time.Now().String(),
		Replies:   []*model.Comment{},
	}

	if _, ok := cr.data[pID]; !ok {
		cr.data[pID] = []*model.Comment{}
	}
	cr.data[pID] = append(cr.data[pID], newComment)

	return newComment, nil
}

func (cr *CommentRepository) GetCommentsByPostID(ctx context.Context, pID string) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	comms, ok := cr.data[pID]
	if !ok {
		return nil, errors.New("not found")
	}

	comments := make(map[string]*model.Comment)

	for _, comm := range comms {
		comments[comm.ID] = comm
	}

	roots := []*model.Comment{}

	for _, comment := range comments {
		if comment.ParentID == nil || *comment.ParentID == "" {
			roots = append(roots, comment)
		} else {
			parent := comments[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	return roots, nil
}

func (cr *CommentRepository) GetCommentsByPostIDPaginated(ctx context.Context, pID string, l, o int) ([]*model.Comment, error) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	comms, ok := cr.data[pID]
	if !ok {
		return nil, errors.New("not found")
	}

	comments := make(map[string]*model.Comment)

	for _, comm := range comms {
		comments[comm.ID] = comm
	}

	roots := []*model.Comment{}

	for _, comment := range comments {
		if comment.ParentID == nil || *comment.ParentID == "" {
			if o > 0 {
				o--
			} else if l > 0 {
				roots = append(roots, comment)
				l--
			}
		} else {
			parent := comments[*comment.ParentID]
			parent.Replies = append(parent.Replies, comment)
		}
	}

	return roots, nil
}
