package data

import (
	"time"
	"github.com/cuanin/emergent-backend/models"
)

// Users contains dummy user data
var Users = []models.User{
	{
		ID:             "1",
		Email:          "user@example.com",
		Password:       "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: password123
		FullName:       "Test User",
		IsAdmin:        false,
		CreatedAt:      time.Now().Add(-7 * 24 * time.Hour),
		LastLogin:      timePtr(time.Now().Add(-1 * time.Hour)),
		EnrolledCourses: []string{"1", "2"},
		Badges:         []string{"Fast Learner"},
		Progress: map[string]int{
			"1": 30, // 30% progress in course 1
			"2": 10, // 10% progress in course 2
		},
	},
	{
		ID:             "2",
		Email:          "admin@example.com",
		Password:       "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: password123
		FullName:       "Admin User",
		IsAdmin:        true,
		CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
		LastLogin:      timePtr(time.Now().Add(-10 * time.Minute)),
		EnrolledCourses: []string{"1", "2", "3"},
		Badges:         []string{"Instructor", "Top Performer"},
		Progress: map[string]int{
			"1": 100,
			"2": 75,
			"3": 20,
		},
	},
}

// Helper function to get a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

// Courses contains dummy course data
var Courses = []models.Course{
	{
		ID:              "1",
		Title:          "Introduction to Personal Finance",
		Description:    "Learn the basics of managing your personal finances. This course covers essential topics like budgeting, saving, and investing to help you take control of your financial future.",
		Price:          49.99,
		Category:       "Finance",
		Level:          "Beginner",
		MentorName:     "John Doe",
		VideoURL:       "https://example.com/videos/finance-intro.mp4",
		PreviewVideoURL: "https://example.com/videos/finance-preview.mp4",
		Duration:       "2 hours",
		Topics:         []string{"Budgeting", "Saving", "Investing", "Debt Management"},
		CreatedAt:      time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		EnrolledCount:  150,
	},
	{
		ID:             "2",
		Title:          "Stock Market Fundamentals",
		Description:    "Understand how the stock market works and how to start investing. Learn about different investment vehicles and strategies to grow your wealth.",
		Price:          79.99,
		Category:       "Investing",
		Level:          "Intermediate",
		MentorName:     "Jane Smith",
		VideoURL:       "https://example.com/videos/stock-market.mp4",
		PreviewVideoURL: "https://example.com/videos/stock-preview.mp4",
		Duration:       "3 hours",
		Topics:         []string{"Stocks", "Bonds", "ETFs", "Market Analysis"},
		CreatedAt:      time.Now().Add(-15 * 24 * time.Hour), // 15 days ago
		EnrolledCount:  89,
	},
	{
		ID:             "3",
		Title:          "Advanced Investment Strategies",
		Description:    "Take your investing skills to the next level with advanced strategies used by professional investors.",
		Price:          129.99,
		Category:       "Investing",
		Level:          "Advanced",
		MentorName:     "Michael Johnson",
		VideoURL:       "https://example.com/videos/advanced-investing.mp4",
		PreviewVideoURL: "https://example.com/videos/advanced-preview.mp4",
		Duration:       "4 hours",
		Topics:         []string{"Options Trading", "Futures", "Hedging", "Portfolio Management"},
		CreatedAt:      time.Now().Add(-7 * 24 * time.Hour), // 7 days ago
		EnrolledCount:  42,
	},
}

// Payments contains dummy payment data
var Payments = []models.Payment{
	{
		ID:            "1",
		UserID:        "1",
		CourseID:      "1",
		Amount:        49.99,
		PaymentMethod: "credit_card",
		Status:        "completed",
		CreatedAt:     time.Now(),
	},
}

// UserCourses maps user IDs to their enrolled course IDs
var UserCourses = map[string][]string{
	"1": {"1"}, // User 1 is enrolled in Course 1
}
