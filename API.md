# Coding Games API

A comprehensive REST API example for a coding games platform built with Go and Fiber v3.

## Features

- **Games Management**: Full CRUD operations for coding games
- **Player Management**: Player registration, profile management, and score tracking
- **Leaderboards**: Global and game-specific leaderboards
- **Comprehensive Error Handling**: Structured error responses with proper HTTP status codes
- **Request Validation**: Input validation for all endpoints
- **Pagination**: Support for paginated responses
- **Filtering**: Query parameter support for filtering results

## API Endpoints

### Health Check
- `GET /health` - API health check

### Games API
- `GET /api/v1/games` - List all games with filtering and pagination
- `POST /api/v1/games` - Create a new game
- `GET /api/v1/games/:id` - Get game details
- `PUT /api/v1/games/:id` - Update game
- `DELETE /api/v1/games/:id` - Delete game

### Players API
- `GET /api/v1/players` - List all players with sorting and pagination
- `POST /api/v1/players` - Create a new player
- `GET /api/v1/players/:id` - Get player details
- `PUT /api/v1/players/:id` - Update player
- `DELETE /api/v1/players/:id` - Delete player
- `POST /api/v1/players/:id/scores` - Submit a score for a game
- `GET /api/v1/players/:id/scores` - Get player's scores

### Leaderboards API
- `GET /api/v1/leaderboards/global` - Get global leaderboard
- `GET /api/v1/leaderboards/games/:gameId` - Get game-specific leaderboard

### Legacy Endpoints
- `GET /game/test` - Legacy test endpoint (backward compatibility)

## Example Usage

### 1. Get all games
```bash
curl http://localhost:8080/api/v1/games
```

### 2. Create a new game
```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Fibonacci Sequence",
    "description": "Calculate the nth Fibonacci number efficiently",
    "difficulty": "medium",
    "language": "javascript",
    "max_score": 150
  }'
```

### 3. Get games with filtering
```bash
# Filter by difficulty
curl "http://localhost:8080/api/v1/games?difficulty=easy"

# Filter by language with pagination
curl "http://localhost:8080/api/v1/games?language=go&limit=5&offset=0"
```

### 4. Create a new player
```bash
curl -X POST http://localhost:8080/api/v1/players \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_coder",
    "email": "john@example.com"
  }'
```

### 5. Submit a score
```bash
curl -X POST http://localhost:8080/api/v1/players/1/scores \
  -H "Content-Type: application/json" \
  -d '{
    "game_id": "1",
    "score": 95
  }'
```

### 6. Get leaderboards
```bash
# Global leaderboard
curl http://localhost:8080/api/v1/leaderboards/global

# Game-specific leaderboard
curl http://localhost:8080/api/v1/leaderboards/games/1
```

## Query Parameters

### Games API
- `difficulty` - Filter by difficulty (easy, medium, hard)
- `language` - Filter by programming language
- `limit` - Number of results per page (default: 10)
- `offset` - Number of results to skip (default: 0)

### Players API
- `sort` - Sort by field (rank, username, total_score, default: rank)
- `limit` - Number of results per page (default: 10)
- `offset` - Number of results to skip (default: 0)

### Leaderboards API
- `limit` - Number of results per page (default: 10)
- `offset` - Number of results to skip (default: 0)

## Data Models

### Game
```json
{
  "id": "string",
  "title": "string",
  "description": "string", 
  "difficulty": "easy|medium|hard",
  "language": "string",
  "max_score": 0,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Player
```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "total_score": 0,
  "rank": 0,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Player Score
```json
{
  "player_id": "string",
  "game_id": "string", 
  "score": 0,
  "completed": true,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Leaderboard Entry
```json
{
  "rank": 1,
  "player_id": "string",
  "username": "string",
  "score": 0
}
```

## Error Responses

All errors follow a consistent structure:

```json
{
  "error": "error_type",
  "message": "Human readable error message",
  "code": 400
}
```

Common error types:
- `invalid_request` - Malformed request body or missing required fields
- `validation_error` - Input validation failure
- `not_found` - Resource not found
- `conflict` - Resource conflict (e.g., duplicate username)
- `internal_error` - Server error

## Running the API

### Development
```bash
make dev
```

### Building
```bash
make build
./main
```

### Testing
```bash
make test
```

The API runs on port 8080 by default. You can change this by setting the `PORT` environment variable.

## Architecture

- **Fiber v3**: High-performance Go web framework
- **Google Wire**: Dependency injection
- **In-Memory Storage**: Simple storage layer for demonstration
- **Structured Logging**: JSON-formatted logs
- **Graceful Shutdown**: Proper server shutdown handling

This example demonstrates modern Go API development patterns including:
- Clean architecture with separated concerns
- Dependency injection
- Proper error handling
- Input validation
- RESTful design principles
- Comprehensive testing