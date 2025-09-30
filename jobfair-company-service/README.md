# Company Service

Service untuk mengelola profil perusahaan, job posting, dan aplikasi lamaran kerja dalam platform JobFair.

## Features

### Company Management
- Create, Read, Update company profile
- Upload company logo, banner, video, dan gallery
- Company profile dengan informasi lengkap (alamat, social media, dll)
- Company verification status

### Job Management
- Create, Read, Update, Delete job postings
- Job status management (draft, active, paused, closed, expired)
- Job filtering (by status, type, location)
- Job view counter

### Application Management
- View job applications
- Update application status (applied, shortlisted, interview, hired, rejected)
- Filter applications by status
- Application statistics

### Dashboard & Analytics
- Dashboard statistics (jobs posted, applicants, views, positions)
- Company analytics (booth visits, profile views, job views, applications)
- Applicant status breakdown

## Tech Stack

- **Language**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT

## API Endpoints

### Public Endpoints
- `GET /api/v1/companies` - List all companies
- `GET /api/v1/companies/:id` - Get company by ID
- `GET /api/v1/jobs/:id` - Get job by ID

### Protected Endpoints (Requires JWT)

#### Company
- `GET /api/v1/my-company` - Get current user's company
- `POST /api/v1/companies` - Create company
- `PUT /api/v1/companies/:id` - Update company
- `POST /api/v1/companies/:id/logo` - Upload logo
- `POST /api/v1/companies/:id/banner` - Upload banner
- `POST /api/v1/companies/:id/videos` - Upload video
- `POST /api/v1/companies/:id/gallery` - Upload gallery
- `GET /api/v1/companies/:id/analytics` - Get analytics
- `GET /api/v1/dashboard` - Get dashboard stats

#### Jobs
- `GET /api/v1/jobs` - List jobs
- `POST /api/v1/jobs` - Create job
- `PUT /api/v1/jobs/:id` - Update job
- `DELETE /api/v1/jobs/:id` - Delete job
- `POST /api/v1/jobs/:id/publish` - Publish job
- `POST /api/v1/jobs/:id/close` - Close job

#### Applications
- `GET /api/v1/applications` - List applications
- `GET /api/v1/applications/:id` - Get application
- `GET /api/v1/jobs/:job_id/applications` - Get applications by job
- `PUT /api/v1/applications/:id/status` - Update application status
- `GET /api/v1/applications/stats` - Get application stats

## Database Schema

### Tables
1. `companies` - Company profiles
2. `company_analytics` - Analytics data
3. `company_media` - Media files
4. `jobs` - Job postings
5. `job_applications` - Job applications

## Setup & Installation

1. Copy environment file:
```bash
cp .env.example .env
```

2. Update `.env` with your configuration

3. Run migrations:
```bash
# Migration files are in /migrations directory
# Run migrations using your preferred tool
```

4. Install dependencies:
```bash
go mod download
```

5. Run the service:
```bash
go run cmd/main.go
```

## Project Structure

```
jobfair-company-service/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── handlers/               # HTTP handlers
│   │   ├── company_handler.go
│   │   ├── job_handler.go
│   │   └── application_handler.go
│   ├── middleware/             # Middlewares
│   ├── models/                 # Data models
│   │   ├── company.go
│   │   ├── job.go
│   │   └── api_response.go
│   ├── repository/             # Database repositories
│   │   ├── company_repository.go
│   │   ├── job_repository.go
│   │   └── application_repository.go
│   ├── services/               # Business logic
│   │   ├── company_service.go
│   │   ├── job_service.go
│   │   └── application_service.go
│   └── utils/                  # Utility functions
├── migrations/                 # Database migrations
├── pkg/
│   └── database/              # Database connection
├── go.mod
└── go.sum
```

## Development

### Adding New Features
1. Create model in `internal/models/`
2. Create repository in `internal/repository/`
3. Create service in `internal/services/`
4. Create handler in `internal/handlers/`
5. Register routes in `cmd/main.go`
6. Create migration in `migrations/`

### Testing
```bash
go test ./...
```

## Docker

Build image:
```bash
docker build -t jobfair-company-service .
```

Run container:
```bash
docker run -p 8081:8081 --env-file .env jobfair-company-service
```

## Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request