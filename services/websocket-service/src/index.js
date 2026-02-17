const express = require('express');
const http = require('http');
const { Server } = require('socket.io');
const jwt = require('jsonwebtoken');

const isProduction = process.env.NODE_ENV === 'production';
const configuredJwtSecret = process.env.JWT_SECRET_KEY;

if (!configuredJwtSecret && isProduction) {
  throw new Error('JWT_SECRET_KEY must be set when NODE_ENV=production');
}

if (!configuredJwtSecret) {
  console.warn('[security] JWT_SECRET_KEY not set. Using a development-only fallback secret.');
}

const jwtSecret = configuredJwtSecret || 'dev-only-secret-change-me';
const allowedOrigins = (process.env.CORS_ORIGINS || '*')
  .split(',')
  .map((origin) => origin.trim())
  .filter(Boolean);

const app = express();
const server = http.createServer(app);
const io = new Server(server, {
  cors: {
    origin: allowedOrigins.includes('*') ? true : allowedOrigins,
    methods: ['GET', 'POST'],
    credentials: true
  },
  transports: ['websocket'],
  pingInterval: 25000,
  pingTimeout: 20000
});

// Authentication middleware
io.use((socket, next) => {
  const token = socket.handshake.auth.token;
  if (!token) {
    return next(new Error('Authentication error'));
  }
  
  try {
    const decoded = jwt.verify(token, jwtSecret);
    socket.userId = decoded.user_id;
    next();
  } catch (err) {
    next(new Error('Authentication error'));
  }
});

io.on('connection', (socket) => {
  console.log(`User connected: ${socket.userId}`);

  socket.on('join-room', (roomId) => {
    socket.join(roomId);
    socket.to(roomId).emit('user-joined', socket.userId);
  });

  socket.on('leave-room', (roomId) => {
    socket.leave(roomId);
    socket.to(roomId).emit('user-left', socket.userId);
  });

  socket.on('chat-message', (data) => {
    if (!data || typeof data.roomId !== 'string' || typeof data.message !== 'string') {
      return;
    }

    if (data.message.length > 2000) {
      return;
    }

    io.to(data.roomId).emit('chat-message', {
      userId: socket.userId,
      message: data.message,
      timestamp: new Date()
    });
  });

  socket.on('disconnect', () => {
    console.log(`User disconnected: ${socket.userId}`);
  });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

const PORT = process.env.PORT || 8080;
server.listen(PORT, () => {
  console.log(`WebSocket service started on port ${PORT}`);
});
