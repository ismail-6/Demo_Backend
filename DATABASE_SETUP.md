# LearnHub Database Setup Guide

## Database Schema Overview

### Tables and Relationships

```
users (1) ────< progresses (M)
                     │
chapters (1) ────────┤
    │                │
    ├──── videos (1) [One-to-One]
    │
    └────< quiz_questions (M) [One-to-Many]
```

### Table Structures

**users**
- id (PK)
- user_id (unique)
- username
- timestamps

**chapters**
- id (PK)
- title
- description
- order_index
- timestamps

**videos** (FK: chapter_id → chapters.id)
- id (PK)
- chapter_id (FK)
- title
- video_url
- duration_seconds
- timestamps

**quiz_questions** (FK: chapter_id → chapters.id)
- id (PK)
- chapter_id (FK)
- question_text
- option_a, option_b, option_c, option_d
- correct_answer (A/B/C/D)
- order_index
- timestamps

**progresses** (FK: chapter_id → chapters.id)
- id (PK)
- user_id
- chapter_id (FK)
- content_type (video/quiz)
- video_timestamp
- quiz_question_index
- is_completed
- timestamps

## Setup Instructions

### Option 1: Full PostgreSQL Setup (Recommended for Production)

```bash
# Set your database URL
export DATABASE_URL="postgresql://username:password@host:port/database"

# Run the complete setup (creates schema + seeds data)
./setup_database.sh
```

This will:
- Drop existing tables (if any)
- Create all tables with proper foreign key relationships
- Seed initial data (chapters, videos, quizzes, sample users)

### Option 2: Schema Only

```bash
# Create tables without seed data
psql $DATABASE_URL -f schema.sql
```

### Option 3: Seed Data Only (After Schema Exists)

```bash
# Seed data with relationships
psql $DATABASE_URL -f seed_with_relations.sql
```

### Option 4: Using Go Seed Script (Works with SQLite/PostgreSQL)

```bash
# For development (SQLite)
./seed.sh

# For production (PostgreSQL)
export DATABASE_URL="postgresql://..."
./seed.sh
```

## What Gets Seeded

### Chapters (5)
1. Introduction to Programming
2. Variables and Data Types
3. Control Flow
4. Functions and Methods
5. Object-Oriented Programming

### Videos (5)
- One video per chapter with real YouTube tutorial links
- Duration ranges from 13 to 62 minutes

### Quiz Questions (15)
- 3 multiple-choice questions per chapter
- Each has 4 options (A, B, C, D)

### Sample Users (3)
- user_001 (john_doe)
- user_002 (jane_smith)
- user_003 (alex_dev)

### Sample Progress Data
- Demonstrates user progress tracking
- Shows video timestamps and quiz progress

## Database Relationship Details

### Foreign Key Constraints

All foreign keys use `ON DELETE CASCADE` and `ON UPDATE CASCADE`:

1. **videos.chapter_id → chapters.id**
   - Each chapter has exactly ONE video
   - Deleting a chapter deletes its video

2. **quiz_questions.chapter_id → chapters.id**
   - Each chapter has MULTIPLE quiz questions
   - Deleting a chapter deletes all its quiz questions

3. **progresses.chapter_id → chapters.id**
   - Tracks which chapter the progress is for
   - Deleting a chapter deletes all progress records

### Indexes

Created for optimal query performance:
- `user_id` lookups
- `chapter_id` foreign key joins
- Combined `(user_id, chapter_id)` for progress queries
- Soft delete support via `deleted_at` indexes

## Testing the Setup

### Check All Tables

```sql
SELECT
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM chapters) as chapters_count,
    (SELECT COUNT(*) FROM videos) as videos_count,
    (SELECT COUNT(*) FROM quiz_questions) as quiz_questions_count,
    (SELECT COUNT(*) FROM progresses) as progresses_count;
```

### Verify Relationships

```sql
-- Get chapters with their videos and quiz counts
SELECT
    c.id, c.title,
    v.title as video_title,
    COUNT(DISTINCT q.id) as quiz_count
FROM chapters c
LEFT JOIN videos v ON c.id = v.chapter_id
LEFT JOIN quiz_questions q ON c.id = q.chapter_id
GROUP BY c.id, c.title, v.title
ORDER BY c.order_index;
```

### Check User Progress

```sql
-- Get user progress summary
SELECT
    u.username,
    c.title as chapter_title,
    p.content_type,
    p.is_completed
FROM progresses p
JOIN users u ON p.user_id = u.user_id
JOIN chapters c ON p.chapter_id = c.id
ORDER BY u.username, c.order_index;
```

## Troubleshooting

### Error: relation "chapters" does not exist
Run the schema creation first:
```bash
psql $DATABASE_URL -f schema.sql
```

### Error: duplicate key value violates unique constraint
The database already has data. Either:
1. Use `seed_with_relations.sql` which includes TRUNCATE
2. Manually clear tables first

### Foreign Key Constraint Violation
Make sure you're creating data in the correct order:
1. chapters (parent)
2. videos (child)
3. quiz_questions (child)
4. progresses (child)

## Files Reference

- `schema.sql` - Complete database schema with relationships
- `seed_with_relations.sql` - Seed data respecting foreign keys
- `seed_data.sql` - Simple seed data (original)
- `setup_database.sh` - Automated setup script
- `seed.sh` - Go-based seeding (works with Go app)
- `cmd/seed/main.go` - Go seed program
