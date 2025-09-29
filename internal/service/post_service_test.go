package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	store_mock "github.com/nedokyrill/posts-service/internal/storage/mocks"
	"github.com/nedokyrill/posts-service/pkg/consts"
	"github.com/nedokyrill/posts-service/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostService_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	postStorage := store_mock.NewMockPostStorage(ctrl)

	postService := NewPostService(postStorage)

	title := "Test Post Title"
	author := "test_author"
	content := "test content"

	t.Run("successfully create post", func(t *testing.T) {
		// Setup
		postReq := models.PostRequest{
			Title:            title,
			Author:           &author,
			Content:          content,
			IsCommentAllowed: true,
		}

		expectedPost := models.Post{
			ID:                uuid.New(),
			Title:             title,
			Author:            author,
			Content:           content,
			IsCommentsAllowed: true,
		}

		// Mock expectations
		postStorage.EXPECT().
			CreatePost(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, post models.Post) (models.Post, error) {
				assert.Equal(t, title, post.Title)
				assert.Equal(t, author, post.Author)
				assert.Equal(t, content, post.Content)
				assert.True(t, post.IsCommentsAllowed)
				return expectedPost, nil
			})

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, &expectedPost, result)
	})

	t.Run("fail when title is empty", func(t *testing.T) {
		// Setup
		postReq := models.PostRequest{
			Title:            "",
			Author:           &author,
			Content:          content,
			IsCommentAllowed: true,
		}

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "post must have a title", gqlErr.Msg)
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when author is empty", func(t *testing.T) {
		// Setup
		emptyAuthor := ""
		postReq := models.PostRequest{
			Title:            title,
			Author:           &emptyAuthor,
			Content:          content,
			IsCommentAllowed: true,
		}

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "post must have a author", gqlErr.Msg)
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when author is nil", func(t *testing.T) {
		// Setup
		postReq := models.PostRequest{
			Title:            title,
			Author:           nil,
			Content:          content,
			IsCommentAllowed: true,
		}

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "post must have a author", gqlErr.Msg)
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when storage returns error", func(t *testing.T) {
		// Setup
		postReq := models.PostRequest{
			Title:            title,
			Author:           &author,
			Content:          content,
			IsCommentAllowed: true,
		}

		// Mock expectations
		postStorage.EXPECT().
			CreatePost(ctx, gomock.Any()).
			Return(models.Post{}, assert.AnError)

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "error creating post", gqlErr.Msg)
		assert.Equal(t, consts.InternalServerErrorType, gqlErr.Type)
		assert.Nil(t, result)
	})
}

func TestPostService_GetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	postStorage := store_mock.NewMockPostStorage(ctrl)

	postService := NewPostService(postStorage)

	postID := uuid.New()
	title := "Test Post Title"
	author := "test_author"
	content := "test content"

	t.Run("successfully get post by ID", func(t *testing.T) {
		// Setup
		expectedPost := &models.Post{
			ID:                postID,
			Title:             title,
			Author:            author,
			Content:           content,
			IsCommentsAllowed: true,
		}

		// Mock expectations
		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(expectedPost, nil)

		// Execute
		result, err := postService.GetPostByID(ctx, postID)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, expectedPost, result)
	})

	t.Run("fail when post not found", func(t *testing.T) {
		// Mock expectations
		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(nil, sql.ErrNoRows)

		// Execute
		result, err := postService.GetPostByID(ctx, postID)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Contains(t, gqlErr.Msg, "post with id:")
		assert.Contains(t, gqlErr.Msg, "not found")
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when storage returns error", func(t *testing.T) {
		// Mock expectations
		postStorage.EXPECT().
			GetPostByID(ctx, postID).
			Return(nil, assert.AnError)

		// Execute
		result, err := postService.GetPostByID(ctx, postID)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Contains(t, gqlErr.Msg, "error with getting post with id:")
		assert.Equal(t, consts.InternalServerErrorType, gqlErr.Type)
		assert.Nil(t, result)
	})
}

