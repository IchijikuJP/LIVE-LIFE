# LiveLife

LiveLife MVP project for Tokyo livehouse event information, CD/shop content, articles, social links, and join/contact intake.

## Local Preview

The local preview currently uses the Go backend with a temporary static page so the project can run even before Node/npm is installed.

```bash
cd backend
go run ./cmd/server
```

Open:

```text
http://localhost:8080
```

Health check:

```text
http://localhost:8080/api/health
```

## Planned App Stack

- Frontend: React + Vite + TypeScript + Tailwind CSS
- Backend: Go
- Database: SQLite for MVP, PostgreSQL later
- Deployment: Docker Compose, Nginx, Alibaba Cloud Tokyo server

## Project Layout

```text
backend/      Go API and local static preview
frontend/     React/Vite/Tailwind source skeleton
docs/         Deployment and server notes
```
