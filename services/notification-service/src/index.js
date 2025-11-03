const express = require('express');
const app = express();

app.use(express.json());

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE');
  res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  next();
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

// Notification routes - Issue #20: Routes updated to match requirements
app.post('/notifications/send', async (req, res) => {
  const { user_id, channel, template, context } = req.body;
  
  // TODO: Send notification via appropriate channel (FCM, Email, SMS)
  // TODO: Render template with context using Handlebars
  // TODO: Track delivery status
  
  res.json({ message: 'Notification sent', notification_id: '123' });
});

app.get('/notifications/:user_id', async (req, res) => {
  const { user_id } = req.params;
  
  // TODO: Retrieve notification history from database
  
  res.json({ notifications: [] });
});

app.put('/notifications/:user_id/preferences', async (req, res) => {
  const { user_id } = req.params;
  const preferences = req.body;
  
  // TODO: Update notification preferences
  // TODO: Store in database
  
  res.json({ message: 'Preferences updated' });
});

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`Notification service started on port ${PORT}`);
});

