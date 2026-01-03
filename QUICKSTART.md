# Quick Start Guide

## Prerequisites

1. **Install Go** (version 1.21 or higher)
   ```bash
   # Check if Go is installed
   go version

   # If not installed, download from:
   # https://golang.org/dl/
   ```

## Setup & Run

1. **Navigate to the project directory**
   ```bash
   cd /home/ismail/Documents/Project/Demo_Backend
   ```

2. **Download dependencies**
   ```bash
   go mod download
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

   Or use the helper script:
   ```bash
   ./run.sh
   ```

4. **Server will start on**: `http://localhost:8080`

## Quick Test

Open another terminal and test the API:

```bash
# Health check
curl http://localhost:8080/health

# Login with a user ID
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "john_doe"}'

# Get all chapters
curl http://localhost:8080/api/chapters

# Run all tests
./test_api.sh
```

## Key Features

### 1. Simple Login (No Password)
Just provide a `user_id` - user will be created automatically if it doesn't exist.

### 2. Resume Functionality
- **Video**: Saves timestamp in seconds
- **Quiz**: Saves current question index (0-based)

### 3. Sample Data
- 3 Chapters with videos and quizzes
- 5 questions per chapter
- Public video URLs for testing

## Common API Flows

### Login Flow
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123"}'
```

### Get Chapters
```bash
curl http://localhost:8080/api/chapters
```

### Save Video Progress (e.g., at 33 seconds)
```bash
curl -X POST http://localhost:8080/api/progress \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "chapter_id": 1,
    "content_type": "video",
    "video_timestamp": 33,
    "is_completed": false
  }'
```

### Save Quiz Progress (e.g., at question 3)
```bash
curl -X POST http://localhost:8080/api/progress \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "chapter_id": 1,
    "content_type": "quiz",
    "quiz_question_index": 3,
    "is_completed": false
  }'
```

### Get User's Latest Progress (for "Continue" feature)
```bash
curl http://localhost:8080/api/progress/user/user123
```

## Project Structure

```
Demo_Backend/
├── main.go              # Entry point with routes
├── models/              # Data models
├── handlers/            # Request handlers
├── database/            # DB setup & seeding
├── middleware/          # CORS, etc.
├── run.sh              # Helper to start server
├── test_api.sh         # API test script
└── README.md           # Full documentation
```

## Database

- **Type**: SQLite (file-based)
- **File**: `learning_app.db` (auto-created on first run)
- **Seeded**: Automatically populated with 3 chapters, videos, and quizzes

## Troubleshooting

**"go: command not found"**
- Install Go from https://golang.org/dl/

**Port 8080 already in use**
- Stop other services on port 8080 or change the port in `main.go`

**Database errors**
- Delete `learning_app.db` and restart the server to recreate

## Next Steps

1. Start the backend server
2. Build the Flutter frontend (coming next)
3. Connect Flutter app to this API
4. Test the resume functionality end-to-end

For detailed API documentation, see [README.md](README.md)
