# LMS Service

Learning Management System (LMS) Service for the Zplus SaaS Platform. This service handles course management, student enrollment, progress tracking, assessments, and learning analytics.

## Features

### 📚 Course Management
- Course creation and management
- Course categories and organization
- Course sections and lessons
- Video lectures and content management
- Course pricing and access control
- Instructor management

### 👥 Student Enrollment
- Student enrollment system
- Course access management
- Enrollment status tracking
- Payment integration
- Certificate management

### 📈 Progress Tracking
- Lesson completion tracking
- Course progress monitoring
- Time spent analytics
- Learning path recommendations
- Achievement badges

### 📝 Assessments & Quizzes
- Quiz creation and management
- Multiple question types
- Automatic grading
- Quiz attempts tracking
- Score and feedback management

### 📋 Assignment Management
- Assignment creation and submission
- File upload support
- Grading and feedback system
- Due date management
- Plagiarism detection (future)

### ⭐ Reviews & Ratings
- Course review system
- Rating and feedback
- Instructor ratings
- Review moderation

### 📊 Learning Analytics
- Student performance analytics
- Course completion statistics
- Learning progress reports
- Instructor dashboards
- System-wide metrics

## API Endpoints

### Health Check
- `GET /health` - Service health check
- `GET /api/v1/lms/health` - LMS service health check

### Course Categories
- `GET /api/v1/lms/categories` - List course categories
- `POST /api/v1/lms/categories` - Create course category

### Courses
- `GET /api/v1/lms/courses` - List courses with filtering
- `POST /api/v1/lms/courses` - Create new course
- `GET /api/v1/lms/courses/:id` - Get course details
- `PUT /api/v1/lms/courses/:id` - Update course
- `DELETE /api/v1/lms/courses/:id` - Delete course
- `GET /api/v1/lms/courses/:id/sections` - Get course sections
- `POST /api/v1/lms/courses/:id/sections` - Create course section
- `GET /api/v1/lms/courses/:id/lessons` - Get course lessons
- `POST /api/v1/lms/courses/:id/lessons` - Create course lesson

### Enrollments
- `GET /api/v1/lms/enrollments` - List enrollments
- `POST /api/v1/lms/enrollments` - Create enrollment
- `GET /api/v1/lms/enrollments/:id` - Get enrollment details

### Student Progress
- `GET /api/v1/lms/progress` - Get student progress
- `GET /api/v1/lms/progress/course/:courseId` - Get course progress
- `POST /api/v1/lms/progress/lesson/:lessonId` - Update lesson progress

### Quizzes
- `GET /api/v1/lms/quizzes` - List quizzes
- `POST /api/v1/lms/quizzes` - Create quiz
- `GET /api/v1/lms/quizzes/:id` - Get quiz details
- `POST /api/v1/lms/quizzes/:id/attempt` - Start quiz attempt

### Assignments
- `GET /api/v1/lms/assignments` - List assignments
- `POST /api/v1/lms/assignments` - Create assignment
- `GET /api/v1/lms/assignments/:id` - Get assignment details
- `POST /api/v1/lms/assignments/:id/submit` - Submit assignment

### Reviews
- `GET /api/v1/lms/reviews/course/:courseId` - Get course reviews
- `POST /api/v1/lms/reviews` - Create review

### Analytics
- `GET /api/v1/lms/analytics` - Get LMS analytics
- `GET /api/v1/lms/analytics/dashboard` - Get dashboard analytics

## Database Schema

The service uses PostgreSQL with the following main tables:

- `course_categories` - Course categories
- `courses` - Course catalog
- `course_sections` - Course sections/modules
- `course_lessons` - Individual lessons/lectures
- `enrollments` - Student enrollments
- `lesson_progress` - Student lesson progress
- `quizzes` - Quizzes and assessments
- `quiz_questions` - Quiz questions
- `quiz_attempts` - Student quiz attempts
- `course_reviews` - Course reviews and ratings
- `assignments` - Course assignments
- `assignment_submissions` - Assignment submissions

## Configuration

### Environment Variables

- `PORT` - Service port (default: 8085)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password

## Multi-tenant Support

The service supports multi-tenancy through:
- Tenant ID in all database operations
- Tenant-specific data isolation
- Header-based tenant identification (`X-Tenant-ID`)

## Running the Service

### Development
```bash
go run apps/backend/lms-service/cmd/main.go
```

### Docker
```bash
docker build -t lms-service -f apps/backend/lms-service/Dockerfile .
docker run -p 8085:8085 lms-service
```

### With Docker Compose
The service is included in the main docker-compose.yml file:
```bash
docker-compose up lms-service
```

## Request/Response Format

### Standard Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Paginated Response Format
```json
{
  "success": true,
  "message": "Results retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "code": 400
}
```

## Example Usage

### Create Course
```bash
curl -X POST http://localhost:8085/api/v1/lms/courses \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant123" \
  -d '{
    "title": "Introduction to Programming",
    "description": "Learn the basics of programming",
    "level": "beginner",
    "price": 99.99,
    "instructor_name": "John Doe"
  }'
```

### Enroll Student
```bash
curl -X POST http://localhost:8085/api/v1/lms/enrollments \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant123" \
  -d '{
    "course_id": 1,
    "student_id": "student123",
    "student_name": "Jane Smith",
    "student_email": "jane@example.com"
  }'
```

### Get Course Progress
```bash
curl http://localhost:8085/api/v1/lms/progress/course/1 \
  -H "X-Tenant-ID: tenant123" \
  -H "X-User-ID: student123"
```

### Get Analytics
```bash
curl http://localhost:8085/api/v1/lms/analytics \
  -H "X-Tenant-ID: tenant123"
```

## Development Status

- ✅ Database schema design (12 tables)
- ✅ Basic service structure
- ✅ API endpoint stubs
- ✅ Health check endpoints
- 🚧 Repository layer implementation
- 🚧 Service layer implementation
- 🚧 Handler implementation
- 🚧 Database integration
- 🚧 Authentication middleware
- ⏳ File upload management
- ⏳ Video streaming integration
- ⏳ Certificate generation
- ⏳ Testing suite
- ⏳ API documentation
- ⏳ Frontend integration

## Integration

The LMS service integrates with:
- **API Gateway** - Routing and authentication
- **Auth Service** - User authentication and authorization
- **Tenant Service** - Multi-tenant support
- **File Service** - Course content and assignment file management
- **Payment Service** - Course payment processing
- **Notification Service** - Student and instructor notifications (future)

## Learning Features

### Content Types Supported
- Video lectures
- Text-based lessons
- Interactive quizzes
- File downloads
- External links
- Live sessions (future)

### Assessment Types
- Multiple choice questions
- True/false questions
- Short answer questions
- Essay questions
- File submissions
- Practical assignments

### Progress Tracking
- Lesson completion percentage
- Time spent on content
- Quiz scores and attempts
- Assignment grades
- Overall course progress
- Learning streaks

## License

Part of the Zplus SaaS Platform - Internal Development Project
