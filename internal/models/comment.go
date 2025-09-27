package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID  `json:"id"`
	Author          string     `json:"author"`
	Content         string     `json:"content"`
	PostID          uuid.UUID  `json:"postId"`
	ParentCommentID *uuid.UUID `json:"parentCommentId,omitempty"`
	//Replies         []*Comment `json:"replies,omitempty"`
	CreatedAt *time.Time `json:"createdAt"`
}
