package i18n

func init() {
	register("en", map[string]string{
		// status labels
		"status.thinking":    "Thinking...",
		"status.processing":  "Processing...",
		"status.interrupted": "[Response was interrupted]",

		// web
		"status.searching":     "Searching...",
		"status.searching_q":   "Searching... (%s)",
		"status.fetching_page": "Fetching page...",
		"status.fetching_q":    "Fetching page... (%s)",

		// file operations
		"status.reading_file":     "Reading file...",
		"status.reading_file_q":   "Reading file... (%s)",
		"status.writing_file":     "Writing file...",
		"status.writing_file_q":   "Writing file... (%s)",
		"status.editing_file":     "Editing file...",
		"status.editing_file_q":   "Editing file... (%s)",
		"status.appending_file":   "Appending to file...",
		"status.appending_file_q": "Appending to file... (%s)",

		// directory
		"status.listing_dir":   "Checking folder...",
		"status.listing_dir_q": "Checking folder... (%s)",

		// exec
		"status.running_command":   "Running command...",
		"status.running_command_q": "Running command... (%s)",

		// memory
		"status.memory_read":         "Loading memory...",
		"status.memory_read_daily":   "Loading today's memo...",
		"status.memory_write":        "Writing memory...",
		"status.memory_append_daily": "Appending to today's memo...",
		"status.memory_default":      "Memory operation...",

		// skill
		"status.skill_list":    "Getting skill list...",
		"status.skill_read":    "Loading skill...",
		"status.skill_read_q":  "Loading skill... (%s)",
		"status.skill_default": "Skill operation...",

		// cron
		"status.cron_add":     "Setting reminder...",
		"status.cron_list":    "Getting schedule...",
		"status.cron_remove":  "Removing schedule...",
		"status.cron_default": "Updating schedule...",

		// message
		"status.sending_message": "Sending message...",

		// spawn/subagent
		"status.spawn":      "Starting subtask...",
		"status.spawn_q":    "Starting subtask... (%s)",
		"status.subagent":   "Running subtask...",
		"status.subagent_q": "Running subtask... (%s)",

		// android
		"status.android_search_apps":  "Searching apps...",
		"status.android_app_info":     "Getting app info...",
		"status.android_app_info_q":   "Getting app info... (%s)",
		"status.android_launch_app":   "Launching app...",
		"status.android_launch_app_q": "Launching app... (%s)",
		"status.android_screenshot":   "Taking screenshot...",
		"status.android_get_ui_tree":  "Getting UI elements...",
		"status.android_tap":          "Tapping...",
		"status.android_swipe":        "Swiping...",
		"status.android_text":         "Entering text...",
		"status.android_keyevent":     "Key operation...",
		"status.android_keyevent_q":   "Key operation... (%s)",
		"status.android_broadcast":    "Sending broadcast...",
		"status.android_intent":       "Sending intent...",
		"status.android_default":      "Device operation...",

		// exit
		"status.exit": "Shutting down assistant...",

		// mcp
		"status.mcp_list":    "Getting MCP server list...",
		"status.mcp_tools":   "Getting MCP tools...",
		"status.mcp_tools_q": "Getting MCP tools... (%s)",
		"status.mcp_call":    "Running MCP tool...",
		"status.mcp_call_q":  "Running MCP tool... (%s)",
		"status.mcp_call_sq": "Running MCP tool... (%s/%s)",
		"status.mcp_default": "MCP operation...",
	})

	register("ja", map[string]string{
		// status labels
		"status.thinking":    "思考中...",
		"status.processing":  "処理中...",
		"status.interrupted": "[応答は中断されました]",

		// web
		"status.searching":     "検索中...",
		"status.searching_q":   "検索中...（%s）",
		"status.fetching_page": "ページ取得中...",
		"status.fetching_q":    "ページ取得中...（%s）",

		// file operations
		"status.reading_file":     "ファイル読み取り中...",
		"status.reading_file_q":   "ファイル読み取り中...（%s）",
		"status.writing_file":     "ファイル書き込み中...",
		"status.writing_file_q":   "ファイル書き込み中...（%s）",
		"status.editing_file":     "ファイル編集中...",
		"status.editing_file_q":   "ファイル編集中...（%s）",
		"status.appending_file":   "ファイル追記中...",
		"status.appending_file_q": "ファイル追記中...（%s）",

		// directory
		"status.listing_dir":   "フォルダ確認中...",
		"status.listing_dir_q": "フォルダ確認中...（%s）",

		// exec
		"status.running_command":   "コマンド実行中...",
		"status.running_command_q": "コマンド実行中...（%s）",

		// memory
		"status.memory_read":         "メモリ読み込み中...",
		"status.memory_read_daily":   "今日のメモ読み込み中...",
		"status.memory_write":        "メモリ書き込み中...",
		"status.memory_append_daily": "今日のメモ追記中...",
		"status.memory_default":      "メモリ操作中...",

		// skill
		"status.skill_list":    "スキル一覧取得中...",
		"status.skill_read":    "スキル読み込み中...",
		"status.skill_read_q":  "スキル読み込み中...（%s）",
		"status.skill_default": "スキル操作中...",

		// cron
		"status.cron_add":     "リマインダー設定中...",
		"status.cron_list":    "スケジュール一覧取得中...",
		"status.cron_remove":  "スケジュール削除中...",
		"status.cron_default": "スケジュール変更中...",

		// message
		"status.sending_message": "メッセージ送信中...",

		// spawn/subagent
		"status.spawn":      "サブタスク開始中...",
		"status.spawn_q":    "サブタスク開始中...（%s）",
		"status.subagent":   "サブタスク実行中...",
		"status.subagent_q": "サブタスク実行中...（%s）",

		// android
		"status.android_search_apps":  "アプリ検索中...",
		"status.android_app_info":     "アプリ情報取得中...",
		"status.android_app_info_q":   "アプリ情報取得中...（%s）",
		"status.android_launch_app":   "アプリ起動中...",
		"status.android_launch_app_q": "アプリ起動中...（%s）",
		"status.android_screenshot":   "スクリーンショット撮影中...",
		"status.android_get_ui_tree":  "UI要素取得中...",
		"status.android_tap":          "タップ中...",
		"status.android_swipe":        "スワイプ中...",
		"status.android_text":         "テキスト入力中...",
		"status.android_keyevent":     "キー操作中...",
		"status.android_keyevent_q":   "キー操作中...（%s）",
		"status.android_broadcast":    "ブロードキャスト送信中...",
		"status.android_intent":       "インテント送信中...",
		"status.android_default":      "デバイス操作中...",

		// exit
		"status.exit": "アシスタント終了中...",

		// mcp
		"status.mcp_list":    "MCPサーバー一覧取得中...",
		"status.mcp_tools":   "MCPツール取得中...",
		"status.mcp_tools_q": "MCPツール取得中...（%s）",
		"status.mcp_call":    "MCPツール実行中...",
		"status.mcp_call_q":  "MCPツール実行中...（%s）",
		"status.mcp_call_sq": "MCPツール実行中...（%s/%s）",
		"status.mcp_default": "MCP操作中...",
	})
}
