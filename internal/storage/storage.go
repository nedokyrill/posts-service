package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostStorage interface {
	GetAllPosts(ctx context.Context, offset, limit int) ([]*models.Post, error)
	GetPostByID(ctx context.Context, postId uuid.UUID) (*models.Post, error)
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)
}

type CommentStorage interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, offset, limit int) ([]*models.Comment, error)
	GetRepliesByComment(ctx context.Context, parentCommentID uuid.UUID) ([]*models.Comment, error)
}
