package blocks

import (
	"errors"
	"goblocks/app/config"
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
