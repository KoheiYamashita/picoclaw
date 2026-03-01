package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestDefaultConfig_HeartbeatEnabled verifies heartbeat is enabled by default
func TestDefaultConfig_HeartbeatEnabled(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.Heartbeat.Enabled {
		t.Error("Heartbeat should be enabled by default")
	}
}

// TestDefaultConfig_WorkspacePath verifies workspace path is correctly set
func TestDefaultConfig_WorkspacePath(t *testing.T) {
	cfg := DefaultConfig()

	// Just verify the workspace is set, don't compare exact paths
	// since expandHome behavior may differ based on environment
	if cfg.Agents.Defaults.Workspace == "" {
		t.Error("Workspace should not be empty")
	}
}

// TestDefaultConfig_Model verifies model default is empty (user must configure)
func TestDefaultConfig_Model(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LLM.Model != "" {
		t.Errorf("LLM.Model should be empty by default, got %q", cfg.LLM.Model)
	}
}

// TestDefaultConfig_MaxTokens verifies max tokens has default value
func TestDefaultConfig_MaxTokens(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agents.Defaults.MaxTokens == 0 {
		t.Error("MaxTokens should not be zero")
	}
}

// TestDefaultConfig_MaxToolIterations verifies max tool iterations has default value
func TestDefaultConfig_MaxToolIterations(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agents.Defaults.MaxToolIterations == 0 {
		t.Error("MaxToolIterations should not be zero")
	}
}

// TestDefaultConfig_Temperature verifies temperature has expected default value (0 = deterministic)
func TestDefaultConfig_Temperature(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agents.Defaults.Temperature != 0 {
		t.Errorf("Temperature should be 0 by default, got %v", cfg.Agents.Defaults.Temperature)
	}
}

// TestDefaultConfig_Gateway verifies gateway defaults
func TestDefaultConfig_Gateway(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Gateway.Port == 0 {
		t.Error("Gateway port should have default value")
	}
}

// TestDefaultConfig_LLM verifies LLM config defaults
func TestDefaultConfig_LLM(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LLM.APIKey != "" {
		t.Error("LLM API key should be empty by default")
	}
	if cfg.LLM.BaseURL != "" {
		t.Error("LLM BaseURL should be empty by default")
	}
	if cfg.LLM.Model != "" {
		t.Errorf("LLM Model should be empty by default, got %q", cfg.LLM.Model)
	}
}

// TestDefaultConfig_Channels verifies channels are disabled by default
func TestDefaultConfig_Channels(t *testing.T) {
	cfg := DefaultConfig()

	// Verify all channels are disabled by default
	if cfg.Channels.WhatsApp.Enabled {
		t.Error("WhatsApp should be disabled by default")
	}
	if cfg.Channels.Telegram.Enabled {
		t.Error("Telegram should be disabled by default")
	}
	if cfg.Channels.Discord.Enabled {
		t.Error("Discord should be disabled by default")
	}
	if cfg.Channels.Slack.Enabled {
		t.Error("Slack should be disabled by default")
	}
}

// TestDefaultConfig_ExecToolDisabled verifies exec tool is disabled by default
func TestDefaultConfig_ExecToolDisabled(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Tools.Exec.Enabled {
		t.Error("Exec tool should be disabled by default")
	}
}

// TestDefaultConfig_WebTools verifies web tools config
func TestDefaultConfig_WebTools(t *testing.T) {
	cfg := DefaultConfig()

	// Verify web tools defaults
	if cfg.Tools.Web.Brave.MaxResults != 5 {
		t.Error("Expected Brave MaxResults 5, got ", cfg.Tools.Web.Brave.MaxResults)
	}
	if cfg.Tools.Web.Brave.APIKey != "" {
		t.Error("Brave API key should be empty by default")
	}
	if cfg.Tools.Web.DuckDuckGo.MaxResults != 5 {
		t.Error("Expected DuckDuckGo MaxResults 5, got ", cfg.Tools.Web.DuckDuckGo.MaxResults)
	}
}

func TestSaveConfig_FilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission bits are not enforced on Windows")
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	cfg := DefaultConfig()
	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("config file has permission %04o, want 0600", perm)
	}
}

