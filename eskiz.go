package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	Message   string `json:"message"`
	Data      struct {
		Token   string `json:"token"`
	} `json:"data"`
	TokenType string `json:"token_type"`
}

var token string

func LoginHandler(c *gin.Context) {
    url := "https://notify.eskiz.uz/api/auth/login"
    method := "POST"

    email := c.PostForm("email")
    password := c.PostForm("password")

    payload := &bytes.Buffer{}
    writer := multipart.NewWriter(payload)
    _ = writer.WriteField("email", email)
    _ = writer.WriteField("password", password)
    _ = writer.Close()

    client := &http.Client{}
    req, err := http.NewRequest(method, url, payload)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    res, err := client.Do(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var loginResponse LoginResponse
    if err := json.Unmarshal(body, &loginResponse); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    token = loginResponse.Data.Token

    c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}


func SendSMSHandler(c *gin.Context) {
	url := "https://notify.eskiz.uz/api/message/sms/send"
	method := "POST"

	// mobile_phone := c.PostForm("mobile_phone")
	// message := c.PostForm("message")
  	// from := c.PostForm("from")
	// callback_url := c.PostForm("callback_url")

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("mobile_phone", phoneNumbers[0])
	_ = writer.WriteField("message", "This is test from Eskiz")
	_ = writer.WriteField("from", "4546")
	_ = writer.WriteField("callback_url", "localhost:8080/api/message/sms/send")
	err := writer.Close()

	// _ = writer.WriteField("mobile_phone", mobile_phone)
	// _ = writer.WriteField("message", message)
	// _ = writer.WriteField("from", from)
	// _ = writer.WriteField("callback_url", callback_url)
	// err := writer.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(res.StatusCode, res.Header.Get("Content-Type"), body)
}