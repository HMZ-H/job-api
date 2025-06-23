# Job Application API

A RESTful API for a job application platform built with Go, Gin, GORM, and PostgreSQL.

## Features

- **User Authentication**: JWT-based authentication with role-based access control
- **Two User Roles**: 
  - Companies: Can post, update, delete jobs and manage applications
  - Applicants: Can browse jobs, apply, and track applications
- **Job Management**: Full CRUD operations for job postings
- **Application System**: Apply to jobs with resume upload and cover letters
- **Search & Filtering**: Search jobs by title, location, and company name
- **Pagination**: All list endpoints support pagination
- **File Upload**: Resume upload integration with Cloudinary

## Technology Stack

- **Backend**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **File Storage**: Cloudinary
- **Validation**: go-playground/validator

## API Endpoints

### Authentication
- `POST /api/auth/signup` - User registration
- `POST /api/auth/login` - User login

### Jobs (Company Only)
- `POST /api/jobs` - Create job posting
- `PUT /api/jobs/:id` - Update job posting
- `DELETE /api/jobs/:id` - Delete job posting
- `GET /api/jobs/my-jobs` - Get company's job postings
- `GET /api/jobs/:id/applications` - Get applications for a job

### Jobs (Applicant Only)
- `GET /api/jobs` - Browse available jobs (with filters)
- `POST /api/jobs/:id/apply` - Apply to a job

### Jobs (Both Roles)
- `GET /api/jobs/:id` - Get job details

### Applications
- `GET /api/applications/my-applications` - Get applicant's applications (Applicant only)
- `PUT /api/applications/:id/status` - Update application status (Company only)

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database
- Cloudinary account (for file uploads)

### Installation

1. **Clone the repository**
   \`\`\`bash
   git clone <repository-url>
   cd job-api
   \`\`\`

2. **Install dependencies**
   \`\`\`bash
   go mod download
   \`\`\`

3. **Set up environment variables**
   \`\`\`bash
   cp .env.example .env
   \`\`\`
   
   Edit `.env` file with your configuration:
   \`\`\`env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=your-password
   DB_NAME=job_api
   DB_PORT=5432
   JWT_SECRET=your-super-secret-jwt-key
   CLOUDINARY_CLOUD_NAME=your-cloud-name
   CLOUDINARY_API_KEY=your-api-key
   CLOUDINARY_API_SECRET=your-api-secret
   PORT=8080
   \`\`\`

4. **Create PostgreSQL database**
   \`\`\`sql
   CREATE DATABASE job_api;
   \`\`\`

5. **Run the application**
   \`\`\`bash
   go run main.go
   \`\`\`

The server will start on `http://localhost:8080`

### Database Migration

The application automatically creates the required tables on startup using GORM's AutoMigrate feature.

## API Usage Examples

### User Signup
\`\`\`bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "role": "applicant"
  }'
\`\`\`

### User Login
\`\`\`bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
\`\`\`

### Create Job (Company)
\`\`\`bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Software Engineer",
    "description": "We are looking for a skilled software engineer to join our team...",
    "location": "Remote"
  }'
\`\`\`

### Browse Jobs (Applicant)
\`\`\`bash
curl -X GET "http://localhost:8080/api/jobs?page=1&page_size=10&title=engineer" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
\`\`\`

### Apply for Job (Applicant)
\`\`\`bash
curl -X POST http://localhost:8080/api/jobs/JOB_ID/apply \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "resume_link": "https://res.cloudinary.com/your-cloud/raw/upload/resume.pdf",
    "cover_letter": "I am very interested in this position..."
  }'
\`\`\`

## Response Format

### Base Response
\`\`\`json
{
  "success": true,
  "message": "Operation successful",
  "object": { /* response data */ },
  "errors": null
}
\`\`\`

### Paginated Response
\`\`\`json
{
  "success": true,
  "message": "Data retrieved successfully",
  "object": [ /* array of items */ ],
  "page_number": 1,
  "page_size": 10,
  "total_size": 50,
  "errors": null
}
\`\`\`

## Validation Rules

### User Registration
- **Name**: Required, alphabets only
- **Email**: Required, valid email format, unique
- **Password**: Required, minimum 8 characters, must contain:
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one digit
  - At least one special character
- **Role**: Required, must be "applicant" or "company"

### Job Creation
- **Title**: Required, 1-100 characters
- **Description**: Required, 20-2000 characters
- **Location**: Optional

### Job Application
- **Resume Link**: Required, valid URL
- **Cover Letter**: Optional, maximum 200 characters

## Security Features

- Password hashing using bcrypt
- JWT token-based authentication
- Role-based access control
- Input validation and sanitization
- CORS support
- SQL injection prevention through GORM

## Error Handling

The API returns appropriate HTTP status codes and error messages:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource (e.g., email already exists)
- `500 Internal Server Error`: Server-side errors

## Development

### Project Structure
\`\`\`
job-api/
├── config/          # Database configuration
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication and authorization middleware
├── models/          # Database models and response structures
├── utils/           # Utility functions (JWT, validation, etc.)
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
└── README.md        # This file
\`\`\`

### Running Tests
\`\`\`bash
go test ./...
\`\`\`

## Deployment

For production deployment:

1. Set `GIN_MODE=release` in environment variables
2. Use a production-grade PostgreSQL database
3. Set strong JWT secret
4. Configure proper CORS settings
5. Use HTTPS in production
6. Set up proper logging and monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
