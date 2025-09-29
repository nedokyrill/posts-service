package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	store_mock "github.com/nedokyrill/posts-service/internal/storage/mocks"
	"github.com/nedokyrill/posts-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	// инициализация логгера в тестовой среде
	log, _ := zap.NewDevelopment()
	logger.Logger = log.Sugar()
	m.Run()
}

func TestCommentService_CreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	postStorage := store_mock.NewMockPostStorage(ctrl)
	commentStorage := store_mock.NewMockCommentStorage(ctrl)

	commentService := NewCommentService(commentStorage, postStorage)

	postID := uuid.New()
	author := "test_author"
	content := "test content"

	t.Run("successfully create root comment", func(t *testing.T) {
		post := &models.Post{
			ID:                postID,
			Title:             "Test Post",
			Author:            "post_author",
			Content:           "post content",
			IsCommentsAllowed: true,
		}

		commentReq := models.CommentRequest{
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		expectedComment := models.Comment{
			ID:        uuid.New(),
			Author:    author,
			Content:   content,
			PostID:    postID,
			CreatedAt: &[]time.Time{time.Now()}[0],
		}

		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(post, nil)

		commentStorage.EXPECT().
			CreateComment(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, comm models.Comment) (models.Comment, error) {
				assert.Equal(t, author, comm.Author)
				assert.Equal(t, content, comm.Content)
				assert.Equal(t, postID, comm.PostID)
				assert.Nil(t, comm.ParentCommentID)
				return expectedComment, nil
			})

		result, err := commentService.CreateComment(ctx, commentReq)

		require.NoError(t, err)
		assert.Equal(t, &expectedComment, result)
	})

	t.Run("successfully create reply comment", func(t *testing.T) {
		parentCommentID := uuid.New()
		post := &models.Post{
			ID:                postID,
			Title:             "Test Post",
			Author:            "post_author",
			Content:           "post content",
			IsCommentsAllowed: true,
		}

		commentReq := models.CommentRequest{
			Author:          author,
			Content:         content,
			PostID:          postID,
			ParentCommentID: &parentCommentID,
		}

		expectedComment := models.Comment{
			ID:              uuid.New(),
			Author:          author,
			Content:         content,
			PostID:          postID,
			ParentCommentID: &parentCommentID,
			CreatedAt:       &[]time.Time{time.Now()}[0],
		}

		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(post, nil)

		commentStorage.EXPECT().
			CreateComment(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, comm models.Comment) (models.Comment, error) {
				assert.Equal(t, author, comm.Author)
				assert.Equal(t, content, comm.Content)
				assert.Equal(t, postID, comm.PostID)
				assert.Equal(t, &parentCommentID, comm.ParentCommentID)
				return expectedComment, nil
			})

		result, err := commentService.CreateComment(ctx, commentReq)

		require.NoError(t, err)
		assert.Equal(t, &expectedComment, result)
	})

	t.Run("fail when post not found", func(t *testing.T) {
		commentReq := models.CommentRequest{
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(nil, assert.AnError)

		result, err := commentService.CreateComment(ctx, commentReq)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("fail when post comments are not allowed", func(t *testing.T) {
		post := &models.Post{
			ID:                postID,
			Title:             "Test Post",
			Author:            "post_author",
			Content:           "post content",
			IsCommentsAllowed: false,
		}

		commentReq := models.CommentRequest{
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(post, nil)

		result, err := commentService.CreateComment(ctx, commentReq)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not allowed to create comment")
		assert.Nil(t, result)
	})

	t.Run("fail when storage fails to create comment", func(t *testing.T) {
		post := &models.Post{
			ID:                postID,
			Title:             "Test Post",
			Author:            "post_author",
			Content:           "post content",
			IsCommentsAllowed: true,
		}

		commentReq := models.CommentRequest{
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(post, nil)

		commentStorage.EXPECT().
			CreateComment(ctx, gomock.Any()).
			Return(models.Comment{}, assert.AnError)

		result, err := commentService.CreateComment(ctx, commentReq)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
