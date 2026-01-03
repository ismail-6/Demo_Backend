# Learning App Backend - Go REST API

A RESTful API backend for a learning app with resume functionality (Netflix-style). Built with Go, Gin framework, and SQLite.

## Features

- Simple userId-based authentication (no password required)
- 3 chapters with video lessons and quizzes
- Progress tracking with resume capability
- Auto-saves video timestamp and quiz question index
- Cross-platform support with CORS enabled

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: SQLite
- **CORS**: gin-contrib/cors

## Project Structure

```
Demo_Backend/
├── main.go                 # Entry point and routes
├── models/                 # Data models
│   ├── user.go
│   ├── chapter.go
│   └── progress.go
├── handlers/               # Request handlers
│   ├── auth.go
│   ├── chapters.go
│   └── progress.go
├── database/               # Database setup and seeding
│   └── db.go
├── middleware/             # Middleware (CORS, etc.)
│   └── cors.go
└── learning_app.db        # SQLite database (auto-created)
```

## Installation

1. **Install Go** (1.21 or higher)
   - Download from: https://golang.org/dl/

2. **Navigate to project directory**
   ```bash
   cd /home/ismail/Documents/Project/Demo_Backend
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run the server**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /health
```
Check if the API is running.

### Authentication

#### Login (Create/Get User)
```
POST /api/auth/login
Content-Type: application/json

{
  "user_id": "user123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": 1,
    "user_id": "user123",
    "username": "user123",
    "created_at": "2024-01-02T10:00:00Z"
  }
}
```

#### Logout
```
POST /api/auth/logout

Response:
{
  "success": true,
  "message": "Logout successful"
}
```

#### Get User Details
```
GET /api/auth/user/:userId

Response:
{
  "success": true,
  "user": {
    "id": 1,
    "user_id": "user123",
    "username": "user123"
  }
}
```

### Chapters

#### Get All Chapters
```
GET /api/chapters

Response:
{
  "success": true,
  "chapters": [
    {
      "id": 1,
      "title": "Introduction to Programming",
      "description": "Learn the basics of programming concepts and logic",
      "order_index": 1
    },
    ...
  ]
}
```

#### Get Chapter by ID
```
GET /api/chapters/:id

Response:
{
  "success": true,
  "chapter": {
    "id": 1,
    "title": "Introduction to Programming",
    "description": "Learn the basics...",
    "order_index": 1
  }
}
```

#### Get Chapter Video
```
GET /api/chapters/:id/video

Response:
{
  "success": true,
  "video": {
    "id": 1,
    "chapter_id": 1,
    "title": "What is Programming?",
    "video_url": "https://...",
    "duration_seconds": 596
  }
}
```

#### Get Chapter Quiz
```
GET /api/chapters/:id/quiz

Response:
{
  "success": true,
  "questions": [
    {
      "id": 1,
      "chapter_id": 1,
      "question_text": "What is a variable in programming?",
      "option_a": "A fixed value that never changes",
      "option_b": "A container for storing data values",
      "option_c": "A type of loop",
      "option_d": "A programming language",
      "correct_answer": "B",
      "order_index": 0
    },
    ...
  ]
}
```

#### Get Chapter Content (Video + Quiz)
```
GET /api/chapters/:id/content

Response:
{
  "success": true,
  "chapter": {
    "id": 1,
    "title": "Introduction to Programming",
    "video": {...},
    "quiz_questions": [...]
  }
}
```

### Progress Tracking

#### Save Progress
```
POST /api/progress
Content-Type: application/json

For Video:
{
  "user_id": "user123",
  "chapter_id": 1,
  "content_type": "video",
  "video_timestamp": 33,
  "is_completed": false
}

For Quiz:
{
  "user_id": "user123",
  "chapter_id": 1,
  "content_type": "quiz",
  "quiz_question_index": 2,
  "is_completed": false
}

Response:
{
  "success": true,
  "message": "Progress saved successfully",
  "progress": {
    "id": 1,
    "user_id": "user123",
    "chapter_id": 1,
    "content_type": "video",
    "video_timestamp": 33,
    "is_completed": false,
    "last_updated": "2024-01-02T10:30:00Z"
  }
}
```

