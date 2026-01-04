# LearnHub Backend - Go REST API

A RESTful API backend for a learning platform with video lessons, quizzes, and progress tracking. Built with Go, Gin framework, and PostgreSQL with **100% raw SQL queries** (no ORM).

## Features

- 🔐 Simple userId-based authentication (no password required)
- 📚 5 chapters with video lessons and quizzes
- 📊 Progress tracking with resume capability
- 🎯 Quiz answer history tracking
- 💾 Auto-saves video timestamp and quiz progress
- 🔄 Resume from where you left off (Netflix-style)
- 🌐 Cross-platform support with CORS enabled
- ⚡ 100% query-driven (raw SQL, no ORM overhead)

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (Production) / SQLite (Development)
- **Queries**: Raw SQL (100% query-driven)
- **CORS**: gin-contrib/cors

## Project Structure

```
Demo_Backend/
├── main.go                 # Entry point and routes
├── handlers/               # Request handlers (all raw SQL)
│   ├── auth.go            # Authentication handlers
│   ├── chapters.go        # Chapter/video/quiz handlers
│   ├── progress.go        # Progress tracking handlers
│   ├── quiz_answers.go    # Quiz answer history handlers
│   └── quiz_with_history.go # Quiz resume/state handlers
├── database/              # Database connection
│   └── db.go             # GORM for connection only (queries are raw SQL)
├── config/               # Configuration
│   └── config.go
├── middleware/           # Middleware
│   └── cors.go          # CORS configuration
└── deploy_production.sh # Production deployment script
```

## Installation & Setup

### Prerequisites

- Go 1.21+
- PostgreSQL (for production) or SQLite (for development)

### 1. Install Dependencies

```bash
cd /home/ismail/Documents/Project/Demo_Backend
go mod download
```

### 2. Database Setup (PostgreSQL Production)

Set your database URL:

```bash
export DATABASE_URL="postgresql://production_db_em4b_user:aMDjUPiChzv5mfw2O70NyICuqf0IvqWc@dpg-d5ca45shg0os73e4ruu0-a.oregon-postgres.render.com/production_db_em4b"
```

The application will automatically connect to PostgreSQL when `DATABASE_URL` is set.

### 3. Run the Server

**Development (SQLite):**

```bash
go run main.go
```

**Production (PostgreSQL):**

```bash
./deploy_production.sh
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

```
GET /health
```

### Authentication

#### Login (Create/Get User)

```
POST /api/auth/login
Content-Type: application/json

{
  "user_id": "user_001"
}
```

#### Logout

```
POST /api/auth/logout
```

#### Get User Details

```
GET /api/auth/user/:userId
```

### Chapters

#### Get All Chapters

```
GET /api/chapters
```

#### Get Chapter by ID

```
GET /api/chapters/:id
```

#### Get Chapter Video

```
GET /api/chapters/:id/video
```

#### Get Chapter Quiz Questions

```
GET /api/chapters/:id/quiz
```

#### Get Chapter Content (Video + Quiz)

```
GET /api/chapters/:id/content
```

### Progress Tracking

#### Save Progress

```
POST /api/progress
Content-Type: application/json

For Video:
{
  "user_id": "user_001",
  "chapter_id": 1,
  "content_type": "video",
  "video_timestamp": 120,
  "is_completed": false
}

For Quiz:
{
  "user_id": "user_001",
  "chapter_id": 1,
  "content_type": "quiz",
  "quiz_question_index": 2,
  "is_completed": false
}
```

#### Get User's Latest Progress

```
GET /api/progress/user/:userId
```

#### Get All User Progress

```
GET /api/progress/user/:userId/all
```

#### Get Chapter-Specific Progress

```
GET /api/progress/user/:userId/chapter/:chapterId
```

#### Reset User Progress

```
DELETE /api/progress/user/:userId/reset
```

### Quiz Answer History

#### Submit Quiz Answer

```
POST /api/quiz/submit
Content-Type: application/json

{
  "user_id": "user_001",
  "chapter_id": 1,
  "quiz_question_id": 5,
  "user_answer": "B"
}

Response includes:
- is_correct: true/false
- correct_answer: "B"
- Full answer details
```

#### Get Quiz History for Chapter

```
GET /api/quiz/history/user/:userId/chapter/:chapterId
```

#### Get All Quiz History

```
GET /api/quiz/history/user/:userId
```

#### Get Quiz Scores Summary

```
GET /api/quiz/score/user/:userId

