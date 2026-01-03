package main

import (
	"learning-app-backend/config"
	"learning-app-backend/database"
	"learning-app-backend/models"
	"log"
)

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDatabase(cfg)

	// Seed chapters
	chapters := []models.Chapter{
		{
			Title:       "Introduction to Programming",
			Description: "Learn the basics of programming and computer science fundamentals",
			OrderIndex:  1,
		},
		{
			Title:       "Variables and Data Types",
			Description: "Understanding variables, data types, and how to work with them",
			OrderIndex:  2,
		},
		{
			Title:       "Control Flow",
			Description: "Master if statements, loops, and conditional logic",
			OrderIndex:  3,
		},
		{
			Title:       "Functions and Methods",
			Description: "Learn how to write reusable code with functions",
			OrderIndex:  4,
		},
		{
			Title:       "Object-Oriented Programming",
			Description: "Introduction to OOP concepts and principles",
			OrderIndex:  5,
		},
	}

	for i := range chapters {
		result := database.DB.Create(&chapters[i])
		if result.Error != nil {
			log.Printf("Error creating chapter '%s': %v", chapters[i].Title, result.Error)
		} else {
			log.Printf("✓ Created chapter: %s (ID: %d)", chapters[i].Title, chapters[i].ID)
		}
	}

	// Seed videos
	videos := []models.Video{
		{ChapterID: chapters[0].ID, Title: "Programming for Beginners - Learn to Code", VideoURL: "https://www.youtube.com/watch?v=zOjov-2OZ0E", DurationSeconds: 3720},
		{ChapterID: chapters[1].ID, Title: "Variables and Data Types Explained", VideoURL: "https://www.youtube.com/watch?v=OH86oLzVzzc", DurationSeconds: 780},
		{ChapterID: chapters[2].ID, Title: "Control Flow - If Statements and Loops", VideoURL: "https://www.youtube.com/watch?v=PYTyhXKSMzE", DurationSeconds: 900},
		{ChapterID: chapters[3].ID, Title: "Functions in Programming - Complete Guide", VideoURL: "https://www.youtube.com/watch?v=nX34jsB_v3E", DurationSeconds: 1020},
		{ChapterID: chapters[4].ID, Title: "Object Oriented Programming (OOP) Explained", VideoURL: "https://www.youtube.com/watch?v=pTB0EiLXUC8", DurationSeconds: 1380},
	}

	for i := range videos {
		result := database.DB.Create(&videos[i])
		if result.Error != nil {
			log.Printf("Error creating video '%s': %v", videos[i].Title, result.Error)
		} else {
			log.Printf("✓ Created video: %s (ID: %d)", videos[i].Title, videos[i].ID)
		}
	}

	// Seed quiz questions
	quizQuestions := []models.QuizQuestion{
		// Chapter 1 questions
		{ChapterID: chapters[0].ID, QuestionText: "What is a program?", OptionA: "A set of instructions for a computer", OptionB: "A type of computer hardware", OptionC: "A programming language", OptionD: "A database", CorrectAnswer: "A", OrderIndex: 1},
		{ChapterID: chapters[0].ID, QuestionText: "Which of these is a programming language?", OptionA: "HTML", OptionB: "Python", OptionC: "CSS", OptionD: "JSON", CorrectAnswer: "B", OrderIndex: 2},

		// Chapter 2 questions
		{ChapterID: chapters[1].ID, QuestionText: "What is a variable?", OptionA: "A fixed value", OptionB: "A container for storing data", OptionC: "A type of loop", OptionD: "A function", CorrectAnswer: "B", OrderIndex: 1},
		{ChapterID: chapters[1].ID, QuestionText: "Which data type stores whole numbers?", OptionA: "String", OptionB: "Float", OptionC: "Integer", OptionD: "Boolean", CorrectAnswer: "C", OrderIndex: 2},

		// Chapter 3 questions
		{ChapterID: chapters[2].ID, QuestionText: "What does an if statement do?", OptionA: "Stores data", OptionB: "Makes decisions based on conditions", OptionC: "Creates loops", OptionD: "Defines functions", CorrectAnswer: "B", OrderIndex: 1},
		{ChapterID: chapters[2].ID, QuestionText: "Which loop runs at least once?", OptionA: "for loop", OptionB: "while loop", OptionC: "do-while loop", OptionD: "if loop", CorrectAnswer: "C", OrderIndex: 2},

		// Chapter 4 questions
		{ChapterID: chapters[3].ID, QuestionText: "What is a function?", OptionA: "A variable", OptionB: "A reusable block of code", OptionC: "A data type", OptionD: "A loop", CorrectAnswer: "B", OrderIndex: 1},
		{ChapterID: chapters[3].ID, QuestionText: "What keyword is used to return a value?", OptionA: "give", OptionB: "return", OptionC: "send", OptionD: "output", CorrectAnswer: "B", OrderIndex: 2},

		// Chapter 5 questions
		{ChapterID: chapters[4].ID, QuestionText: "What does OOP stand for?", OptionA: "Online Operating Program", OptionB: "Object-Oriented Programming", OptionC: "Ordered Operation Process", OptionD: "Optimized Output Protocol", CorrectAnswer: "B", OrderIndex: 1},
		{ChapterID: chapters[4].ID, QuestionText: "What is encapsulation?", OptionA: "Combining data and methods", OptionB: "A type of loop", OptionC: "A variable type", OptionD: "A function call", CorrectAnswer: "A", OrderIndex: 2},
	}

	for i := range quizQuestions {
		result := database.DB.Create(&quizQuestions[i])
		if result.Error != nil {
			log.Printf("Error creating quiz question: %v", result.Error)
		} else {
			log.Printf("✓ Created quiz question (ID: %d)", quizQuestions[i].ID)
		}
	}

	log.Println("\n✅ Database seeding completed successfully!")
	log.Printf("Summary: %d chapters, %d videos, %d quiz questions", len(chapters), len(videos), len(quizQuestions))
}
