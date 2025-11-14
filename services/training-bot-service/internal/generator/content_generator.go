package generator

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/training-bot-service/pkg/ai"
)

type ContentGenerator struct {
	aiClient ai.AIClient
}

type ManualRequest struct {
	Category string   `json:"category"` // end-user, creator, admin, developer
	Topics   []string `json:"topics"`
	Format   string   `json:"format"` // markdown, pdf, html
	Language string   `json:"language"`
}

type VideoScriptRequest struct {
	Topic    string `json:"topic"`
	Duration int    `json:"duration"` // in seconds
	Audience string `json:"audience"`
	Style    string `json:"style"` // tutorial, quick-tip, deep-dive
}

type QuizRequest struct {
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"` // beginner, intermediate, advanced
	NumQuestions int  `json:"numQuestions"`
}

type OnboardingRequest struct {
	UserType string `json:"userType"`
	Goals    []string `json:"goals"`
}

func NewContentGenerator(aiClient ai.AIClient) *ContentGenerator {
	return &ContentGenerator{
		aiClient: aiClient,
	}
}

func (cg *ContentGenerator) GenerateManual(c *gin.Context) {
	var req ManualRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := cg.buildManualPrompt(req)
	manual, err := cg.aiClient.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"manual":   manual,
		"category": req.Category,
		"format":   req.Format,
	})
}

func (cg *ContentGenerator) GenerateVideoScript(c *gin.Context) {
	var req VideoScriptRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := cg.buildVideoScriptPrompt(req)
	script, err := cg.aiClient.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"script":   script,
		"topic":    req.Topic,
		"duration": req.Duration,
		"scenes":   cg.parseScenes(script),
	})
}

func (cg *ContentGenerator) GenerateQuiz(c *gin.Context) {
	var req QuizRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := cg.buildQuizPrompt(req)
	quiz, err := cg.aiClient.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz":       quiz,
		"topic":      req.Topic,
		"difficulty": req.Difficulty,
	})
}

func (cg *ContentGenerator) GenerateOnboarding(c *gin.Context) {
	var req OnboardingRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := cg.buildOnboardingPrompt(req)
	onboarding, err := cg.aiClient.Generate(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"onboarding": onboarding,
		"userType":   req.UserType,
		"steps":      cg.parseSteps(onboarding),
	})
}

func (cg *ContentGenerator) buildManualPrompt(req ManualRequest) string {
	return fmt.Sprintf(`Generate a comprehensive training manual for %s users of the StreamVerse platform.

Topics to cover:
%v

Format: %s
Language: %s

The manual should include:
1. Introduction and overview
2. Step-by-step instructions with screenshots descriptions
3. Best practices and tips
4. Common issues and troubleshooting
5. FAQ section
6. Glossary of terms

Make it beginner-friendly, well-structured, and include practical examples.`,
		req.Category, req.Topics, req.Format, req.Language)
}

func (cg *ContentGenerator) buildVideoScriptPrompt(req VideoScriptRequest) string {
	return fmt.Sprintf(`Create a video training script for StreamVerse on the topic: %s

Duration: %d seconds
Target Audience: %s
Style: %s

The script should include:
1. Opening hook (5-10 seconds)
2. Introduction (10-15 seconds)
3. Main content with clear sections
4. Visual cues and on-screen text suggestions
5. Transitions between sections
6. Call-to-action at the end
7. Closing (5-10 seconds)

Format the script with:
- [SCENE X] markers
- [VISUAL: description] for what should appear on screen
- [VOICEOVER: text] for narration
- [TEXT: content] for on-screen text
- Timestamps for each section

Make it engaging, concise, and educational.`, req.Topic, req.Duration, req.Audience, req.Style)
}

func (cg *ContentGenerator) buildQuizPrompt(req QuizRequest) string {
	return fmt.Sprintf(`Generate a quiz for StreamVerse training on topic: %s

Difficulty: %s
Number of Questions: %d

Create questions in the following format:
1. Multiple choice (with 4 options, 1 correct)
2. True/False
3. Fill in the blank
4. Scenario-based questions

For each question, provide:
- Question text
- Options (if applicable)
- Correct answer
- Explanation of why the answer is correct
- Learning objective addressed

Ensure questions test practical understanding and real-world application.`,
		req.Topic, req.Difficulty, req.NumQuestions)
}

func (cg *ContentGenerator) buildOnboardingPrompt(req OnboardingRequest) string {
	return fmt.Sprintf(`Create a personalized onboarding flow for a %s user of StreamVerse.

User Goals:
%v

Generate a step-by-step onboarding plan that includes:
1. Welcome and introduction
2. Account setup steps
3. Platform familiarization
4. Key features to explore
5. First actions to take
6. Resources for learning more
7. Success milestones

For each step, provide:
- Step number and title
- Clear instructions
- Expected outcome
- Estimated time
- Tips and best practices

Make it engaging and help the user achieve quick wins early.`, req.UserType, req.Goals)
}

func (cg *ContentGenerator) parseScenes(script string) []map[string]string {
	// Parse script into scenes
	// This is a placeholder - implement actual parsing
	return []map[string]string{}
}

func (cg *ContentGenerator) parseSteps(onboarding string) []map[string]interface{} {
	// Parse onboarding into steps
	// This is a placeholder - implement actual parsing
	return []map[string]interface{}{}
}
