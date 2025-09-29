package mem

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentsStorageMem_CreateComment(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	author := "Test_author"
	content := "Test_test_test"

	t.Run("without ID", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		comment := models.Comment{
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		created, err := storage.CreateComment(ctx, comment)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.NotNil(t, created.CreatedAt)
		assert.WithinDuration(t, time.Now(), *created.CreatedAt, time.Second)
		assert.Equal(t, author, created.Author)
		assert.Equal(t, content, created.Content)
		assert.Equal(t, postID, created.PostID)
		assert.Nil(t, created.ParentCommentID)
	})

	t.Run("with predefined ID", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		predefinedID := uuid.New()
		comment := models.Comment{
			ID:      predefinedID,
			Author:  author,
			Content: content,
			PostID:  postID,
		}

		created, err := storage.CreateComment(ctx, comment)
		require.NoError(t, err)
		assert.Equal(t, predefinedID, created.ID)
	})
}

func TestCommentsStorageMem_GetCommentsByPostID(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	author := "Test_author"
	content := "Test_test_test"

	t.Run("with invalid arguments", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		_, err := storage.GetCommentsByPostID(ctx, postID, -1, 10)
		assert.Error(t, err)

		_, err = storage.GetCommentsByPostID(ctx, postID, 0, -5)
		assert.Error(t, err)
	})

	t.Run("with non-existent post", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		comments, err := storage.GetCommentsByPostID(ctx, uuid.New(), 0, 10)
		require.NoError(t, err)
		assert.Empty(t, comments)
	})

	t.Run("with pagination", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		for i := 0; i < 5; i++ {
			comment := models.Comment{
				Author:  author,
				Content: content,
				PostID:  postID,
			}
			_, err := storage.CreateComment(ctx, comment)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		retrieved, err := storage.GetCommentsByPostID(ctx, postID, 0, 10)
		require.NoError(t, err)
		assert.Len(t, retrieved, 5)

		assert.True(t, retrieved[0].CreatedAt.After(*retrieved[1].CreatedAt))

		retrieved, err = storage.GetCommentsByPostID(ctx, postID, 0, 3)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		retrieved, err = storage.GetCommentsByPostID(ctx, postID, 2, 10)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		retrieved, err = storage.GetCommentsByPostID(ctx, postID, 1, 2)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
	})

	t.Run("only returns root comments", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		rootComment, err := storage.CreateComment(ctx, models.Comment{
			Author:  author,
			Content: "Root comment",
			PostID:  postID,
		})
		require.NoError(t, err)

		_, err = storage.CreateComment(ctx, models.Comment{
			Author:          author,
			Content:         "Reply 1",
			PostID:          postID,
			ParentCommentID: &rootComment.ID,
		})
		require.NoError(t, err)

		_, err = storage.CreateComment(ctx, models.Comment{
			Author:          author,
			Content:         "Reply 2",
			PostID:          postID,
			ParentCommentID: &rootComment.ID,
		})
		require.NoError(t, err)

		rootComments, err := storage.GetCommentsByPostID(ctx, postID, 0, 10)
		require.NoError(t, err)
		assert.Len(t, rootComments, 1)
		assert.Equal(t, rootComment.ID, rootComments[0].ID)
	})
}

func TestCommentsStorageMem_GetRepliesByParentCommentID(t *testing.T) {
	ctx := context.Background()
	postID := uuid.New()
	author := "Test_author"
	content := "Test_test_test"

	t.Run("with existing parent comment", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		parentComment, err := storage.CreateComment(ctx, models.Comment{
			Author:  author,
			Content: content,
			PostID:  postID,
		})
		require.NoError(t, err)

		reply1, err := storage.CreateComment(ctx, models.Comment{
			Author:          author,
			Content:         "First reply",
			PostID:          postID,
			ParentCommentID: &parentComment.ID,
		})
		require.NoError(t, err)

		reply2, err := storage.CreateComment(ctx, models.Comment{
			Author:          author,
			Content:         "Second reply",
			PostID:          postID,
			ParentCommentID: &parentComment.ID,
		})
		require.NoError(t, err)

		replies, err := storage.GetRepliesByParentCommentID(ctx, parentComment.ID)
		require.NoError(t, err)
		assert.Len(t, replies, 2)

		assert.True(t, replies[0].CreatedAt.After(*replies[1].CreatedAt))

		replyIDs := []uuid.UUID{replies[0].ID, replies[1].ID}
		assert.Contains(t, replyIDs, reply1.ID)
		assert.Contains(t, replyIDs, reply2.ID)
	})

	t.Run("with non-existent parent", func(t *testing.T) {
		storage := NewCommentsStorageMem()

		replies, err := storage.GetRepliesByParentCommentID(ctx, uuid.New())
		require.NoError(t, err)
		assert.Empty(t, replies)
	})
}

func TestCommentsStorageMem_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	storage := NewCommentsStorageMem()

	postID := uuid.New()
	author := "Test_author"
	content := "Test_test_test"

	goroutines := 10
	commentsPerRoutine := 5

	errChan := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			for j := 0; j < commentsPerRoutine; j++ {
				_, err := storage.CreateComment(ctx, models.Comment{
					Author:  author,
					Content: content,
					PostID:  postID,
				})
				if err != nil {
					errChan <- err
					return
				}
			}
			errChan <- nil
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		assert.NoError(t, <-errChan)
	}

	comments, err := storage.GetCommentsByPostID(ctx, postID, 0, 50)
	require.NoError(t, err)
	assert.Len(t, comments, goroutines*commentsPerRoutine)
}
