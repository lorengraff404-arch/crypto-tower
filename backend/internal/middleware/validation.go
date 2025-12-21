package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// InputValidator provides sanitization and validation
type InputValidator struct{}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// SanitizeString removes potentially dangerous characters
func (v *InputValidator) SanitizeString(input string) string {
	// Remove HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	sanitized := htmlTagRegex.ReplaceAllString(input, "")

	// Remove script tags specifically (double check)
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized = scriptRegex.ReplaceAllString(sanitized, "")

	// Remove potentially dangerous characters
	sanitized = strings.ReplaceAll(sanitized, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#x27;")
	sanitized = strings.ReplaceAll(sanitized, "/", "&#x2F;")

	return sanitized
}

// ValidateStringLength checks if string is within limits
func (v *InputValidator) ValidateStringLength(input string, min, max int) bool {
	length := len(input)
	return length >= min && length <= max
}

// ValidateEthereumAddress checks if valid Ethereum address format
func (v *InputValidator) ValidateEthereumAddress(address string) bool {
	// Ethereum addresses are 42 characters (0x + 40 hex chars)
	if len(address) != 42 {
		return false
	}

	if !strings.HasPrefix(address, "0x") {
		return false
	}

	// Check if remaining 40 chars are valid hex
	hexRegex := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	return hexRegex.MatchString(address)
}

// ValidateInteger checks if integer is within range
func (v *InputValidator) ValidateInteger(value, min, max int) bool {
	return value >= min && value <= max
}

// ValidatePositiveInteger checks if value is positive
func (v *InputValidator) ValidatePositiveInteger(value int) bool {
	return value > 0
}

// SanitizationMiddleware applies input sanitization to request body
func SanitizationMiddleware() gin.HandlerFunc {
	validator := NewInputValidator()

	return func(c *gin.Context) {
		// Get content type
		contentType := c.ContentType()

		// Only sanitize JSON requests
		if strings.Contains(contentType, "application/json") {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err == nil {
				// Sanitize string fields
				sanitized := sanitizeMap(body, validator)
				c.Set("sanitized_body", sanitized)
			}
		}

		c.Next()
	}
}

// sanitizeMap recursively sanitizes map values
func sanitizeMap(data map[string]interface{}, validator *InputValidator) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			sanitized[key] = validator.SanitizeString(v)
		case map[string]interface{}:
			sanitized[key] = sanitizeMap(v, validator)
		case []interface{}:
			sanitized[key] = sanitizeSlice(v, validator)
		default:
			sanitized[key] = value
		}
	}

	return sanitized
}

// sanitizeSlice recursively sanitizes slice values
func sanitizeSlice(data []interface{}, validator *InputValidator) []interface{} {
	sanitized := make([]interface{}, len(data))

	for i, value := range data {
		switch v := value.(type) {
		case string:
			sanitized[i] = validator.SanitizeString(v)
		case map[string]interface{}:
			sanitized[i] = sanitizeMap(v, validator)
		case []interface{}:
			sanitized[i] = sanitizeSlice(v, validator)
		default:
			sanitized[i] = value
		}
	}

	return sanitized
}

// MaxBodySizeMiddleware limits request body size
func MaxBodySizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// ValidateJSONMiddleware ensures valid JSON in request
func ValidateJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.ContentType()
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
