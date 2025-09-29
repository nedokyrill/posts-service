package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostStorage interface {
	GetAllPosts(ctx context.Context, offset, limit int) ([]*models.Post, error) // получение списка всех постов
	GetPostByID(ctx context.Context, postId uuid.UUID) (*models.Post, error)    // получение поста по его id
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)      // создание поста
}

type CommentStorage interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)                       // создание комментария
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, offset, limit int) ([]*models.Comment, error) // получение комментариев по id поста
	GetRepliesByParentCommentID(ctx context.Context, parentCommentID uuid.UUID) ([]*models.Comment, error)   // получение ответов на комментарий по его id
}
