package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// InputValidationMiddleware validates and sanitizes input
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if !isValidInput(value) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Invalid input detected",
					})
					c.Abort()
					return
				}
			}
		}

		// Validate path parameters
		for _, param := range c.Params {
			if !isValidInput(param.Value) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid path parameter",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// isValidInput checks if input contains potentially dangerous content
func isValidInput(input string) bool {
	if len(input) > 1000 {
		return false
	}

	// Check for SQL injection patterns
	sqlPatterns := []string{
		"';", "--", "/*", "*/", "xp_", "sp_", "exec", "execute",
		"union", "select", "insert", "update", "delete", "drop", "create",
	}

	for _, pattern := range sqlPatterns {
		if strings.Contains(strings.ToLower(input), pattern) {
			return false
		}
	}

	// Check for XSS patterns
	xssPatterns := []string{
		"<script", "javascript:", "onload=", "onerror=", "onclick=",
		"<iframe", "<object", "<embed", "<form",
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(strings.ToLower(input), pattern) {
			return false
		}
	}

	// Check for command injection patterns
	cmdPatterns := []string{
		";", "|", "&", "`", "$(", "&&", "||",
	}

	for _, pattern := range cmdPatterns {
		if strings.Contains(input, pattern) {
			return false
		}
	}

	return true
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")

	// Limit length
	if len(input) > 1000 {
		input = input[:1000]
	}

	return input
}
