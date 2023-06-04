package posts

import (
	"errors"
	"math"

	"github.com/pedromspeixoto/posts-api/internal/data"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PostId  string
	Content string
}

// postRepository is a repository for dealing with the post object.
type PostRepository interface {
	// List lists posts from the database with pagination.
	List(limit, page int) ([]Post, *data.Pagination, error)
	// GetByUUID gets a post from the database by uuid.
	GetByUUID(uuid string) (*Post, error)
	// Get gets a post from the database by id.
	Get(id uint) (*Post, error)
	// Create creates a post in the database.
	Create(post *Post) error
	// Upsert creates or updates a post if the post already exists.
	Upsert(post *Post) error
	// Update updates a post config in the database. Should be paired with Get
	// to retrieve the existing object, then the object modified and passed to this
	// method.
	Update(post *Post) error
	// SoftDelete soft deletes a post record from the database.
	SoftDelete(post *Post) error
	// HardDelete hard deletes a post record from the database.
	HardDelete(post *Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (p postRepository) List(limit, page int) ([]Post, *data.Pagination, error) {
	var posts []Post

	// pagination object
	pagination := &data.Pagination{
		Limit: limit,
		Page:  page,
	}
	p.db.Scopes(pagination.Paginate()).Find(&posts)

	// pagination details
	p.db.Model(&Post{}).Count(&pagination.TotalRows)
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))

	return posts, pagination, nil
}

func (p postRepository) GetByUUID(uuid string) (*Post, error) {
	post := Post{}
	result := p.db.Unscoped().Where("post_id = ?", uuid).Find(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &post, nil
}

func (p postRepository) Get(id uint) (*Post, error) {
	post := Post{}
	result := p.db.First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (p postRepository) Create(post *Post) error {
	result := p.db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p postRepository) Upsert(upsertPost *Post) error {
	post, err := p.Get(upsertPost.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.Create(upsertPost)
		}
		return err
	}
	post.Content = upsertPost.Content
	return p.Update(upsertPost)
}

func (p postRepository) Update(post *Post) error {
	result := p.db.Save(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p postRepository) SoftDelete(post *Post) error {
	result := p.db.Delete(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p postRepository) HardDelete(post *Post) error {
	result := p.db.Unscoped().Delete(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
