package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/edwardmccormick/gowebserver/internal/store"
)

func TestChatMessageServiceGetPagedMessages(t *testing.T) {
	now := time.Now()
	chatStore := &fakeChatHistoryStore{
		conversation: store.Conversation{
			MatchID: 5,
			Messages: []store.ChatMessage{
				{ID: 2, MatchID: 5, Message: "second", Time: now},
				{ID: 1, MatchID: 5, Message: "first", Time: now.Add(-time.Minute)},
				{ID: 3, MatchID: 5, Message: "third", Time: now.Add(time.Minute)},
			},
		},
	}
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) { return domain.Match{}, nil },
		nil, nil, nil, nil,
	)

	messages, total, err := service.GetPagedMessages(context.Background(), 5, 2, 1)
	if err != nil {
		t.Fatalf("GetPagedMessages returned error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(messages) != 2 || messages[0].Message != "second" || messages[1].Message != "third" {
		t.Fatalf("unexpected paged messages: %+v", messages)
	}
}

func TestChatMessageServiceGetPagedMessagesNotFound(t *testing.T) {
	chatStore := &fakeChatHistoryStore{loadErr: store.ErrNotFound}
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) { return domain.Match{}, nil },
		nil, nil, nil, nil,
	)

	messages, total, err := service.GetPagedMessages(context.Background(), 5, 10, 0)
	if err != nil {
		t.Fatalf("GetPagedMessages returned error: %v", err)
	}
	if total != 0 || len(messages) != 0 {
		t.Fatalf("expected empty result, got total=%d messages=%+v", total, messages)
	}
}

func TestChatMessageServiceEnsureIntroductionSavesWhenEmpty(t *testing.T) {
	chatStore := &fakeChatHistoryStore{loadErr: store.ErrNotFound}
	loadCalls := 0
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) {
			loadCalls++
			return domain.Match{Offered: 1, Accepted: 2}, nil
		},
		func(match domain.Match) (*domain.ChatMessage, error) {
			return &domain.ChatMessage{MatchID: 6, Message: "intro"}, nil
		},
		nil, nil, nil,
	)

	message, err := service.EnsureIntroduction(context.Background(), 6)
	if err != nil {
		t.Fatalf("EnsureIntroduction returned error: %v", err)
	}
	if loadCalls != 1 {
		t.Fatalf("expected match to be loaded once, got %d", loadCalls)
	}
	if message == nil || message.Message != "intro" {
		t.Fatalf("unexpected intro message: %+v", message)
	}
	if chatStore.saved.MatchID != 6 || len(chatStore.saved.Messages) != 1 || chatStore.saved.Messages[0].Message != "intro" {
		t.Fatalf("expected intro conversation to be saved, got %+v", chatStore.saved)
	}
}

func TestChatMessageServiceEnsureIntroductionSkipsWhenMessagesExist(t *testing.T) {
	chatStore := &fakeChatHistoryStore{
		conversation: store.Conversation{
			MatchID: 2,
			Messages: []store.ChatMessage{
				{ID: 1, MatchID: 2, Message: "existing"},
			},
		},
	}
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) {
			t.Fatal("did not expect match loader to be called")
			return domain.Match{}, nil
		},
		func(match domain.Match) (*domain.ChatMessage, error) {
			t.Fatal("did not expect intro generator to be called")
			return nil, nil
		},
		nil, nil, nil,
	)

	message, err := service.EnsureIntroduction(context.Background(), 2)
	if err != nil {
		t.Fatalf("EnsureIntroduction returned error: %v", err)
	}
	if message != nil {
		t.Fatalf("expected nil message when conversation already exists, got %+v", message)
	}
}

func TestChatMessageServiceGenerateDateMessageAppendsAndBroadcasts(t *testing.T) {
	now := time.Now()
	chatStore := &fakeChatHistoryStore{
		conversation: store.Conversation{
			MatchID: 8,
			Messages: []store.ChatMessage{
				{ID: 2, MatchID: 8, Message: "later", Time: now},
				{ID: 1, MatchID: 8, Message: "earlier", Time: now.Add(-time.Minute)},
			},
		},
	}

	var broadcastMatchID uint
	var broadcastMessage *domain.ChatMessage
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) {
			return domain.Match{Offered: 1, Accepted: 2}, nil
		},
		nil,
		func(match domain.Match, recentMessages []domain.ChatMessage) (*domain.ChatMessage, error) {
			if len(recentMessages) != 2 || recentMessages[0].Message != "earlier" || recentMessages[1].Message != "later" {
				t.Fatalf("expected sorted message history, got %+v", recentMessages)
			}
			return &domain.ChatMessage{MatchID: 8, Message: "date idea"}, nil
		},
		nil,
		func(matchID uint, message *domain.ChatMessage) {
			broadcastMatchID = matchID
			broadcastMessage = message
		},
	)

	message, err := service.GenerateDateMessage(context.Background(), 8)
	if err != nil {
		t.Fatalf("GenerateDateMessage returned error: %v", err)
	}
	if message == nil || message.Message != "date idea" {
		t.Fatalf("unexpected generated message: %+v", message)
	}
	if chatStore.appendMatch != 8 || chatStore.appended.Message != "date idea" {
		t.Fatalf("expected generated message to be appended, got match=%d message=%+v", chatStore.appendMatch, chatStore.appended)
	}
	if broadcastMatchID != 8 || broadcastMessage == nil || broadcastMessage.Message != "date idea" {
		t.Fatalf("expected generated message to be broadcast, got match=%d message=%+v", broadcastMatchID, broadcastMessage)
	}
}

func TestChatMessageServiceGenerateVibeMessageReturnsGeneratorError(t *testing.T) {
	chatStore := &fakeChatHistoryStore{loadErr: store.ErrNotFound}
	service := NewChatMessageService(
		NewChatService(&fakeMatchStore{}, chatStore, nil, nil),
		func(matchID uint) (domain.Match, error) { return domain.Match{}, nil },
		nil,
		nil,
		func(match domain.Match, recentMessages []domain.ChatMessage) (*domain.ChatMessage, error) {
			return nil, errors.New("generator failed")
		},
		nil,
	)

	_, err := service.GenerateVibeMessage(context.Background(), 3)
	if err == nil || err.Error() != "generator failed" {
		t.Fatalf("expected generator error, got %v", err)
	}
}