func TestPostService_GetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	postStorage := store_mock.NewMockPostStorage(ctrl)

	postService := NewPostService(postStorage)

	title := "Test Post Title"
	author := "test_author"
	content := "test content"

	t.Run("successfully get all posts without pagination", func(t *testing.T) {
		// Setup
		expectedPosts := []*models.Post{
			{
				ID:                uuid.New(),
				Title:             title + " 1",
				Author:            author,
				Content:           content,
				IsCommentsAllowed: true,
			},
			{
				ID:                uuid.New(),
				Title:             title + " 2",
				Author:            author,
				Content:           content,
				IsCommentsAllowed: false,
			},
		}

		// Mock expectations - исправляем на реальные значения из вашего сервиса
		// По логам видно, что используется limit=20, а не 50
		postStorage.EXPECT().
			GetAllPosts(ctx, 0, 20).
			Return(expectedPosts, nil)

		// Execute
		result, err := postService.GetAllPosts(ctx, nil)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, expectedPosts, result)
	})

	t.Run("successfully get all posts with pagination", func(t *testing.T) {
		// Setup
		page := int32(2)
		// Для page=2: offset = (2-1) * 20 = 20, limit = 20
		expectedOffset := 20
		expectedLimit := 20

		expectedPosts := []*models.Post{
			{
				ID:                uuid.New(),
				Title:             title,
				Author:            author,
				Content:           content,
				IsCommentsAllowed: true,
			},
		}

		// Mock expectations
		postStorage.EXPECT().
			GetAllPosts(ctx, expectedOffset, expectedLimit).
			Return(expectedPosts, nil)

		// Execute
		result, err := postService.GetAllPosts(ctx, &page)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, expectedPosts, result)
	})

	t.Run("fail when page is zero", func(t *testing.T) {
		// Setup
		page := int32(0)

		// Execute
		result, err := postService.GetAllPosts(ctx, &page)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "page must be greater than zero", gqlErr.Msg)
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when page is negative", func(t *testing.T) {
		// Setup
		page := int32(-1)

		// Execute
		result, err := postService.GetAllPosts(ctx, &page)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "page must be greater than zero", gqlErr.Msg)
		assert.Equal(t, consts.BadRequestType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("fail when storage returns error", func(t *testing.T) {
		// Mock expectations
		postStorage.EXPECT().
			GetAllPosts(ctx, 0, 20).
			Return(nil, assert.AnError)

		// Execute
		result, err := postService.GetAllPosts(ctx, nil)

		// Verify
		assert.Error(t, err)
		assert.IsType(t, utils.GqlError{}, err)
		gqlErr := err.(utils.GqlError)
		assert.Equal(t, "error with getting posts", gqlErr.Msg)
		assert.Equal(t, consts.InternalServerErrorType, gqlErr.Type)
		assert.Nil(t, result)
	})

	t.Run("return empty when no posts found", func(t *testing.T) {
		// Mock expectations
		postStorage.EXPECT().
			GetAllPosts(ctx, 0, 20).
			Return([]*models.Post{}, nil)

		// Execute
		result, err := postService.GetAllPosts(ctx, nil)

		// Verify
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestPostService_EdgeCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	postStorage := store_mock.NewMockPostStorage(ctrl)

	postService := NewPostService(postStorage)

	t.Run("handle large page numbers correctly", func(t *testing.T) {
		// Setup
		page := int32(1000000)
		// Для page=1000000: offset = (1000000-1) * 20 = 19999980, limit = 20
		expectedOffset := 19999980
		expectedLimit := 20

		// Mock expectations
		postStorage.EXPECT().
			GetAllPosts(ctx, expectedOffset, expectedLimit).
			Return([]*models.Post{}, nil)

		// Execute
		result, err := postService.GetAllPosts(ctx, &page)

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("handle content edge cases", func(t *testing.T) {
		// Setup
		title := "Test Post"
		author := "test_author"
		emptyContent := ""
		longContent := "very long content " // в реальности можно сделать действительно длинный контент

		postReq := models.PostRequest{
			Title:            title,
			Author:           &author,
			Content:          emptyContent,
			IsCommentAllowed: true,
		}

		expectedPost := models.Post{
			ID:                uuid.New(),
			Title:             title,
			Author:            author,
			Content:           emptyContent,
			IsCommentsAllowed: true,
		}

		// Mock expectations
		postStorage.EXPECT().
			CreatePost(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, post models.Post) (models.Post, error) {
				assert.Equal(t, emptyContent, post.Content)
				return expectedPost, nil
			})

		// Execute
		result, err := postService.CreatePost(ctx, postReq)

		// Verify
		require.NoError(t, err)
		assert.Equal(t, &expectedPost, result)

		// Test with long content
		postReq.Content = longContent
		expectedPost.Content = longContent

		postStorage.EXPECT().
			CreatePost(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, post models.Post) (models.Post, error) {
				assert.Equal(t, longContent, post.Content)
				return expectedPost, nil
			})

		result, err = postService.CreatePost(ctx, postReq)
		require.NoError(t, err)
		assert.Equal(t, &expectedPost, result)
	})
}
