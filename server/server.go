package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	router := gin.Default()

	// Use the CORS middleware
	router.Use(cors.Default())

	// GET route for root URL
	router.GET("/", get_root)

	// POST route for /send-sms
	router.POST("/send-sms", sendSMS)

	// Start the HTTP server on port 8080
	router.Run("localhost:8080")
}

func get_root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "HTTP is Online!",
	})
}

func sendSMS(c *gin.Context) {
	var smsReq SMSRequest
	if err := c.ShouldBindJSON(&smsReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(smsReq.To)
	params.SetFrom(fromPhone)
	params.SetBody(smsReq.Message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending SMS: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS sent successfully"})
}