// TestDefaultConfig_DataDir verifies data dir default value
func TestDefaultConfig_DataDir(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agents.Defaults.DataDir == "" {
		t.Error("DataDir should not be empty")
	}
	if cfg.Agents.Defaults.DataDir != "~/.clawdroid/data" {
		t.Errorf("DataDir should be '~/.clawdroid/data', got '%s'", cfg.Agents.Defaults.DataDir)
	}
}

// TestConfig_DataPath verifies DataPath expands home directory
func TestConfig_DataPath(t *testing.T) {
	cfg := DefaultConfig()

	path := cfg.DataPath()
	if path == "" {
		t.Error("DataPath should not be empty")
	}
	if path == "~/.clawdroid/data" {
		t.Error("DataPath should expand ~ to home directory")
	}
	if path[0] == '~' {
		t.Error("DataPath should not start with ~")
	}
}

// TestDefaultConfig_GatewayAPIKey verifies Gateway APIKey is empty by default
func TestDefaultConfig_GatewayAPIKey(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Gateway.APIKey != "" {
		t.Error("Gateway APIKey should be empty by default")
	}
}

// TestConfig_Complete verifies all config fields are set
func TestConfig_Complete(t *testing.T) {
	cfg := DefaultConfig()

	// Verify complete config structure
	if cfg.Agents.Defaults.Workspace == "" {
		t.Error("Workspace should not be empty")
	}
	if cfg.Agents.Defaults.MaxTokens == 0 {
		t.Error("MaxTokens should not be zero")
	}
	if cfg.Agents.Defaults.MaxToolIterations == 0 {
		t.Error("MaxToolIterations should not be zero")
	}
	if cfg.Gateway.Port == 0 {
		t.Error("Gateway port should have default value")
	}
	if !cfg.Heartbeat.Enabled {
		t.Error("Heartbeat should be enabled by default")
	}
}

// --- FlexibleStringSlice.UnmarshalJSON ---

func TestFlexibleStringSlice_StringArray(t *testing.T) {
	var f FlexibleStringSlice
	if err := json.Unmarshal([]byte(`["a","b","c"]`), &f); err != nil {
		t.Fatal(err)
	}
	if len(f) != 3 || f[0] != "a" || f[1] != "b" || f[2] != "c" {
		t.Errorf("got %v", f)
	}
}

func TestFlexibleStringSlice_NumberArray(t *testing.T) {
	var f FlexibleStringSlice
	if err := json.Unmarshal([]byte(`[123, 456]`), &f); err != nil {
		t.Fatal(err)
	}
	if len(f) != 2 || f[0] != "123" || f[1] != "456" {
		t.Errorf("got %v", f)
	}
}

func TestFlexibleStringSlice_MixedArray(t *testing.T) {
	var f FlexibleStringSlice
	if err := json.Unmarshal([]byte(`["user1", 12345, "user2"]`), &f); err != nil {
		t.Fatal(err)
	}
	if len(f) != 3 || f[0] != "user1" || f[1] != "12345" || f[2] != "user2" {
		t.Errorf("got %v", f)
	}
}

func TestFlexibleStringSlice_EmptyArray(t *testing.T) {
	var f FlexibleStringSlice
	if err := json.Unmarshal([]byte(`[]`), &f); err != nil {
		t.Fatal(err)
	}
	if len(f) != 0 {
		t.Errorf("expected empty, got %v", f)
	}
}

func TestFlexibleStringSlice_InvalidJSON(t *testing.T) {
	var f FlexibleStringSlice
	err := json.Unmarshal([]byte(`"not an array"`), &f)
	if err == nil {
		t.Error("expected error for non-array JSON")
	}
}

// --- expandHome ---

func TestExpandHome_WithTilde(t *testing.T) {
	result := expandHome("~/Documents")
	if strings.HasPrefix(result, "~") {
		t.Errorf("~ should be expanded, got %q", result)
	}
	if !strings.HasSuffix(result, "/Documents") {
		t.Errorf("expected path to end with /Documents, got %q", result)
	}
}

func TestExpandHome_WithoutTilde(t *testing.T) {
	result := expandHome("/absolute/path")
	if result != "/absolute/path" {
		t.Errorf("expected unchanged path, got %q", result)
	}
}

