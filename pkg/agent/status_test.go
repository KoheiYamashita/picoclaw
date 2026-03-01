package agent

import (
	"testing"
)

// --- strArg ---

func TestStrArg(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		key  string
		want string
	}{
		{"existing key", map[string]interface{}{"k": "v"}, "k", "v"},
		{"missing key", map[string]interface{}{"k": "v"}, "other", ""},
		{"non-string value", map[string]interface{}{"k": 123}, "k", ""},
		{"nil map", nil, "k", ""},
		{"empty string value", map[string]interface{}{"k": ""}, "k", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strArg(tt.args, tt.key)
			if got != tt.want {
				t.Errorf("strArg(%v, %q) = %q, want %q", tt.args, tt.key, got, tt.want)
			}
		})
	}
}

// --- truncLabel ---

func TestTruncLabel(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxRunes int
		want     string
	}{
		{"short string", "hello", 10, "hello"},
		{"exact max", "hello", 5, "hello"},
		{"exceeds max", "hello world", 5, "hello..."},
		{"unicode within", "日本語テスト", 6, "日本語テスト"},
		{"unicode exceeds", "日本語テスト", 3, "日本語..."},
		{"empty", "", 5, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncLabel(tt.s, tt.maxRunes)
			if got != tt.want {
				t.Errorf("truncLabel(%q, %d) = %q, want %q", tt.s, tt.maxRunes, got, tt.want)
			}
		})
	}
}

// --- hostFromURL ---

func TestHostFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"normal URL", "https://example.com/path", "example.com"},
		{"URL with port", "https://example.com:8080/path", "example.com:8080"},
		{"invalid URL", "://bad", "://bad"},
		{"empty string", "", ""},
		{"no scheme", "not-a-url", "not-a-url"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hostFromURL(tt.url)
			if got != tt.want {
				t.Errorf("hostFromURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

// --- statusLabel ---

func TestStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		args     map[string]interface{}
		contains string
	}{
		{"web_search with query", "web_search", map[string]interface{}{"query": "golang"}, "golang"},
		{"web_search no query", "web_search", map[string]interface{}{}, "検索中..."},
		{"web_fetch with url", "web_fetch", map[string]interface{}{"url": "https://example.com/page"}, "example.com"},
		{"web_fetch no url", "web_fetch", map[string]interface{}{}, "ページ取得中..."},
		{"read_file with path", "read_file", map[string]interface{}{"path": "/home/user/file.txt"}, "file.txt"},
		{"read_file no path", "read_file", map[string]interface{}{}, "ファイル読み取り中..."},
		{"write_file", "write_file", map[string]interface{}{"path": "/tmp/out.txt"}, "out.txt"},
		{"edit_file", "edit_file", map[string]interface{}{}, "ファイル編集中..."},
		{"append_file", "append_file", map[string]interface{}{}, "ファイル追記中..."},
		{"list_dir with path", "list_dir", map[string]interface{}{"path": "/home/user/docs"}, "docs/"},
		{"list_dir no path", "list_dir", map[string]interface{}{}, "フォルダ確認中..."},
		{"exec with command", "exec", map[string]interface{}{"command": "ls -la"}, "ls -la"},
		{"exec no command", "exec", map[string]interface{}{}, "コマンド実行中..."},
		{"memory", "memory", map[string]interface{}{"action": "read_long_term"}, "メモリ読み込み中..."},
		{"skill", "skill", map[string]interface{}{"action": "skill_list"}, "スキル一覧取得中..."},
		{"cron", "cron", map[string]interface{}{"action": "add"}, "リマインダー設定中..."},
		{"message", "message", map[string]interface{}{}, "メッセージ送信中..."},
		{"spawn with label", "spawn", map[string]interface{}{"label": "task1"}, "task1"},
		{"spawn no label", "spawn", map[string]interface{}{}, "サブタスク開始中..."},
		{"subagent with label", "subagent", map[string]interface{}{"label": "sub1"}, "sub1"},
		{"subagent no label", "subagent", map[string]interface{}{}, "サブタスク実行中..."},
		{"android", "android", map[string]interface{}{"action": "screenshot"}, "スクリーンショット撮影中..."},
		{"exit", "exit", map[string]interface{}{}, "アシスタント終了中..."},
		{"mcp", "mcp", map[string]interface{}{"action": "mcp_list"}, "MCPサーバー一覧取得中..."},
		{"unknown tool", "unknown_tool", map[string]interface{}{}, "処理中..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusLabel(tt.toolName, tt.args)
			if got == "" {
				t.Error("statusLabel returned empty string")
			}
			if !containsStr(got, tt.contains) {
				t.Errorf("statusLabel(%q, %v) = %q, want to contain %q", tt.toolName, tt.args, got, tt.contains)
			}
		})
	}
}

