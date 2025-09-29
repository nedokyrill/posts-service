package mem

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostStorageMem_CreatePost(t *testing.T) {
	ctx := context.Background()
	title := "test_title"
	author := "test_author"
	content := "test_content"

	t.Run("with auto-generated fields", func(t *testing.T) {
		storage := NewPostStorageMem()

		post := models.Post{
			Title:             title,
			Author:            author,
			Content:           content,
			IsCommentsAllowed: true,
		}

		created, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, created.ID)

		assert.NotNil(t, created.CreatedAt)
		assert.WithinDuration(t, time.Now(), *created.CreatedAt, time.Second)

		assert.Equal(t, title, created.Title)
		assert.Equal(t, author, created.Author)
		assert.Equal(t, content, created.Content)
		assert.True(t, created.IsCommentsAllowed)
	})

	t.Run("field persistence with different values", func(t *testing.T) {
		storage := NewPostStorageMem()

		testCases := []struct {
			name              string
			title             string
			author            string
			content           string
			isCommentsAllowed bool
		}{
			{
				name:              "with comments allowed",
				title:             "Post with comments",
				author:            "author1",
				content:           "Content 1",
				isCommentsAllowed: true,
			},
			{
				name:              "with comments disabled",
				title:             "Post without comments",
				author:            "author2",
				content:           "Content 2",
				isCommentsAllowed: false,
			},
			{
				name:              "empty content",
				title:             "Empty Content Post",
				author:            "author3",
				content:           "",
				isCommentsAllowed: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				post := models.Post{
					Title:             tc.title,
					Author:            tc.author,
					Content:           tc.content,
					IsCommentsAllowed: tc.isCommentsAllowed,
				}

				created, err := storage.CreatePost(ctx, post)
				require.NoError(t, err)

				retrieved, err := storage.GetPostByID(ctx, created.ID)
				require.NoError(t, err)

				assert.Equal(t, tc.title, retrieved.Title)
				assert.Equal(t, tc.author, retrieved.Author)
				assert.Equal(t, tc.content, retrieved.Content)
				assert.Equal(t, tc.isCommentsAllowed, retrieved.IsCommentsAllowed)
			})
		}
	})
}

func TestPostStorageMem_GetPostByID(t *testing.T) {
	ctx := context.Background()
	title := "test_title"
	author := "test_author"
	content := "test_content"

	t.Run("with existing post", func(t *testing.T) {
		storage := NewPostStorageMem()

		post := models.Post{
			Title:             title,
			Author:            author,
			Content:           content,
			IsCommentsAllowed: true,
		}

		created, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		retrieved, err := storage.GetPostByID(ctx, created.ID)
		require.NoError(t, err)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, title, retrieved.Title)
		assert.Equal(t, author, retrieved.Author)
		assert.Equal(t, content, retrieved.Content)
		assert.True(t, retrieved.IsCommentsAllowed)
		assert.Equal(t, created.CreatedAt, retrieved.CreatedAt)
	})

	t.Run("with non-existent post", func(t *testing.T) {
		storage := NewPostStorageMem()

		nonExistentID := uuid.New()
		post, err := storage.GetPostByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), nonExistentID.String())
		assert.NotNil(t, post)
		assert.Equal(t, uuid.Nil, post.ID)
		assert.Empty(t, post.Title)
	})
}

func TestPostStorageMem_GetAllPosts(t *testing.T) {
	ctx := context.Background()
	title := "test_title"
	author := "test_author"
	content := "test_content"

	t.Run("with empty storage", func(t *testing.T) {
		storage := NewPostStorageMem()

		posts, err := storage.GetAllPosts(ctx, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("with single post", func(t *testing.T) {
		storage := NewPostStorageMem()

		post := models.Post{
			Title:             title,
			Author:            author,
			Content:           content,
			IsCommentsAllowed: true,
		}

		created, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		posts, err := storage.GetAllPosts(ctx, 0, 10)
		require.NoError(t, err)
		assert.Len(t, posts, 1)
		assert.Equal(t, created.ID, posts[0].ID)
	})

	t.Run("with multiple posts", func(t *testing.T) {
		storage := NewPostStorageMem()

		postCount := 5
		for i := 0; i < postCount; i++ {
			post := models.Post{
				Title:             fmt.Sprintf("%s %d", title, i),
				Author:            author,
				Content:           content,
				IsCommentsAllowed: i%2 == 0, // чередуем разрешение комментариев
			}
			_, err := storage.CreatePost(ctx, post)
			require.NoError(t, err)
			time.Sleep(time.Millisecond) // для разных времен создания
		}

		posts, err := storage.GetAllPosts(ctx, 0, 10)
		require.NoError(t, err)
		assert.Len(t, posts, postCount)

		for i := 1; i < len(posts); i++ {
			assert.True(t, posts[i-1].CreatedAt.Before(*posts[i].CreatedAt) ||
				posts[i-1].CreatedAt.Equal(*posts[i].CreatedAt))
		}
	})

	t.Run("with pagination", func(t *testing.T) {
		storage := NewPostStorageMem()

		postCount := 10
		for i := 0; i < postCount; i++ {
			post := models.Post{
				Title:             fmt.Sprintf("%s %d", title, i),
				Author:            author,
				Content:           content,
				IsCommentsAllowed: true,
			}
			_, err := storage.CreatePost(ctx, post)
			require.NoError(t, err)
		}

		posts, err := storage.GetAllPosts(ctx, 0, 3)
		require.NoError(t, err)
		assert.Len(t, posts, 3)

		posts, err = storage.GetAllPosts(ctx, 5, 10)
		require.NoError(t, err)
		assert.Len(t, posts, 5)

		posts, err = storage.GetAllPosts(ctx, 2, 4)
		require.NoError(t, err)
		assert.Len(t, posts, 4)

		posts, err = storage.GetAllPosts(ctx, 15, 10)
		require.NoError(t, err)
		assert.Nil(t, posts)
	})

	t.Run("with boundary conditions", func(t *testing.T) {
		storage := NewPostStorageMem()

		for i := 0; i < 3; i++ {
			post := models.Post{
				Title:             fmt.Sprintf("%s %d", title, i),
				Author:            author,
				Content:           content,
				IsCommentsAllowed: true,
			}
			_, err := storage.CreatePost(ctx, post)
			require.NoError(t, err)
		}

		posts, err := storage.GetAllPosts(ctx, 3, 10)
		require.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Empty(t, posts)

		posts, err = storage.GetAllPosts(ctx, 1, 5)
		require.NoError(t, err)
		assert.Len(t, posts, 2)
	})
}

func TestPostStorageMem_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	storage := NewPostStorageMem()

	goroutines := 10
	postsPerRoutine := 5

	errChan := make(chan error, goroutines)
	createdIDs := make(chan uuid.UUID, goroutines*postsPerRoutine)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			for j := 0; j < postsPerRoutine; j++ {
				post := models.Post{
					Title:             fmt.Sprintf("Concurrent Post %d-%d", index, j),
					Author:            "test_author",
					Content:           "test content",
					IsCommentsAllowed: true,
				}
				created, err := storage.CreatePost(ctx, post)
				if err != nil {
					errChan <- err
					return
				}
				createdIDs <- created.ID
			}
			errChan <- nil
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		assert.NoError(t, <-errChan)
	}
	close(createdIDs)

	posts, err := storage.GetAllPosts(ctx, 0, goroutines*postsPerRoutine)
	require.NoError(t, err)
	assert.Len(t, posts, goroutines*postsPerRoutine)

	for id := range createdIDs {
		post, err := storage.GetPostByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, post.ID)
	}
}
