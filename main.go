package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go/twiml"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	"github.com/cohere-ai/cohere-go/v2/finetuning"
)

func loadEnv() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type MyReader struct {
	io.Reader
	name string
}

func (m *MyReader) Name() string {
	return m.name
}
func main() {
	// Load Environment Variables
	loadEnv()

	co := client.NewClient(client.WithToken(os.Getenv("COHERE_API_KEY")))

	// Read dataset file
	datasetFile := "test.jsonl"
	datasetContent, err := os.ReadFile(datasetFile)
	if err != nil {
		log.Fatalf("Failed to read dataset file: %v", err)
	}

	// Create dataset
	resp, err := co.Datasets.Create(
		context.TODO(),
		&MyReader{Reader: strings.NewReader(string(datasetContent))},
		&MyReader{Reader: strings.NewReader("")},
		&cohere.DatasetsCreateRequest{
			Name: "test-dataset",
			Type: cohere.DatasetTypeSingleLabelClassificationFinetuneInput,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create dataset: %v", err)
	}
	log.Printf("Dataset created: %+v", resp)

	// Wait for dataset validation
	for {
		datasetStatus, err := co.Datasets.Get(context.TODO(), *resp.Id)
		if err != nil {
			log.Fatal("Error getting dataset status:", err)
		}

		if datasetStatus.Dataset.ValidationStatus == cohere.DatasetValidationStatusValidated {
			break
		}

		log.Printf("Dataset is still processing. Waiting for validation...")
		time.Sleep(5 * time.Second)
	}

	// Test question
	question := "What car does Shernan drive?"

	// Create a finetuned model
	modelRes, err := co.Finetuning.CreateFinetunedModel(
		context.TODO(),
		&finetuning.FinetunedModel{
			Name: "test-finetuned-model",
			Settings: &finetuning.Settings{
				DatasetId: *resp.Id,
				BaseModel: &finetuning.BaseModel{
					BaseType: finetuning.BaseTypeBaseTypeClassification,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create finetuned model: %v", err)
	}
	log.Printf("Finetuned model created: %+v", modelRes.FinetunedModel)

	// Chat with the finetuned model
	testRes, err := co.ChatStream(
		context.TODO(),
		&cohere.ChatStreamRequest{
			Model:   modelRes.FinetunedModel.Id,
			Message: question,
		},
	)

	if err != nil {
		log.Fatalf("Failed to start chat stream: %v", err)
	}

	defer testRes.Close()

	for {
		message, err := testRes.Recv()

		if errors.Is(err, io.EOF) {
			// An io.EOF error means the server is done sending messages
			// and should be treated as a success.
			break
		}

		if err != nil {
			log.Fatal("Error receiving message:", err)
		}

		if message.TextGeneration != nil {
			log.Printf("Received message: %+v", message.TextGeneration)
		}
	}

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
