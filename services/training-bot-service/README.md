# StreamVerse Training Bot Service

AI-powered training and learning platform for StreamVerse users, creators, administrators, and developers.

## Overview

The Training Bot Service is an intelligent, conversational training system that provides:
- 🤖 **AI-Powered Chat Bot** - Interactive training assistant
- 📚 **Automated Manual Generation** - Create comprehensive guides
- 🎬 **Video Script Generation** - Produce training video scripts
- 📝 **Quiz & Assessment Creation** - Generate interactive quizzes
- 🎯 **Personalized Learning Paths** - Customized onboarding flows
- 📊 **Progress Tracking** - Monitor learning achievements

## Features

### Conversational Training Bot

- **Context-Aware**: Understands user role (end-user, creator, admin, developer)
- **Multi-Turn Conversations**: Maintains conversation context
- **Action Suggestions**: Recommends relevant guides, videos, docs
- **24/7 Availability**: Always ready to help

### Content Generation

- **User Manuals**: Comprehensive guides for all user types
- **Video Scripts**: Detailed training video scripts with scenes
- **Interactive Quizzes**: Knowledge assessments with explanations
- **Onboarding Flows**: Step-by-step personalized paths

### Training Library

- **50+ Video Tutorials**: Covering all platform features
- **25+ Interactive Modules**: Hands-on exercises and simulations
- **15+ Certifications**: Achievement badges and completion certificates
- **Multi-Language Support**: Content in 10+ languages

## Architecture

```
┌─────────────────────────┐
│   Training Bot UI       │
│   (React Component)     │
└───────────┬─────────────┘
            │
┌───────────▼─────────────┐
│  Training Bot Service   │
│  (Go, Port 8096)        │
├─────────────────────────┤
│  • Chat Handler         │
│  • Content Generator    │
│  • Progress Tracker     │
│  • Analytics Engine     │
└───────────┬─────────────┘
            │
┌───────────▼─────────────┐
│   Gemini AI API         │
│   (Content Generation)  │
└─────────────────────────┘
```

## API Endpoints

### Training Bot

```
POST   /api/v1/bot/chat              # Send message to bot
POST   /api/v1/bot/session           # Create chat session
GET    /api/v1/bot/session/:id       # Get session history
```

### Content Generation

```
POST   /api/v1/generate/manual       # Generate user manual
POST   /api/v1/generate/video-script # Generate video script
POST   /api/v1/generate/quiz         # Generate quiz
POST   /api/v1/generate/onboarding   # Generate onboarding flow
```

### Training Content

```
GET    /api/v1/training/manuals      # List all manuals
GET    /api/v1/training/videos       # List video trainings
GET    /api/v1/training/modules      # List interactive modules
```

### Progress Tracking

```
GET    /api/v1/progress/:userId      # Get user progress
POST   /api/v1/progress/:userId/complete  # Mark module complete
```

## Quick Start

### Prerequisites

- Go 1.21+
- Gemini API key
- PostgreSQL (for progress tracking)
- Redis (for session management)

### Installation

```bash
# Clone repository
git clone https://github.com/streamverse/training-bot-service
cd training-bot-service

# Install dependencies
go mod download

# Set environment variables
export GEMINI_API_KEY=your_api_key
export DATABASE_URL=postgresql://...
export REDIS_URL=redis://...

# Run service
go run cmd/main.go
```

### Docker Deployment

```bash
# Build image
docker build -t streamverse/training-bot:latest .

# Run container
docker run -p 8096:8096 \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e DATABASE_URL=$DATABASE_URL \
  streamverse/training-bot:latest
```

### Kubernetes Deployment

```bash
kubectl apply -f k8s/training-bot-deployment.yaml
```

## Usage Examples

### Creating a Chat Session

```javascript
// Create session
const response = await fetch('https://api.streamverse.io/v1/bot/session', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_TOKEN'
  },
  body: JSON.stringify({
    userId: 'user123',
    userType: 'end-user' // or 'creator', 'admin', 'developer'
  })
});

const { sessionId, message } = await response.json();
console.log(message);
// "👋 Welcome to StreamVerse! I'm your training assistant..."
```

### Chatting with Bot

```javascript
const response = await fetch('https://api.streamverse.io/v1/bot/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_TOKEN'
  },
  body: JSON.stringify({
    sessionId: 'session123',
    userId: 'user123',
    userType: 'end-user',
    message: 'How do I download videos for offline viewing?'
  })
});

const { message, actions } = await response.json();
console.log(message);
// "To download videos for offline viewing, you'll need to use our mobile app..."

console.log(actions);
// [
//   { type: 'guide', label: '📖 View Download Guide', data: {...} },
//   { type: 'video', label: '🎥 Watch Download Tutorial', data: {...} }
// ]
```

### Generating a Manual

