package user

import (
	"github.com/Babatunde13/golangr/api/database"
	"github.com/gin-gonic/gin"
)

type Login struct {
	User database.User `json:"user"`
	Token string `json:"token"`
}

type SuggestRequest struct {
	Message string `json:"message"`
}

type IUserController interface {
	GetUsers() gin.HandlerFunc
	CreateUser() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	LoginUser() gin.HandlerFunc
	AuthMiddleware() gin.HandlerFunc
	AuthUser() gin.HandlerFunc
	SuggestWithGPT() gin.HandlerFunc
	TextToSpeech() gin.HandlerFunc
	Upload() gin.HandlerFunc
	Download() gin.HandlerFunc
}

type UserController struct{}

type OpenAIResponse struct {
	ID               string       `json:"id"`
	Object           string       `json:"object"`
	Created          float64      `json:"created"`
	Model            string       `json:"model"`
	Choices          []Choice     `json:"choices"`
	Usage            UsageDetails `json:"usage,omitempty"`
	SystemFingerprint string      `json:"system_fingerprint,omitempty"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
}

type ChatMessage struct {
	Content  string  `json:"content"`
	Role     string  `json:"role"`
	Refusal  string `json:"refusal,omitempty"`
}

type UsageDetails struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type TTSResponse struct {
	AudioContent	string	`json:"audioContent"`
}