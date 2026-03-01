package channels

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/KarakuriAgent/clawdroid/pkg/config"
)

// helper to build a minimal LINEChannel for testing pure functions (no Start needed)
func newTestLINEChannel(secret, botUserID, botDisplayName string) *LINEChannel {
	return &LINEChannel{
		BaseChannel: &BaseChannel{name: "line"},
		config: config.LINEConfig{
			ChannelSecret: secret,
		},
		botUserID:      botUserID,
		botDisplayName: botDisplayName,
	}
}

// --- buildTextMessage ---

func TestBuildTextMessage_WithQuoteToken(t *testing.T) {
	msg := buildTextMessage("hello", "qt123")
	if msg["type"] != "text" {
		t.Errorf("type = %q", msg["type"])
	}
	if msg["text"] != "hello" {
		t.Errorf("text = %q", msg["text"])
	}
	if msg["quoteToken"] != "qt123" {
		t.Errorf("quoteToken = %q", msg["quoteToken"])
	}
}

func TestBuildTextMessage_WithoutQuoteToken(t *testing.T) {
	msg := buildTextMessage("hello", "")
	if msg["type"] != "text" {
		t.Errorf("type = %q", msg["type"])
	}
	if msg["text"] != "hello" {
		t.Errorf("text = %q", msg["text"])
	}
	if _, ok := msg["quoteToken"]; ok {
		t.Error("quoteToken should not be present")
	}
}

// --- verifySignature ---

func computeSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func TestVerifySignature_Valid(t *testing.T) {
	c := newTestLINEChannel("my-secret", "", "")
	body := []byte(`{"events":[]}`)
	sig := computeSignature("my-secret", body)

	if !c.verifySignature(body, sig) {
		t.Error("expected valid signature")
	}
}

func TestVerifySignature_Invalid(t *testing.T) {
	c := newTestLINEChannel("my-secret", "", "")
	body := []byte(`{"events":[]}`)

	if c.verifySignature(body, "invalid-signature") {
		t.Error("expected invalid signature")
	}
}

func TestVerifySignature_Empty(t *testing.T) {
	c := newTestLINEChannel("my-secret", "", "")
	body := []byte(`{"events":[]}`)

	if c.verifySignature(body, "") {
		t.Error("expected false for empty signature")
	}
}

func TestVerifySignature_WrongSecret(t *testing.T) {
	c := newTestLINEChannel("my-secret", "", "")
	body := []byte(`{"events":[]}`)
	sig := computeSignature("wrong-secret", body)

	if c.verifySignature(body, sig) {
		t.Error("expected false for wrong secret")
	}
}

// --- isBotMentioned ---

func TestIsBotMentioned_AllMention(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "@All test",
		Mention: &struct {
			Mentionees []lineMentionee `json:"mentionees"`
		}{
			Mentionees: []lineMentionee{
				{Type: "all", Index: 0, Length: 4},
			},
		},
	}
	if !c.isBotMentioned(msg) {
		t.Error("expected true for @all mention")
	}
}

func TestIsBotMentioned_UserIDMatch(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "@BotName hello",
		Mention: &struct {
			Mentionees []lineMentionee `json:"mentionees"`
		}{
			Mentionees: []lineMentionee{
				{Type: "user", UserID: "botUID", Index: 0, Length: 8},
			},
		},
	}
	if !c.isBotMentioned(msg) {
		t.Error("expected true for bot userId match")
	}
}

func TestIsBotMentioned_DisplayNameMatch(t *testing.T) {
	c := newTestLINEChannel("", "", "BotName")
	msg := lineMessage{
		Text: "@BotName hello",
		Mention: &struct {
			Mentionees []lineMentionee `json:"mentionees"`
		}{
			Mentionees: []lineMentionee{
				{Type: "user", UserID: "otherUID", Index: 0, Length: 8},
			},
		},
	}
	if !c.isBotMentioned(msg) {
		t.Error("expected true for display name match in mentionee text")
	}
}

func TestIsBotMentioned_TextFallback(t *testing.T) {
	c := newTestLINEChannel("", "", "BotName")
	msg := lineMessage{
		Text: "hello @BotName how are you",
	}
	if !c.isBotMentioned(msg) {
		t.Error("expected true for text-based @BotName fallback")
	}
}

func TestIsBotMentioned_NoMention(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "hello everyone",
	}
	if c.isBotMentioned(msg) {
		t.Error("expected false when no mention")
	}
}

func TestIsBotMentioned_OtherUserMentioned(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "@OtherUser hello",
		Mention: &struct {
			Mentionees []lineMentionee `json:"mentionees"`
		}{
			Mentionees: []lineMentionee{
				{Type: "user", UserID: "otherUID", Index: 0, Length: 10},
			},
		},
	}
	if c.isBotMentioned(msg) {
		t.Error("expected false when other user mentioned")
	}
}

// --- stripBotMention ---

func TestStripBotMention_WithMentionMetadata(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "@BotName hello world",
		Mention: &struct {
			Mentionees []lineMentionee `json:"mentionees"`
		}{
			Mentionees: []lineMentionee{
				{Type: "user", UserID: "botUID", Index: 0, Length: 8},
			},
		},
	}
	result := c.stripBotMention("@BotName hello world", msg)
	if result != "hello world" {
		t.Errorf("got %q, want %q", result, "hello world")
	}
}

func TestStripBotMention_TextFallback(t *testing.T) {
	c := newTestLINEChannel("", "", "BotName")
	msg := lineMessage{
		Text: "hey @BotName do stuff",
	}
	result := c.stripBotMention("hey @BotName do stuff", msg)
	// Fallback path uses strings.ReplaceAll which leaves a double space in the middle.
	// TrimSpace only trims leading/trailing whitespace, not internal.
	if result != "hey  do stuff" {
		t.Errorf("got %q, want %q", result, "hey  do stuff")
	}
}

func TestStripBotMention_NoMention(t *testing.T) {
	c := newTestLINEChannel("", "botUID", "BotName")
	msg := lineMessage{
		Text: "hello world",
	}
	result := c.stripBotMention("hello world", msg)
	if result != "hello world" {
		t.Errorf("got %q, want %q", result, "hello world")
	}
}

// --- resolveChatID ---

func TestResolveChatID_Group(t *testing.T) {
	c := newTestLINEChannel("", "", "")
	source := lineSource{Type: "group", GroupID: "G123", UserID: "U456"}
	if got := c.resolveChatID(source); got != "G123" {
		t.Errorf("got %q, want G123", got)
	}
}

func TestResolveChatID_Room(t *testing.T) {
	c := newTestLINEChannel("", "", "")
	source := lineSource{Type: "room", RoomID: "R789", UserID: "U456"}
	if got := c.resolveChatID(source); got != "R789" {
		t.Errorf("got %q, want R789", got)
	}
}

func TestResolveChatID_User(t *testing.T) {
	c := newTestLINEChannel("", "", "")
	source := lineSource{Type: "user", UserID: "U456"}
	if got := c.resolveChatID(source); got != "U456" {
		t.Errorf("got %q, want U456", got)
	}
}
