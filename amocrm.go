package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var accessToken string = "AmoCRM token"
var accountSubdomain = "Subdomain"

type ContactResponse struct {
	Embedded struct {
		Contacts []Contact `json:"contacts"`
	} `json:"_embedded"`
}

type Contact struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	ResponsibleUserID  int    `json:"responsible_user_id"`
	CustomFieldsValues []struct {
		FieldID   int    `json:"field_id"`
		FieldName string `json:"field_name"`
		FieldType string `json:"field_type"`
		Values    []struct {
			Value    string `json:"value"`
			EnumID   int    `json:"enum_id"`
			EnumCode string `json:"enum_code"`
		} `json:"values"`
	} `json:"custom_fields_values"`
}

var phoneNumbers []string

func getPhoneNumbers(c *gin.Context) {
	phoneNumbers = nil

	contactsURL := fmt.Sprintf("https://%s.amocrm.ru/api/v4/contacts", accountSubdomain)

	client := &http.Client{}

	req, err := http.NewRequest("GET", contactsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var contactResponse ContactResponse
	err = json.Unmarshal(body, &contactResponse)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse contacts response"})
		return
	}

	for _, contact := range contactResponse.Embedded.Contacts {
		for _, field := range contact.CustomFieldsValues {
			if field.FieldName == "Телефон" {
				for _, value := range field.Values {
					phoneNumbers = append(phoneNumbers, value.Value)
				}
			}
		}
	}

	for i, number := range phoneNumbers {
		phoneNumbers[i] = strings.TrimPrefix(number, "+")
	}

	c.JSON(http.StatusOK, gin.H{"phone_numbers": phoneNumbers})
}
