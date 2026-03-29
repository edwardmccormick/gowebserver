package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/edwardmccormick/gowebserver/internal/domain"
	"github.com/edwardmccormick/gowebserver/internal/realtime"
	"github.com/edwardmccormick/gowebserver/internal/store"
)

type ChatService struct {
	matches store.MatchStore
	chat    store.ChatHistoryStore
	hub     *realtime.Hub
	notify  *realtime.NotificationCenter
}

func NewChatService(matches store.MatchStore, chat store.ChatHistoryStore, hub *realtime.Hub, notify *realtime.NotificationCenter) *ChatService {
	return &ChatService{
		matches: matches,
		chat:    chat,
		hub:     hub,
		notify:  notify,
	}
}

func (s *ChatService) UpdateUnreadCounts(ctx context.Context, matchID uint, senderID uint) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	var recipientID uint
	if senderID == match.Offered {
		recipientID = match.Accepted
		match.UnreadAccepted++
	} else {
		recipientID = match.Offered
		match.UnreadOffered++
	}
	match.LastMessageTime = time.Now()

	if err := s.matches.Update(ctx, match); err != nil {
		return err
	}

	if s.hub != nil && s.notify != nil && !s.hub.HasActiveConnection(recipientID) {
		unreadCount := match.UnreadAccepted
		if recipientID == match.Offered {
			unreadCount = match.UnreadOffered
		}

		s.notify.Broadcast(
			uintToString(recipientID),
			realtime.NotificationEvent{
				Type:      realtime.NotificationTypeMessage,
				MatchID:   matchID,
				Count:     unreadCount,
				Timestamp: time.Now(),
			},
		)
	}

	return nil
}

func (s *ChatService) ResetUnreadCount(ctx context.Context, matchID uint, userID uint) error {
	match, err := s.matches.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	if userID == match.Offered {
		match.UnreadOffered = 0
	} else if userID == match.Accepted {
		match.UnreadAccepted = 0
	} else {
		return nil
	}

	return s.matches.Update(ctx, match)
}

func (s *ChatService) LoadConversation(ctx context.Context, matchID uint) (domain.Conversation, error) {
	conversation, err := s.chat.LoadByMatchID(ctx, matchID)
	if err != nil {
		return domain.Conversation{}, err
	}

	result := domain.Conversation{
		MatchID:  conversation.MatchID,
		Messages: make([]domain.ChatMessage, 0, len(conversation.Messages)),
	}
	for _, msg := range conversation.Messages {
		result.Messages = append(result.Messages, domain.ChatMessage{
			ID:          msg.ID,
			MatchID:     msg.MatchID,
			Time:        msg.Time,
			Who:         msg.Who,
			MessageType: msg.MessageType,
			Message:     msg.Message,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		})
	}

	return result, nil
}

func (s *ChatService) SaveConversation(ctx context.Context, conversation domain.Conversation) error {
	payload := store.Conversation{
		MatchID:  conversation.MatchID,
		Messages: make([]store.ChatMessage, 0, len(conversation.Messages)),
	}
	for _, msg := range conversation.Messages {
		payload.Messages = append(payload.Messages, store.ChatMessage{
			ID:          msg.ID,
			MatchID:     msg.MatchID,
			Time:        msg.Time,
			Who:         msg.Who,
			MessageType: msg.MessageType,
			Message:     msg.Message,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
		})
	}
	return s.chat.SaveConversation(ctx, payload)
}

func (s *ChatService) AppendMessage(ctx context.Context, matchID uint, message domain.ChatMessage) error {
	return s.chat.AppendMessage(ctx, matchID, store.ChatMessage{
		ID:          message.ID,
		MatchID:     message.MatchID,
		Time:        message.Time,
		Who:         message.Who,
		MessageType: message.MessageType,
		Message:     message.Message,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	})
}

func (s *ChatService) MarkConversationRead(ctx context.Context, matchID uint, userID uint) error {
	lastMessageID, err := s.chat.LatestMessageID(ctx, matchID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return err
	}

	if err := s.ResetUnreadCount(ctx, matchID, userID); err != nil {
		return err
	}

	if errors.Is(err, store.ErrNotFound) {
		return nil
	}

	return s.chat.SaveReadState(ctx, store.MatchReadState{
		MatchID:           matchID,
		UserID:            userID,
		LastReadMessageID: lastMessageID,
		UpdatedAt:         time.Now(),
	})
}

func uintToString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
