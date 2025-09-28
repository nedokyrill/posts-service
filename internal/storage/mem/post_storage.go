package mem

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/pkg/consts"
)

type PostStorageMem struct {
	posts []*models.Post
	mu    sync.RWMutex
}

func NewPostStorageMem() *PostStorageMem {
	return &PostStorageMem{
		posts: make([]*models.Post, consts.InitPostsSizeInMem),
	}
}

func (s *PostStorageMem) GetAllPosts(_ context.Context, offset, limit int) ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if offset > len(s.posts) {
		return nil, nil
	}

	if offset+limit > len(s.posts) {
		return s.posts[offset:], nil
	}

	return s.posts[offset : offset+limit], nil
}

func (s *PostStorageMem) CreatePost(_ context.Context, post models.Post) (models.Post, error) {
	now := time.Now()
	post.ID = uuid.New()
	post.CreatedAt = &now

	s.mu.Lock()
	defer s.mu.Unlock()

	s.posts = append(s.posts, &post)
	return post, nil
}

func (s *PostStorageMem) GetPostByID(_ context.Context, postId uuid.UUID) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.posts {
		if s.posts[i].ID == postId {
			return s.posts[i], nil
		}
	}

	return &models.Post{}, errors.New(fmt.Sprintf("Post with id: %s not found", postId.String()))
}
