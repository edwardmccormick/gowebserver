package service

import (
	"context"
	"errors"
	"sort"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/edwardmccormick/gowebserver/internal/store"
)

type IntroGenerator func(match domain.Match) (*domain.ChatMessage, error)
type AIMessageGenerator func(match domain.Match, recentMessages []domain.ChatMessage) (*domain.ChatMessage, error)
type MessageBroadcaster func(matchID uint, message *domain.ChatMessage)
type MatchLoader func(matchID uint) (domain.Match, error)

type ChatMessageService struct {
	chat          *ChatService
	loadMatch     MatchLoader
	intro         IntroGenerator
	dateGenerator AIMessageGenerator
	vibeGenerator AIMessageGenerator
	broadcast     MessageBroadcaster
}

func NewChatMessageService(
	chat *ChatService,
	loadMatch MatchLoader,
	intro IntroGenerator,
	dateGenerator AIMessageGenerator,
	vibeGenerator AIMessageGenerator,
	broadcast MessageBroadcaster,
) *ChatMessageService {
	return &ChatMessageService{
		chat:          chat,
		loadMatch:     loadMatch,
		intro:         intro,
		dateGenerator: dateGenerator,
		vibeGenerator: vibeGenerator,
		broadcast:     broadcast,
	}
}

func (s *ChatMessageService) GetPagedMessages(ctx context.Context, matchID uint, limit, offset int) ([]domain.ChatMessage, int, error) {
	conversation, err := s.chat.LoadConversation(ctx, matchID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return []domain.ChatMessage{}, 0, nil
		}
		return nil, 0, err
	}

	messages := conversation.Messages
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})

	total := len(messages)
	start := offset
	end := offset + limit
	if start >= total {
		return []domain.ChatMessage{}, total, nil
	}
	if end > total {
		end = total
	}

	return messages[start:end], total, nil
}

func (s *ChatMessageService) EnsureIntroduction(ctx context.Context, matchID uint) (*domain.ChatMessage, error) {
	conversation, err := s.chat.LoadConversation(ctx, matchID)
	if err == nil && len(conversation.Messages) > 0 {
		return nil, nil
	}
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	match, err := s.loadMatch(matchID)
	if err != nil {
		return nil, err
	}

	message, err := s.intro(match)
	if err != nil {
		return nil, err
	}

	if err := s.chat.SaveConversation(ctx, domain.Conversation{
		MatchID:  matchID,
		Messages: []domain.ChatMessage{*message},
	}); err != nil {
		return message, nil
	}

	return message, nil
}

func (s *ChatMessageService) GenerateDateMessage(ctx context.Context, matchID uint) (*domain.ChatMessage, error) {
	conversation, err := s.chat.LoadConversation(ctx, matchID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			conversation = domain.Conversation{MatchID: matchID, Messages: []domain.ChatMessage{}}
		} else {
			return nil, err
		}
	}

	messages := conversation.Messages
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})

	match, err := s.loadMatch(matchID)
	if err != nil {
		return nil, err
	}

	message, err := s.dateGenerator(match, messages)
	if err != nil {
		return nil, err
	}

	_ = s.chat.AppendMessage(ctx, matchID, *message)
	if s.broadcast != nil {
		s.broadcast(matchID, message)
	}

	return message, nil
}

func (s *ChatMessageService) GenerateVibeMessage(ctx context.Context, matchID uint) (*domain.ChatMessage, error) {
	conversation, err := s.chat.LoadConversation(ctx, matchID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			conversation = domain.Conversation{MatchID: matchID, Messages: []domain.ChatMessage{}}
		} else {
			return nil, err
		}
	}

	messages := conversation.Messages
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Time.Before(messages[j].Time)
	})

	match, err := s.loadMatch(matchID)
	if err != nil {
		return nil, err
	}

	message, err := s.vibeGenerator(match, messages)
	if err != nil {
		return nil, err
	}

	_ = s.chat.AppendMessage(ctx, matchID, *message)
	if s.broadcast != nil {
		s.broadcast(matchID, message)
	}

	return message, nil
}
