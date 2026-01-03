# API Summary - Learning App Backend

## Server Info
- **URL**: `http://localhost:8080`
- **Framework**: Gin (Go)
- **Database**: SQLite
- **CORS**: Enabled for all origins

## Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/login` | Login/create user (no password) |
| POST | `/api/auth/logout` | Logout user |
| GET | `/api/auth/user/:userId` | Get user details |

## Chapter Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/chapters` | Get all chapters (ordered) |
| GET | `/api/chapters/:id` | Get chapter by ID |
| GET | `/api/chapters/:id/video` | Get chapter video |
| GET | `/api/chapters/:id/quiz` | Get chapter quiz questions |
| GET | `/api/chapters/:id/content` | Get video + quiz together |

## Progress Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/progress` | Save/update progress |
| GET | `/api/progress/user/:userId` | Get latest progress (for resume) |
| GET | `/api/progress/user/:userId/all` | Get all progress for user |
| GET | `/api/progress/user/:userId/chapter/:chapterId` | Get chapter progress |
| DELETE | `/api/progress/user/:userId/reset` | Reset all progress |

## Data Models

### User
```json
{
  "id": 1,
  "user_id": "john_doe",
  "username": "john_doe",
  "created_at": "2024-01-02T10:00:00Z"
}
```

### Chapter
```json
{
  "id": 1,
  "title": "Introduction to Programming",
  "description": "Learn the basics...",
  "order_index": 1
}
```

### Video
```json
{
  "id": 1,
  "chapter_id": 1,
  "title": "What is Programming?",
  "video_url": "https://...",
  "duration_seconds": 596
}
```

### Quiz Question
```json
{
  "id": 1,
  "chapter_id": 1,
  "question_text": "What is a variable?",
  "option_a": "...",
  "option_b": "...",
  "option_c": "...",
  "option_d": "...",
  "correct_answer": "B",
  "order_index": 0
}
```

### Progress
```json
{
  "id": 1,
  "user_id": "john_doe",
  "chapter_id": 1,
  "content_type": "video",
  "video_timestamp": 33,
  "quiz_question_index": null,
  "is_completed": false,
  "last_updated": "2024-01-02T10:30:00Z"
}
```

## Resume Logic Implementation

### For Video Resume:
1. Client saves progress every 5-10 seconds:
   ```json
   POST /api/progress
   {
     "user_id": "user123",
     "chapter_id": 1,
     "content_type": "video",
     "video_timestamp": 33,
     "is_completed": false
   }
   ```

2. On app start, get latest progress:
   ```json
   GET /api/progress/user/user123

   Response:
   {
     "has_progress": true,
     "last_chapter_id": 1,
     "last_content_type": "video",
     "last_video_time": 33,
     "chapter_title": "Introduction to Programming"
   }
   ```

3. Navigate to that chapter and seek to timestamp 33

### For Quiz Resume:
1. Save after each question:
   ```json
   POST /api/progress
   {
     "user_id": "user123",
     "chapter_id": 1,
     "content_type": "quiz",
     "quiz_question_index": 2,
     "is_completed": false
   }
   ```

2. On return, navigate to question at index 2

## Sample Content

### Chapters (3)
1. Introduction to Programming
2. Data Structures Fundamentals
3. Algorithms and Problem Solving

### Videos (3)
- Uses Google's public sample videos
- Durations: 15s to 653s

### Quiz Questions (15)
- 5 per chapter
- Multiple choice format
- Covers chapter topics

## Running the Server

```bash
# Install dependencies
go mod download

# Run server
go run main.go

# Or use helper script
./run.sh
```

## Testing

```bash
# Run all API tests
./test_api.sh

# Or test individual endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/chapters
```

## Flutter Integration Points

1. **On App Launch**: Call `/api/progress/user/:userId` to check for existing progress
2. **Video Player**:
   - Fetch video via `/api/chapters/:id/video`
   - Auto-save progress every 5-10 seconds to `/api/progress`
   - On dispose, save final timestamp
3. **Quiz Screen**:
   - Fetch questions via `/api/chapters/:id/quiz`
   - Save progress after each answer to `/api/progress`
4. **Home Screen**:
   - Show "Continue" card if `has_progress: true`
   - Display chapter title and progress type
   - On tap, navigate to exact position

## Error Responses

All errors follow this format:
```json
{
  "success": false,
  "message": "Error description"
}
```

Common HTTP status codes:
- `200`: Success
- `400`: Bad request (invalid data)
- `404`: Resource not found
- `500`: Server error
