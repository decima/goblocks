package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"goblocks/app/services/blocks"
	"io"
	"net/http"
)

type GetBlockController struct {
	*BaseController
	blockManager blocks.BlockManager
}

func NewGetBlockController(blockManager blocks.BlockManager) *GetBlockController {
	return &GetBlockController{
		NewBaseRoute("GET /blocks/{path...}"),
		blockManager,
	}
}

func (c *GetBlockController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	raw := r.URL.Query().Has("raw") || r.URL.Query().Get("format") == "raw"
	block, err := c.blockManager.Get(path, raw)
	status := Ok
	if err != nil {
		status = blockErrorToStatus(err)
	}

	if raw {
		if err != nil {
			c.Error(w, err.Error(), status)
			return
		}
		w.Header().Set("Content-Type", block.Type)
		w.WriteHeader(http.StatusOK)
		io.Copy(w, bytes.NewBuffer(block.Content))
		return
	}

	block.Children, err = c.blockManager.List(path)
	c.JSON(w, block, status)

}

type WriteBlockController struct {
	*BaseController
	blockManager blocks.BlockManager
}

func NewWriteBlockController(blockManager blocks.BlockManager) *WriteBlockController {
	return &WriteBlockController{
		NewBaseRoute("PUT /blocks/{path...}"),
		blockManager,
	}
}

func (c *WriteBlockController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	defer r.Body.Close()
	content, err := io.ReadAll(r.Body)
	if err != nil {
		c.Error(w, err.Error(), blockErrorToStatus(err))
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	err = c.blockManager.Set(path, content, contentType)
	if err != nil {
		c.Error(w, err.Error(), blockErrorToStatus(err))
		return
	}

	block, err := c.blockManager.Get(path, false)
	if err != nil {
		c.Error(w, err.Error(), blockErrorToStatus(err))
		return
	}
	c.JSON(w, block, Accepted)
}

type DeleteBlockController struct {
	*BaseController
	blockManager blocks.BlockManager
}

func NewDeleteBlockController(blockManager blocks.BlockManager) *DeleteBlockController {
	return &DeleteBlockController{
		NewBaseRoute("DELETE /blocks/{path...}"),
		blockManager,
	}
}

func (c *DeleteBlockController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	err := c.blockManager.Delete(path)
	if err != nil {
		c.Error(w, err.Error(), blockErrorToStatus(err))
		return
	}
	c.JSON(w, nil, NoContent)

}

func blockErrorToStatus(err error) Option {
	if errors.Is(err, blocks.ErrNotFound) {
		return NotFound
	} else if errors.Is(err, blocks.ErrForbidden) {
		return Forbidden
	} else {
		fmt.Println(err.Error())
		return Unprocessable
	}
}
