package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	twilio "github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Printf("Error loading credentials: %v", envErr)
	}

	// routes declaration
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/users", getUsers)
	router.GET("/user/:id", getUserId)
	router.POST("/sendmessage", sendMessage)
	router.GET("/listmessages", listMessages)
	router.Run()
}

var users = []User{{Id: 0, First_name: "Mayank", Last_name: "Sonu", Phone_number: "+917037414934"},
	{Id: 1, First_name: "Chetan", Last_name: "Rana", Phone_number: "+917037414934"},
	{Id: 3, First_name: "Zeeshan", Last_name: "Zama", Phone_number: "+917037414934"},
}

type User struct {
	Id           int    `json:"id"`
	First_name   string `json:"firstName"`
	Last_name    string `json:"lastName"`
	Phone_number string `json:"phoneNumber"`
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}
func getUserId(c *gin.Context) {
	response := make(map[string]interface{})
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		response["Error"] = "Error during conversion"
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, users[userId])
}

func sendMessage(c *gin.Context) {

	var newMessage messageFields
	response := make(map[string]interface{})

	// Getting params from query
	if err := c.BindJSON(&newMessage); err != nil {
		fmt.Println(err)
		response["Error"] = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// validate object
	if newMessage.From == "" || newMessage.Message == "" || newMessage.To == "" {
		response["Error"] = "Invalid Properties"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWLIO_AUTH_TOKEN")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &api.CreateMessageParams{}
	params.SetFrom(newMessage.From)
	params.SetBody(newMessage.Message)
	params.SetTo(newMessage.To)
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		response["Error"] = "Server Error"
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		fmt.Println(resp.Sid)
		response["Message"] = "Message Sent"
		response["Data"] = resp
		c.JSON(http.StatusOK, response)
	}
}

func listMessages(c *gin.Context) {
	response := make(map[string]interface{})
	client := twilio.NewRestClient()

	params := &api.ListMessageParams{}
	params.SetLimit(20)
	resp, err := client.Api.ListMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		response["Error"] = "Server Error"
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		response["Messages"] = resp
		fmt.Println("Response: ", response)
		c.JSON(http.StatusOK, response)
	}
}

type messageFields struct {
	From    string `json:"sender"`
	To      string `json:"receiver"`
	Message string `json:"message"`
}