```javascript
const response = await fetch('https://api.streamverse.io/v1/generate/manual', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer ADMIN_TOKEN'
  },
  body: JSON.stringify({
    category: 'end-user',
    topics: [
      'Finding content',
      'Creating profiles',
      'Parental controls',
      'Downloads'
    ],
    format: 'markdown',
    language: 'en'
  })
});

const { manual } = await response.json();
// Returns comprehensive markdown manual
```

### Generating a Video Script

```javascript
const response = await fetch('https://api.streamverse.io/v1/generate/video-script', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer ADMIN_TOKEN'
  },
  body: JSON.stringify({
    topic: 'Getting Started with StreamVerse',
    duration: 90, // seconds
    audience: 'new-users',
    style: 'tutorial'
  })
});

const { script, scenes } = await response.json();
// Returns detailed video script with scenes, voiceover, visuals
```

## Training Content Structure

```
training/
├── manuals/
│   ├── end-users/
│   │   └── END_USER_MANUAL.md
│   ├── creators/
│   │   └── CREATOR_MANUAL.md
│   ├── admins/
│   │   └── ADMIN_MANUAL.md
│   └── developers/
│       └── DEVELOPER_MANUAL.md
├── video-scripts/
│   └── VIDEO_TRAINING_SCRIPTS.md
├── interactive-modules/
│   └── INTERACTIVE_MODULES.md
├── assessments/
│   └── quizzes.json
└── onboarding/
    └── onboarding-flows.json
```

## Configuration

### Environment Variables

```bash
# Required
GEMINI_API_KEY=your_gemini_api_key
DATABASE_URL=postgresql://user:pass@localhost/training_bot
REDIS_URL=redis://localhost:6379

# Optional
PORT=8096
LOG_LEVEL=info
SESSION_TIMEOUT=3600
MAX_CONTEXT_LENGTH=10
ENABLE_ANALYTICS=true
```

### AI Model Configuration

```go
// Customize AI behavior
model.SetTemperature(0.7)  // Creativity (0.0-1.0)
model.SetTopP(0.9)         // Diversity
model.SetMaxOutputTokens(2048)  // Response length
```

## Analytics & Insights

### Training Metrics

Track key metrics:
- **Total Sessions**: Number of bot conversations
- **Completion Rate**: % of users completing modules
- **Average Session Duration**: Time spent in training
- **Popular Topics**: Most asked questions
- **User Satisfaction**: Ratings and feedback

### Dashboard

Access analytics at: `https://admin.streamverse.io/training/analytics`

## Integration

### Embedding in Web App

```jsx
import { TrainingBot } from '@streamverse/training-bot-react';

function App() {
  return (
    <TrainingBot
      apiUrl="https://api.streamverse.io/v1"
      userType="end-user"
      userId={currentUser.id}
      theme="dark"
    />
  );
}
```

### Mobile Integration

```dart
// Flutter
import 'package:streamverse_training/training_bot.dart';

TrainingBotWidget(
  apiUrl: 'https://api.streamverse.io/v1',
  userType: UserType.endUser,
  userId: currentUser.id,
)
```

## Localization

Supported languages:
- English (en)
- Spanish (es)
- French (fr)
- German (de)
- Portuguese (pt)
- Japanese (ja)
- Korean (ko)
- Chinese (zh)
- Hindi (hi)
- Arabic (ar)

Add new language:
```bash
# Generate translations
./scripts/generate-translations.sh --lang=it
```

## Testing

### Unit Tests

```bash
go test ./...
```

### Integration Tests

```bash
go test ./tests/integration -v
```

### Load Testing

```bash
k6 run tests/load/training-bot-load-test.js
```

## Monitoring

### Health Check

```bash
curl http://localhost:8096/health
```

### Metrics

Prometheus metrics available at: `/metrics`

Key metrics:
- `training_bot_sessions_total`
- `training_bot_messages_total`
- `training_bot_response_duration_seconds`
- `content_generation_duration_seconds`

### Logging

Structured JSON logs:
```json
{
  "level": "info",
  "timestamp": "2025-01-20T10:00:00Z",
  "service": "training-bot",
  "event": "chat_message",
  "session_id": "sess_123",
  "user_id": "user_456",
  "duration_ms": 234
}
```

## Troubleshooting

### Common Issues

**Bot not responding:**
1. Check Gemini API key is valid
2. Verify API quota not exceeded
3. Check network connectivity

**Slow responses:**
1. Reduce context length
2. Optimize prompts
3. Check AI model load

**Session not found:**
1. Verify Redis connection
2. Check session expiration
3. Ensure session ID is correct

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- **Documentation**: [docs.streamverse.io/training](https://docs.streamverse.io/training)
- **API Reference**: [api.streamverse.io/training](https://api.streamverse.io/training)
- **Issues**: [GitHub Issues](https://github.com/streamverse/training-bot-service/issues)
- **Email**: training-support@streamverse.io

---

**Training Bot Service Version**: 2.0
**Last Updated**: 2025
**Status**: Production Ready

**Built with ❤️ by the StreamVerse Team**
