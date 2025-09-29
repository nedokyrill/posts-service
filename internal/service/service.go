package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostService interface {
	GetAllPosts(ctx context.Context, page *int32) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
	CreatePost(ctx context.Context, postReq models.PostRequest) (*models.Post, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, commReq models.CommentRequest) (*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, page *int32) ([]*models.Comment, error)
	GetRepliesByComment(ctx context.Context, commentID uuid.UUID) ([]*models.Comment, error)
}

type ViewerService interface {
	CreateViewer(ctx context.Context, postId uuid.UUID) (int, chan *models.Comment, error)
	DeleteViewer(ctx context.Context, postId uuid.UUID, id int) error
	NotifyViewers(ctx context.Context, postId uuid.UUID, comment models.Comment) error
}
