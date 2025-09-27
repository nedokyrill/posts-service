package storage

import (
	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostStorage interface {
	GetAllPosts(offset, limit int) ([]*models.Post, error)
	GetPostByID(postId uuid.UUID) (models.Post, error)
	CreatePost(post models.Post) (models.Post, error)
}

type CommentStorage interface {
	CreateComment(comment models.Comment) (models.Comment, error)
	GetCommentsByPostID(postID uuid.UUID, offset, limit int) ([]*models.Comment, error)
	GetRepliesByComment(parentCommentID uuid.UUID) ([]*models.Comment, error)
}
