package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go/twiml"
)

func loadEnv() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {

	// Load Environment Variables
	loadEnv()

	router := gin.Default()

	router.POST("/answer", incomingCallHandler)

	router.Run(":8080")
}

// Handles Twilio Webhook for incoming calls
func incomingCallHandler(c *gin.Context) {
	say := &twiml.VoiceSay{
		Message: "This is a test message. Shernan is a really cool guy wink wink",
	}

	twimlResult, err := twiml.Voice([]twiml.Element{say})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Header("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}
}
