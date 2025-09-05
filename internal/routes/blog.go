package routes

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlogRoutes sets up blog-related routes
func BlogRoutes(router *gin.Engine, db *gorm.DB) {
	blog := router.Group("/api/blog")
	{
		// Get all published blog posts
		blog.GET("/posts", func(c *gin.Context) {
			var posts []struct {
				ID           string `json:"id"`
				Title        string `json:"title"`
				Slug         string `json:"slug"`
				Excerpt      string `json:"excerpt"`
				Image        string `json:"image"`
				Category     string `json:"category"`
				ReadTime     int    `json:"read_time"`
				ViewCount    int    `json:"view_count"`
				CreatedAt    string `json:"created_at"`
				AuthorName   string `json:"author_name"`
				AuthorTitle  string `json:"author_title"`
				AuthorAvatar string `json:"author_avatar"`
			}

			query := db.Table("blog_posts").
				Select(`
					blog_posts.id, blog_posts.title, blog_posts.slug, blog_posts.excerpt,
					blog_posts.image, blog_posts.category, blog_posts.read_time,
					blog_posts.view_count, blog_posts.created_at,
					users.name as author_name, users.title as author_title, users.avatar as author_avatar
				`).
				Joins("LEFT JOIN users ON blog_posts.author_id = users.id").
				Where("blog_posts.status = ?", "published").
				Order("blog_posts.published_at DESC")

			// Apply category filter
			if category := c.Query("category"); category != "" {
				query = query.Where("blog_posts.category = ?", category)
			}

			// Apply search filter
			if search := c.Query("search"); search != "" {
				query = query.Where("blog_posts.title LIKE ? OR blog_posts.excerpt LIKE ?", "%"+search+"%", "%"+search+"%")
			}

			// Apply pagination
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "9"))
			offset := (page - 1) * limit

			var total int64
			query.Count(&total)

			result := query.Offset(offset).Limit(limit).Find(&posts)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blog posts"})
				return
			}

			// Transform the data to match frontend expectations
			var transformedPosts []map[string]interface{}
			for _, post := range posts {
				transformedPosts = append(transformedPosts, map[string]interface{}{
					"id":        post.ID,
					"title":     post.Title,
					"slug":      post.Slug,
					"excerpt":   post.Excerpt,
					"image":     post.Image,
					"category":  post.Category,
					"readTime":  post.ReadTime,
					"viewCount": post.ViewCount,
					"date":      post.CreatedAt,
					"author": map[string]interface{}{
						"name":   post.AuthorName,
						"title":  post.AuthorTitle,
						"avatar": post.AuthorAvatar,
					},
				})
			}

			c.JSON(http.StatusOK, gin.H{
				"posts": transformedPosts,
				"total": total,
				"page":  page,
				"limit": limit,
			})
		})

		// Get a single blog post by slug
		blog.GET("/posts/:slug", func(c *gin.Context) {
			slug := c.Param("slug")

			var post struct {
				ID           string `json:"id"`
				Title        string `json:"title"`
				Slug         string `json:"slug"`
				Excerpt      string `json:"excerpt"`
				Content      string `json:"content"`
				Image        string `json:"image"`
				Category     string `json:"category"`
				Tags         string `json:"tags"`
				ReadTime     int    `json:"read_time"`
				ViewCount    int    `json:"view_count"`
				PublishedAt  string `json:"published_at"`
				CreatedAt    string `json:"created_at"`
				AuthorID     string `json:"author_id"`
				AuthorName   string `json:"author_name"`
				AuthorTitle  string `json:"author_title"`
				AuthorAvatar string `json:"author_avatar"`
				AuthorBio    string `json:"author_bio"`
			}

			// First get the raw data to see what we're working with
			var rawPost struct {
				ID           string `json:"id"`
				Title        string `json:"title"`
				Slug         string `json:"slug"`
				Excerpt      string `json:"excerpt"`
				Content      string `json:"content"`
				Image        string `json:"image"`
				Category     string `json:"category"`
				Tags         string `json:"tags"`
				ReadTime     int    `json:"read_time"`
				ViewCount    int    `json:"view_count"`
				PublishedAt  string `json:"published_at"`
				CreatedAt    string `json:"created_at"`
				AuthorID     string `json:"author_id"`
				AuthorName   string `json:"author_name"`
				AuthorTitle  string `json:"author_title"`
				AuthorAvatar string `json:"author_avatar"`
				AuthorBio    string `json:"author_bio"`
			}

			result := db.Table("blog_posts").
				Select(`
					blog_posts.id, blog_posts.title, blog_posts.slug, blog_posts.excerpt,
					blog_posts.content, blog_posts.image, blog_posts.category, blog_posts.tags,
					blog_posts.read_time, blog_posts.view_count, blog_posts.published_at,
					blog_posts.created_at,
					users.id as author_id, users.name as author_name, users.title as author_title,
					users.avatar as author_avatar, users.bio as author_bio
				`).
				Joins("LEFT JOIN users ON blog_posts.author_id = users.id").
				Where("blog_posts.slug = ? AND blog_posts.status = ?", slug, "published").
				First(&rawPost)

			// Copy to post struct
			post = rawPost

			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
				return
			}

			// Debug: Log the raw tags data
			log.Printf("DEBUG: Raw tags from DB: %q", post.Tags)
			log.Printf("DEBUG: Raw tags length: %d", len(post.Tags))

			// Increment view count
			db.Model(&struct{ ID string }{ID: post.ID}).
				Table("blog_posts").
				Update("view_count", gorm.Expr("view_count + 1"))

			// Get comments for this post
			var comments []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Content   string `json:"content"`
				CreatedAt string `json:"created_at"`
			}

			db.Table("blog_comments").
				Select("id, name, content, created_at").
				Where("post_id = ? AND status = ?", post.ID, "approved").
				Order("created_at DESC").
				Find(&comments)

			// Get related posts (same category, exclude current post)
			var relatedPosts []struct {
				ID        string `json:"id"`
				Title     string `json:"title"`
				Slug      string `json:"slug"`
				Excerpt   string `json:"excerpt"`
				Image     string `json:"image"`
				Category  string `json:"category"`
				ReadTime  int    `json:"read_time"`
				CreatedAt string `json:"created_at"`
			}

			db.Table("blog_posts").
				Select("id, title, slug, excerpt, image, category, read_time, created_at").
				Where("category = ? AND slug != ? AND status = ?", post.Category, slug, "published").
				Order("created_at DESC").
				Limit(3).
				Find(&relatedPosts)

			// Parse tags from JSON string
			var tags []string
			if post.Tags != "" {
				// Try to parse as JSON array first
				if err := json.Unmarshal([]byte(post.Tags), &tags); err != nil {
					// If JSON parsing fails, try to handle as comma-separated string
					// or if it's already a single string, wrap it in an array
					if post.Tags[0] == '[' && post.Tags[len(post.Tags)-1] == ']' {
						// It looks like a JSON array but failed to parse, try to clean it up
						cleaned := post.Tags[1 : len(post.Tags)-1] // Remove brackets
						// Split by comma and clean up quotes
						parts := strings.Split(cleaned, ",")
						for _, part := range parts {
							cleanedPart := strings.TrimSpace(part)
							// Remove surrounding quotes if present
							if len(cleanedPart) >= 2 && cleanedPart[0] == '"' && cleanedPart[len(cleanedPart)-1] == '"' {
								cleanedPart = cleanedPart[1 : len(cleanedPart)-1]
							}
							if cleanedPart != "" {
								tags = append(tags, cleanedPart)
							}
						}
					} else {
						// Treat as single string
						tags = []string{post.Tags}
					}
				}
			}

			// Transform the data
			transformedPost := map[string]interface{}{
				"id":          post.ID,
				"title":       post.Title,
				"slug":        post.Slug,
				"excerpt":     post.Excerpt,
				"content":     post.Content,
				"image":       post.Image,
				"category":    post.Category,
				"tags":        tags,
				"readTime":    post.ReadTime,
				"viewCount":   post.ViewCount,
				"publishedAt": post.PublishedAt,
				"date":        post.CreatedAt,
				"author": map[string]interface{}{
					"id":     post.AuthorID,
					"name":   post.AuthorName,
					"title":  post.AuthorTitle,
					"avatar": post.AuthorAvatar,
					"bio":    post.AuthorBio,
				},
				"comments":     comments,
				"relatedPosts": relatedPosts,
			}

			c.JSON(http.StatusOK, transformedPost)
		})

		// Get comments for a blog post
		blog.GET("/posts/:slug/comments", func(c *gin.Context) {
			slug := c.Param("slug")

			// Get the post ID
			var post struct {
				ID string `json:"id"`
			}
			if err := db.Table("blog_posts").Select("id").Where("slug = ?", slug).First(&post).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
				return
			}

			// Apply pagination
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
			offset := (page - 1) * limit

			// Get total count
			var total int64
			if err := db.Table("blog_comments").
				Where("post_id = ? AND status = ?", post.ID, "approved").
				Count(&total).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
				return
			}

			// Get comments for this post with pagination
			var comments []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Content   string `json:"content"`
				Status    string `json:"status"`
				CreatedAt string `json:"created_at"`
			}

			if err := db.Table("blog_comments").
				Select("id, name, content, status, created_at").
				Where("post_id = ? AND status = ?", post.ID, "approved").
				Order("created_at DESC, id DESC").
				Offset(offset).
				Limit(limit).
				Find(&comments).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"comments": comments,
				"total":    total,
				"page":     page,
				"limit":    limit,
				"pages":    int(math.Ceil(float64(total) / float64(limit))),
			})
		})

		// Add a comment to a blog post
		blog.POST("/posts/:slug/comments", func(c *gin.Context) {
			slug := c.Param("slug")

			// Get the post ID
			var post struct {
				ID string `json:"id"`
			}
			if err := db.Table("blog_posts").Select("id").Where("slug = ?", slug).First(&post).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
				return
			}

			var comment struct {
				Name    string `json:"name" binding:"required"`
				Email   string `json:"email" binding:"required,email"`
				Content string `json:"content" binding:"required"`
			}

			if err := c.ShouldBindJSON(&comment); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment data"})
				return
			}

			// Create the comment
			newComment := map[string]interface{}{
				"id":         uuid.New().String(),
				"post_id":    post.ID,
				"name":       comment.Name,
				"email":      comment.Email,
				"content":    comment.Content,
				"status":     "pending", // Comments need approval by default
				"ip_address": c.ClientIP(),
				"user_agent": c.GetHeader("User-Agent"),
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}

			if err := db.Table("blog_comments").Create(newComment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully and pending approval"})
		})

		// Get blog categories
		blog.GET("/categories", func(c *gin.Context) {
			var categories []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Slug        string `json:"slug"`
				Description string `json:"description"`
				Color       string `json:"color"`
			}

			if err := db.Table("blog_categories").Find(&categories).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
				return
			}

			c.JSON(http.StatusOK, categories)
		})

		// Get blog tags
		blog.GET("/tags", func(c *gin.Context) {
			var tags []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			}

			if err := db.Table("blog_tags").Find(&tags).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
				return
			}

			c.JSON(http.StatusOK, tags)
		})

		// Get comment statistics
		blog.GET("/comments/stats", func(c *gin.Context) {
			// Get comment stats
			var totalComments int64
			var pendingComments int64
			var approvedComments int64

			db.Table("blog_comments").Count(&totalComments)
			db.Table("blog_comments").Where("status = ?", "pending").Count(&pendingComments)
			db.Table("blog_comments").Where("status = ?", "approved").Count(&approvedComments)

			c.JSON(http.StatusOK, gin.H{
				"total_comments":    totalComments,
				"pending_comments":  pendingComments,
				"approved_comments": approvedComments,
				"auto_approval": gin.H{
					"enabled": true,
					"delay":   "10m",
				},
			})
		})
	}
}
