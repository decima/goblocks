package controllers

import (
	"bytes"
	"goblocks/app/config"
	"goblocks/app/services/blocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBlockController(t *testing.T) {
	// Setup
	manager := blocks.NewInMemoryBlockManager()
	manager.Set("test/block", []byte("Hello, World!"), "text/plain")

	controller := NewGetBlockController(manager)

	tests := []struct {
		name           string
		path           string
		raw            bool
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "get existing block metadata",
			path:           "test/block",
			raw:            false,
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
		},
		{
			name:           "get existing block raw",
			path:           "test/block",
			raw:            true,
			expectedStatus: http.StatusOK,
			expectedType:   "text/plain",
		},
		{
			name:           "get nonexistent block",
			path:           "nonexistent",
			raw:            false,
			expectedStatus: http.StatusNotFound,
			expectedType:   "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/blocks/" + tt.path
			if tt.raw {
				url += "?raw"
			}

			req := httptest.NewRequest("GET", url, nil)
			req.SetPathValue("path", tt.path)

			w := httptest.NewRecorder()
			controller.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != tt.expectedType {
				t.Errorf("Content-Type = %s, want %s", contentType, tt.expectedType)
			}

			if tt.raw && tt.expectedStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				if string(body) != "Hello, World!" {
					t.Errorf("Body = %s, want Hello, World!", string(body))
				}
			}
		})
	}
}

func TestGetBlockController_PathValidation(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	controller := NewGetBlockController(manager)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "path traversal attempt",
			path:           "../../../etc/passwd",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "absolute path",
			path:           "/etc/passwd",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "path too deep",
			path:           "a/b/c/d/e/f/g/h/i/j/k",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/blocks/"+tt.path, nil)
			req.SetPathValue("path", tt.path)

			w := httptest.NewRecorder()
			controller.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

func TestWriteBlockController(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	cfg := &config.Config{
		Http: config.Http{
			MaxUploadSize: 10 * 1024 * 1024, // 10MB
		},
	}
	controller := NewWriteBlockController(manager, cfg)

	tests := []struct {
		name           string
		path           string
		content        string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "create new block",
			path:           "new/block",
			content:        "Hello, World!",
			contentType:    "text/plain",
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "create with default content type",
			path:           "default/type",
			content:        "binary data",
			contentType:    "",
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "invalid content type",
			path:           "invalid/type",
			content:        "content",
			contentType:    "invalid",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/blocks/"+tt.path, bytes.NewBufferString(tt.content))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			req.SetPathValue("path", tt.path)

			w := httptest.NewRecorder()
			controller.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}

			// If successful, verify the block was created
			if tt.expectedStatus == http.StatusAccepted {
				block, err := manager.Get(tt.path, true)
				if err != nil {
					t.Errorf("Failed to get created block: %v", err)
				}
				if string(block.Content) != tt.content {
					t.Errorf("Block content = %s, want %s", string(block.Content), tt.content)
				}
			}
		})
	}
}

func TestWriteBlockController_SizeLimit(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	cfg := &config.Config{
		Http: config.Http{
			MaxUploadSize: 10, // Only 10 bytes
		},
	}
	controller := NewWriteBlockController(manager, cfg)

	// Try to upload more than the limit
	largeContent := bytes.Repeat([]byte("a"), 20)
	req := httptest.NewRequest("PUT", "/blocks/large", bytes.NewReader(largeContent))
	req.Header.Set("Content-Type", "text/plain")
	req.SetPathValue("path", "large")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Should fail with unprocessable entity (the error comes from io.ReadAll)
	if resp.StatusCode == http.StatusAccepted {
		t.Errorf("Large upload should fail, got status %d", resp.StatusCode)
	}
}

func TestDeleteBlockController(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	manager.Set("test/block", []byte("content"), "text/plain")

	controller := NewDeleteBlockController(manager)

	// Delete the block
	req := httptest.NewRequest("DELETE", "/blocks/test/block", nil)
	req.SetPathValue("path", "test/block")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	// Verify block is deleted
	_, err := manager.Get("test/block", false)
	if err != blocks.ErrNotFound {
		t.Errorf("Block should be deleted, got error: %v", err)
	}
}

func TestDeleteBlockController_NotFound(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	controller := NewDeleteBlockController(manager)

	req := httptest.NewRequest("DELETE", "/blocks/nonexistent", nil)
	req.SetPathValue("path", "nonexistent")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Delete of nonexistent should succeed (idempotent)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestDeleteBlockController_PathValidation(t *testing.T) {
	manager := blocks.NewInMemoryBlockManager()
	controller := NewDeleteBlockController(manager)

	req := httptest.NewRequest("DELETE", "/blocks/../../../etc/passwd", nil)
	req.SetPathValue("path", "../../../etc/passwd")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Status = %d, want %d (Forbidden)", resp.StatusCode, http.StatusForbidden)
	}
}
