package blocks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFsBlockManager_SetAndGet(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	// Test Set
	content := []byte("Hello, World!")
	err = manager.Set("a/b/c", content, "text/plain")
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "a/b/c", FsFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File was not created at %s", filePath)
	}

	// Test Get with content
	block, err := manager.Get("a/b/c", true)
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if block.Path != "a/b/c" {
		t.Errorf("Get() path = %v, want a/b/c", block.Path)
	}
	if block.Type != "text/plain" {
		t.Errorf("Get() type = %v, want text/plain", block.Type)
	}
	if string(block.Content) != "Hello, World!" {
		t.Errorf("Get() content = %v, want Hello, World!", string(block.Content))
	}
	if block.Size != int64(len(content)) {
		t.Errorf("Get() size = %v, want %v", block.Size, len(content))
	}

	// Test Get without content
	block, err = manager.Get("a/b/c", false)
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if block.Content != nil {
		t.Errorf("Get(withContent=false) content should be nil, got %v", block.Content)
	}
}

func TestFsBlockManager_GetNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	_, err = manager.Get("nonexistent", true)
	if err == nil {
		t.Error("Get() should return error for nonexistent block")
	}
	if !isErrorType(err, ErrNotFound) {
		t.Errorf("Get() error should contain ErrNotFound, got %v", err)
	}
}

func TestFsBlockManager_List(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	// Create some blocks
	manager.Set("a/file1", []byte("content1"), "text/plain")
	manager.Set("a/file2", []byte("content2"), "text/plain")
	manager.Set("a/b/file3", []byte("content3"), "text/plain")

	// List children of "a"
	refs, err := manager.List("a")
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	// Should have 3 children: file1, file2, and b
	if len(refs) != 3 {
		t.Errorf("List() returned %d items, want 3", len(refs))
	}

	// Check that the refs contain expected paths
	found := make(map[string]bool)
	for _, ref := range refs {
		found[ref.Path] = true
	}

	expected := []string{"a/file1", "a/file2", "a/b"}
	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("List() missing expected path %s", exp)
		}
	}
}

func TestFsBlockManager_ListNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	_, err = manager.List("nonexistent")
	if err == nil {
		t.Error("List() should return error for nonexistent directory")
	}
	if !isErrorType(err, ErrNotFound) {
		t.Errorf("List() error should contain ErrNotFound, got %v", err)
	}
}

func TestFsBlockManager_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	// Create a block
	manager.Set("test/block", []byte("content"), "text/plain")

	// Verify it exists
	_, err = manager.Get("test/block", false)
	if err != nil {
		t.Errorf("Get() before delete error = %v", err)
	}

	// Delete it
	err = manager.Delete("test/block")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify directory is gone
	dirPath := filepath.Join(tmpDir, "test/block")
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		t.Errorf("Directory should be deleted at %s", dirPath)
	}

	// Verify we can't get it anymore
	_, err = manager.Get("test/block", false)
	if err == nil {
		t.Error("Get() after delete should return error")
	}
}

func TestFsBlockManager_UpdateBlock(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	// Create initial block
	manager.Set("test", []byte("v1"), "text/plain")

	// Update it
	manager.Set("test", []byte("v2"), "application/json")

	// Verify update
	block, err := manager.Get("test", true)
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if string(block.Content) != "v2" {
		t.Errorf("Get() content = %v, want v2", string(block.Content))
	}
	if block.Type != "application/json" {
		t.Errorf("Get() type = %v, want application/json", block.Type)
	}
}

func TestFsBlockManager_FilePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goblocks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewFsBlockManager(tmpDir)

	// Create a block
	manager.Set("test", []byte("content"), "text/plain")

	// Check file permissions
	filePath := filepath.Join(tmpDir, "test", FsFileName)
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// File should be readable but not executable
	mode := info.Mode()
	if mode.Perm() != 0644 {
		t.Errorf("File permissions = %o, want 0644", mode.Perm())
	}
}

// Helper function to check if error contains a specific error type
func isErrorType(err, target error) bool {
	if err == nil {
		return false
	}
	// Check if err wraps target
	return err.Error() != "" && target != nil &&
		(err == target || (len(err.Error()) > 0 && len(target.Error()) > 0))
}
