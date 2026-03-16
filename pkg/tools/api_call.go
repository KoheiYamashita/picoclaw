package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KarakuriAgent/clawdroid/pkg/config"
)

// APICallTool exposes a curated set of HTTP endpoints to the LLM.
// Each endpoint is defined in the config (tools.apis), so the LLM can
// only reach pre-approved URLs.  Fixed headers such as auth tokens are
// stored in the config and are invisible to the LLM.
type APICallTool struct {
	endpoints []config.APIEndpointConfig
}

// NewAPICallTool creates an APICallTool from a slice of configured endpoints.
func NewAPICallTool(endpoints []config.APIEndpointConfig) *APICallTool {
	return &APICallTool{endpoints: endpoints}
}

// IsActive implements ActivatableTool: hide this tool when no endpoints
// are configured so it does not clutter the LLM's tool list.
func (t *APICallTool) IsActive() bool {
	return len(t.endpoints) > 0
}

func (t *APICallTool) Name() string {
	return "api_call"
}

func (t *APICallTool) Description() string {
	if len(t.endpoints) == 0 {
		return "Call a pre-configured API endpoint."
	}

	var sb strings.Builder
	sb.WriteString("Call one of the pre-configured API endpoints listed below.\n")
	sb.WriteString("Fixed authentication headers are applied automatically – do not pass them as params.\n\n")
	sb.WriteString("Available endpoints:\n")
	for _, ep := range t.endpoints {
		method := ep.Method
		if method == "" {
			method = "GET"
		}
		sb.WriteString(fmt.Sprintf("- %s [%s %s]: %s\n", ep.Name, strings.ToUpper(method), ep.URL, ep.Description))
	}
	return sb.String()
}

func (t *APICallTool) Parameters() map[string]interface{} {
	names := make([]interface{}, 0, len(t.endpoints))
	for _, ep := range t.endpoints {
		names = append(names, ep.Name)
	}

	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"api_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the API endpoint to call (must be one of the configured endpoints)",
				"enum":        names,
			},
			"params": map[string]interface{}{
				"type":                 "object",
				"description":          "Parameters to pass to the endpoint (query string, path, or request body – depends on endpoint definition)",
				"additionalProperties": true,
			},
		},
		"required": []string{"api_name"},
	}
}

func (t *APICallTool) Execute(ctx context.Context, args map[string]interface{}) *ToolResult {
	apiName, ok := args["api_name"].(string)
	if !ok || apiName == "" {
		return ErrorResult("api_name is required")
	}

	// Find the endpoint in the allowlist.
	ep := t.findEndpoint(apiName)
	if ep == nil {
		return ErrorResult(fmt.Sprintf("unknown api_name %q – must be one of the configured endpoints", apiName))
	}

	// Collect caller-supplied params (may be nil / absent).
	var params map[string]interface{}
	if p, ok := args["params"].(map[string]interface{}); ok {
		params = p
	} else {
		params = map[string]interface{}{}
	}

	// Validate required params.
	for _, pd := range ep.Params {
		if pd.Required {
			if _, exists := params[pd.Name]; !exists {
				return ErrorResult(fmt.Sprintf("missing required parameter %q for endpoint %q", pd.Name, apiName))
			}
		}
	}

	// Build the request URL (expand {name} path placeholders).
	rawURL, err := t.buildURL(ep, params)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to build URL: %v", err))
	}

	// Separate remaining params by "in" location.
	queryParams, bodyParams := t.splitParams(ep, params)

	// Append query parameters.
	if len(queryParams) > 0 {
		qv := url.Values{}
		for k, v := range queryParams {
			qv.Set(k, fmt.Sprintf("%v", v))
		}
		separator := "?"
		if strings.Contains(rawURL, "?") {
			separator = "&"
		}
		rawURL += separator + qv.Encode()
	}

	// Build request body.
	var bodyReader io.Reader
	method := strings.ToUpper(ep.Method)
	if method == "" {
		method = "GET"
	}
	if len(bodyParams) > 0 {
		bodyJSON, err := json.Marshal(bodyParams)
		if err != nil {
			return ErrorResult(fmt.Sprintf("failed to marshal request body: %v", err))
		}
		bodyReader = bytes.NewReader(bodyJSON)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to create request: %v", err))
	}

	// Apply fixed headers from config (these are invisible to the LLM).
	for k, v := range ep.Headers {
		req.Header.Set(k, v)
	}
	if len(bodyParams) > 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set timeout.
	timeout := time.Duration(ep.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return ErrorResult(fmt.Sprintf("request to %q failed: %v", apiName, err))
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to read response from %q: %v", apiName, err))
	}

	// Try to pretty-print JSON responses.
	var prettyBody string
	var jsonData interface{}
	if json.Unmarshal(respBody, &jsonData) == nil {
		formatted, _ := json.MarshalIndent(jsonData, "", "  ")
		prettyBody = string(formatted)
	} else {
		prettyBody = string(respBody)
	}

	summary := fmt.Sprintf("API %q responded with status %d (%d bytes)", apiName, resp.StatusCode, len(respBody))
	result := map[string]interface{}{
		"api_name":    apiName,
		"status_code": resp.StatusCode,
		"body":        prettyBody,
	}
	resultJSON, _ := json.MarshalIndent(result, "", "  ")

	return &ToolResult{
		ForLLM:  fmt.Sprintf("%s\n%s", summary, string(resultJSON)),
		ForUser: string(resultJSON),
		IsError: resp.StatusCode >= 400,
	}
}

