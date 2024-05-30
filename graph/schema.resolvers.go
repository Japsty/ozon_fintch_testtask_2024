package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"Ozon_testtask/internal/model"
	"context"
	"errors"
	"fmt"
)

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, post model.NewPost) (*model.Post, error) {
	if len(post.Content) > 2000 {
		return nil, errors.New("maximum content length is 2000")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	createdPost, err := r.PostService.AddPost(ctxWithTimeout, post.Title, post.Content, post.CommentsAllowed)
	if err != nil {
		r.Logger.Error("AddPost Resolver Service Error: ", err)
		return nil, err
	}

	return &createdPost, nil
}

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, comment model.NewComment) (*model.Post, error) {
	if len(comment.Content) > 2000 {
		return nil, errors.New("maximum content length is 2000")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	postComments, err := r.CommentService.CommentPost(ctxWithTimeout, comment.ParentID, comment.Content, comment.ParentID)
	if err != nil {
		r.Logger.Error("AddComment Resolver CommentPost Error: ", err)
		return nil, err
	}

	post, err := r.PostService.GetPostByPostID(ctxWithTimeout, comment.PostID)
	if err != nil {
		r.Logger.Error("AddComment Resolver GetPostByPostID Error: ", err)
		return nil, err
	}
	post.Comments = postComments

	return &post, nil
}

// ToggleComments is the resolver for the toggleComments field.
func (r *mutationResolver) ToggleComments(ctx context.Context, postID string, allowed bool) (*model.Post, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	post, err := r.PostService.UpdatePostCommentsStatus(ctxWithTimeout, postID, allowed)
	if err != nil {
		r.Logger.Error("ToggleComments UpdatePostCommentsStatus Error: ", err)
		return nil, err
	}

	return &post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	posts, err := r.PostService.GetAllPosts(ctxWithTimeout)
	if err != nil {
		r.Logger.Error("Posts GetAllPosts Error: ", err)
		return nil, err
	}

	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string, limit int, offset int) (*model.Post, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, TimeoutTime)
	defer cancel()

	post, err := r.PostService.GetPostByPostID(ctxWithTimeout, id)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			return nil, err
		}
		r.Logger.Error("Post GetPostByPostID Error: ", err)
		return nil, err
	}

	comments, err := r.CommentService.GetCommentsByPostID(ctxWithTimeout, id, limit, offset)
	if err != nil {
		r.Logger.Error("Post GetCommentsByPostID Error: ", err)
		return nil, err
	}

	post.Comments = comments

	return &post, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
