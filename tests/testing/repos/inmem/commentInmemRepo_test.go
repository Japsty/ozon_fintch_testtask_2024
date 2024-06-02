package repos

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/repos/inmem"
	"context"
	"reflect"
	"testing"
	"time"
)

var createdAtTime = time.Now()

var mockComment = model.Comment{
	ID:       "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:  "Test Content",
	AuthorID: "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:   "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID: nil,
	Replies: []*model.Comment{
		&mockChildComment,
	},
	CreatedAt: createdAtTime.String(),
}

var childParentID = "14ad7024-7c45-4453-9fac-2dfae1ad2c96"
var mockChildComment = model.Comment{
	ID:        "24ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  &childParentID,
	Replies:   []*model.Comment{&mockChildChildComment},
	CreatedAt: createdAtTime.String(),
}

var childChildParentID = "24ad7024-7c45-4453-9fac-2dfae1ad2c96"
var mockChildChildComment = model.Comment{
	ID:        "34ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  &childChildParentID,
	Replies:   nil,
	CreatedAt: createdAtTime.String(),
}

type CommentWithoutPointers struct {
	ID        string
	Content   string
	AuthorID  string
	PostID    string
	ParentID  string
	Replies   []CommentWithoutPointers
	CreatedAt string
}

func deletePointer(comment *model.Comment) CommentWithoutPointers {
	var parentIDStr string
	if comment.ParentID != nil {
		parentIDStr = *comment.ParentID
	}

	replies := []CommentWithoutPointers{}
	for _, reply := range comment.Replies {
		replies = append(replies, deletePointer(reply))
	}

	return CommentWithoutPointers{
		ID:        comment.ID,
		Content:   comment.Content,
		AuthorID:  comment.AuthorID,
		PostID:    comment.PostID,
		ParentID:  parentIDStr,
		Replies:   replies,
		CreatedAt: comment.CreatedAt,
	}
}

func TestCreateComment(t *testing.T) {
	repo := inmem.NewCommentInMemoryRepository()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	comment, err := repo.CreateComment(
		ctxWithTimeout,
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		"",
	)
	if err != nil {
		return
	}

	expectedComment := mockComment
	expectedComment.CreatedAt = comment.CreatedAt
	expectedComment.Replies = []*model.Comment{}

	if !reflect.DeepEqual(deletePointer(comment), deletePointer(&expectedComment)) {
		t.Errorf(
			"Unexpected comment data. Got %+v, expected %+v",
			deletePointer(comment),
			deletePointer(&expectedComment),
		)
	}
}

func TestGetCommentByPostID(t *testing.T) {
	repo := inmem.NewCommentInMemoryRepository()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	root, err := repo.CreateComment(
		ctxWithTimeout,
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		"",
	)
	if err != nil {
		return
	}

	child, err := repo.CreateComment(
		ctxWithTimeout,
		mockChildComment.ID,
		mockChildComment.Content,
		mockChildComment.AuthorID,
		mockChildComment.PostID,
		*mockChildComment.ParentID,
	)
	if err != nil {
		return
	}

	childChild, err := repo.CreateComment(
		ctxWithTimeout,
		mockChildChildComment.ID,
		mockChildChildComment.Content,
		mockChildChildComment.AuthorID,
		mockChildChildComment.PostID,
		*mockChildChildComment.ParentID,
	)
	if err != nil {
		return
	}

	comments, err := repo.GetCommentsByPostID(
		ctxWithTimeout,
		mockComment.PostID,
	)
	if err != nil {
		return
	}

	mockComment.CreatedAt = root.CreatedAt
	mockComment.Replies[0].CreatedAt = child.CreatedAt
	mockComment.Replies[0].Replies[0].CreatedAt = childChild.CreatedAt

	expectedComment := mockComment

	for _, comm := range comments {
		if !reflect.DeepEqual(deletePointer(comm), deletePointer(&expectedComment)) {
			t.Errorf(
				"Unexpected comment data. Got %+v, expected %+v",
				deletePointer(comm),
				deletePointer(&expectedComment),
			)
		}
	}
}

func TestGetCommentByPostIDPaginated(t *testing.T) {
	repo := inmem.NewCommentInMemoryRepository()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	root, err := repo.CreateComment(
		ctxWithTimeout,
		mockComment.ID,
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		"",
	)
	if err != nil {
		return
	}

	child, err := repo.CreateComment(
		ctxWithTimeout,
		mockChildComment.ID,
		mockChildComment.Content,
		mockChildComment.AuthorID,
		mockChildComment.PostID,
		*mockChildComment.ParentID,
	)
	if err != nil {
		return
	}

	childChild, err := repo.CreateComment(
		ctxWithTimeout,
		mockChildChildComment.ID,
		mockChildChildComment.Content,
		mockChildChildComment.AuthorID,
		mockChildChildComment.PostID,
		*mockChildChildComment.ParentID,
	)
	if err != nil {
		return
	}

	mockComment.CreatedAt = root.CreatedAt
	mockComment.Replies[0].CreatedAt = child.CreatedAt
	mockComment.Replies[0].Replies[0].CreatedAt = childChild.CreatedAt

	expectedComment := mockComment

	comments, err := repo.GetCommentsByPostIDPaginated(
		ctxWithTimeout,
		mockComment.PostID,
		10,
		0,
	)
	if err != nil {
		return
	}

	for _, comm := range comments {
		if !reflect.DeepEqual(deletePointer(comm), deletePointer(&expectedComment)) {
			t.Errorf(
				"Unexpected comment data. Got %+v, expected %+v",
				deletePointer(comm),
				deletePointer(&expectedComment),
			)
		}
	}
}
