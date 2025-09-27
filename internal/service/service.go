package service

import (
	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostService interface {
	GetAllPosts(page *int32) ([]*models.Post, error)
	GetPostByID(id *uuid.UUID) (*models.Post, error)
	CreatePost(title string, author *string, content string, isCommentAllowed bool) (*models.Post, error)
}

type CommentService interface {
	CreateComment(author string, content string, postID uuid.UUID, parentCommentID *uuid.UUID) (*models.Comment, error)
	GetCommentsByPostID(postID uuid.UUID, page *int32) ([]*models.Comment, error)
	GetRepliesByComment(commentID uuid.UUID) ([]*models.Comment, error)
}

type ViewerService interface {
	CreateViewer(postId uuid.UUID) (int, chan *models.Comment, error)
	DeleteViewer(postId uuid.UUID, id int) error
	NotifyViewers(postId uuid.UUID, comment models.Comment) error
}
