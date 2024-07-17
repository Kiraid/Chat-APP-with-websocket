// routes/register.go
package routes

import (
	"log"
	"net/http"
	"strconv"
	"example.com/chat/models"
	"github.com/gin-gonic/gin"
)


func getEvents(context *gin.Context) {
	events, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Couldnt fetch events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	log.Printf("User data: %+v\n", user)

	if err := user.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}


func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := models.Authenticate(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	log.Printf("Setting user ID %d in the cookie", user.ID)
	c.SetCookie("user_id", strconv.Itoa(user.ID), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/chat")
}

func GetUsername(c *gin.Context) {
    userIDstr, err := c.Cookie("user_id")
    if err != nil {
        log.Printf("Error retrieving user ID from cookie: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    userID, err := strconv.Atoi(userIDstr)
    if err != nil {
        log.Printf("Error converting user ID to integer: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := models.GetUserByID(userID)
    if err != nil {
        log.Printf("Error fetching user details for user ID %d: %v", userID, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "username": user.Username})
}