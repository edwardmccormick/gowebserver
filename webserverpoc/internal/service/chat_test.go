package service

import (
	"context"
	"testing"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/edwardmccormick/gowebserver/internal/realtime"
	"github.com/edwardmccormick/gowebserver/internal/store"
)

type fakeMatchStore struct {
	match       store.Match
	getErr      error
	updateErr   error
	updated     store.Match
	updateCalls int
}

func (f *fakeMatchStore) GetByID(ctx context.Context, id uint) (store.Match, error) {
	return f.match, f.getErr
}

func (f *fakeMatchStore) Update(ctx context.Context, match store.Match) error {
	f.updated = match
	f.updateCalls++
	return f.updateErr
}

type fakeChatHistoryStore struct {
	conversation store.Conversation
	loadErr      error
	saveErr      error
	appendErr    error
	saved        store.Conversation
	appended     store.ChatMessage
	appendMatch  uint
}

func (f *fakeChatHistoryStore) LoadByMatchID(ctx context.Context, matchID uint) (store.Conversation, error) {
	return f.conversation, f.loadErr
}

func (f *fakeChatHistoryStore) SaveConversation(ctx context.Context, conversation store.Conversation) error {
	f.saved = conversation
	return f.saveErr
}

func (f *fakeChatHistoryStore) AppendMessage(ctx context.Context, matchID uint, message store.ChatMessage) error {
	f.appendMatch = matchID
	f.appended = message
	return f.appendErr
}

func (f *fakeChatHistoryStore) LatestMessageID(ctx context.Context, matchID uint) (int64, error) {
	if len(f.conversation.Messages) == 0 {
		return 0, store.ErrNotFound
	}
	return f.conversation.Messages[len(f.conversation.Messages)-1].ID, nil
}

func (f *fakeChatHistoryStore) SaveReadState(ctx context.Context, state store.MatchReadState) error {
	return nil
}

func TestChatServiceUpdateUnreadCountsOfferedSender(t *testing.T) {
	matchStore := &fakeMatchStore{
		match: store.Match{ID: 9, Offered: 10, Accepted: 20},
	}
	notifications := realtime.NewNotificationCenter()
	hub := realtime.NewHub()
	service := NewChatService(matchStore, &fakeChatHistoryStore{}, hub, notifications)

	ch := make(chan realtime.NotificationEvent, 1)
	notifications.Register("20", ch)

	if err := service.UpdateUnreadCounts(context.Background(), 9, 10); err != nil {
		t.Fatalf("UpdateUnreadCounts returned error: %v", err)
	}

	if matchStore.updateCalls != 1 {
		t.Fatalf("expected 1 update call, got %d", matchStore.updateCalls)
	}
	if matchStore.updated.UnreadAccepted != 1 {
		t.Fatalf("expected unread accepted increment, got %+v", matchStore.updated)
	}
	if matchStore.updated.LastMessageTime.IsZero() {
		t.Fatal("expected last message time to be set")
	}

	select {
	case event := <-ch:
		if event.MatchID != 9 || event.Count != 1 {
			t.Fatalf("unexpected notification event: %+v", event)
		}
	default:
		t.Fatal("expected notification for inactive recipient")
	}
}

func TestChatServiceUpdateUnreadCountsNoNotificationWhenRecipientActive(t *testing.T) {
	matchStore := &fakeMatchStore{
		match: store.Match{ID: 11, Offered: 1, Accepted: 2},
	}
	notifications := realtime.NewNotificationCenter()
	hub := realtime.NewHub()
	service := NewChatService(matchStore, &fakeChatHistoryStore{}, hub, notifications)

	hub.SetActiveConnection(2, nil)

	ch := make(chan realtime.NotificationEvent, 1)
	notifications.Register("2", ch)

	if err := service.UpdateUnreadCounts(context.Background(), 11, 1); err != nil {
		t.Fatalf("UpdateUnreadCounts returned error: %v", err)
	}

	select {
	case event := <-ch:
		t.Fatalf("did not expect notification, got %+v", event)
	default:
	}
}

func TestChatServiceResetUnreadCount(t *testing.T) {
	matchStore := &fakeMatchStore{
		match: store.Match{ID: 7, Offered: 5, Accepted: 8, UnreadOffered: 3, UnreadAccepted: 4},
	}
	service := NewChatService(matchStore, &fakeChatHistoryStore{}, nil, nil)

	if err := service.ResetUnreadCount(context.Background(), 7, 8); err != nil {
		t.Fatalf("ResetUnreadCount returned error: %v", err)
	}

	if matchStore.updated.UnreadAccepted != 0 {
		t.Fatalf("expected accepted unread count reset, got %+v", matchStore.updated)
	}
	if matchStore.updated.UnreadOffered != 3 {
		t.Fatalf("expected offered unread count unchanged, got %+v", matchStore.updated)
	}
}

func TestChatServiceLoadSaveAndAppendConversation(t *testing.T) {
	now := time.Now()
	chatStore := &fakeChatHistoryStore{
		conversation: store.Conversation{
			MatchID: 3,
			Messages: []store.ChatMessage{
				{ID: 1, MatchID: 3, Who: 9, Message: "hello", Time: now},
			},
		},
	}
	service := NewChatService(&fakeMatchStore{}, chatStore, nil, nil)

	conversation, err := service.LoadConversation(context.Background(), 3)
	if err != nil {
		t.Fatalf("LoadConversation returned error: %v", err)
	}
	if conversation.MatchID != 3 || len(conversation.Messages) != 1 || conversation.Messages[0].Message != "hello" {
		t.Fatalf("unexpected conversation: %+v", conversation)
	}

	savePayload := domain.Conversation{
		MatchID: 4,
		Messages: []domain.ChatMessage{
			{ID: 2, MatchID: 4, Who: 7, Message: "saved"},
		},
	}
	if err := service.SaveConversation(context.Background(), savePayload); err != nil {
		t.Fatalf("SaveConversation returned error: %v", err)
	}
	if chatStore.saved.MatchID != 4 || len(chatStore.saved.Messages) != 1 || chatStore.saved.Messages[0].Message != "saved" {
		t.Fatalf("unexpected saved payload: %+v", chatStore.saved)
	}

	appendPayload := domain.ChatMessage{ID: 3, MatchID: 4, Who: 8, Message: "append"}
	if err := service.AppendMessage(context.Background(), 4, appendPayload); err != nil {
		t.Fatalf("AppendMessage returned error: %v", err)
	}
	if chatStore.appendMatch != 4 || chatStore.appended.Message != "append" {
		t.Fatalf("unexpected append payload: match=%d message=%+v", chatStore.appendMatch, chatStore.appended)
	}
}