Returns score per chapter with percentages
```

#### Get Question Answer History

```
GET /api/quiz/history/user/:userId/question/:questionId
```

#### Clear Quiz History

```
DELETE /api/quiz/history/user/:userId/clear?chapter_id=1
```

### Quiz Resume Feature

#### Get Quiz with User's Answer History (Preserves State)

```
GET /api/quiz/chapter/:id/with-history?user_id=USER_ID

Returns all questions with:
- has_answered: true/false
- user_answer: "A"/"B"/"C"/"D" (if answered)
- is_correct: true/false (if answered)
- correct_answer: shown only if user has answered
- times_attempted: number of attempts
```

#### Get Quiz Resume Point

```
GET /api/quiz/resume/user/:userId/chapter/:chapterId

Returns first unanswered question or completion status
```

## Database Schema

### Tables

**users**

- id, user_id (unique), username
- created_at, updated_at, deleted_at

**chapters**

- id, title, description, order_index
- created_at, updated_at, deleted_at

**videos** (One-to-One with chapters)

- id, chapter_id (FK), title, video_url, duration_seconds
- created_at, updated_at, deleted_at

**quiz_questions** (One-to-Many with chapters)

- id, chapter_id (FK), question_text
- option_a, option_b, option_c, option_d
- correct_answer, order_index
- created_at, updated_at, deleted_at

**progresses**

- id, user_id, chapter_id (FK), content_type
- video_timestamp, quiz_question_index
- is_completed, last_updated
- created_at, updated_at, deleted_at

**quiz_answers** (Quiz history tracking)

- id, user_id, chapter_id (FK), quiz_question_id (FK)
- user_answer, is_correct, answered_at
- created_at, updated_at, deleted_at

### Relationships

- chapters (1) ──── videos (1) [One-to-One]
- chapters (1) ────< quiz_questions (M) [One-to-Many]
- chapters (1) ────< progresses (M) [One-to-Many]
- quiz_questions (1) ────< quiz_answers (M) [One-to-Many]

## Sample Data

### Chapters (5)

1. Introduction to Programming
2. Variables and Data Types
3. Control Flow
4. Functions and Methods
5. Object-Oriented Programming

### Videos (5)

- Real YouTube tutorial links
- One video per chapter
- Duration: 13-62 minutes

### Quiz Questions (15)

- 3 questions per chapter
- Multiple choice (A, B, C, D)
- Covers chapter topics

## Development

### Build

```bash
go build -o learnhub-server main.go
```

### Run

```bash
./learnhub-server
```

### Deploy to Production

```bash
./deploy_production.sh
```

## Architecture Highlights

### 🚀 100% Query-Driven

- **No ORM models** - All queries are raw SQL
- **Direct PostgreSQL** - Using database/sql with GORM connection pool
- **Better Performance** - No ORM overhead
- **Clear Intent** - SQL queries show exactly what's happening

### 📊 Quiz State Preservation

When users reopen a quiz:

- ✅ See which questions they've answered
- ✅ See if answers were correct/incorrect
- ✅ Resume from first unanswered question
- ✅ Review past answers and see correct answers
- ✅ Full history of all attempts

### 💾 Progress Tracking

- Auto-saves video playback position
- Tracks quiz question progress
- Resume functionality (Netflix-style)
- Per-chapter and global progress views

## Environment Variables

```bash
# PostgreSQL (Production)
export DATABASE_URL="postgresql://production_db_em4b_user:aMDjUPiChzv5mfw2O70NyICuqf0IvqWc@dpg-d5ca45shg0os73e4ruu0-a.oregon-postgres.render.com/production_db_em4b"
export PORT=8080
export ENVIRONMENT=production

# SQLite (Development) - Auto-detected if DATABASE_URL not set
# No environment variables needed
```

## Testing with cURL

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

**Submit Quiz Answer:**

```bash
curl -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user",
    "chapter_id": 1,
    "quiz_question_id": 1,
    "user_answer": "B"
  }'
```

**Get Quiz with History:**

```bash
curl "http://localhost:8080/api/quiz/chapter/1/with-history?user_id=test_user"
```

## Production Deployment

The app is configured for production deployment with:

- PostgreSQL database
- Silent logging (no SQL query logs)
- Environment-based configuration
- CORS enabled for cross-origin requests
- Port 8080 (configurable via PORT env var)

## Notes

- Database schema managed via external SQL files
- All queries use raw SQL with parameterized statements ($1, $2, etc.)
- Foreign keys enforce referential integrity
- Soft deletes via `deleted_at` column
- Indexes on frequently queried columns
- User progress preserved across sessions
- Quiz history tracks all attempts

## License

MIT
