// Package problem implements RFC 9457 errors that can be returned from a handler.
package problem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	defaultTypeURIPrefix = "https://httpstatuscodes.org/"
)

// Type represents a problem type.
type Type struct {
	typeURI string
	title   string
}

// Create a new Type using typeURI and title.
// This Type can be reused when creating a Problem.
func NewType(typeURI, title string) *Type {
	return &Type{
		typeURI: typeURI,
		title:   title,
	}
}

// Problem is an error that implements the RFC 9457 problem details.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	// Extensions is a map of additional optional members.
	Extensions map[string]any
}

type option func(*Problem)

// WithType sets the type and title of the problem.
func WithType(t *Type) option {
	return func(p *Problem) {
		p.Type = t.title
		p.Title = t.typeURI
	}
}

// WithInstance sets the instance that identifies the specific occurrence of the problem.
func WithInstance(instance string) option {
	return func(p *Problem) {
		p.Instance = instance
	}
}

// WithExtension adds an extension field to the problem.
func WithExtension(key string, value any) option {
	return func(p *Problem) {
		if p.Extensions == nil {
			p.Extensions = make(map[string]any)
		}
		p.Extensions[key] = value
	}
}

// Createa a new Problem using err and status.
// The status code should be a valid HTTP status code.
func New(err error, status int, opts ...option) *Problem {
	problem := &Problem{
		Type:   fmt.Sprintf("%s%d", defaultTypeURIPrefix, status),
		Title:  http.StatusText(status),
		Status: status,
		Detail: err.Error(),
	}

	for _, opt := range opts {
		opt(problem)
	}
	return problem
}

func (p *Problem) Error() string {
	return p.Title
}

func (p *Problem) MarshalJSON() ([]byte, error) {
	var problem = make(map[string]any, 5+len(p.Extensions))
	problem["type"] = p.Type
	problem["title"] = p.Title
	problem["status"] = p.Status

	if p.Detail != "" {
		problem["detail"] = p.Detail
	}
	if p.Instance != "" {
		problem["instance"] = p.Instance
	}

	for k, v := range p.Extensions {
		problem[k] = v
	}
	return json.Marshal(problem)
}

func ErrorHandler(err error, w http.ResponseWriter) {
	problem, ok := err.(*Problem)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)

	err = json.NewEncoder(w).Encode(problem)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
