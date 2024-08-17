package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ConvertStringToObjectId(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func ConvertObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func NewError(message string) error {
	return errors.New(message)
}

func HandleError(err error, message string) error {
	if err == nil && message != ""{
		return NewError(message)
	}

	if err == nil {
		return nil
	}

	if message == "" {
		message = "Data not found"
	}

	if err.Error() == "mongo: no documents in result" {
		return NewError(message)
	}

	return err
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); if err != nil {
		return ""
	}

	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println(err)
	return err == nil
}

type PromptInput struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type GPTInput struct {
	Model          string        `json:"model"`
	Messages       []PromptInput `json:"messages"`
	Temperature    float64       `json:"temperature"`
}

type Input struct {
	Text string `json:"text"`
}

type Voice struct {
	LanguageCode string `json:"languageCode"`
	Name string `json:"name"`
	SsmlGender string `json:"ssmlGender"`
}

type AudioConfig struct {
	AudioEncoding string `json:"audioEncoding"`
}

type TTSInput struct {
	Input       	Input	     `json:"input"`
	Voice       	Voice 		 `json:"voice"`
	AudioConfig 	AudioConfig  `json:"audioConfig"`
}

func ComputeGPTPrompt(message string) GPTInput {
	return GPTInput{
		Model: "gpt-3.5-turbo", // Use the correct model
		Messages: []PromptInput{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: message,
			},
		},
		Temperature: 0.7,
	}
}

func ComputeTTSInput(message string) TTSInput {
	return TTSInput{
		Input: Input{
			Text: message,
		},
		Voice: Voice{
			LanguageCode: "en-US",
			Name: "en-GB-Standard-A",
			SsmlGender: "FEMALE",
		},
		AudioConfig: AudioConfig{
			AudioEncoding: "MP3",
		},
	}
}

func SaveBase64ToDisk(encodedString, filePath string) error {
	// Decode the base64 encoded string into bytes
	data, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Write the decoded data to the specified file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Println("File saved successfully to:", filePath)
	return nil
}
