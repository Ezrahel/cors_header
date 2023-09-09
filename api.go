package main

import (
	//"encoding/json"
	"net/http"
	"time"
	"os"
	"log"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Slack          string `json:"slack_name"`
	DayOfWeek      string `json:"current_day"`
	CurrentUTCTime string `json:"utc_time"`
	Track          string `json:"track"`
	Github_file_url string `json:"github_file_url"`
	Github_repo_url string `json:"github_repo_url"`
	StatusMessage  int    `json:"status_code"`
}

func infoHandler(c *gin.Context) {
    // Parse query parameters
    slack := c.DefaultQuery("slack_name", "") // Use DefaultQuery to provide a default value if the parameter is missing
    track := c.DefaultQuery("track", "")

    // Get the current day of the week
    dayOfWeek := time.Now().Weekday().String()

    // Get the current UTC time in Nigeria
    currentTime := nigeriaTime()

    // Calculate the time difference in hours between Nigeria time and UTC time
    timeDiff := currentTime.Sub(time.Now().UTC()).Hours()

    // Determine the HTTP status code based on the time difference
    var statusCode int
    if timeDiff >= -2 && timeDiff <= 2 {
        statusCode = http.StatusOK
    } else {
        statusCode = http.StatusInternalServerError
    }

    // Create the response struct
    response := Response{
        Slack:          slack,
        DayOfWeek:      dayOfWeek,
        CurrentUTCTime: currentTime.Format("2006-01-02 15:04:05"),
        Track:          track,
        Github_file_url: "https://github.com/Ezrahel/Go_API/blob/main/GoEndpoint.go",
        Github_repo_url: "https://github.com/Ezrahel/Go_API",
        StatusMessage:  statusCode,
    }

    // Set the content type header to JSON
    c.Header("Content-Type", "application/json")

    // Set the HTTP status code
    c.JSON(statusCode, response)
}

func nigeriaTime() time.Time {
	// Create a fixed time zone for Nigeria (UTC+1)
	nigeriaTimeZone := time.FixedZone("UTC+1", 3600)
	return time.Now().In(nigeriaTimeZone)
}


func CORSMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router:= gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/api", infoHandler)
	//http.HandleFunc("/api", infoHandler)
	port := os.Getenv("PORT")
	if port == ""{
		port ="8080"
	}

	if err := router.Run(":"+ port); err!=nil{
		log.Panicf("error: %s", err)
	}


}