// --- fileStatusLabel ---

func TestFileStatusLabel(t *testing.T) {
	got := fileStatusLabel("ファイル読み取り中...", map[string]interface{}{"path": "/home/user/test.go"})
	if got != "ファイル読み取り中...（test.go）" {
		t.Errorf("got %q", got)
	}

	got = fileStatusLabel("ファイル読み取り中...", map[string]interface{}{})
	if got != "ファイル読み取り中..." {
		t.Errorf("got %q", got)
	}
}

// --- memoryStatusLabel ---

func TestMemoryStatusLabel(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{"read_long_term", "メモリ読み込み中..."},
		{"read_daily", "今日のメモ読み込み中..."},
		{"write_long_term", "メモリ書き込み中..."},
		{"append_daily", "今日のメモ追記中..."},
		{"unknown", "メモリ操作中..."},
	}
	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got := memoryStatusLabel(map[string]interface{}{"action": tt.action})
			if got != tt.want {
				t.Errorf("memoryStatusLabel(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}

// --- skillStatusLabel ---

func TestSkillStatusLabel(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		want string
	}{
		{"skill_list", map[string]interface{}{"action": "skill_list"}, "スキル一覧取得中..."},
		{"skill_read with name", map[string]interface{}{"action": "skill_read", "name": "github"}, "スキル読み込み中...（github）"},
		{"skill_read no name", map[string]interface{}{"action": "skill_read"}, "スキル読み込み中..."},
		{"unknown", map[string]interface{}{"action": "other"}, "スキル操作中..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := skillStatusLabel(tt.args)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- cronStatusLabel ---

func TestCronStatusLabel(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{"add", "リマインダー設定中..."},
		{"list", "スケジュール一覧取得中..."},
		{"remove", "スケジュール削除中..."},
		{"unknown", "スケジュール変更中..."},
	}
	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got := cronStatusLabel(map[string]interface{}{"action": tt.action})
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- androidStatusLabel ---

func TestAndroidStatusLabel(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		want string
	}{
		{"search_apps", map[string]interface{}{"action": "search_apps"}, "アプリ検索中..."},
		{"app_info with pkg", map[string]interface{}{"action": "app_info", "package_name": "com.example"}, "アプリ情報取得中...（com.example）"},
		{"app_info no pkg", map[string]interface{}{"action": "app_info"}, "アプリ情報取得中..."},
		{"launch_app with pkg", map[string]interface{}{"action": "launch_app", "package_name": "com.test"}, "アプリ起動中...（com.test）"},
		{"launch_app no pkg", map[string]interface{}{"action": "launch_app"}, "アプリ起動中..."},
		{"screenshot", map[string]interface{}{"action": "screenshot"}, "スクリーンショット撮影中..."},
		{"get_ui_tree", map[string]interface{}{"action": "get_ui_tree"}, "UI要素取得中..."},
		{"tap", map[string]interface{}{"action": "tap"}, "タップ中..."},
		{"swipe", map[string]interface{}{"action": "swipe"}, "スワイプ中..."},
		{"text", map[string]interface{}{"action": "text"}, "テキスト入力中..."},
		{"keyevent with key", map[string]interface{}{"action": "keyevent", "key": "BACK"}, "キー操作中...（BACK）"},
		{"keyevent no key", map[string]interface{}{"action": "keyevent"}, "キー操作中..."},
		{"broadcast", map[string]interface{}{"action": "broadcast"}, "ブロードキャスト送信中..."},
		{"intent", map[string]interface{}{"action": "intent"}, "インテント送信中..."},
		{"unknown", map[string]interface{}{"action": "other"}, "デバイス操作中..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := androidStatusLabel(tt.args)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- mcpStatusLabel ---

func TestMcpStatusLabel(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		want string
	}{
		{"mcp_list", map[string]interface{}{"action": "mcp_list"}, "MCPサーバー一覧取得中..."},
		{"mcp_tools with server", map[string]interface{}{"action": "mcp_tools", "server": "myserver"}, "MCPツール取得中...（myserver）"},
		{"mcp_tools no server", map[string]interface{}{"action": "mcp_tools"}, "MCPツール取得中..."},
		{"mcp_call tool+server", map[string]interface{}{"action": "mcp_call", "tool": "mytool", "server": "srv"}, "MCPツール実行中...（srv/mytool）"},
		{"mcp_call tool only", map[string]interface{}{"action": "mcp_call", "tool": "mytool"}, "MCPツール実行中...（mytool）"},
		{"mcp_call no args", map[string]interface{}{"action": "mcp_call"}, "MCPツール実行中..."},
		{"unknown", map[string]interface{}{"action": "other"}, "MCP操作中..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mcpStatusLabel(tt.args)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func containsStr(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
