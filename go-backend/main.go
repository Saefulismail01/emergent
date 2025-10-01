package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cuanin/emergent-backend/data"
	"github.com/cuanin/emergent-backend/handlers"
	"github.com/cuanin/emergent-backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Set JWT secret key
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-secret-key-2024")
	}

	r := gin.Default()

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	// Routes
	setupRoutes(r)

	// Start server
	port := ":8080"
	log.Printf("Server running on port %s", port)
	log.Fatal(r.Run(port))
}

func setupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Courses routes
		courses := v1.Group("/courses")
		{
			courses.GET("", handlers.GetCourses)
			courses.GET("/:id", handlers.GetCourse)
			courses.POST("", handlers.CreateCourse)
		}

		// Payment routes
		payment := v1.Group("/payment")
		payment.Use(authMiddleware())
		{
			payment.POST("", handlers.PurchaseCourse)
		}

		// User dashboard
		user := v1.Group("/user")
		user.Use(authMiddleware())
		{
			user.GET("/dashboard", handlers.GetUserDashboard)
		}

		// Categories
		v1.GET("/categories", handlers.GetCategories)
	}
}

// Auth middleware to validate JWT token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid authorization format"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if token is expired
			exp, ok := claims["exp"].(float64)
			if !ok || float64(time.Now().Unix()) > exp {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Token has expired"})
				c.Abort()
				return
			}

			// Add user ID to context
			userID, ok := claims["user_id"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token claims"})
				c.Abort()
				return
			}

			// Check if user exists
			userExists := false
			for _, user := range data.Users {
				if user.ID == userID {
					userExists = true
					break
				}
			}

			if !userExists {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
				c.Abort()
				return
			}

			// Set user ID and is_admin flag in context
			c.Set("user_id", userID)
			isAdmin, _ := claims["is_admin"].(bool)
			c.Set("is_admin", isAdmin)

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}
	}
}
