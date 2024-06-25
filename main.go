package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/contacts", getPhoneNumbers)

	router.POST("/api/auth/login", LoginHandler)
 	
	router.POST("/api/message/sms/send", SendSMSHandler)

	router.Run(":8080")
}