package bot

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streamverse/training-bot-service/pkg/ai"
)

type TrainingBot struct {
	aiClient ai.AIClient
	sessions map[string]*ChatSession
}

type ChatSession struct {
	ID        string
	UserID    string
	UserType  string
	Context   []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatRequest struct {
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	UserType  string `json:"userType"` // end-user, creator, admin, developer
	Message   string `json:"message"`
}

type ChatResponse struct {
	SessionID string    `json:"sessionId"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Actions   []Action  `json:"actions,omitempty"`
}

type Action struct {
	Type  string                 `json:"type"`
	Label string                 `json:"label"`
	Data  map[string]interface{} `json:"data"`
}

func NewTrainingBot(aiClient ai.AIClient) *TrainingBot {
	return &TrainingBot{
		aiClient: aiClient,
		sessions: make(map[string]*ChatSession),
	}
}

func (tb *TrainingBot) CreateSession(c *gin.Context) {
	var req struct {
		UserID   string `json:"userId"`
		UserType string `json:"userType"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := &ChatSession{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		UserType:  req.UserType,
		Context:   []Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Add system prompt based on user type
	systemPrompt := tb.getSystemPrompt(req.UserType)
	session.Context = append(session.Context, Message{
		Role:      "system",
		Content:   systemPrompt,
		Timestamp: time.Now(),
	})

	tb.sessions[session.ID] = session

	c.JSON(http.StatusOK, gin.H{
		"sessionId": session.ID,
		"message":   tb.getWelcomeMessage(req.UserType),
	})
}

func (tb *TrainingBot) HandleChat(c *gin.Context) {
	var req ChatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, exists := tb.sessions[req.SessionID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Add user message to context
	session.Context = append(session.Context, Message{
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	})

	// Generate response using AI
	response, err := tb.generateResponse(session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add assistant response to context
	session.Context = append(session.Context, Message{
		Role:      "assistant",
		Content:   response,
		Timestamp: time.Now(),
	})

	session.UpdatedAt = time.Now()

	// Detect if actions should be suggested
	actions := tb.detectActions(req.Message, session.UserType)

	c.JSON(http.StatusOK, ChatResponse{
		SessionID: session.ID,
		Message:   response,
		Timestamp: time.Now(),
		Actions:   actions,
	})
}

func (tb *TrainingBot) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	session, exists := tb.sessions[sessionID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (tb *TrainingBot) generateResponse(session *ChatSession) (string, error) {
	// Build conversation context
	var messages []ai.Message
	for _, msg := range session.Context {
		messages = append(messages, ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Call AI to generate response
	response, err := tb.aiClient.GenerateResponse(messages)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (tb *TrainingBot) getSystemPrompt(userType string) string {
	prompts := map[string]string{
		"end-user": `You are StreamVerse Training Assistant, an expert guide for end users of the StreamVerse streaming platform.
Your role is to help users learn how to:
- Navigate the platform and find content
- Create and manage profiles
- Use features like watchlists, downloads, and parental controls
- Troubleshoot common issues
- Optimize their viewing experience

Be friendly, patient, and provide step-by-step instructions with examples.`,

		"creator": `You are StreamVerse Creator Training Assistant, specialized in helping content creators succeed on the platform.
Your role is to help creators learn how to:
- Upload and manage content
- Use the Creator Studio
- Optimize content for discovery
- Understand analytics and audience insights
- Monetize their content
- Follow best practices for content quality

Be professional, encouraging, and provide actionable advice.`,

		"admin": `You are StreamVerse Admin Training Assistant, an expert in platform administration and management.
Your role is to help administrators learn how to:
- Manage users and permissions
- Configure platform settings
- Monitor system health and performance
- Handle content moderation
- Generate reports and analytics
- Manage subscriptions and billing

Be precise, technical, and security-conscious.`,

		"developer": `You are StreamVerse Developer Training Assistant, specialized in technical integration and development.
Your role is to help developers learn how to:
- Integrate with StreamVerse APIs
- Implement authentication and authorization
- Use SDKs for different platforms
- Handle streaming protocols and DRM
- Debug and troubleshoot integrations
- Follow best practices for performance and security

Be technical, detailed, and provide code examples.`,
	}

	if prompt, ok := prompts[userType]; ok {
		return prompt
	}

	return prompts["end-user"] // Default
}

func (tb *TrainingBot) getWelcomeMessage(userType string) string {
	messages := map[string]string{
		"end-user":  "👋 Welcome to StreamVerse! I'm your training assistant. I can help you learn how to use the platform, find your favorite content, and make the most of your streaming experience. What would you like to learn about?",
		"creator":   "👋 Welcome to StreamVerse Creator Training! I'm here to help you succeed as a content creator on our platform. Whether you're uploading your first video or optimizing for millions of views, I've got you covered. What can I help you with?",
		"admin":     "👋 Welcome to StreamVerse Admin Training. I'm your guide for platform administration and management. I can help you with user management, system configuration, analytics, and more. How can I assist you today?",
		"developer": "👋 Welcome to StreamVerse Developer Training. I'm here to help you integrate with our platform APIs, SDKs, and services. Whether you're building a mobile app, TV app, or backend integration, I'm here to help. What are you working on?",
	}

	if msg, ok := messages[userType]; ok {
		return msg
	}

	return messages["end-user"]
}

func (tb *TrainingBot) detectActions(message string, userType string) []Action {
	// Detect common intents and suggest actions
	actions := []Action{}

	// Check for keywords and suggest relevant actions
	keywords := map[string][]Action{
		"upload": {
			{Type: "guide", Label: "📖 View Upload Guide", Data: map[string]interface{}{"guide": "content-upload"}},
			{Type: "video", Label: "🎥 Watch Upload Tutorial", Data: map[string]interface{}{"video": "upload-101"}},
		},
		"profile": {
			{Type: "guide", Label: "📖 Profile Management Guide", Data: map[string]interface{}{"guide": "profile-setup"}},
		},
		"api": {
			{Type: "docs", Label: "📚 View API Documentation", Data: map[string]interface{}{"url": "/docs/api"}},
			{Type: "code", Label: "💻 View Code Examples", Data: map[string]interface{}{"examples": "api-integration"}},
		},
	}

	for keyword, keywordActions := range keywords {
		if contains(message, keyword) {
			actions = append(actions, keywordActions...)
		}
	}

	return actions
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) && (str[:len(substr)] == substr ||
		str[len(str)-len(substr):] == substr)))
}
