package posts

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"net/http"

	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/data/models/posts"
	"github.com/pedromspeixoto/posts-api/internal/dto"
	postsdto "github.com/pedromspeixoto/posts-api/internal/dto/posts"
	"github.com/pedromspeixoto/posts-api/internal/pkg/logger"
	"go.uber.org/fx"
)

// PostService provides methods pertaining to managing posts.
type PostService interface {
	// CreatePost creates a post entry
	CreatePost(ctx context.Context, Post *postsdto.PostRequest) (int, *postsdto.PostResponse, error)
	// ListPosts retrieves all posts with pagination.
	ListPosts(ctx context.Context, pagination *dto.PaginationRequest) (int, *dto.PaginationResponse, error)
	// UpdatePost updates a post entry by uuid
	UpdatePost(ctx context.Context, uuid string, request *postsdto.PostRequest) (int, *postsdto.PostResponse, error)
	// UpsertPost updates or creates a post entry by uuid
	UpsertPost(ctx context.Context, uuid string, request *postsdto.PostRequest) (int, *postsdto.PostResponse, error)
	// GetPost retrieves a post entry by uuid
	GetPost(ctx context.Context, uuid string) (int, *postsdto.PostResponse, error)
	// DeletePost hard deletes a post entry by uuid
	DeletePost(ctx context.Context, uuid string) (int, *postsdto.PostResponse, error)
}

type PostServiceDeps struct {
	fx.In

	Config         *config.Config
	Logger         *logger.LoggingClient
	PostRepository posts.PostRepository
}

type postService struct {
	PostServiceDeps
	logger.Logger
}

func NewPostService(deps PostServiceDeps) PostService {
	return &postService{
		PostServiceDeps: deps,
		Logger:          deps.Logger.GetLogger(),
	}
}

func (p *postService) CreatePost(ctx context.Context, request *postsdto.PostRequest) (int, *postsdto.PostResponse, error) {
	model := postsdto.ModelFromPostRequest(request)
	err := p.PostRepository.Create(model)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("error creating new post: %v", err)
	}

	return http.StatusCreated, postsdto.NewPostResponse(model), nil
}

func (p *postService) ListPosts(ctx context.Context, paginationRequest *dto.PaginationRequest) (int, *dto.PaginationResponse, error) {
	posts, pageEnv, err := p.PostRepository.List(paginationRequest.Limit, paginationRequest.Page)
	if err != nil {
		return http.StatusNotFound, nil, fmt.Errorf("error fetching posts: %v", err)
	}

	pageEnv.Data = postsdto.NewPostListResponse(posts)
	return http.StatusOK, dto.NewPaginationResponse(pageEnv), nil
}

func (p *postService) UpdatePost(ctx context.Context, uuid string, request *postsdto.PostRequest) (int, *postsdto.PostResponse, error) {
	post, err := p.PostRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	post.Content = request.Content
	err = p.PostRepository.Update(post)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error updating post: %v", err)
	}

	return http.StatusOK, postsdto.NewPostResponse(post), nil
}

func (p *postService) UpsertPost(ctx context.Context, uuid string, request *postsdto.PostRequest) (int, *postsdto.PostResponse, error) {
	_, err := p.PostRepository.GetByUUID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return p.CreatePost(ctx, request)
		}
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error fetching post: %v", err)
	}

	return p.UpdatePost(ctx, uuid, request)
}

func (p *postService) GetPost(ctx context.Context, uuid string) (int, *postsdto.PostResponse, error) {
	post, err := p.PostRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	return http.StatusOK, postsdto.NewPostResponse(post), nil
}

func (p *postService) DeletePost(ctx context.Context, uuid string) (int, *postsdto.PostResponse, error) {
	post, err := p.PostRepository.GetByUUID(uuid)
	if err != nil {
		return http.StatusNotFound, nil, err
	}

	err = p.PostRepository.HardDelete(post)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("unexpected error deleting post: %v", err)
	}

	return http.StatusOK, nil, nil
}
