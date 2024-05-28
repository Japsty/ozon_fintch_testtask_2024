package inmem

import (
	"Ozon_testtask/internal/models"
	"context"
	"sync"
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
		Replies:         []*models.Comment{},
	}

	cr.data[postID] = append(cr.data[postID], newComment)

	return cr.data[postID], nil
}

func (cr *CommentInMemoryRepository) DeleteCommentByID(_ context.Context, userID, postID, commentID string) ([]models.Comment, error) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	postComments := cr.data[postID]
	for idx, val := range postComments {
		if val.ID == commentID {
			if val.AuthorID != userID {
				return []models.Comment{}, errNotFound
			}
			postComments = append(postComments[:idx], postComments[idx+1:]...)
			cr.data[postID] = postComments
			return cr.data[postID], nil
		}
	}

	return []models.Comment{}, errNotFound
}
