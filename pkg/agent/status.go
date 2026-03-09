package agent

import (
	"net/url"
	"path/filepath"
	"unicode/utf8"

	"github.com/KarakuriAgent/clawdroid/pkg/i18n"
)

// statusLabel generates a human-readable localized status label for a tool call.
func statusLabel(toolName string, args map[string]interface{}, locale string) string {
	switch toolName {
	case "web_search":
		if q := strArg(args, "query"); q != "" {
			return i18n.Tf(locale, "status.searching_q", truncLabel(q, 20))
		}
		return i18n.T(locale, "status.searching")
	case "web_fetch":
		if u := strArg(args, "url"); u != "" {
			return i18n.Tf(locale, "status.fetching_q", hostFromURL(u))
		}
		return i18n.T(locale, "status.fetching_page")
	case "read_file":
		return fileStatusLabel(locale, "status.reading_file", "status.reading_file_q", args)
	case "write_file":
		return fileStatusLabel(locale, "status.writing_file", "status.writing_file_q", args)
	case "edit_file":
		return fileStatusLabel(locale, "status.editing_file", "status.editing_file_q", args)
	case "append_file":
		return fileStatusLabel(locale, "status.appending_file", "status.appending_file_q", args)
	case "list_dir":
		if p := strArg(args, "path"); p != "" {
			return i18n.Tf(locale, "status.listing_dir_q", filepath.Base(p)+"/")
		}
		return i18n.T(locale, "status.listing_dir")
	case "exec":
		if c := strArg(args, "command"); c != "" {
			return i18n.Tf(locale, "status.running_command_q", truncLabel(c, 30))
		}
		return i18n.T(locale, "status.running_command")
	case "memory":
		return memoryStatusLabel(args, locale)
	case "skill":
		return skillStatusLabel(args, locale)
	case "cron":
		return cronStatusLabel(args, locale)
	case "message":
		return i18n.T(locale, "status.sending_message")
	case "spawn":
		if l := strArg(args, "label"); l != "" {
			return i18n.Tf(locale, "status.spawn_q", truncLabel(l, 20))
		}
		return i18n.T(locale, "status.spawn")
	case "subagent":
		if l := strArg(args, "label"); l != "" {
			return i18n.Tf(locale, "status.subagent_q", truncLabel(l, 20))
		}
		return i18n.T(locale, "status.subagent")
	case "android":
		return androidStatusLabel(args, locale)
	case "exit":
		return i18n.T(locale, "status.exit")
	case "mcp":
		return mcpStatusLabel(args, locale)
	default:
		return i18n.T(locale, "status.processing")
	}
}

func fileStatusLabel(locale, baseKey, fmtKey string, args map[string]interface{}) string {
	if p := strArg(args, "path"); p != "" {
		return i18n.Tf(locale, fmtKey, filepath.Base(p))
	}
	return i18n.T(locale, baseKey)
}

func memoryStatusLabel(args map[string]interface{}, locale string) string {
	switch strArg(args, "action") {
	case "read_long_term":
		return i18n.T(locale, "status.memory_read")
	case "read_daily":
		return i18n.T(locale, "status.memory_read_daily")
	case "write_long_term":
		return i18n.T(locale, "status.memory_write")
	case "append_daily":
		return i18n.T(locale, "status.memory_append_daily")
	default:
		return i18n.T(locale, "status.memory_default")
	}
}

func skillStatusLabel(args map[string]interface{}, locale string) string {
	switch strArg(args, "action") {
	case "skill_list":
		return i18n.T(locale, "status.skill_list")
	case "skill_read":
		if n := strArg(args, "name"); n != "" {
			return i18n.Tf(locale, "status.skill_read_q", n)
		}
		return i18n.T(locale, "status.skill_read")
	default:
		return i18n.T(locale, "status.skill_default")
	}
}

func cronStatusLabel(args map[string]interface{}, locale string) string {
	switch strArg(args, "action") {
	case "add":
		return i18n.T(locale, "status.cron_add")
	case "list":
		return i18n.T(locale, "status.cron_list")
	case "remove":
		return i18n.T(locale, "status.cron_remove")
	default:
		return i18n.T(locale, "status.cron_default")
	}
}

func androidStatusLabel(args map[string]interface{}, locale string) string {
	switch strArg(args, "action") {
	case "search_apps":
		return i18n.T(locale, "status.android_search_apps")
	case "app_info":
		if p := strArg(args, "package_name"); p != "" {
			return i18n.Tf(locale, "status.android_app_info_q", truncLabel(p, 25))
		}
		return i18n.T(locale, "status.android_app_info")
	case "launch_app":
		if p := strArg(args, "package_name"); p != "" {
			return i18n.Tf(locale, "status.android_launch_app_q", truncLabel(p, 25))
		}
		return i18n.T(locale, "status.android_launch_app")
	case "screenshot":
		return i18n.T(locale, "status.android_screenshot")
	case "get_ui_tree":
		return i18n.T(locale, "status.android_get_ui_tree")
	case "tap":
		return i18n.T(locale, "status.android_tap")
	case "swipe":
		return i18n.T(locale, "status.android_swipe")
	case "text":
		return i18n.T(locale, "status.android_text")
	case "keyevent":
		if k := strArg(args, "key"); k != "" {
			return i18n.Tf(locale, "status.android_keyevent_q", k)
		}
		return i18n.T(locale, "status.android_keyevent")
	case "broadcast":
		return i18n.T(locale, "status.android_broadcast")
	case "intent":
		return i18n.T(locale, "status.android_intent")
	default:
		return i18n.T(locale, "status.android_default")
	}
}

func mcpStatusLabel(args map[string]interface{}, locale string) string {
	switch strArg(args, "action") {
	case "mcp_list":
		return i18n.T(locale, "status.mcp_list")
	case "mcp_tools":
		if s := strArg(args, "server"); s != "" {
			return i18n.Tf(locale, "status.mcp_tools_q", s)
		}
		return i18n.T(locale, "status.mcp_tools")
	case "mcp_call":
		if t := strArg(args, "tool"); t != "" {
			if s := strArg(args, "server"); s != "" {
				return i18n.Tf(locale, "status.mcp_call_sq", s, t)
			}
			return i18n.Tf(locale, "status.mcp_call_q", t)
		}
		return i18n.T(locale, "status.mcp_call")
	default:
		return i18n.T(locale, "status.mcp_default")
	}
}

// strArg extracts a string argument from a tool arguments map.
func strArg(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// truncLabel truncates a string to maxRunes runes, appending "..." if truncated.
func truncLabel(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "..."
}

// hostFromURL extracts the hostname from a URL string.
func hostFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return truncLabel(rawURL, 30)
	}
	return u.Host
}
