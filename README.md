# Chess Puzzle API

A Go-powered REST API to solve chess puzzles, with a simplified Glicko rating system that rewards good decisions, and penalises bad ones - because confidence has never been a valid chess strategy. Inspired by Lichess Puzzles.

**Base URL:** `https://chess-puzzles-api-production.up.railway.app`  
**Auth:** Bearer token via `Authorization` header  
**Content-Type:** `application/json`

---

## Authentication

### Register
```
POST /auth/register
```
**Body:**
```json
{
  "username": "string (3-30 chars, alphanumeric + underscore)",
  "email": "string (valid email)",
  "password": "string (8-72 chars)"
}
```
**Response `201`:**
```json
{
  "token": "eyJ...",
  "user": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "rating": 1000.0,
    "rating_deviation": 350.0,
    "created_at": "timestamp"
  }
}
```
Sets `refresh_token` httpOnly cookie.

---

### Login
```
POST /auth/login
```
**Body:**
```json
{
  "email": "string",
  "password": "string"
}
```
**Response `200`:**
```json
{
  "token": "eyJ..."
}
```
Sets `refresh_token` httpOnly cookie.

---

### Refresh Token
```
POST /auth/refresh
```
Reads `refresh_token` cookie automatically.  
**Response `200`:**
```json
{
  "token": "eyJ..."
}
```
Rotates refresh token — sets new `refresh_token` cookie.

---

### Logout
```
POST /auth/logout
```
🔒 Protected  
Invalidates refresh token, clears cookie.  
**Response `200`:**
```json
{
  "message": "logged out successfully"
}
```

---

## Puzzles

### Get Random Puzzle
```
GET /puzzles/random
```
🔒 Protected  
**Response `200`:**
```json
{
  "id": "00008",
  "fen": "r6k/pp2r2p/...",
  "moves": ["f6h5", "e6e7", "h5f4"],
  "rating": 1862.0,
  "rating_deviation": 73.0,
  "popularity": 96,
  "nb_plays": 36672,
  "themes": ["crushing", "middlegame"],
  "created_at": "timestamp"
}
```
> `moves[0]` is the opponent's setup move — apply it automatically.  
> User plays `moves[1]` onward.

---

### Get Puzzle by ID
```
GET /puzzles/{id}
```
🔒 Protected  
**Response `200`:** same as above  
**Response `404`:** puzzle not found

---

### Solve Puzzle
```
POST /puzzles/{id}/solve
```
🔒 Protected  
**Body:**
```json
{
  "moves": ["e2e4", "e7e5"]
}
```
> Send only the user's moves — not `moves[0]` (the setup move).  
> Moves must be valid UCI format: `[a-h][1-8][a-h][1-8][qrbn]?`

**Response `200`:**
```json
{
  "success": true,
  "new_rating": 1087.3,
  "rating_change": 87.3
}
```

---

## Users

### Get Profile
```
GET /users/me
```
🔒 Protected  
**Response `200`:**
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "rating": 1087.3,
  "rating_deviation": 245.1,
  "created_at": "timestamp"
}
```

---

## Leaderboard

### Get Top Users
```
GET /leaderboard
```
Public — no auth required.  
**Response `200`:**
```json
[
  {"username": "magnus", "rating": 2341.0},
  {"username": "hikaru", "rating": 2187.5}
]
```

---

## Error Responses

All errors follow this shape:

| Status | Meaning |
|--------|---------|
| 400 | Bad request — invalid input |
| 401 | Unauthorized — missing/invalid/expired token |
| 404 | Resource not found |
| 429 | Too many requests — rate limited |
| 500 | Internal server error |

---

## Rate Limiting
Auth endpoints (`/auth/register`, `/auth/login`) are rate limited per IP.

## Security
- Passwords hashed with bcrypt (cost 10)
- JWT access tokens expire in 15 minutes
- Refresh tokens expire in 7 days with automatic rotation
- Refresh token stored as httpOnly cookie (XSS-safe)
- All inputs validated server-side
