-- Seed Data for LearnHub Backend
-- Run this with: psql $DATABASE_URL -f seed_data.sql

BEGIN;

-- Insert sample chapters
INSERT INTO chapters (title, description, order_index, created_at, updated_at) VALUES
('Introduction to Programming', 'Learn the basics of programming and computer science fundamentals', 1, NOW(), NOW()),
('Variables and Data Types', 'Understanding variables, data types, and how to work with them', 2, NOW(), NOW()),
('Control Flow', 'Master if statements, loops, and conditional logic', 3, NOW(), NOW()),
('Functions and Methods', 'Learn how to write reusable code with functions', 4, NOW(), NOW()),
('Object-Oriented Programming', 'Introduction to OOP concepts and principles', 5, NOW(), NOW());

-- Insert sample videos (assuming chapters have IDs 1-5)
INSERT INTO videos (chapter_id, title, video_url, duration_seconds, created_at, updated_at) VALUES
(1, 'Programming for Beginners - Learn to Code', 'https://www.youtube.com/watch?v=zOjov-2OZ0E', 3720, NOW(), NOW()),
(2, 'Variables and Data Types Explained', 'https://www.youtube.com/watch?v=OH86oLzVzzc', 780, NOW(), NOW()),
(3, 'Control Flow - If Statements and Loops', 'https://www.youtube.com/watch?v=PYTyhXKSMzE', 900, NOW(), NOW()),
(4, 'Functions in Programming - Complete Guide', 'https://www.youtube.com/watch?v=nX34jsB_v3E', 1020, NOW(), NOW()),
(5, 'Object Oriented Programming (OOP) Explained', 'https://www.youtube.com/watch?v=pTB0EiLXUC8', 1380, NOW(), NOW());

-- Insert sample quiz questions
INSERT INTO quiz_questions (chapter_id, question_text, option_a, option_b, option_c, option_d, correct_answer, order_index, created_at, updated_at) VALUES
(1, 'What is a program?', 'A set of instructions for a computer', 'A type of computer hardware', 'A programming language', 'A database', 'A', 1, NOW(), NOW()),
(1, 'Which of these is a programming language?', 'HTML', 'Python', 'CSS', 'JSON', 'B', 2, NOW(), NOW()),

(2, 'What is a variable?', 'A fixed value', 'A container for storing data', 'A type of loop', 'A function', 'B', 1, NOW(), NOW()),
(2, 'Which data type stores whole numbers?', 'String', 'Float', 'Integer', 'Boolean', 'C', 2, NOW(), NOW()),

(3, 'What does an if statement do?', 'Stores data', 'Makes decisions based on conditions', 'Creates loops', 'Defines functions', 'B', 1, NOW(), NOW()),
(3, 'Which loop runs at least once?', 'for loop', 'while loop', 'do-while loop', 'if loop', 'C', 2, NOW(), NOW()),

(4, 'What is a function?', 'A variable', 'A reusable block of code', 'A data type', 'A loop', 'B', 1, NOW(), NOW()),
(4, 'What keyword is used to return a value?', 'give', 'return', 'send', 'output', 'B', 2, NOW(), NOW()),

(5, 'What does OOP stand for?', 'Online Operating Program', 'Object-Oriented Programming', 'Ordered Operation Process', 'Optimized Output Protocol', 'B', 1, NOW(), NOW()),
(5, 'What is encapsulation?', 'Combining data and methods', 'A type of loop', 'A variable type', 'A function call', 'A', 2, NOW(), NOW());

COMMIT;
