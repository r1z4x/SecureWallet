import { apiClient } from './auth.js'

export const blogService = {
  // Get all blog posts with optional filters
  async getPosts(params = {}) {
    try {
      const response = await apiClient.get('/blog/posts', { params })
      return response.data
    } catch (error) {
      console.error('Error fetching blog posts:', error)
      throw error
    }
  },

  // Get a single blog post by slug
  async getPost(slug) {
    try {
      const response = await apiClient.get(`/blog/posts/${slug}`)
      return response.data
    } catch (error) {
      console.error('Error fetching blog post:', error)
      throw error
    }
  },

  // Get comments for a blog post
        async getComments(slug, params = {}) {
    try {
      const response = await apiClient.get(`/blog/posts/${slug}/comments`, { params })
      return response.data
    } catch (error) {
      console.error('Error fetching comments:', error)
      throw error
    }
  },

  // Add a comment to a blog post
  async addComment(slug, commentData) {
    try {
      const response = await apiClient.post(`/blog/posts/${slug}/comments`, commentData)
      return response.data
    } catch (error) {
      console.error('Error adding comment:', error)
      throw error
    }
  },

  // Get blog categories
  async getCategories() {
    try {
      const response = await apiClient.get('/blog/categories')
      return response.data
    } catch (error) {
      console.error('Error fetching categories:', error)
      throw error
    }
  },

  // Get blog tags
  async getTags() {
    try {
      const response = await apiClient.get('/blog/tags')
      return response.data
    } catch (error) {
      console.error('Error fetching tags:', error)
      throw error
    }
  }
}

export default blogService
