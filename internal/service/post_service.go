package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/internal/storage"
	"github.com/nedokyrill/posts-service/pkg/consts"
	"github.com/nedokyrill/posts-service/pkg/logger"
	"github.com/nedokyrill/posts-service/pkg/utils"
)

type PostServiceImpl struct {
	store storage.PostStorage
}

func NewPostService(store storage.PostStorage) *PostServiceImpl {
	return &PostServiceImpl{
		store: store,
	}
}

func (s *PostServiceImpl) GetAllPosts(ctx context.Context, page *int32) ([]*models.Post, error) {
	if page != nil && *page <= 0 {
		logger.Logger.Error("page must be greater than zero")
		return nil, utils.GqlError{Msg: "page must be greater than zero", Type: consts.BadRequestType}

	}

	offset, limit := utils.GetOffsetNLimit(page, consts.PageSize)

	posts, err := s.store.GetAllPosts(ctx, offset, limit)
	if err != nil {
		logger.Logger.Error("error with getting posts: ", err)
		return nil, utils.GqlError{Msg: "error with getting posts", Type: consts.InternalServerErrorType}
	}

	logger.Logger.Info("get all posts successfully")
	return posts, nil
}
func (s *PostServiceImpl) GetPostByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	post, err := s.store.GetPostByID(ctx, id)

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("post with id: %s not found, err: %v", id.String(), err))
		return nil, utils.GqlError{Msg: fmt.Sprintf("post with id: %s not found", id.String()),
			Type: consts.InternalServerErrorType}
	}

	logger.Logger.Info(fmt.Sprintf("get post with id: %s successfully", post.ID.String()))
	return post, nil
}
func (s *PostServiceImpl) CreatePost(ctx context.Context, title string, author *string, content string,
	isCommentAllowed bool) (*models.Post, error) {
	if len(title) == 0 {
		logger.Logger.Error("post must have a title")
		return nil, utils.GqlError{Msg: "post must have a title", Type: consts.BadRequestType}
	}

	if len(*author) == 0 {
		logger.Logger.Error("post must have a author")
		return nil, utils.GqlError{Msg: "post must have a author", Type: consts.BadRequestType}
	}

	newPost, err := s.store.CreatePost(ctx, models.Post{Title: title, Author: *author, Content: content,
		IsCommentsAllowed: isCommentAllowed})
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error with creating post: %v", err))
		return nil, utils.GqlError{Msg: "error creating post", Type: consts.InternalServerErrorType}
	}

	logger.Logger.Info(fmt.Sprintf("create post with id: %s successfully", newPost.ID.String()))
	return &newPost, nil
}
