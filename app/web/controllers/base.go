package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseController struct {
	pattern string
}

func NewBaseRoute(pattern string) *BaseController {
	return &BaseController{pattern: pattern}
}

func (b *BaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "default handler, don't forget to implement ServeHTTP route", http.StatusNotImplemented)
}

func (b *BaseController) Pattern() string {
	return b.pattern
}

type response struct {
	status  int
	headers map[string]string
	content any
}

func newResponse(content any) *response {
	return &response{
		status:  200,
		headers: map[string]string{},
		content: content,
	}
}

func (resp *response) write(w http.ResponseWriter) error {
	w.WriteHeader(resp.status)
	for k, v := range resp.headers {
		w.Header().Set(k, v)
	}

	_, err := fmt.Fprint(w, resp.content)
	return err
}

type Option func(r *response)

func JsonWrite(w http.ResponseWriter, data any, opts ...Option) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp := newResponse(string(content))
	AsJson(resp)
	for _, opt := range opts {
		opt(resp)
	}

	return resp.write(w)
}
func Error(w http.ResponseWriter, msg string, opts ...Option) {
	JsonWrite(w, H{"error": msg}, opts...)
}

func (b *BaseController) Error(w http.ResponseWriter, msg string, opts ...Option) {
	Error(w, msg, opts...)
}

func (b *BaseController) JSON(w http.ResponseWriter, data any, opts ...Option) {
	err := JsonWrite(w, data, opts...)
	if err != nil {
		Error(w, err.Error())
	}
}

func WithStatus(status int) Option {
	return func(r *response) {
		r.status = status
	}
}

var Ok = WithStatus(http.StatusOK)
var Created = WithStatus(http.StatusCreated)
var NoContent = WithStatus(http.StatusNoContent)
var NotFound = WithStatus(http.StatusNotFound)
var Forbidden = WithStatus(http.StatusForbidden)
var Unauthorized = WithStatus(http.StatusUnauthorized)
var AsJson = WithHeader("Content-Type", "application/json")

func WithHeader(h string, v string) Option {
	return func(r *response) {
		r.headers[h] = v
	}
}

type H map[string]any
