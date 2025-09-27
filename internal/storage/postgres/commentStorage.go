package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedokyrill/posts-service/internal/models"
)

type CommentsStorePgx struct {
	db *pgxpool.Pool
}

func NewCommentsStorePgx(db *pgxpool.Pool) *CommentsStorePgx {
	return &CommentsStorePgx{
		db: db,
	}
}

func (s *CommentsStorePgx) CreateComment(comment models.Comment) (models.Comment, error) {
	query := `INSERT INTO comments (author, content, post_id, parent_comment_id)
				VALUES ($1, $2, $3, $4) RETURNING id, created_at;`

	var id uuid.UUID
	var createdAt time.Time

	err := s.db.QueryRow(context.Background(), query, comment.Author, comment.Content, comment.PostID,
		comment.ParentCommentID).Scan(&id, &createdAt)
	if err != nil {
		return models.Comment{}, err
	}

	comment.ID = id
	comment.CreatedAt = &createdAt
	return comment, nil
}

func (s *CommentsStorePgx) GetCommentsByPostID(postID uuid.UUID, offset, limit int) ([]models.Comment, error) {
	var comments []models.Comment
	query := `SELECT * FROM comments WHERE post_id = $1 AND parent_comment_id IS NULL 
         		ORDER BY created_at DESC LIMIT $2 OFFSET $3;`

	rows, err := s.db.Query(context.Background(), query, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id, postId uuid.UUID
		var parentId *uuid.UUID
		var author, content string
		var createdAt *time.Time

		if err = rows.Scan(&id, &author, &content, &postId, &parentId, &createdAt); err != nil {
			return nil, err
		}
		comments = append(comments, models.Comment{ID: id, Author: author, Content: content, PostID: postId, ParentCommentID: parentId, CreatedAt: createdAt})
	}
	return comments, nil
}

// слишком много одинакового кода, но решил не выносить в отдельный метод
func (s *CommentsStorePgx) GetRepliesByComment(parentCommentID uuid.UUID) ([]models.Comment, error) {
	var replies []models.Comment
	query := `SELECT * FROM comments WHERE parent_comment_id = $1 ORDER BY created_at;`

	rows, err := s.db.Query(context.Background(), query, parentCommentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id, postId uuid.UUID
		var parentId *uuid.UUID
		var author, content string
		var createdAt *time.Time

		if err = rows.Scan(&id, &author, &content, &postId, &parentId, &createdAt); err != nil {
			return nil, err
		}
		replies = append(replies, models.Comment{ID: id, Author: author, Content: content, PostID: postId, ParentCommentID: parentId, CreatedAt: createdAt})
	}
	return replies, nil
}
