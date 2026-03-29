package gowebserver

import (
	"fmt"
	"os"

	"github.com/edwardmccormick/gowebserver/internal/httpapi"
	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
}

func LoadRuntimeConfig() (*Config, error) {
	if configPath := os.Getenv("CONFIG_FILE"); configPath != "" {
		return LoadConfig(configPath)
	}

	configPath := "./configlocal.json"
	if isRunningInDockerContainer() {
		configPath = "./config.json"
	}

	return LoadConfig(configPath)
}

func NewApp(config *Config) (*App, error) {
	var err error

	db, err = ConnectToPostgresWithConfig(config)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	db = db.Debug()

	if err := MigrateSchema(db); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	if err := InitializeAIClient(); err != nil {
		return nil, fmt.Errorf("initialize ai client: %w", err)
	}

	mongoClient, err = ConnectToMongoDBWithConfig(config)
	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	if err := initializeStores(); err != nil {
		return nil, fmt.Errorf("initialize stores: %w", err)
	}
	if err := initializeServices(); err != nil {
		return nil, fmt.Errorf("initialize services: %w", err)
	}

	return &App{router: newRouter()}, nil
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func newRouter() *gin.Engine {
	return httpapi.NewRouter(httpapi.Handlers{
		Signup:                   Signup,
		Login:                    Login,
		Logout:                   Logout,
		PostPeople:               PostPeople,
		GetUsers:                 GetUsers,
		GetPeople:                GetPeople,
		SearchPeople:             SearchPeople,
		GetPeopleByLocation:      GetPeopleByLocation,
		GetPeopleByID:            GetPeopleByID,
		GetPhotosByID:            GetPhotosByID,
		GetMatches:               GetMatches,
		GetMatchByPersonID:       GetMatchByPersonID,
		PostMatch:                PostMatch,
		WebsocketListener:        WebsocketListener,
		GetFaviconIco:            GetFaviconIco,
		GreetUser:                GreetUser,
		GreetUserByName:          GreetUserByName,
		ChatMessagesFromSQL:      ChatMessagesFromSQL,
		ChatMessagesFromMongo:    ChatMessagesFromMongo,
		GetChatMessagesForMatch:  GetChatMessagesForMatch,
		GenerateVibeChatForMatch: GenerateVibeChatForMatch,
		GenerateDateSuggestion:   GenerateDateSuggestionForMatch,
		SSEHandler:               SSEHandler,
		MarkMessagesAsRead:       MarkMessagesAsRead,
		AdminGetAllUsers:         AdminGetAllUsers,
		AdminGetAllChats:         AdminGetAllChats,
		AdminGetRecentUpdates:    AdminGetRecentUpdates,
		AdminSetUserAdmin:        AdminSetUserAdmin,
		JwtMiddleware:            JwtMiddleware,
		AdminMiddleware:          AdminMiddleware,
	})
}
