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

	// Parses the incoming Twilio request

	// var msg = &api.CreateMessageParams{}
	// err := c.Bind(&msg)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Test Message

	// testMessage := "Hello, this is a test message from Mockcim."

	// accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	// authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	// phoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: accountSid,
	// 	Password: authToken,
	// })

	// params := &api.CreateMessageParams{}
	// params.SetFrom(phoneNumber)
	// params.SetBody(testMessage + "\n Shernan is also really cool and you should hire him if you're interested in hiring a software engineer.")

	// // Send the message
	// _, err = client.Api.CreateMessage(params)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	say := &twiml.VoiceSay{
		Message: "Hello, this is a test message from Mockcim. Shernan is also really cool and you should hire him if you're interested in hiring a software engineer.",
	}

	twimlResult, err := twiml.Voice([]twiml.Element{say})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.Header("Content-Type", "text/xml")
		c.String(http.StatusOK, twimlResult)
	}
}
