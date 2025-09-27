package service

import (
	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostService interface {
	GetAllPosts(page *int) ([]models.Post, error)
	GetPostByID(id uuid.UUID) (*models.Post, error)
	CreatePost(post models.Post) (*models.Post, error)
}

type CommentService interface {
	CreateComment(comment models.Comment) (*models.Comment, error)
	GetCommentsByPostID(postID uuid.UUID, page *int) ([]models.Comment, error)
	GetRepliesByComment(commentID uuid.UUID) ([]models.Comment, error)
}

type ViewerService interface {
	CreateViewer(postId uuid.UUID) (int, chan *models.Comment, error)
	DeleteViewer(postId uuid.UUID, chanId int) error
	NotifyViewers(postId uuid.UUID, comment models.Comment) error
}
