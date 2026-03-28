package gowebserver

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/edwardmccormick/gowebserver/internal/service"
)

var (
	authService        *service.AuthService
	chatService        *service.ChatService
	matchService       *service.MatchService
	profileService     *service.ProfileService
	chatMessageService *service.ChatMessageService
)

func initializeServices() error {
	if userStore == nil || matchStore == nil {
		return errors.New("stores are not initialized")
	}

	var s3Client *s3.S3
	if region := os.Getenv("AWS_REGION"); region != "" {
		awsSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err == nil {
			s3Client = s3.New(awsSession)
		}
	}

	authService = service.NewAuthService(userStore)
	matchService = service.NewMatchService(db)
	profileService = service.NewProfileService(db, s3Client, os.Getenv("AWS_S3_BUCKET"))
	if chatHistoryStore != nil {
		chatService = service.NewChatService(matchStore, chatHistoryStore, realtimeHub, notificationCenter)
		chatMessageService = service.NewChatMessageService(
			chatService,
			matchService.GetMatchDetails,
			CreateInitialChatMessage,
			CreateDateSuggestionMessage,
			CreateVibeChatMessage,
			broadcastMessageToMatch,
		)
	}

	return nil
}
