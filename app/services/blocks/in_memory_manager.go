package blocks

import (
	"path"
	"sync"
)

type InMemoryBlockManager struct {
	blocks sync.Map
}

func (i *InMemoryBlockManager) List(p string) ([]BlockReference, error) {
	response := []BlockReference{}
	if p == "" {
		p = "."
	}
	i.blocks.Range(func(k, v any) bool {
		if p != path.Dir(k.(string)) {
			return true
		}
		response = append(response, BlockReference{Path: k.(string)})
		return true

	})

	return response, nil
}

func (i *InMemoryBlockManager) Get(path string, withContent bool) (Block, error) {
	item, ok := i.blocks.Load(path)
	if !ok {
		return Block{}, ErrNotFound
	}
	if ite, ok := item.(Block); ok {
		if withContent == false {
			ite.Content = nil
		}
		return ite, nil
	}
	return Block{}, ErrNotFound
}

func (i *InMemoryBlockManager) Set(p string, content []byte, contentType string) error {
	i.blocks.Store(p, Block{
		Path:    p,
		Content: content,
		Type:    contentType,
		Size:    int64(len(content)),
	})
	for {
		p = path.Dir(p)
		if p == "." {
			break
		}
		i.blocks.Store(p, nil)
	}

	return nil
}

func (i *InMemoryBlockManager) Delete(path string) error {
	i.blocks.Delete(path)
	return nil
}

func NewInMemoryBlockManager() *InMemoryBlockManager {
	t := &InMemoryBlockManager{
		blocks: sync.Map{},
	}

	t.blocks.Store("a/b/c", Block{Path: "a/b/c"})

	return t
}
