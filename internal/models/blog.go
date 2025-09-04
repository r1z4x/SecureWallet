package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlogPost represents a blog post in the system
type BlogPost struct {
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	Title       string     `json:"title" gorm:"type:varchar(500);not null"`
	Slug        string     `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Excerpt     string     `json:"excerpt" gorm:"type:text"`
	Content     string     `json:"content" gorm:"type:longtext"`
	Image       string     `json:"image"`
	Category    string     `json:"category" gorm:"not null"`
	Tags        string     `json:"tags" gorm:"type:text"` // JSON string of tags
	ReadTime    int        `json:"read_time"`
	AuthorID    uuid.UUID  `json:"author_id" gorm:"type:char(36)"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	Status      string     `json:"status" gorm:"default:'published'"` // draft, published, archived
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BlogComment represents a comment on a blog post
type BlogComment struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	PostID    uuid.UUID `json:"post_id" gorm:"type:char(36);not null"`
	Post      BlogPost  `json:"post" gorm:"foreignKey:PostID"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Status    string    `json:"status" gorm:"default:'approved'"` // pending, approved, rejected
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BlogCategory represents a category for blog posts
type BlogCategory struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"default:'#3B82F6'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BlogTag represents a tag for blog posts
type BlogTag struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Slug      string    `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID for blog models
func (b *BlogPost) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *BlogComment) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *BlogCategory) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *BlogTag) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