#### Get User's Latest Progress
```
GET /api/progress/user/:userId

Response:
{
  "success": true,
  "progress": {
    "has_progress": true,
    "last_chapter_id": 2,
    "last_content_type": "video",
    "last_video_time": 33,
    "chapter_title": "Data Structures Fundamentals",
    "last_updated": "2024-01-02T10:30:00Z"
  }
}
```

#### Get All User Progress
```
GET /api/progress/user/:userId/all

Response:
{
  "success": true,
  "progress": [
    {
      "id": 1,
      "user_id": "user123",
      "chapter_id": 1,
      "content_type": "video",
      "video_timestamp": 596,
      "is_completed": true
    },
    ...
  ]
}
```

#### Get Chapter-Specific Progress
```
GET /api/progress/user/:userId/chapter/:chapterId

Response:
{
  "success": true,
  "progress": [
    {
      "id": 1,
      "user_id": "user123",
      "chapter_id": 1,
      "content_type": "video",
      "video_timestamp": 150
    },
    {
      "id": 2,
      "user_id": "user123",
      "chapter_id": 1,
      "content_type": "quiz",
      "quiz_question_index": 3
    }
  ]
}
```

#### Reset User Progress
```
DELETE /api/progress/user/:userId/reset

Response:
{
  "success": true,
  "message": "Progress reset successfully",
  "deleted": 5
}
```

## Sample Data

The database is auto-seeded with:

### Chapters (3)
1. Introduction to Programming
2. Data Structures Fundamentals
3. Algorithms and Problem Solving

### Videos (3)
- Each chapter has one video (using public sample videos)
- Videos range from 15 to 653 seconds

### Quiz Questions (15)
- 5 questions per chapter
- Multiple choice (A, B, C, D)
- Questions cover chapter topics

## Development

### Build
```bash
go build -o learning-app-server
```

### Run
```bash
./learning-app-server
```

### Test with cURL

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test_user"}'
```

**Get Chapters:**
```bash
curl http://localhost:8080/api/chapters
```

**Save Video Progress:**
```bash
curl -X POST http://localhost:8080/api/progress \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user",
    "chapter_id": 1,
    "content_type": "video",
    "video_timestamp": 120,
    "is_completed": false
  }'
```

**Get User Progress:**
```bash
curl http://localhost:8080/api/progress/user/test_user
```

## Database Schema

### Users
- `id` (primary key)
- `user_id` (unique)
- `username`
- `created_at`, `updated_at`

### Chapters
- `id` (primary key)
- `title`
- `description`
- `order_index`

### Videos
- `id` (primary key)
- `chapter_id` (foreign key)
- `title`
- `video_url`
- `duration_seconds`

### Quiz Questions
- `id` (primary key)
- `chapter_id` (foreign key)
- `question_text`
- `option_a`, `option_b`, `option_c`, `option_d`
- `correct_answer`
- `order_index`

### Progress
- `id` (primary key)
- `user_id` (indexed)
- `chapter_id` (indexed)
- `content_type` (video/quiz)
- `video_timestamp` (nullable)
- `quiz_question_index` (nullable)
- `is_completed`
- `last_updated`

## Notes

- The database file `learning_app.db` is created automatically on first run
- User IDs are auto-created on first login (no pre-registration needed)
- Progress is automatically updated with latest timestamp/index
- CORS is enabled for all origins (suitable for development)
- Video URLs use Google's public sample videos

## Next Steps

For production deployment:
1. Change CORS settings to specific origins
2. Add authentication middleware for protected routes
3. Switch to PostgreSQL or MySQL for production
4. Add rate limiting
5. Add logging middleware
6. Add input validation and sanitization
7. Add unit tests

## License

MIT
