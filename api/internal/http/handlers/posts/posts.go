package posts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pedromspeixoto/posts-api/internal/config"
	"github.com/pedromspeixoto/posts-api/internal/domain/posts"
	"github.com/pedromspeixoto/posts-api/internal/dto"
	postsdto "github.com/pedromspeixoto/posts-api/internal/dto/posts"
	"github.com/pedromspeixoto/posts-api/internal/http/handlers/common"
	"github.com/pedromspeixoto/posts-api/internal/http/middlewares"
	"github.com/pedromspeixoto/posts-api/internal/pkg/logger"
	"go.uber.org/fx"
)

type PostServiceHandler interface {
	Routes() chi.Router
}

type postServiceDeps struct {
	fx.In

	Config      *config.Config
	Logger      *logger.LoggingClient
	Validator   *validator.Validate
	PostService posts.PostService
}

type postServiceHandler struct {
	postServiceDeps
	logger.Logger
}

func NewPostServiceHandler(deps postServiceDeps) PostServiceHandler {
	return &postServiceHandler{
		postServiceDeps: deps,
		Logger:          deps.Logger.GetLogger(),
	}
}

func (h postServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// posts
	r.With(middlewares.Paginate).Get("/", h.ListPosts)
	r.Post("/", h.CreatePost)
	r.Get("/{postId}", h.GetPost)
	r.Put("/{postId}", h.UpdatePost)
	r.Delete("/{postId}", h.DeletePost)

	return r
}

// CreatePost Handler - Handles posts requests creation
// @Summary Create a new post request.
// @Description This API is used to create a new post request
// @Param request body postsdto.PostRequest true "Post Payload"
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /v1/posts [post]
func (h postServiceHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	post := postsdto.PostRequest{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		common.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Validator.Struct(post)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, postResponse, err := h.PostService.CreatePost(r.Context(), &post)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "new post created", postResponse)
}

// ListPosts - Handles posts requests creation
// @Summary Gets all post requests.
// @Description This API is used to list all post request created
// @Param limit query int false "Limit"
// @Param page  query int false "Page"
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /v1/posts [get]
func (h postServiceHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(middlewares.LimitKey).(int)
	page := r.Context().Value(middlewares.PageKey).(int)
	sort := r.Context().Value(middlewares.SortKey).(string)
	filter := r.Context().Value(middlewares.FilterKey).(map[string]string)
	search := r.Context().Value(middlewares.SearchKey).(map[string]string)

	pageRequest, err := dto.NewPaginationRequest(limit, page, sort, filter, search)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
	}

	statusCode, env, err := h.postServiceDeps.PostService.ListPosts(r.Context(), pageRequest)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "posts retrieved", env)
}

// GetPost - Handles posts requests creation
// @Summary Get an post request.
// @Description This API is used to get post request created
// @Param post_id path string true "Post Id"
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /v1/posts/{post_id} [get]
func (h postServiceHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	statusCode, post, err := h.postServiceDeps.PostService.GetPost(r.Context(), postId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "post retrieved", post)
}

// UpdatePost - Handles posts requests updates
// @Summary Updates an post request.
// @Description This API is used to update an post request
// @Param post_id path string true "Post Id"
// @Param request body postsdto.PostRequest true "Post Update Payload"
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /v1/posts/{post_id} [put]
func (h postServiceHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	post := postsdto.PostRequest{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		common.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Validator.Struct(post)
	if err != nil {
		common.Err(w, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, postResponse, err := h.PostService.UpdatePost(r.Context(), postId, &post)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "post updated", postResponse)
}

// DeletePost - Handles posts requests creation
// @Summary Delete an post request.
// @Description This API is used to delete an post request created
// @Param post_id path string true "Post Id"
// @Tags posts
// @Accept  json
// @Produce  json
// @Router /v1/posts/{post_id} [delete]
func (h postServiceHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	statusCode, env, err := h.postServiceDeps.PostService.DeletePost(r.Context(), postId)
	if err != nil {
		common.Err(w, statusCode, err.Error())
		return
	}

	common.Json(w, statusCode, "", env)
}
