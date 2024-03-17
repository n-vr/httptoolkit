package problem_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/n-vr/httptoolkit/problem"
)

func TestProblem_newWithDefaults(t *testing.T) {
	problem := problem.New(errors.New("problem error"), 418)

	if problem.Type != "https://httpstatuscodes.org/418" {
		t.Errorf("problem.Type = %s; want https://httpstatuscodes.org/418", problem.Type)
	}

	if problem.Title != "I'm a teapot" {
		t.Errorf("problem.Title = %q; want \"I'm a teapot\"", problem.Title)
	}

	if problem.Status != 418 {
		t.Errorf("problem.Status = %d; want 418", problem.Status)
	}

	if problem.Detail != "problem error" {
		t.Errorf("problem.Detail = %q; want \"problem error\"", problem.Detail)
	}
}

func TestProblem_newWithType(t *testing.T) {
	userNotFoundProblemType := problem.NewType("https://example.com/user-not-found", "User not found")

	problem := problem.New(errors.New("user not found"), 404, problem.WithType(userNotFoundProblemType))

	if problem.Type != "User not found" {
		t.Errorf("problem.Type = %q; want \"User not found\"", problem.Type)
	}

	if problem.Title != "https://example.com/user-not-found" {
		t.Errorf("problem.Title = %q; want \"https://example.com/user-not-found\"", problem.Title)
	}
}

func TestProblem_newWithInstance(t *testing.T) {
	problem := problem.New(errors.New("user not found"), 404, problem.WithInstance("/users/123"))

	if problem.Instance != "/users/123" {
		t.Errorf("problem.Instance = %q; want \"/users/123\"", problem.Instance)
	}
}

func TestProblem_newWithExtension(t *testing.T) {
	p := problem.New(errors.New("user not found"), 404, problem.WithExtension("user", "123"))

	if p.Extensions["user"] != "123" {
		t.Errorf("problem.Extensions[\"user\"] = %q; want \"123\"", p.Extensions["user"])
	}

	recorder := httptest.NewRecorder()
	problem.HTTPErrorHandler(p, recorder)
	response := recorder.Result()

	if response.StatusCode != 404 {
		t.Errorf("response.StatusCode = %d; want 404", response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/problem+json" {
		t.Errorf("response header \"Content-Type\") = %q; want \"application/problem+json\"", response.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var problemResponse map[string]any
	err = json.Unmarshal(body, &problemResponse)
	if err != nil {
		t.Fatal(err)
	}

	if problemResponse["user"] != "123" {
		t.Errorf("problemResponse[\"user\"] = %q; want \"123\"", problemResponse["user"])
	}
}
