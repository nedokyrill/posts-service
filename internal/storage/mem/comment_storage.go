package mem

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/pkg/consts"
)

type CommentsStorageMem struct {
	comms []models.Comment
	mu    sync.RWMutex
}

func NewCommentsStorageMem() *CommentsStorageMem {
	return &CommentsStorageMem{
		comms: make([]models.Comment, 0, consts.InitCommentsSizeInMem),
	}
}

func (s *CommentsStorageMem) CreateComment(_ context.Context, comment models.Comment) (models.Comment, error) {
	now := time.Now()
	comment.ID = uuid.New()
	comment.CreatedAt = &now

	s.mu.Lock()
	defer s.mu.Unlock()

	s.comms = append(s.comms, comment)
	return comment, nil
}
func (s *CommentsStorageMem) GetCommentsByPostID(_ context.Context, postID uuid.UUID,
	offset, limit int) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := len(s.comms) - 1; i >= 0; i-- { // order by created_at desc
		if s.comms[i].PostID == postID && s.comms[i].ParentCommentID == nil {
			comments = append(comments, &s.comms[i])
		}
	}

	if offset > len(s.comms) {
		return nil, nil
	}

	if offset+limit > len(s.comms) {
		return comments[offset:], nil
	}

	return comments[offset : offset+limit], nil
}
func (s *CommentsStorageMem) GetRepliesByParentCommentID(_ context.Context, parentCommentID uuid.UUID) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := len(s.comms) - 1; i >= 0; i-- {
		if s.comms[i].ParentCommentID == &parentCommentID {
			comments = append(comments, &s.comms[i])
		}
	}

	return comments, nil
}
