-- LearnHub Backend Database Schema
-- PostgreSQL Database Schema with Relationships

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS progresses CASCADE;
DROP TABLE IF EXISTS quiz_questions CASCADE;
DROP TABLE IF EXISTS videos CASCADE;
DROP TABLE IF EXISTS chapters CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_user_id ON users(user_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Chapters Table
CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_chapters_order_index ON chapters(order_index);
CREATE INDEX idx_chapters_deleted_at ON chapters(deleted_at);

-- Videos Table (One-to-One with Chapter)
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    video_url TEXT NOT NULL,
    duration_seconds INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    -- Foreign Key Constraint
    CONSTRAINT fk_videos_chapter FOREIGN KEY (chapter_id)
        REFERENCES chapters(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_videos_chapter_id ON videos(chapter_id);
CREATE INDEX idx_videos_deleted_at ON videos(deleted_at);

-- Quiz Questions Table (One-to-Many with Chapter)
CREATE TABLE quiz_questions (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    option_a VARCHAR(500) NOT NULL,
    option_b VARCHAR(500) NOT NULL,
    option_c VARCHAR(500) NOT NULL,
    option_d VARCHAR(500) NOT NULL,
    correct_answer VARCHAR(1) NOT NULL CHECK (correct_answer IN ('A', 'B', 'C', 'D')),
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    -- Foreign Key Constraint
    CONSTRAINT fk_quiz_questions_chapter FOREIGN KEY (chapter_id)
        REFERENCES chapters(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_quiz_questions_chapter_id ON quiz_questions(chapter_id);
CREATE INDEX idx_quiz_questions_deleted_at ON quiz_questions(deleted_at);

-- Progress Table (Tracks user progress through chapters)
CREATE TABLE progresses (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    chapter_id INTEGER NOT NULL,
    content_type VARCHAR(10) NOT NULL CHECK (content_type IN ('video', 'quiz')),
    video_timestamp INTEGER DEFAULT NULL,
    quiz_question_index INTEGER DEFAULT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    -- Foreign Key Constraints
    CONSTRAINT fk_progresses_chapter FOREIGN KEY (chapter_id)
        REFERENCES chapters(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_progresses_user_id ON progresses(user_id);
CREATE INDEX idx_progresses_chapter_id ON progresses(chapter_id);
CREATE INDEX idx_progresses_user_chapter ON progresses(user_id, chapter_id);
CREATE INDEX idx_progresses_deleted_at ON progresses(deleted_at);

-- Comments explaining the relationships:
COMMENT ON TABLE users IS 'Stores user account information';
COMMENT ON TABLE chapters IS 'Learning chapters/modules';
COMMENT ON TABLE videos IS 'One video per chapter - teaching content';
COMMENT ON TABLE quiz_questions IS 'Multiple quiz questions per chapter - assessment';
COMMENT ON TABLE progresses IS 'Tracks user progress through videos and quizzes';

-- Relationship Summary:
-- users (1) ----< progresses (M)
-- chapters (1) ---- videos (1)         [One-to-One]
-- chapters (1) ----< quiz_questions (M) [One-to-Many]
-- chapters (1) ----< progresses (M)     [One-to-Many]
