package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go/twiml"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
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

	client := cohereclient.NewClient(cohereclient.WithToken(os.Getenv("COHERE_API_KEY")))

	// Temporary Dataset
	dataset, err := client.Datasets.Create(
		context.TODO(),
		&MyReader{Reader: strings.NewReader(""), name: "test.jsonl"},
		&MyReader{Reader: strings.NewReader(""), name: "a.jsonl"},
		&cohere.DatasetsCreateRequest{
			Name: "prompt-completion-dataset",
			Type: cohere.DatasetTypeEmbedResult,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", dataset)

	datasetResult, err := client.Datasets.Get(context.TODO(), "prompt-completion-dataset")

	log.Printf("Dataset: %v", datasetResult)

	// router := gin.Default()

	// router.POST("/answer", incomingCallHandler)
	// router.POST("/handle-user-input", handleUserInput)

	// router.Run(":8080")
}

// Handles Twilio Webhook for incoming calls
func incomingCallHandler(c *gin.Context) {
	msg := "Welcome! This is a test call, Try to ask me questions and I will try to answer them. Example, ask me how the weather is like in your location."

	gather := &twiml.VoiceGather{
		Input:         "speech",
		Language:      "en-US",
		Action:        "/handle-user-input",
		SpeechTimeout: "1",
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

type MyReader struct {
	io.Reader
	name string
}

func (m *MyReader) Name() string {
	return m.name
}

// Handles Twilio Webhook for user input
func handleUserInput(c *gin.Context) {
	userInput := c.PostForm("SpeechResult")

	client := cohereclient.NewClient(cohereclient.WithToken(os.Getenv("COHERE_API_KEY")))

	// Temporary Dataset
	dataset, err := client.Datasets.Create(
		context.TODO(),
		&MyReader{Reader: strings.NewReader(""), name: "test.jsonl"},
		&MyReader{Reader: strings.NewReader(""), name: "a.jsonl"},
		&cohere.DatasetsCreateRequest{
			Name: "prompt-completion-dataset",
			Type: cohere.DatasetTypeEmbedResult,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", dataset)

	datasetResult, err := client.Datasets.Get(context.TODO(), "prompt-completion-dataset")

	log.Printf("Dataset: %v", datasetResult)

	response, err := client.Chat(
		context.TODO(),
		&cohere.ChatRequest{
			Message: userInput,
		},
	)

	say := &twiml.VoiceSay{
		Message: response.Text,
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