func TestExpandHome_Empty(t *testing.T) {
	result := expandHome("")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestExpandHome_TildeOnly(t *testing.T) {
	result := expandHome("~")
	if result == "~" {
		t.Error("~ should be expanded to home directory")
	}
	if result == "" {
		t.Error("result should not be empty")
	}
}

// --- CopyFrom ---

func TestCopyFrom(t *testing.T) {
	src := DefaultConfig()
	src.LLM.Model = "test-model"
	src.LLM.APIKey = "test-key"
	src.Gateway.Port = 9999

	dst := &Config{}
	dst.CopyFrom(src)

	if dst.LLM.Model != "test-model" {
		t.Errorf("LLM.Model = %q, want %q", dst.LLM.Model, "test-model")
	}
	if dst.LLM.APIKey != "test-key" {
		t.Errorf("LLM.APIKey = %q, want %q", dst.LLM.APIKey, "test-key")
	}
	if dst.Gateway.Port != 9999 {
		t.Errorf("Gateway.Port = %d, want %d", dst.Gateway.Port, 9999)
	}
	if dst.Heartbeat.Enabled != src.Heartbeat.Enabled {
		t.Error("Heartbeat.Enabled not copied")
	}
}

// --- LoadConfig ---

func TestLoadConfig_NonExistentFile(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should return default config - verify multiple fields
	defaults := DefaultConfig()
	if cfg.Heartbeat.Interval != defaults.Heartbeat.Interval {
		t.Errorf("Heartbeat.Interval = %d, want %d", cfg.Heartbeat.Interval, defaults.Heartbeat.Interval)
	}
	if cfg.Heartbeat.Enabled != defaults.Heartbeat.Enabled {
		t.Errorf("Heartbeat.Enabled = %v, want %v", cfg.Heartbeat.Enabled, defaults.Heartbeat.Enabled)
	}
	if cfg.Gateway.Port != defaults.Gateway.Port {
		t.Errorf("Gateway.Port = %d, want %d", cfg.Gateway.Port, defaults.Gateway.Port)
	}
	if cfg.Agents.Defaults.MaxToolIterations != defaults.Agents.Defaults.MaxToolIterations {
		t.Errorf("MaxToolIterations = %d, want %d", cfg.Agents.Defaults.MaxToolIterations, defaults.Agents.Defaults.MaxToolIterations)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	data := `{"llm":{"model":"gpt-4"},"gateway":{"port":5555}}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LLM.Model != "gpt-4" {
		t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, "gpt-4")
	}
	if cfg.Gateway.Port != 5555 {
		t.Errorf("Gateway.Port = %d, want %d", cfg.Gateway.Port, 5555)
	}
}

// --- WorkspacePath ---

func TestWorkspacePath_ExpandsTilde(t *testing.T) {
	cfg := DefaultConfig()
	path := cfg.WorkspacePath()
	if path == "" {
		t.Error("WorkspacePath should not be empty")
	}
	if path[0] == '~' {
		t.Error("WorkspacePath should expand ~")
	}
	if !strings.HasSuffix(path, "/.clawdroid/workspace") {
		t.Errorf("WorkspacePath = %q, expected to end with /.clawdroid/workspace", path)
	}
}

// --- SaveConfig → LoadConfig round-trip ---

func TestSaveConfig_LoadConfig_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "roundtrip.json")

	cfg := DefaultConfig()
	cfg.LLM.Model = "test-model-rt"
	cfg.LLM.APIKey = "test-key-rt"
	cfg.Gateway.Port = 7777
	cfg.Heartbeat.Interval = 15

	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.LLM.Model != "test-model-rt" {
		t.Errorf("LLM.Model = %q, want %q", loaded.LLM.Model, "test-model-rt")
	}
	if loaded.LLM.APIKey != "test-key-rt" {
		t.Errorf("LLM.APIKey = %q", loaded.LLM.APIKey)
	}
	if loaded.Gateway.Port != 7777 {
		t.Errorf("Gateway.Port = %d, want 7777", loaded.Gateway.Port)
	}
	if loaded.Heartbeat.Interval != 15 {
		t.Errorf("Heartbeat.Interval = %d, want 15", loaded.Heartbeat.Interval)
	}
}
