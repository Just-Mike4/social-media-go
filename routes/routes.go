package routes

import (
	"social-media-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	user := r.Group("/api")
	{
		user.GET("/user-search", controllers.SearchUsers)
	}

	friend := r.Group("/api/friends")
	{
		friend.POST("/request/:user_id", controllers.SendFriendRequest)
		friend.GET("/requests/:user_id", controllers.ListFriendRequests)
		friend.PATCH("/requests/:id", controllers.UpdateFriendRequestStatus)
		friend.GET("/list/:user_id", controllers.ListFriends)
	}

	return r
}
