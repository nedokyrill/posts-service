package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedokyrill/posts-service/internal/models"
)

type PostStorePgx struct {
	db *pgxpool.Pool
}

func NewPostStorePgx(db *pgxpool.Pool) *PostStorePgx {
	return &PostStorePgx{
		db: db,
	}
}

func (s *PostStorePgx) GetAllPosts(ctx context.Context, offset, limit int) ([]*models.Post, error) {
	var posts []*models.Post

	query := `SELECT * FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2;`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var title, author, content string
		var isCommentsAllowed bool
		var createdAt *time.Time

		if err = rows.Scan(&id, &title, &content, &author, &isCommentsAllowed, &createdAt); err != nil {
			return nil, err
		}
		posts = append(posts, &models.Post{ID: id, Title: title, Content: content, Author: author,
			IsCommentsAllowed: isCommentsAllowed, CreatedAt: createdAt})
	}
	return posts, nil
}

func (s *PostStorePgx) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	query := `INSERT INTO posts (title, content, author, is_comments_allowed) 
				VALUES ($1, $2, $3, $4) RETURNING id, created_at;`

	var id uuid.UUID
	var createdAt time.Time

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.Author,
		post.IsCommentsAllowed).Scan(&id, &createdAt)
	if err != nil {
		return models.Post{}, err
	}

	post.ID = id
	post.CreatedAt = &createdAt
	return post, nil
}

func (s *PostStorePgx) GetPostByID(ctx context.Context, postId uuid.UUID) (*models.Post, error) {
	var post models.Post

	query := `SELECT * FROM posts WHERE id = $1;`
	err := s.db.QueryRow(ctx, query, postId).Scan(&post.ID, &post.Title, &post.Content, &post.Author,
		&post.IsCommentsAllowed, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
