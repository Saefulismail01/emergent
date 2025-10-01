package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/cuanin/emergent-backend/data"
	"github.com/cuanin/emergent-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Users returns a copy of the users slice to avoid direct modification
func getUsers() []models.User {
	return append([]models.User(nil), data.Users...)
}

// Dummy courses data
var courses = []models.Course{
	{
		ID:          "1",
		Title:       "Introduction to Personal Finance",
		Description: "Learn the basics of managing your personal finances.",
		Price:       49.99,
		Category:    "Finance",
		Level:       "Beginner",
		MentorName:  "John Doe",
		Duration:    "2 hours",
		Topics:      []string{"Budgeting", "Saving", "Investing"},
		CreatedAt:   time.Now(),
	},
	// Add more dummy courses as needed
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generateJWT(user models.User) (string, error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Validate email and password
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		c.JSON(400, models.ErrorResponse{Error: "Email, password, and full name are required"})
		return
	}

	// Check if user already exists
	for _, user := range data.Users {
		if user.Email == req.Email {
			c.JSON(409, models.ErrorResponse{Error: "Email already registered"})
			return
		}
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: "Failed to process password"})
		return
	}

	// Create new user with UUID
	newUser := models.User{
		ID:              uuid.New().String(),
		Email:          req.Email,
		Password:       hashedPassword,
		FullName:       req.FullName,
		IsAdmin:        false,
		CreatedAt:      time.Now(),
		EnrolledCourses: []string{},
		Badges:         []string{"New User"}, // Add default badge
		Progress:       make(map[string]int),
	}

	// Add to users slice
	data.Users = append(data.Users, newUser)

	// Generate JWT token
	token, err := generateJWT(newUser)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: "Failed to generate authentication token"})
		return
	}

	// Return user data without password
	newUser.Password = ""
	c.JSON(201, gin.H{
		"token": token,
		"user":  newUser,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		c.JSON(400, models.ErrorResponse{Error: "Email and password are required"})
		return
	}

	// Find user by email
	var foundUser *models.User
	fmt.Printf("Debug: Searching for user with email: %s\n", req.Email)
	fmt.Printf("Debug: Available users: %+v\n", data.Users)
	
	for i, user := range data.Users {
		fmt.Printf("Debug: Checking user %d - Email: %s\n", i, user.Email)
		if user.Email == req.Email {
			foundUser = &data.Users[i] // Get pointer to the actual user in the slice
			fmt.Printf("Debug: Found user: %+v\n", *foundUser)
			break
		}
	}

	// Check if user exists and password is correct
	if foundUser == nil {
		// Don't reveal if user exists or not for security
		fmt.Println("Debug: User not found")
		c.JSON(401, models.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Verify password
	fmt.Printf("Debug: Verifying password for user: %s\n", foundUser.Email)
	fmt.Printf("Debug: Stored hash: %s\n", foundUser.Password)
	
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		fmt.Printf("Debug: Password verification failed: %v\n", err)
		c.JSON(401, models.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(*foundUser)
	if err != nil {
		c.JSON(500, models.ErrorResponse{Error: "Failed to generate authentication token"})
		return
	}

	// Update last login time (in a real app, you'd save this to the database)
	now := time.Now()
	foundUser.LastLogin = &now

	// Return user data without password
	responseUser := *foundUser
	responseUser.Password = ""
	
	c.JSON(200, gin.H{
		"token": token,
		"user":  responseUser,
	})
}
