package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {

	router := gin.Default()

	router.POST("/incoming-call", incomingCallHandler)

	router.Run(":8080")
}

// Handles Twilio Webhook for incoming calls
func incomingCallHandler(c *gin.Context) {

	// Parses the incoming Twilio request

	msg := &api.CreateMessageParams{}
	err := c.Bind(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
