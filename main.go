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
	router.POST("/handle-user-input", handleUserInput)

	router.Run(":8080")
}

// Handles Twilio Webhook for incoming calls
func incomingCallHandler(c *gin.Context) {
	msg := "Welcome! This is a test call, say anything and I will say it back to you"

	gather := &twiml.VoiceGather{
		Input:    "speech",
		Language: "en-US",
		Action:   "/handle-user-input",
		Timeout:  "3",
	}

	say := &twiml.VoiceSay{
		Message: msg,
	}

	twimlResult, err := twiml.Voice([]twiml.Element{say, gather})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Header("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}
}

// Handles Twilio Webhook for user input
func handleUserInput(c *gin.Context) {
	userInput := c.PostForm("SpeechResult")

	// Process the user's input
	responseMsg := "You said: " + userInput

	say := &twiml.VoiceSay{
		Message: responseMsg,
	}

	// Send the TwiML response back to Twilio
	twimlResult, err := twiml.Voice([]twiml.Element{say})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Header("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}
}
