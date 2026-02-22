package blocks

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

const (
	FsFileName = ".content"
)

type FsBlockManager struct {
	baseDir string
}

func NewFsBlockManager(baseDir string) *FsBlockManager {
	return &FsBlockManager{baseDir: baseDir}
}

func (f *FsBlockManager) Get(path string, withContent bool) (Block, error) {
	file, err := os.Open(f.getAbsoluteFilePath(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Block{}, errors.Join(err, ErrNotFound)
		}
		if errors.Is(err, os.ErrPermission) {
			return Block{}, errors.Join(err, ErrForbidden)
		}

		return Block{}, errors.Join(err, ErrUnknown)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return Block{}, errors.Join(err, ErrUnknown)
	}
	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil {
		return Block{}, errors.Join(err, ErrForbidden)
	}

	var fileContent = FileContent{}
	err = json.Unmarshal(content, &fileContent)
	if err != nil {
		return Block{}, errors.Join(err, ErrUnknown)
	}

	block := Block{
		Path: path,
		Type: fileContent.ContentType,
		Size: fileContent.Size,
	}
	if withContent {
		block.Content = fileContent.Content
	}

	return block, nil
}

func (f *FsBlockManager) Set(path string, content []byte, contentType string) error {
	//refFormat := "::ref(/a/b/c)"
	err := os.MkdirAll(f.getAbsolutePath(path), 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		if errors.Is(err, os.ErrPermission) {
			return errors.Join(err, ErrForbidden)
		}
		return errors.Join(err, ErrUnknown)
	}
	augmentedContent := FileContent{
		Content:     content,
		ContentType: contentType,
		Size:        int64(len(content)),
	}

	jsonContent, err := json.Marshal(augmentedContent)
	if err != nil {
		return errors.Join(err, ErrUnknown)
	}
	err = os.WriteFile(f.getAbsoluteFilePath(path), jsonContent, 0755)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return errors.Join(err, ErrForbidden)
		}
		return errors.Join(err, ErrUnknown)
	}

	return nil
}

func (f *FsBlockManager) Delete(path string) error {
	err := os.RemoveAll(f.getAbsolutePath(path))
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return errors.Join(err, ErrForbidden)
		}
		return errors.Join(err, ErrUnknown)
	}
	return nil

}

func (f *FsBlockManager) List(path string) ([]BlockReference, error) {
	entries, err := os.ReadDir(f.getAbsolutePath(path))
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, errors.Join(err, ErrForbidden)
		}
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.Join(err, ErrNotFound)
		}
		return nil, errors.Join(err, ErrUnknown)
	}
	references := []BlockReference{}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if !e.IsDir() {
			continue
		}

		references = append(references, BlockReference{Path: path + "/" + e.Name()})

	}
	return references, nil
}

func (f *FsBlockManager) getAbsolutePath(path string) string {
	return f.baseDir + path
}
func (f *FsBlockManager) getAbsoluteFilePath(path string) string {
	return f.baseDir + path + "/" + FsFileName
}

type FileContent struct {
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}
