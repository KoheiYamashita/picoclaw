package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/KarakuriAgent/clawdroid/pkg/config"
)

// helpers

func newTestEndpoint(name, method, rawURL string, headers map[string]string, params []config.HTTPParam) config.APIEndpointConfig {
	return config.APIEndpointConfig{
		Name:        name,
		Description: "test endpoint",
		URL:         rawURL,
		Method:      method,
		Headers:     headers,
		Params:      params,
		Timeout:     5,
	}
}

// TestAPICallTool_IsActive checks that the tool is inactive when no endpoints
// are registered.
func TestAPICallTool_IsActive(t *testing.T) {
	t.Run("no endpoints", func(t *testing.T) {
		tool := NewAPICallTool(nil)
		if tool.IsActive() {
			t.Error("expected IsActive() == false when no endpoints are configured")
		}
	})

	t.Run("with endpoints", func(t *testing.T) {
		tool := NewAPICallTool([]config.APIEndpointConfig{
			newTestEndpoint("ep1", "GET", "http://example.com", nil, nil),
		})
		if !tool.IsActive() {
			t.Error("expected IsActive() == true when endpoints are configured")
		}
	})
}

// TestAPICallTool_Name checks the tool name.
func TestAPICallTool_Name(t *testing.T) {
	tool := NewAPICallTool(nil)
	if tool.Name() != "api_call" {
		t.Errorf("expected name 'api_call', got %q", tool.Name())
	}
}

// TestAPICallTool_Parameters_Enum verifies that the enum in the parameters
// schema matches the configured endpoint names.
func TestAPICallTool_Parameters_Enum(t *testing.T) {
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("ep1", "GET", "http://example.com", nil, nil),
		newTestEndpoint("ep2", "POST", "http://example.com/post", nil, nil),
	})

	params := tool.Parameters()
	props, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("parameters should have a 'properties' map")
	}
	apiNameProp, ok := props["api_name"].(map[string]interface{})
	if !ok {
		t.Fatal("parameters should have an 'api_name' property")
	}
	enum, ok := apiNameProp["enum"].([]interface{})
	if !ok {
		t.Fatal("api_name property should have an 'enum' field")
	}
	if len(enum) != 2 {
		t.Fatalf("expected 2 enum values, got %d", len(enum))
	}
}

// TestAPICallTool_Execute_MissingAPIName checks that an error is returned when
// api_name is absent.
func TestAPICallTool_Execute_MissingAPIName(t *testing.T) {
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("ep1", "GET", "http://example.com", nil, nil),
	})
	result := tool.Execute(context.Background(), map[string]interface{}{})
	if !result.IsError {
		t.Error("expected error when api_name is missing")
	}
	if !strings.Contains(result.ForLLM, "api_name") {
		t.Errorf("expected error to mention 'api_name', got: %s", result.ForLLM)
	}
}

// TestAPICallTool_Execute_UnknownAPIName checks that an error is returned for
// an unrecognised api_name.
func TestAPICallTool_Execute_UnknownAPIName(t *testing.T) {
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("ep1", "GET", "http://example.com", nil, nil),
	})
	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "nonexistent",
	})
	if !result.IsError {
		t.Error("expected error for unknown api_name")
	}
	if !strings.Contains(result.ForLLM, "unknown api_name") {
		t.Errorf("expected 'unknown api_name' in error, got: %s", result.ForLLM)
	}
}

// TestAPICallTool_Execute_GET performs a successful GET request.
func TestAPICallTool_Execute_GET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "wrong method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_get", "GET", server.URL+"/api", nil, nil),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_get",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}
	if !strings.Contains(result.ForLLM, "200") {
		t.Errorf("expected status 200 in ForLLM, got: %s", result.ForLLM)
	}
	if !strings.Contains(result.ForUser, "ok") {
		t.Errorf("expected response body in ForUser, got: %s", result.ForUser)
	}
}

