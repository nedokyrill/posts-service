package service

import (
	"errors"
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

func (s *PostServiceImpl) GetAllPosts(page *int) ([]models.Post, error) {
	if page == nil || *page <= 0 {
		logger.Logger.Error("page must be greater than zero")
		return nil, errors.New("page must be greater than zero")
	}

	offset, limit := utils.GetOffsetNLimit(page, consts.PageSize)

	posts, err := s.store.GetAllPosts(offset, limit)
	if err != nil {
		logger.Logger.Error("error with getting posts: ", err)
		return nil, err
	}

	logger.Logger.Info("get all posts successfully")
	return posts, nil
}
func (s *PostServiceImpl) GetPostByID(id uuid.UUID) (*models.Post, error) {
	post, err := s.store.GetPostByID(id)

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("post with id: %s not found, err: %v", id.String(), err))
		return nil, errors.New(fmt.Sprintf("post with id: %s not found, err: %v", id.String(), err))
	}

	logger.Logger.Info(fmt.Sprintf("get post with id: %s successfully", post.ID.String()))
	return &post, nil
}
func (s *PostServiceImpl) CreatePost(post models.Post) (*models.Post, error) {
	if len(post.Title) == 0 {
		logger.Logger.Error("post must have a title")
		return nil, errors.New("post must have a title")
	}

	if len(post.Author) == 0 {
		logger.Logger.Error("post must have a author")
		return nil, errors.New("post must have a author")
	}

	newPost, err := s.store.CreatePost(post)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error with creating post: %v", err))
		return nil, errors.New(fmt.Sprintf("error creating post"))
	}

	logger.Logger.Info(fmt.Sprintf("create post with id: %s successfully", newPost.ID.String()))
	return &newPost, nil
}
