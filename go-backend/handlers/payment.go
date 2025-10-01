package handlers

import (
	"net/http"
	// "strconv"
	"time"

	"github.com/cuanin/emergent-backend/data"
	"github.com/cuanin/emergent-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getPaymentsCopy returns a copy of the payments slice to avoid direct modification
func getPaymentsCopy() []models.Payment {
	return append([]models.Payment(nil), data.Payments...)
}

// getUserCoursesCopy returns a copy of the userCourses map to avoid direct modification
func getUserCoursesCopy() map[string][]string {
	copied := make(map[string][]string)
	for k, v := range data.UserCourses {
		courses := make([]string, len(v))
		copy(courses, v)
		copied[k] = courses
	}
	return copied
}

// PurchaseCourse handles course purchase
func PurchaseCourse(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	// Find course
	var course *models.Course
	var courseIndex int
	for i, c := range data.Courses {
		if c.ID == req.CourseID {
			course = &data.Courses[i]
			courseIndex = i
			break
		}
	}

	if course == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Course not found"})
		return
	}

	// Check if user is already enrolled
	for _, user := range data.Users {
		if user.ID == userID {
			for _, courseID := range user.EnrolledCourses {
				if courseID == req.CourseID {
					c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "User already enrolled in this course"})
					return
				}
			}
			break
		}
	}

	// Create payment
	payment := models.Payment{
		ID:            uuid.New().String(),
		UserID:        userID.(string),
		CourseID:      req.CourseID,
		Amount:        course.Price, // Use course price instead of request amount for security
		PaymentMethod: req.PaymentMethod,
		Status:        "completed",
		CreatedAt:     time.Now(),
	}

	// Save payment
	data.Payments = append(data.Payments, payment)

	// Update user's enrolled courses and progress
	for i := range data.Users {
		if data.Users[i].ID == userID {
			data.Users[i].EnrolledCourses = append(data.Users[i].EnrolledCourses, req.CourseID)
			if data.Users[i].Progress == nil {
				data.Users[i].Progress = make(map[string]int)
			}
			data.Users[i].Progress[req.CourseID] = 0 // Initialize progress at 0%
			break
		}
	}

	// Update course enrollment count
	data.Courses[courseIndex].EnrolledCount++

	c.JSON(http.StatusCreated, payment)
}

// GetUserDashboard returns user's dashboard data
func GetUserDashboard(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not authenticated"})
		return
	}

	// Find user
	var user *models.User
	for i, u := range data.Users {
		if u.ID == userID {
			user = &data.Users[i]
			break
		}
	}

	if user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	// Get enrolled courses with progress
	var enrolledCourses []models.EnrolledCourse
	totalSpent := 0.0

	for _, courseID := range user.EnrolledCourses {
		for _, course := range data.Courses {
			if course.ID == courseID {
				enrolledCourse := models.EnrolledCourse{
					Course:  course,
					Progress: user.Progress[courseID],
				}
				enrolledCourses = append(enrolledCourses, enrolledCourse)
				totalSpent += course.Price
				break
			}
		}
	}

	// Get user's recent payments
	recentPayments := []models.Payment{}
	for _, payment := range data.Payments {
		if payment.UserID == userID {
			recentPayments = append(recentPayments, payment)
			// Limit to last 5 payments
			if len(recentPayments) >= 5 {
				break
			}
		}
	}

	c.JSON(http.StatusOK, models.DashboardResponse{
		EnrolledCourses: enrolledCourses,
		TotalSpent:      totalSpent,
		Badges:          user.Badges,
		RecentPayments:  recentPayments,
	})
}
