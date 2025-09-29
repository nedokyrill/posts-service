package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Title             string    `json:"title"`
	Author            string    `json:"author"`
	Content           string    `json:"content"`
	IsCommentsAllowed bool      `json:"isCommentsAllowed"`
	//Comments          []*Comment `json:"comments,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type PostRequest struct {
	Title            string
	Author           *string
	Content          string
	IsCommentAllowed bool
}
