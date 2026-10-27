# TodoNest

TodoNest is a full-stack task management application with user authentication, private tasks, statistics, and a modern web interface.

## Features

- User registration and login
- JWT-based authentication
- Private tasks for each user
- Create, update, complete, and delete tasks
- Task statistics
- PostgreSQL database
- REST API
- TypeScript frontend
- Docker-based development environment

## Tech Stack

### Backend

- Go
- PostgreSQL
- pgx
- JWT
- bcrypt
- Docker Compose

### Frontend

- TypeScript
- Vite
- HTML
- CSS

## Requirements

- Go
- Node.js
- npm
- Docker Desktop
- Make

## Getting Started

Clone the repository:

```bash
git clone https://github.com/stfprod17/todoapp.git
cd todoapp
```

Create the environment file:

```bash
cp env.example .env
```

Configure the required values in `.env`:

```env
HTTP_ADDR=:5050
HTTP_SHUTDOWN_TIMEOUT=30s

POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=your_database
POSTGRES_TIMEOUT=10s

LOGGER_LEVEL=DEBUG
LOGGER_FOLDER=./out/logs

JWT_SECRET=your_secret_key
TIME_ZONE=UTC
```

Never commit `.env` to the repository.

## Run the Backend

Start PostgreSQL:

```bash
make env-up
```

Run database migrations:

```bash
make migrate-up
```

Start the backend:

```bash
make todoapp-run
```

The backend API will be available at:

```text
http://localhost:5050
```

## Run the Frontend

Open another terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at:

```text
http://localhost:5173
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/users` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive a JWT token |

### Tasks

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/tasks` | Get the current user's tasks |
| `POST` | `/api/v1/tasks` | Create a new task |
| `GET` | `/api/v1/tasks/{id}` | Get a task by ID |
| `PATCH` | `/api/v1/tasks/{id}` | Update a task |
| `DELETE` | `/api/v1/tasks/{id}` | Delete a task |

### Statistics

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/statistics` | Get the current user's statistics |

Protected endpoints require:

```http
Authorization: Bearer <jwt_token>
```

## Example Requests

Register a user:

```http
POST /api/v1/users
Content-Type: application/json

{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

Login:

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

Create a task:

```http
POST /api/v1/tasks
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "Learn Go",
  "description": "Study interfaces and goroutines"
}
```

## Database Commands

Start PostgreSQL:

```bash
make env-up
```

Stop PostgreSQL:

```bash
make env-down
```

Apply migrations:

```bash
make migrate-up
```

Rollback the latest migration:

```bash
make migrate-down
```

Check running containers:

```bash
make ps
```

## Project Structure

```text
todoapp/
├── cmd/
│   └── todoapp/
├── internal/
│   └── core/
│       ├── auth/
│       ├── config/
│       ├── domain/
│       ├── errors/
│       ├── features/
│       │   ├── statistics/
│       │   ├── tasks/
│       │   └── users/
│       ├── logger/
│       ├── repository/
│       └── transport/
├── migrations/
├── frontend/
├── docker-compose.yaml
├── Makefile
├── env.example
├── go.mod
└── go.sum
```

## Testing

Run backend tests:

```bash
go test ./...
```

Build the frontend:

```bash
cd frontend
npm run build
```

## Development Commands

Create a new migration:

```bash
make migrate-create seq=your_migration_name
```

Clean application logs:

```bash
make logs-cleanup
```

Deploy the application with Docker Compose:

```bash
make todoapp-deploy
```

Stop the deployed application:

```bash
make todoapp-undeploy
```

## Future Improvements

- Password reset
- Email verification
- Refresh tokens
- Task priorities and deadlines
- Task categories
- Search and sorting
- Shared task lists
- Production deployment
- Automated tests

## License

This project is intended for educational and personal use.
