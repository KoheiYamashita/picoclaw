package bus

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewMessageBus(t *testing.T) {
	mb := NewMessageBus()
	if mb == nil {
		t.Fatal("NewMessageBus returned nil")
	}
	if mb.inbound == nil {
		t.Error("inbound channel should be initialized")
	}
	if mb.outbound == nil {
		t.Error("outbound channel should be initialized")
	}
	if mb.handlers == nil {
		t.Error("handlers map should be initialized")
	}
	if mb.closed {
		t.Error("bus should not be closed initially")
	}
	// Verify buffer size
	if cap(mb.inbound) != 100 {
		t.Errorf("inbound buffer size = %d, want 100", cap(mb.inbound))
	}
	if cap(mb.outbound) != 100 {
		t.Errorf("outbound buffer size = %d, want 100", cap(mb.outbound))
	}
}

func TestPublishConsumeInbound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	msg := InboundMessage{
		Channel:  "test",
		SenderID: "user1",
		Content:  "hello",
	}

	mb.PublishInbound(msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, ok := mb.ConsumeInbound(ctx)
	if !ok {
		t.Fatal("ConsumeInbound returned false")
	}
	if got.Channel != msg.Channel {
		t.Errorf("Channel = %q, want %q", got.Channel, msg.Channel)
	}
	if got.SenderID != msg.SenderID {
		t.Errorf("SenderID = %q, want %q", got.SenderID, msg.SenderID)
	}
	if got.Content != msg.Content {
		t.Errorf("Content = %q, want %q", got.Content, msg.Content)
	}
}

func TestPublishSubscribeOutbound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	msg := OutboundMessage{
		Channel: "discord",
		ChatID:  "123",
		Content: "response",
	}

	mb.PublishOutbound(msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, ok := mb.SubscribeOutbound(ctx)
	if !ok {
		t.Fatal("SubscribeOutbound returned false")
	}
	if got.Channel != msg.Channel {
		t.Errorf("Channel = %q, want %q", got.Channel, msg.Channel)
	}
	if got.ChatID != msg.ChatID {
		t.Errorf("ChatID = %q, want %q", got.ChatID, msg.ChatID)
	}
	if got.Content != msg.Content {
		t.Errorf("Content = %q, want %q", got.Content, msg.Content)
	}
}

func TestConsumeInbound_ContextCancel(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, ok := mb.ConsumeInbound(ctx)
	if ok {
		t.Error("ConsumeInbound should return false when context is cancelled")
	}
}

func TestSubscribeOutbound_ContextCancel(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, ok := mb.SubscribeOutbound(ctx)
	if ok {
		t.Error("SubscribeOutbound should return false when context is cancelled")
	}
}

func TestRegisterGetHandler(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	called := false
	handler := func(msg InboundMessage) error {
		called = true
		return nil
	}

	mb.RegisterHandler("discord", handler)

	got, ok := mb.GetHandler("discord")
	if !ok {
		t.Fatal("GetHandler returned false for registered handler")
	}
	if got == nil {
		t.Fatal("GetHandler returned nil handler")
	}

	// Verify handler works
	_ = got(InboundMessage{})
	if !called {
		t.Error("handler was not called")
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	_, ok := mb.GetHandler("nonexistent")
	if ok {
		t.Error("GetHandler should return false for unregistered handler")
	}
}

func TestClose(t *testing.T) {
	mb := NewMessageBus()

	mb.Close()
	if !mb.closed {
		t.Error("bus should be closed after Close()")
	}

	// Double close should not panic
	mb.Close()

	// Publish after close should not panic
	mb.PublishInbound(InboundMessage{Content: "after close"})
	mb.PublishOutbound(OutboundMessage{Content: "after close"})
}

func TestConcurrentPublishConsume(t *testing.T) {
	mb := NewMessageBus()
	defer mb.Close()

	const n = 50
	var wg sync.WaitGroup

	// Publish concurrently
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			mb.PublishInbound(InboundMessage{Content: "msg"})
		}(i)
	}

	// Consume concurrently
	consumed := make(chan struct{}, n)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if _, ok := mb.ConsumeInbound(ctx); ok {
				consumed <- struct{}{}
			}
		}()
	}

	wg.Wait()

	if len(consumed) != n {
		t.Errorf("consumed %d messages, want %d", len(consumed), n)
	}
}
