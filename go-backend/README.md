# Emergent Go Backend

This is a Go-based backend for the Emergent financial education platform, built with Gin framework.

## Features

- User authentication (JWT)
- Course management
- Payment processing
- User dashboard
- Category listing

## Prerequisites

- Go 1.20 or higher
- Git

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/emergent-backend.git
   cd emergent-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   Then edit the `.env` file with your configuration.

4. Run the application:
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Courses
- `GET /api/courses` - Get all courses (filter with query params: `category`, `level`)
- `GET /api/courses/:id` - Get a single course
- `POST /api/courses` - Create a new course (admin only)

### Payment
- `POST /api/payment` - Purchase a course (requires authentication)

### User
- `GET /api/user/dashboard` - Get user dashboard (requires authentication)

### Categories
- `GET /api/categories` - Get all course categories with counts

## Dummy Data

The application comes with some dummy data for testing:

### Users
- **Regular User**
  - Email: user@example.com
  - Password: password123

- **Admin User**
  - Email: admin@example.com
  - Password: admin123

## Environment Variables

- `JWT_SECRET`: Secret key for JWT token signing (default: 'your-secret-key-2024')
- `PORT`: Port to run the server on (default: 8080)

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o emergent-backend
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
