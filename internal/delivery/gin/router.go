package gin

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/Shuhrat55/auth/pkg/client"
	"github.com/Shuhrat55/forum/internal/delivery/gin/handler"
	"github.com/Shuhrat55/forum/internal/usecase"
	"github.com/Shuhrat55/forum/pkg/wsserver"

	_ "github.com/Shuhrat55/forum/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(P usecase.PostUseCase, T usecase.ThreadUseCase, authClient *client.AuthClient, hub *wsserver.Hub) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	forumHandler := NewForumHandler(P, T)
	go hub.Run()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v2 group
	api := router.Group("/api/v2")
	{
		api.GET("/threads", forumHandler.GetAllThread)
		api.GET("/thread/:id", forumHandler.GetThreadByID)

		authGroup := api.Group("")
		authGroup.Use(handler.AuthMiddleware(authClient))
		{
			authGroup.POST("/threads", forumHandler.CreateThread)
			authGroup.POST("/threads/posts", forumHandler.CreatePost)

			authGroup.GET("/threads/user/:id", forumHandler.GetThreadsByUserID)
			authGroup.GET("/posts/user/:id", forumHandler.GetPostsByUserID)
			authGroup.GET("/thread/:id/posts", forumHandler.GetPostsByThreadID)

			authGroup.DELETE("/posts/:id", forumHandler.DeletePostByID)
			authGroup.DELETE("/threads/:id", forumHandler.DeleteTheadByID)

			authGroup.PUT("/threads", forumHandler.EditThread)

			api.GET("/ws/threads/:id", hub.ThreadChat)
		}
	}

	// API v1 group for backward compatibility
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/threads", forumHandler.GetAllThread)
		apiV1.GET("/thread/:id", forumHandler.GetThreadByID)

		authGroupV1 := apiV1.Group("")
		authGroupV1.Use(handler.AuthMiddleware(authClient))
		{
			authGroupV1.POST("/threads", forumHandler.CreateThread)
			authGroupV1.POST("/threads/posts", forumHandler.CreatePost)

			authGroupV1.GET("/threads/user/:id", forumHandler.GetThreadsByUserID)
			authGroupV1.GET("/posts/user/:id", forumHandler.GetPostsByUserID)
			authGroupV1.GET("/thread/:id/posts", forumHandler.GetPostsByThreadID)

			authGroupV1.DELETE("/posts/:id", forumHandler.DeletePostByID)
			authGroupV1.DELETE("/threads/:id", forumHandler.DeleteTheadByID)

			authGroupV1.PUT("/threads", forumHandler.EditThread)

			apiV1.GET("/ws/threads/:id", hub.ThreadChat)
		}
	}

	return router
}
