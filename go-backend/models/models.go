package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID             string         `json:"id" bson:"_id"`
	Email          string         `json:"email" bson:"email"`
	Password       string         `json:"-" bson:"password"`
	FullName       string         `json:"full_name" bson:"full_name"`
	IsAdmin        bool           `json:"is_admin" bson:"is_admin"`
	CreatedAt      time.Time      `json:"created_at" bson:"created_at"`
	LastLogin      *time.Time     `json:"last_login,omitempty" bson:"last_login,omitempty"`
	EnrolledCourses []string       `json:"enrolled_courses" bson:"enrolled_courses"`
	Badges         []string       `json:"badges" bson:"badges"`
	Progress       map[string]int `json:"progress" bson:"progress"`
}

// Course represents a course in the platform
type Course struct {
	ID              string    `json:"id" bson:"_id"`
	Title           string    `json:"title" bson:"title"`
	Description     string    `json:"description" bson:"description"`
	Price           float64   `json:"price" bson:"price"`
	Category        string    `json:"category" bson:"category"`
	Level           string    `json:"level" bson:"level"`
	MentorName      string    `json:"mentor_name" bson:"mentor_name"`
	VideoURL        string    `json:"video_url,omitempty" bson:"video_url,omitempty"`
	PreviewVideoURL string    `json:"preview_video_url,omitempty" bson:"preview_video_url,omitempty"`
	Duration        string    `json:"duration" bson:"duration"`
	Topics          []string  `json:"topics" bson:"topics"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	EnrolledCount   int       `json:"enrolled_count" bson:"enrolled_count"`
}

// Payment represents a payment transaction
type Payment struct {
	ID            string    `json:"id" bson:"_id"`
	UserID        string    `json:"user_id" bson:"user_id"`
	CourseID      string    `json:"course_id" bson:"course_id"`
	Amount        float64   `json:"amount" bson:"amount"`
	PaymentMethod string    `json:"payment_method" bson:"payment_method"`
	Status        string    `json:"status" bson:"status"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

// Request and response models
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	User        User   `json:"user"`
}

type CourseCreateRequest struct {
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Price           float64  `json:"price" binding:"required"`
	Category        string   `json:"category" binding:"required"`
	Level           string   `json:"level" binding:"required"`
	MentorName      string   `json:"mentor_name" binding:"required"`
	VideoURL        string   `json:"video_url,omitempty"`
	PreviewVideoURL string   `json:"preview_video_url,omitempty"`
	Duration        string   `json:"duration" binding:"required"`
	Topics          []string `json:"topics"`
}

type PaymentRequest struct {
	CourseID      string  `json:"course_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

// EnrolledCourse represents a course that a user is enrolled in, including progress
type EnrolledCourse struct {
	Course   Course `json:"course"`
	Progress int    `json:"progress"` // 0-100
}

// DashboardResponse represents the data returned for a user's dashboard
type DashboardResponse struct {
	EnrolledCourses []EnrolledCourse `json:"enrolled_courses"`
	TotalSpent     float64          `json:"total_spent"`
	Badges         []string         `json:"badges"`
	RecentPayments []Payment        `json:"recent_payments"`
}

type CategoryResponse struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}
