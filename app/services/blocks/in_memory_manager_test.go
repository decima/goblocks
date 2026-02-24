package blocks

import (
	"testing"
)

func TestInMemoryBlockManager_SetAndGet(t *testing.T) {
	manager := NewInMemoryBlockManager()

	// Test Set
	content := []byte("Hello, World!")
	err := manager.Set("a/b/c", content, "text/plain")
	if err != nil {
		t.Errorf("Set() error = %v", err)
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

func TestInMemoryBlockManager_GetNotFound(t *testing.T) {
	manager := NewInMemoryBlockManager()

	_, err := manager.Get("nonexistent", true)
	if err != ErrNotFound {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestInMemoryBlockManager_ParentDirectories(t *testing.T) {
	manager := NewInMemoryBlockManager()

	// Set a deep path
	err := manager.Set("a/b/c/d", []byte("test"), "text/plain")
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Check that parent directories were created
	parents := []string{"a", "a/b", "a/b/c"}
	for _, parent := range parents {
		block, err := manager.Get(parent, false)
		if err != nil {
			t.Errorf("Get(%s) error = %v", parent, err)
		}
		if block.Type != "directory" {
			t.Errorf("Get(%s) type = %v, want directory", parent, block.Type)
		}
	}
}

func TestInMemoryBlockManager_List(t *testing.T) {
	manager := NewInMemoryBlockManager()

	// Create some blocks
	manager.Set("a/file1", []byte("content1"), "text/plain")
	manager.Set("a/file2", []byte("content2"), "text/plain")
	manager.Set("a/b/file3", []byte("content3"), "text/plain")
	manager.Set("c/file4", []byte("content4"), "text/plain")

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

func TestInMemoryBlockManager_Delete(t *testing.T) {
	manager := NewInMemoryBlockManager()

	// Create a block
	manager.Set("test/block", []byte("content"), "text/plain")

	// Verify it exists
	_, err := manager.Get("test/block", false)
	if err != nil {
		t.Errorf("Get() before delete error = %v", err)
	}

	// Delete it
	err = manager.Delete("test/block")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify it's gone
	_, err = manager.Get("test/block", false)
	if err != ErrNotFound {
		t.Errorf("Get() after delete error = %v, want ErrNotFound", err)
	}
}

func TestInMemoryBlockManager_UpdateBlock(t *testing.T) {
	manager := NewInMemoryBlockManager()

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
