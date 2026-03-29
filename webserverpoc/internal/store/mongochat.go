package store

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoChatStore struct {
	collection *mongo.Collection
}

func NewMongoChatStore(client *mongo.Client, databaseName string) *MongoChatStore {
	return &MongoChatStore{
		collection: client.Database(databaseName).Collection("chathistory"),
	}
}

func (s *MongoChatStore) LoadByMatchID(ctx context.Context, matchID uint) (Conversation, error) {
	var conversation Conversation
	err := s.collection.FindOne(ctx, bson.M{"match_id": matchID}).Decode(&conversation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Conversation{MatchID: matchID, Messages: []ChatMessage{}}, ErrNotFound
		}
		return Conversation{}, err
	}
	return conversation, nil
}

func (s *MongoChatStore) SaveConversation(ctx context.Context, conversation Conversation) error {
	_, err := s.collection.ReplaceOne(
		ctx,
		bson.M{"match_id": conversation.MatchID},
		conversation,
		options.Replace().SetUpsert(true),
	)
	return err
}

func (s *MongoChatStore) AppendMessage(ctx context.Context, matchID uint, message ChatMessage) error {
	current, err := s.LoadByMatchID(ctx, matchID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if errors.Is(err, ErrNotFound) {
		current = Conversation{
			MatchID:  matchID,
			Messages: []ChatMessage{message},
		}
	} else {
		current.Messages = append(current.Messages, message)
	}

	return s.SaveConversation(ctx, current)
}

func (s *MongoChatStore) LatestMessageID(ctx context.Context, matchID uint) (int64, error) {
	conversation, err := s.LoadByMatchID(ctx, matchID)
	if err != nil {
		return 0, err
	}
	if len(conversation.Messages) == 0 {
		return 0, ErrNotFound
	}
	return conversation.Messages[len(conversation.Messages)-1].ID, nil
}

func (s *MongoChatStore) SaveReadState(ctx context.Context, state MatchReadState) error {
	return nil
}
