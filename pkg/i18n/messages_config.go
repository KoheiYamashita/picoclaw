package i18n

func init() {
	// English labels match the struct tag defaults — no registration needed
	// since the struct tags already provide English labels.

	register("ja", map[string]string{
		// Top-level sections
		"LLM":                "LLM",
		"Agent Defaults":     "エージェント設定",
		"Messaging Channels": "メッセージングチャンネル",
		"Gateway":            "ゲートウェイ",
		"Tool Settings":      "ツール設定",
		"Heartbeat":          "ハートビート",
		"Rate Limits":        "レート制限",

		// LLM
		"Model":    "モデル",
		"API Key":  "APIキー",
		"Base URL": "ベースURL",

		// Agent Defaults
		"Defaults":              "デフォルト",
		"Workspace":             "ワークスペース",
		"Data Directory":        "データディレクトリ",
		"Restrict to Workspace": "ワークスペースに制限",
		"Max Tokens":            "最大トークン数",
		"Context Window":        "コンテキストウィンドウ",
		"Temperature":           "温度",
		"Max Tool Iterations":   "最大ツール反復回数",
		"Queue Messages":        "メッセージキュー",
		"Show Errors":           "エラー表示",
		"Show Warnings":         "警告表示",

		// Channels
		"WhatsApp":             "WhatsApp",
		"Telegram":             "Telegram",
		"Discord":              "Discord",
		"Slack":                "Slack",
		"LINE":                 "LINE",
		"WebSocket":            "WebSocket",
		"Enabled":              "有効",
		"Token":                "トークン",
		"Bot Token":            "ボットトークン",
		"App Token":            "アプリトークン",
		"Proxy":                "プロキシ",
		"Allow From":           "許可リスト",
		"Bridge URL":           "ブリッジURL",
		"Host":                 "ホスト",
		"Port":                 "ポート",
		"Path":                 "パス",
		"Channel Secret":       "チャンネルシークレット",
		"Channel Access Token": "チャンネルアクセストークン",
		"Webhook Host":         "Webhookホスト",
		"Webhook Port":         "Webhookポート",
		"Webhook Path":         "Webhookパス",

		// Heartbeat
		"Interval": "間隔",

		// Rate Limits
		"Max Tool Calls Per Minute": "1分あたりの最大ツール呼び出し数",
		"Max Requests Per Minute":   "1分あたりの最大リクエスト数",

		// Tools
		"Web Search":  "Web検索",
		"Shell Exec":  "シェル実行",
		"Android":     "Android",
		"Memory":      "メモリ",
		"MCP Servers": "MCPサーバー",

		// Web search sub
		"Brave Search": "Brave検索",
		"DuckDuckGo":   "DuckDuckGo",
		"Max Results":  "最大結果数",

		// Android tool categories
		"Alarm":            "アラーム",
		"Calendar":         "カレンダー",
		"Calendar Account": "カレンダーアカウント",
		"Contacts":         "連絡先",
		"Communication":    "コミュニケーション",
		"Media":            "メディア",
		"Navigation":       "ナビゲーション",
		"Device Control":   "デバイス制御",
		"Settings":         "設定",
		"Web":              "Web",
		"Clipboard":        "クリップボード",

		// Android actions
		"Set Alarm":            "アラーム設定",
		"Set Timer":            "タイマー設定",
		"Dismiss Alarm":        "アラーム解除",
		"Show Alarms":          "アラーム表示",
		"Create Event":         "イベント作成",
		"Query Events":         "イベント検索",
		"Update Event":         "イベント更新",
		"Delete Event":         "イベント削除",
		"List Calendars":       "カレンダー一覧",
		"Add Reminder":         "リマインダー追加",
		"Search Contacts":      "連絡先検索",
		"Get Contact Detail":   "連絡先詳細",
		"Add Contact":          "連絡先追加",
		"Dial":                 "電話発信",
		"Compose SMS":          "SMS作成",
		"Compose Email":        "メール作成",
		"Play/Pause":           "再生/一時停止",
		"Next":                 "次へ",
		"Previous":             "前へ",
		"Play Music Search":    "音楽検索再生",
		"Navigate":             "ナビゲーション",
		"Search Nearby":        "周辺検索",
		"Show Map":             "地図表示",
		"Get Current Location": "現在地取得",
		"Flashlight":           "フラッシュライト",
		"Set Volume":           "音量設定",
		"Set Ringer Mode":      "着信モード設定",
		"Set DND":              "おやすみモード設定",
		"Set Brightness":       "画面の明るさ設定",
		"Open Settings":        "設定を開く",
		"Open URL":             "URL を開く",
		"Copy":                 "コピー",
		"Read":                 "読み取り",
	})
}
