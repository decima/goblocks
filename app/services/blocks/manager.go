package blocks

import (
	"errors"
	"goblocks/app/config"
	"path/filepath"
	"strings"
)

type Block struct {
	Path     string           `json:"path"`
	Content  []byte           `json:"content,omitempty"`
	Type     string           `json:"type"`
	Children []BlockReference `json:"children"`
	Size     int64            `json:"size,omitempty"`
}

type BlockReference struct {
	Path string `json:"path"`
}

type BlockManager interface {
	List(path string) ([]BlockReference, error)
	Get(path string, withContent bool) (Block, error)
	Set(path string, content []byte, contentType string) error
	Delete(path string) error
}

func NewBlockManager(c *config.Config) BlockManager {
	switch c.Blocks.Storage.Type {
	case config.Fs:
		return BlockManager(NewFsBlockManager(c.Blocks.Storage.Path))
	case config.InMemory:
		return BlockManager(NewInMemoryBlockManager())
	}

	return nil
}

var ErrNotFound = errors.New("Not Found Error")
var ErrUnknown = errors.New("Unknown Error")
var ErrForbidden = errors.New("Forbidden")
var ErrInvalidPath = errors.New("Invalid Path")
var ErrPathTooDeep = errors.New("Path Too Deep")
var ErrInvalidContentType = errors.New("Invalid Content-Type")

const MaxPathDepth = 10

// ValidatePath validates and sanitizes a path to prevent path traversal attacks
func ValidatePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Clean the path to resolve any .. or . elements
	cleaned := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleaned, "..") || strings.HasPrefix(cleaned, "/") {
		return "", ErrInvalidPath
	}

	// Remove leading ./ if present
	cleaned = strings.TrimPrefix(cleaned, "./")

	// Check path depth
	parts := strings.Split(cleaned, string(filepath.Separator))
	if len(parts) > MaxPathDepth {
		return "", ErrPathTooDeep
	}

	// Reject paths with null bytes or other suspicious characters
	if strings.ContainsAny(cleaned, "\x00") {
		return "", ErrInvalidPath
	}

	return cleaned, nil
}

// ValidateContentType validates that a content type follows the MIME type format
func ValidateContentType(contentType string) error {
	if contentType == "" {
		return ErrInvalidContentType
	}

	// Basic MIME type format: type/subtype
	// May include parameters like: type/subtype; charset=utf-8
	parts := strings.Split(contentType, ";")
	mainType := strings.TrimSpace(parts[0])

	// Check format: must have type/subtype
	typeParts := strings.Split(mainType, "/")
	if len(typeParts) != 2 {
		return ErrInvalidContentType
	}

	// Both type and subtype must be non-empty and contain valid characters
	if typeParts[0] == "" || typeParts[1] == "" {
		return ErrInvalidContentType
	}

	// Reject content types with control characters or null bytes
	if strings.ContainsAny(contentType, "\x00\r\n") {
		return ErrInvalidContentType
	}

	return nil
}
