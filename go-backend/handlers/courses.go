package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cuanin/emergent-backend/data"
	"github.com/cuanin/emergent-backend/models"
	"github.com/gin-gonic/gin"
)

// getCoursesCopy returns a copy of the courses slice to avoid direct modification
func getCoursesCopy() []models.Course {
	return append([]models.Course(nil), data.Courses...)
}

// GetCourses returns a list of courses with optional filtering
func GetCourses(c *gin.Context) {
	// Get query parameters
	category := c.Query("category")
	level := c.Query("level")

	// Get a copy of courses
	courses := getCoursesCopy()
	
	// Filter courses based on query parameters
	filteredCourses := make([]models.Course, 0)
	for _, course := range courses {
		if (category == "" || course.Category == category) &&
			(level == "" || course.Level == level) {
			filteredCourses = append(filteredCourses, course)
		}
	}

	c.JSON(http.StatusOK, filteredCourses)
}

// GetCourse returns a single course by ID
func GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	// Find course by ID
	for _, course := range getCoursesCopy() {
		if course.ID == courseID {
			c.JSON(http.StatusOK, course)
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Course not found"})
}

// CreateCourse creates a new course (admin only)
func CreateCourse(c *gin.Context) {
	// In a real app, check if user is admin
	// user := c.MustGet("user").(models.User)
	// if !user.IsAdmin {
	// 	c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Admin access required"})
	// 	return
	// }

	var req models.CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	// Create new course (in a real app, save to database)
	newCourse := models.Course{
		ID:              strconv.Itoa(len(courses) + 1),
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		Category:        req.Category,
		Level:           req.Level,
		MentorName:      req.MentorName,
		VideoURL:        req.VideoURL,
		PreviewVideoURL: req.PreviewVideoURL,
		Duration:        req.Duration,
		Topics:          req.Topics,
		CreatedAt:       time.Now(),
		EnrolledCount:   0,
	}

	// In a real app, save the course to database here
	data.Courses = append(data.Courses, newCourse)

	c.JSON(http.StatusCreated, newCourse)
}

// GetCategories returns a list of course categories with counts
func GetCategories(c *gin.Context) {
	// Count courses by category
	categoryCount := make(map[string]int)
	for _, course := range getCoursesCopy() {
		categoryCount[course.Category]++
	}

	// Convert to response format
	var response []models.CategoryResponse
	for name, count := range categoryCount {
		response = append(response, models.CategoryResponse{
			Name:  name,
			Count: count,
		})
	}

	c.JSON(http.StatusOK, response)
}
