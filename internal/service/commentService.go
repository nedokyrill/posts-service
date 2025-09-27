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

type CommentServiceImpl struct {
	commStore storage.CommentStorage
	postStore storage.PostStorage
}

func NewCommentService(commStore storage.CommentStorage, postStore storage.PostStorage) *CommentServiceImpl {
	return &CommentServiceImpl{
		commStore: commStore,
		postStore: postStore,
	}
}

func (s *CommentServiceImpl) CreateComment(comment models.Comment) (*models.Comment, error) {
	if len(comment.Author) == 0 {
		logger.Logger.Error("comment must have a author")
		return nil, errors.New("comment must have a author")
	}

	if len(comment.Content) > consts.ContentMaxLen {
		logger.Logger.Error(fmt.Sprintf("comment content should no exceed %v characters", consts.ContentMaxLen))
		return nil, errors.New(fmt.Sprintf("comment content should no exceed %v characters",
			consts.ContentMaxLen))
	}

	post, err := s.postStore.GetPostByID(comment.PostID)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("post with id: %s does not exist, error: %v",
			comment.PostID.String(), err))
		return nil, errors.New(fmt.Sprintf("post with id: %s does not exist", comment.PostID.String()))
	}

	if !post.IsCommentsAllowed {
		logger.Logger.Error(fmt.Sprintf("post with id: %s is not allowed to create comment",
			comment.PostID.String()))
		return nil, errors.New(fmt.Sprintf("post with id: %s is not allowed to create comment",
			comment.PostID.String()))
	}

	newComm, err := s.commStore.CreateComment(comment)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error creating comment: %v", err))
		return nil, errors.New("error creating comment")
	}

	logger.Logger.Info(fmt.Sprintf("create comment with id: %s successfully", newComm.ID.String()))
	return &newComm, nil
}
func (s *CommentServiceImpl) GetCommentsByPostID(postID uuid.UUID, page *int) ([]models.Comment, error) {
	if page == nil || *page <= 0 {
		logger.Logger.Error("page must be greater than zero")
		return nil, errors.New("page must be greater than zero")
	}

	offset, limit := utils.GetOffsetNLimit(page, consts.PageSize)

	comments, err := s.commStore.GetCommentsByPostID(postID, offset, limit)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error getting comments: %v", err))
		return nil, errors.New("error getting comments")
	}

	logger.Logger.Info(fmt.Sprintf("get comments by postId: %s successfully", postID.String()))
	return comments, nil
}
func (s *CommentServiceImpl) GetRepliesByComment(commentID uuid.UUID) ([]models.Comment, error) {
	replies, err := s.commStore.GetRepliesByComment(commentID)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error with getting comments: %v", err))
		return nil, errors.New("error with getting comments")
	}

	logger.Logger.Info(fmt.Sprintf("get comments by commentId: %s successfully", commentID.String()))
	return replies, nil
}