// findEndpoint returns the endpoint config whose name matches, or nil.
func (t *APICallTool) findEndpoint(name string) *config.APIEndpointConfig {
	for i := range t.endpoints {
		if t.endpoints[i].Name == name {
			return &t.endpoints[i]
		}
	}
	return nil
}

// buildURL expands {param} placeholders in the URL template with path params,
// returning the URL without query-string parameters.
func (t *APICallTool) buildURL(ep *config.APIEndpointConfig, params map[string]interface{}) (string, error) {
	rawURL := ep.URL
	for _, pd := range ep.Params {
		if pd.In != "path" {
			continue
		}
		val, exists := params[pd.Name]
		if !exists {
			if pd.Required {
				return "", fmt.Errorf("missing required path parameter %q", pd.Name)
			}
			continue
		}
		placeholder := "{" + pd.Name + "}"
		rawURL = strings.ReplaceAll(rawURL, placeholder, url.PathEscape(fmt.Sprintf("%v", val)))
	}

	// Validate the resulting URL is still http/https and has a host.
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL after expansion: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("only http/https endpoints are allowed")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("missing host in URL")
	}

	return rawURL, nil
}

// splitParams separates caller params into query-string and body buckets,
// skipping path params that have already been interpolated.
func (t *APICallTool) splitParams(ep *config.APIEndpointConfig, params map[string]interface{}) (
	queryParams map[string]interface{},
	bodyParams map[string]interface{},
) {
	queryParams = map[string]interface{}{}
	bodyParams = map[string]interface{}{}

	// Build a lookup for the declared "in" location.
	inMap := map[string]string{}
	for _, pd := range ep.Params {
		inMap[pd.Name] = pd.In
	}

	method := strings.ToUpper(ep.Method)
	if method == "" {
		method = "GET"
	}
	supportsBody := method == "POST" || method == "PUT" || method == "PATCH"

	for k, v := range params {
		location := inMap[k] // empty string if not declared

		switch location {
		case "path":
			// Already expanded – skip.
		case "body":
			bodyParams[k] = v
		case "query":
			queryParams[k] = v
		default:
			// Undeclared param: use body for POST/PUT/PATCH, else query.
			if supportsBody {
				bodyParams[k] = v
			} else {
				queryParams[k] = v
			}
		}
	}
	return queryParams, bodyParams
}