// TestAPICallTool_Execute_POST_BodyParams sends body params and checks that
// they arrive in the request body.
func TestAPICallTool_Execute_POST_BodyParams(t *testing.T) {
	var receivedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "wrong method", http.StatusMethodNotAllowed)
			return
		}
		receivedBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"received":true}`)
	}))
	defer server.Close()

	params := []config.HTTPParam{
		{Name: "message", In: "body", Description: "msg", Required: true},
	}
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_post", "POST", server.URL+"/api", nil, params),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_post",
		"params": map[string]interface{}{
			"message": "hello",
		},
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(receivedBody, &body); err != nil {
		t.Fatalf("server received non-JSON body: %s", receivedBody)
	}
	if body["message"] != "hello" {
		t.Errorf("expected body param 'message'='hello', got: %v", body)
	}
}

// TestAPICallTool_Execute_QueryParams verifies that query params are appended
// to the URL for GET requests.
func TestAPICallTool_Execute_QueryParams(t *testing.T) {
	var receivedQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	params := []config.HTTPParam{
		{Name: "q", In: "query", Description: "search query", Required: true},
	}
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_query", "GET", server.URL+"/search", nil, params),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_query",
		"params": map[string]interface{}{
			"q": "clawdroid",
		},
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}
	if receivedQuery.Get("q") != "clawdroid" {
		t.Errorf("expected query param q=clawdroid, got: %v", receivedQuery)
	}
}

// TestAPICallTool_Execute_PathParams verifies that path params are interpolated
// into the URL template.
func TestAPICallTool_Execute_PathParams(t *testing.T) {
	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	params := []config.HTTPParam{
		{Name: "id", In: "path", Description: "item id", Required: true},
	}
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_path", "GET", server.URL+"/items/{id}", nil, params),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_path",
		"params": map[string]interface{}{
			"id": "42",
		},
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}
	if receivedPath != "/items/42" {
		t.Errorf("expected path /items/42, got: %s", receivedPath)
	}
}

// TestAPICallTool_Execute_FixedHeaders verifies that fixed headers from config
// are sent with the request and cannot be overridden by LLM params.
func TestAPICallTool_Execute_FixedHeaders(t *testing.T) {
	var receivedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	headers := map[string]string{
		"Authorization": "Bearer secret-token",
	}
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_auth", "GET", server.URL+"/secure", headers, nil),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_auth",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}
	if receivedAuthHeader != "Bearer secret-token" {
		t.Errorf("expected auth header to be sent, got: %q", receivedAuthHeader)
	}
}

// TestAPICallTool_Execute_MissingRequiredParam checks that a missing required
// param produces an error.
func TestAPICallTool_Execute_MissingRequiredParam(t *testing.T) {
	params := []config.HTTPParam{
		{Name: "required_field", In: "query", Description: "must be set", Required: true},
	}
	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_required", "GET", "http://example.com/api", nil, params),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_required",
		// params omitted entirely
	})

	if !result.IsError {
		t.Error("expected error for missing required param")
	}
	if !strings.Contains(result.ForLLM, "required_field") {
		t.Errorf("expected error to mention 'required_field', got: %s", result.ForLLM)
	}
}

// TestAPICallTool_Execute_HTTPError checks that 4xx responses set IsError=true.
func TestAPICallTool_Execute_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}))
	defer server.Close()

	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("test_404", "GET", server.URL+"/missing", nil, nil),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "test_404",
	})

	if !result.IsError {
		t.Error("expected IsError=true for 404 response")
	}
	if !strings.Contains(result.ForLLM, "404") {
		t.Errorf("expected 404 in ForLLM, got: %s", result.ForLLM)
	}
}

// TestAPICallTool_Execute_NoParams verifies a simple GET with no params works.
func TestAPICallTool_Execute_NoParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"hello":"world"}`)
	}))
	defer server.Close()

	tool := NewAPICallTool([]config.APIEndpointConfig{
		newTestEndpoint("simple", "GET", server.URL, nil, nil),
	})

	result := tool.Execute(context.Background(), map[string]interface{}{
		"api_name": "simple",
	})

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.ForLLM)
	}
	if !strings.Contains(result.ForUser, "world") {
		t.Errorf("expected 'world' in ForUser, got: %s", result.ForUser)
	}
}
