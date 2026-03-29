package httpapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Signup                   gin.HandlerFunc
	Login                    gin.HandlerFunc
	Logout                   gin.HandlerFunc
	PostPeople               gin.HandlerFunc
	GetUsers                 gin.HandlerFunc
	GetPeople                gin.HandlerFunc
	SearchPeople             gin.HandlerFunc
	GetPeopleByLocation      gin.HandlerFunc
	GetPeopleByID            gin.HandlerFunc
	GetPhotosByID            gin.HandlerFunc
	GetMatches               gin.HandlerFunc
	GetMatchByPersonID       gin.HandlerFunc
	PostMatch                gin.HandlerFunc
	WebsocketListener        gin.HandlerFunc
	GetFaviconIco            gin.HandlerFunc
	GreetUser                gin.HandlerFunc
	GreetUserByName          gin.HandlerFunc
	ChatMessagesFromSQL      gin.HandlerFunc
	ChatMessagesFromMongo    gin.HandlerFunc
	GetChatMessagesForMatch  gin.HandlerFunc
	GenerateVibeChatForMatch gin.HandlerFunc
	GenerateDateSuggestion   gin.HandlerFunc
	SSEHandler               gin.HandlerFunc
	MarkMessagesAsRead       gin.HandlerFunc
	AdminGetAllUsers         gin.HandlerFunc
	AdminGetAllChats         gin.HandlerFunc
	AdminGetRecentUpdates    gin.HandlerFunc
	AdminSetUserAdmin        gin.HandlerFunc
	JwtMiddleware            gin.HandlerFunc
	AdminMiddleware          gin.HandlerFunc
}

func NewRouter(h Handlers) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5172"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/signup", h.Signup)
	router.POST("/login", h.Login)
	router.GET("/logout", h.Logout)
	router.POST("/people", h.PostPeople)
	router.GET("/users", h.GetUsers)

	router.GET("/people", h.GetPeople)
	router.POST("/people/search", h.JwtMiddleware, h.SearchPeople)
	router.POST("/peoplelocation", h.GetPeopleByLocation)
	router.GET("/people/:id", h.GetPeopleByID)
	router.GET("/photos/:id", h.GetPhotosByID)

	router.GET("/matches", h.GetMatches)
	router.GET("/matches/:id", h.GetMatchByPersonID)
	router.POST("/matches", h.PostMatch)
	router.GET("/ws", h.WebsocketListener)

	router.GET("/favicon.ico", h.GetFaviconIco)
	router.GET("/", h.GreetUser)
	router.GET("/greet/:name", h.GreetUserByName)
	router.GET("/chat", h.ChatMessagesFromSQL)
	router.GET("/chat/:id", h.ChatMessagesFromMongo)
	router.GET("/chatmessages/:id", h.GetChatMessagesForMatch)

	router.POST("/vibechat/:id", h.JwtMiddleware, h.GenerateVibeChatForMatch)
	router.POST("/perfectdate/:id", h.JwtMiddleware, h.GenerateDateSuggestion)
	router.GET("/notifications/:id", h.JwtMiddleware, h.SSEHandler)
	router.POST("/chat/markread/:id", h.JwtMiddleware, h.MarkMessagesAsRead)

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(h.AdminMiddleware)
	adminRoutes.GET("/users", h.AdminGetAllUsers)
	adminRoutes.GET("/chats", h.AdminGetAllChats)
	adminRoutes.GET("/recent", h.AdminGetRecentUpdates)
	adminRoutes.POST("/set-admin", h.AdminSetUserAdmin)

	return router
}
