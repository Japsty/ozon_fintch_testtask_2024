package services

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/services"
	mocks "Ozon_testtask/tests/mocks/repos"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var createdAtTime = time.Now()

var mockComment = model.Comment{
	ID:        "14ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Content:   "Test Content",
	AuthorID:  "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	PostID:    "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	ParentID:  nil,
	Replies:   []*model.Comment{},
	CreatedAt: createdAtTime.String(),
}

func TestCommentPostSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	cs := services.NewCommentService(mockCommentRepo, mockPostRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")
	mockPostRepo.On("GetPostByPostID", ctx, mockComment.PostID).Return(mockPost, nil)
	mockCommentRepo.On(
		"CreateComment",
		ctx,
		mock.AnythingOfType("string"),
		mockComment.Content,
		mockComment.AuthorID,
		mockComment.PostID,
		"",
	).Return(&mockComment, nil)
	mockCommentRepo.On("GetCommentsByPostID", ctx, mockComment.PostID).Return(
		[]*model.Comment{&mockComment},
		nil,
	)

	_, comment, err := cs.CommentPost(ctx, mockComment.PostID, mockComment.Content, "")
	if err != nil {
		t.Errorf("CommentPost Error: %s", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, mockComment.Content, comment.Content)
	assert.Equal(t, mockComment.PostID, comment.PostID)
	assert.Equal(t, mockComment.AuthorID, comment.AuthorID)
	assert.Equal(t, mockComment.ParentID, comment.ParentID)
	assert.Equal(t, mockComment.ID, comment.ID)
}

func TestCommentGetCommentsByPostIDSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	cs := services.NewCommentService(mockCommentRepo, mockPostRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")
	mockCommentRepo.On("GetCommentsByPostIDPaginated", ctx, mockComment.PostID, 1, 0).
		Return([]*model.Comment{&mockComment}, nil)

	comments, err := cs.GetCommentsByPostID(ctx, mockComment.PostID, 1, 0)
	if err != nil {
		t.Errorf("GetCommentsByPostID Error: %s", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Equal(t, mockComment.Content, comments[0].Content)
	assert.Equal(t, mockComment.PostID, comments[0].PostID)
	assert.Equal(t, mockComment.AuthorID, comments[0].AuthorID)
	assert.Equal(t, mockComment.ParentID, comments[0].ParentID)
	assert.Equal(t, mockComment.ID, comments[0].ID)
}
