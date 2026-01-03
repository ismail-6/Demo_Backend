-- Seed Data for LearnHub Backend (With Foreign Key Relations)
-- Run this AFTER creating schema with: psql $DATABASE_URL -f schema.sql
-- Then run: psql $DATABASE_URL -f seed_with_relations.sql

BEGIN;

-- Clear existing data (respects foreign key constraints with CASCADE)
TRUNCATE TABLE progresses, quiz_questions, videos, chapters, users RESTART IDENTITY CASCADE;

-- Insert sample users
INSERT INTO users (user_id, username, created_at, updated_at) VALUES
('user_001', 'john_doe', NOW(), NOW()),
('user_002', 'jane_smith', NOW(), NOW()),
('user_003', 'alex_dev', NOW(), NOW());

-- Insert chapters (IDs will be 1, 2, 3, 4, 5)
INSERT INTO chapters (title, description, order_index, created_at, updated_at) VALUES
('Introduction to Programming', 'Learn the basics of programming and computer science fundamentals', 1, NOW(), NOW()),
('Variables and Data Types', 'Understanding variables, data types, and how to work with them', 2, NOW(), NOW()),
('Control Flow', 'Master if statements, loops, and conditional logic', 3, NOW(), NOW()),
('Functions and Methods', 'Learn how to write reusable code with functions', 4, NOW(), NOW()),
('Object-Oriented Programming', 'Introduction to OOP concepts and principles', 5, NOW(), NOW());

-- Insert videos (One video per chapter - uses chapter_id foreign key)
INSERT INTO videos (chapter_id, title, video_url, duration_seconds, created_at, updated_at) VALUES
(1, 'Programming for Beginners - Learn to Code', 'https://www.youtube.com/watch?v=zOjov-2OZ0E', 3720, NOW(), NOW()),
(2, 'Variables and Data Types Explained', 'https://www.youtube.com/watch?v=OH86oLzVzzc', 780, NOW(), NOW()),
(3, 'Control Flow - If Statements and Loops', 'https://www.youtube.com/watch?v=PYTyhXKSMzE', 900, NOW(), NOW()),
(4, 'Functions in Programming - Complete Guide', 'https://www.youtube.com/watch?v=nX34jsB_v3E', 1020, NOW(), NOW()),
(5, 'Object Oriented Programming (OOP) Explained', 'https://www.youtube.com/watch?v=pTB0EiLXUC8', 1380, NOW(), NOW());

-- Insert quiz questions (Multiple questions per chapter - uses chapter_id foreign key)
INSERT INTO quiz_questions (chapter_id, question_text, option_a, option_b, option_c, option_d, correct_answer, order_index, created_at, updated_at) VALUES
-- Chapter 1: Introduction to Programming
(1, 'What is a program?', 'A set of instructions for a computer', 'A type of computer hardware', 'A programming language', 'A database', 'A', 1, NOW(), NOW()),
(1, 'Which of these is a programming language?', 'HTML', 'Python', 'CSS', 'JSON', 'B', 2, NOW(), NOW()),
(1, 'What is the purpose of an algorithm?', 'To store data', 'To provide step-by-step instructions to solve a problem', 'To design user interfaces', 'To connect to databases', 'B', 3, NOW(), NOW()),

-- Chapter 2: Variables and Data Types
(2, 'What is a variable?', 'A fixed value', 'A container for storing data', 'A type of loop', 'A function', 'B', 1, NOW(), NOW()),
(2, 'Which data type stores whole numbers?', 'String', 'Float', 'Integer', 'Boolean', 'C', 2, NOW(), NOW()),
(2, 'What does a Boolean data type represent?', 'Numbers only', 'True or False values', 'Text strings', 'Decimal numbers', 'B', 3, NOW(), NOW()),

-- Chapter 3: Control Flow
(3, 'What does an if statement do?', 'Stores data', 'Makes decisions based on conditions', 'Creates loops', 'Defines functions', 'B', 1, NOW(), NOW()),
(3, 'Which loop runs at least once?', 'for loop', 'while loop', 'do-while loop', 'if loop', 'C', 2, NOW(), NOW()),
(3, 'What is the purpose of a loop?', 'To store variables', 'To repeat code multiple times', 'To define functions', 'To create classes', 'B', 3, NOW(), NOW()),

-- Chapter 4: Functions and Methods
(4, 'What is a function?', 'A variable', 'A reusable block of code', 'A data type', 'A loop', 'B', 1, NOW(), NOW()),
(4, 'What keyword is used to return a value?', 'give', 'return', 'send', 'output', 'B', 2, NOW(), NOW()),
(4, 'What are function parameters?', 'Return values', 'Input values passed to a function', 'Local variables only', 'Global constants', 'B', 3, NOW(), NOW()),

-- Chapter 5: Object-Oriented Programming
(5, 'What does OOP stand for?', 'Online Operating Program', 'Object-Oriented Programming', 'Ordered Operation Process', 'Optimized Output Protocol', 'B', 1, NOW(), NOW()),
(5, 'What is encapsulation?', 'Combining data and methods', 'A type of loop', 'A variable type', 'A function call', 'A', 2, NOW(), NOW()),
(5, 'What is inheritance in OOP?', 'Copying code manually', 'A class acquiring properties from another class', 'Deleting unused code', 'Running loops repeatedly', 'B', 3, NOW(), NOW());

-- Insert sample progress data (Demonstrates the many-to-many relationship)
INSERT INTO progresses (user_id, chapter_id, content_type, video_timestamp, quiz_question_index, is_completed, last_updated, created_at, updated_at) VALUES
-- User 001 progress
('user_001', 1, 'video', 1200, NULL, false, NOW(), NOW(), NOW()),
('user_001', 1, 'quiz', NULL, 2, true, NOW(), NOW(), NOW()),
('user_001', 2, 'video', 300, NULL, false, NOW(), NOW(), NOW()),

-- User 002 progress
('user_002', 1, 'video', 3720, NULL, true, NOW(), NOW(), NOW()),
('user_002', 1, 'quiz', NULL, 3, true, NOW(), NOW(), NOW()),
('user_002', 2, 'video', 780, NULL, true, NOW(), NOW(), NOW()),
('user_002', 2, 'quiz', NULL, 3, true, NOW(), NOW(), NOW()),
('user_002', 3, 'video', 450, NULL, false, NOW(), NOW(), NOW()),

-- User 003 progress
('user_003', 1, 'video', 600, NULL, false, NOW(), NOW(), NOW());

COMMIT;

-- Verify the relationships
SELECT 'Data seeded successfully!' as status;

-- Show counts
SELECT
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM chapters) as chapters_count,
    (SELECT COUNT(*) FROM videos) as videos_count,
    (SELECT COUNT(*) FROM quiz_questions) as quiz_questions_count,
    (SELECT COUNT(*) FROM progresses) as progresses_count;

-- Show sample joined data
SELECT
    c.id as chapter_id,
    c.title as chapter_title,
    v.title as video_title,
    v.video_url,
    COUNT(q.id) as quiz_count
FROM chapters c
LEFT JOIN videos v ON c.id = v.chapter_id
LEFT JOIN quiz_questions q ON c.id = q.chapter_id
GROUP BY c.id, c.title, v.title, v.video_url
ORDER BY c.order_index;
