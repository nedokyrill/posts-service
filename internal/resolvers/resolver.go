package resolvers

import "github.com/nedokyrill/posts-service/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService    service.PostService
	CommentService service.CommentService
	ViewerService  service.ViewerService
}
