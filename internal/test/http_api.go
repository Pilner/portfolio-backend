package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	commons "portfolio-backend/internal/adapters/handlers/http_api/commons"
	authdomain "portfolio-backend/internal/core/domain/auth"
	"testing"
)

type HttpRequest struct {
	T                  *testing.T
	Method             string
	Endpoint           string
	RequestBody        any
	UserCtx            *authdomain.User
	HandlerFunc        func(w http.ResponseWriter, r *http.Request)
	ExpectedStatusCode int
}

func (h *HttpRequest) Execute() *http.Response {
	h.T.Helper()

	w := httptest.NewRecorder()

	// For Request Body
	var body io.Reader
	if h.RequestBody != nil {
		jsonData, err := json.Marshal(h.RequestBody)
		if err != nil {
			h.T.Fatalf("failed to marshal request body: %v", err)
		}
		body = bytes.NewReader(jsonData)
	}
	req, err := http.NewRequest(h.Method, h.Endpoint, body)
	if err != nil {
		h.T.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// For Protected Routes
	if h.UserCtx != nil {
		ctx := context.WithValue(req.Context(), commons.AuthUserCtxKey{}, h.UserCtx)
		req = req.WithContext(ctx)
	}

	h.HandlerFunc(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != h.ExpectedStatusCode {
		h.T.Errorf(string(TestUnexpectedValue), h.ExpectedStatusCode, resp.StatusCode)
	}

	return resp
}
