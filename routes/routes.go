package routes

import (

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.GET("/users", getEvents)
	router.GET("/register", RegisterPage)
	router.GET("/api/get-username", GetUsername)
	router.POST("/register", Register)
	router.GET("/login", LoginPage)
	router.POST("/login", Login)
	router.GET("/chat", ChatPage)
	router.POST("/create-channel", CreateChannel)
	router.GET("/channels/:channel/messages", GetMessages)
	router.POST("/channels/:channel/messages", PutMessage)
	router.POST("/join-channel", joinChannel)

}


