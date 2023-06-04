package posts

import (
	"time"

	postmodel "github.com/pedromspeixoto/posts-api/internal/data/models/posts"
	"github.com/pedromspeixoto/posts-api/internal/pkg/uuid"
)

type Post struct {
}

// request
type PostRequest struct {
	Content string `json:"content" validate:"required"`
}

func ModelFromPostRequest(post *PostRequest) *postmodel.Post {
	model := &postmodel.Post{
		PostId:  uuid.GenerateUUID(),
		Content: post.Content,
	}
	return model
}

// response
type PostResponse struct {
	PostId    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewPostResponse(post *postmodel.Post) *PostResponse {
	resp := &PostResponse{
		PostId:    post.PostId,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}
	return resp
}

type PostListResponse struct {
	Posts []PostResponse `json:"posts,omitempty"`
}

func NewPostListResponse(models []postmodel.Post) *PostListResponse {
	var posts []PostResponse
	for _, m := range models {
		posts = append(posts, PostResponse{
			PostId:    m.PostId,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}
	resp := &PostListResponse{
		Posts: posts,
	}
	return resp
}
