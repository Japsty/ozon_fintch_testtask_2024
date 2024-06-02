package services

import (
	"Ozon_testtask/internal/model"
	"Ozon_testtask/internal/services"
	"Ozon_testtask/tests/mocks/repos"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var mockPost = model.Post{
	ID:              "a4ad7024-7c45-4453-9fac-2dfae1ad2c96",
	Title:           "Test Title",
	Content:         "Test Content",
	UserID:          "5594a70f-ad01-427e-be8a-43bf94fc76fd",
	Comments:        []*model.Comment{},
	CommentsAllowed: true,
	CreatedAt:       "",
}

func TestGetAllPostsSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	ps := services.NewPostService(mockPostRepo, mockCommentRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")
	mockPostRepo.On("GetAllPosts", ctx).
		Return([]*model.Post{&mockPost}, nil)

	mockCommentRepo.On("GetCommentsByPostID", ctx, mockPost.ID).
		Return([]*model.Comment{&mockComment}, nil)

	posts, err := ps.GetAllPosts(ctx)
	if err != nil {
		t.Errorf("CommentPost Error: %s", err)
		return
	}

	post := mockPost
	post.Comments = []*model.Comment{&mockComment}

	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Equal(t, post.Title, posts[0].Title)
	assert.Equal(t, post.Content, posts[0].Content)
	assert.Equal(t, post.UserID, posts[0].UserID)
	assert.Equal(t, post.Comments, posts[0].Comments)
	assert.Equal(t, post.CommentsAllowed, posts[0].CommentsAllowed)
}

func TestAddPostSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	ps := services.NewPostService(mockPostRepo, mockCommentRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")

	mockPostRepo.On(
		"CreatePost",
		ctx,
		mock.AnythingOfType("string"),
		mockPost.Title,
		mockPost.Content,
		mockPost.UserID,
		mockPost.CommentsAllowed,
	).
		Return(mockPost, nil)

	post, err := ps.AddPost(ctx, mockPost.Title, mockPost.Content, mockPost.CommentsAllowed)
	if err != nil {
		t.Errorf("AddPost Error: %s", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, mockPost.Title, post.Title)
	assert.Equal(t, mockPost.Content, post.Content)
	assert.Equal(t, mockPost.UserID, post.UserID)
	assert.Equal(t, mockPost.Comments, post.Comments)
	assert.Equal(t, mockPost.CommentsAllowed, post.CommentsAllowed)
}

func TestGetPostByPostIDSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	ps := services.NewPostService(mockPostRepo, mockCommentRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")

	mockPostRepo.On("GetPostByPostID", ctx, mockPost.ID).
		Return(mockPost, nil)

	post, err := ps.GetPostByPostID(ctx, mockPost.ID)
	if err != nil {
		t.Errorf("CommentPost Error: %s", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, mockPost.Title, post.Title)
	assert.Equal(t, mockPost.Content, post.Content)
	assert.Equal(t, mockPost.UserID, post.UserID)
	assert.Equal(t, mockPost.Comments, post.Comments)
	assert.Equal(t, mockPost.CommentsAllowed, post.CommentsAllowed)
}

func TestUpdatePostCommentsStatusSuccess(t *testing.T) {
	mockCommentRepo := new(mocks.CommentRepoMock)
	mockPostRepo := new(mocks.PostRepoMock)

	ps := services.NewPostService(mockPostRepo, mockCommentRepo)

	ctx := context.WithValue(context.Background(), "userID", "5594a70f-ad01-427e-be8a-43bf94fc76fd")

	mockPostRepo.On("UpdatePostStatus", ctx, mockPost.ID, mockPost.UserID, mockPost.CommentsAllowed).
		Return(mockPost, nil)

	post, err := ps.UpdatePostCommentsStatus(ctx, mockPost.ID, mockPost.CommentsAllowed)
	if err != nil {
		t.Errorf("CommentPost Error: %s", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, mockPost.Title, post.Title)
	assert.Equal(t, mockPost.Content, post.Content)
	assert.Equal(t, mockPost.UserID, post.UserID)
	assert.Equal(t, mockPost.Comments, post.Comments)
	assert.Equal(t, mockPost.CommentsAllowed, post.CommentsAllowed)
}
