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

func (s *CommentServiceImpl) CreateComment(ctx context.Context,
	commReq models.CommentRequest) (*models.Comment, error) {
	if len(commReq.Author) == 0 {
		return nil, utils.GqlError{
			Msg:  "comment must have a author",
			Type: consts.BadRequestType,
		}

	}

	if len(commReq.Content) > consts.ContentMaxLen {
		return nil, utils.GqlError{
			Msg:  fmt.Sprintf("comment content should no exceed %v characters", consts.ContentMaxLen),
			Type: consts.BadRequestType,
		}
	}

	post, err := s.postStore.GetPostByID(ctx, commReq.PostID)
	if err != nil {
		return nil, utils.GqlError{
			Msg:  fmt.Sprintf("post with id: %s does not exist", commReq.PostID.String()),
			Type: consts.BadRequestType,
		}
	}

	if !post.IsCommentsAllowed {
		return nil, utils.GqlError{
			Msg:  fmt.Sprintf("post with id: %s is not allowed to create comment", commReq.PostID.String()),
			Type: consts.BadRequestType,
		}

	}

	newComm, err := s.commStore.CreateComment(ctx, models.Comment{
		Author:          commReq.Author,
		Content:         commReq.Content,
		PostID:          commReq.PostID,
		ParentCommentID: commReq.ParentCommentID,
	})

	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error creating comment: %v", err))
		return nil, utils.GqlError{
			Msg:  "error creating comment",
			Type: consts.InternalServerErrorType,
		}
	}

	logger.Logger.Info(fmt.Sprintf("create comment with id: %s successfully", newComm.ID.String()))
	return &newComm, nil
}
func (s *CommentServiceImpl) GetCommentsByPostID(ctx context.Context, postID uuid.UUID,
	page *int32) ([]*models.Comment, error) {
	if page != nil && *page <= 0 {
		return nil, utils.GqlError{
			Msg:  "page must be greater than zero",
			Type: consts.BadRequestType,
		}
	}

	offset, limit := utils.GetOffsetNLimit(page, consts.PageSize)

	comments, err := s.commStore.GetCommentsByPostID(ctx, postID, offset, limit)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error getting comments: %v", err))
		return nil, utils.GqlError{
			Msg:  "error getting comments",
			Type: consts.InternalServerErrorType,
		}
	}

	logger.Logger.Info(fmt.Sprintf("get comments by postId: %s successfully", postID.String()))
	return comments, nil
}
func (s *CommentServiceImpl) GetRepliesByComment(ctx context.Context, commentID uuid.UUID) ([]*models.Comment, error) {
	replies, err := s.commStore.GetRepliesByParentCommentID(ctx, commentID)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("error getting comments: %v", err))
		return nil, utils.GqlError{
			Msg:  "error getting comments",
			Type: consts.InternalServerErrorType,
		}
	}

	logger.Logger.Info(fmt.Sprintf("get comments by commentId: %s successfully", commentID.String()))
	return replies, nil
}
