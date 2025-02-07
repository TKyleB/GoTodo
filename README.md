# REST API Documentation

## Authentication & User Management

### Register a New User
**Endpoint:** `POST /api/users/register`

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Responses:**
- `201 Created`: User successfully registered.
- `400 Bad Request`: Invalid username or password.
- `409 Conflict`: Username already registered.

---

### Login a User
**Endpoint:** `POST /api/users/login`

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Responses:**
- `200 OK`: Returns JWT and refresh token.
- `401 Unauthorized`: Invalid username or password.

---

### Refresh JWT Token
**Endpoint:** `POST /api/users/refresh`

**Request Body:**
```json
{
  "refreshToken": "string"
}
```

**Responses:**
- `200 OK`: Returns new JWT token.
- `401 Unauthorized`: Expired or invalid token.

---

### Logout User
**Endpoint:** `POST /api/users/logout`

**Request Body:**
```json
{
  "refreshToken": "string"
}
```

**Responses:**
- `204 No Content`: Logout successful.

---

### Get Authenticated User
**Endpoint:** `GET /api/users`

**Headers:**
`Authorization: Bearer <token>`

**Responses:**
- `200 OK`: Returns user data.
- `401 Unauthorized`: Invalid or missing token.

---

## Snippet Management

### Create a New Snippet
**Endpoint:** `POST /api/snippets`

**Headers:**
`Authorization: Bearer <token>`

**Request Body:**
```json
{
  "language": "string",
  "snippet_text": "string",
  "snippet_desc": "string",
  "snippet_title": "string"
}
```

**Responses:**
- `201 Created`: Snippet successfully created.
- `400 Bad Request`: Missing required fields.
- `401 Unauthorized`: Invalid or missing token.

---

### Get All Snippets
**Endpoint:** `GET /api/snippets`

**Query Parameters (Optional):**
- `language` - Filter by programming language.
- `username` - Filter by author.
- `q` - Search query.
- `limit` - Number of snippets per page.
- `offset` - Pagination offset.

**Responses:**
- `200 OK`: Returns a paginated list of snippets.

---

### Get Snippet by ID
**Endpoint:** `GET /api/snippets/{id}`

**Responses:**
- `200 OK`: Returns snippet details.
- `400 Bad Request`: Invalid snippet ID.
- `404 Not Found`: Snippet not found.

---

### Delete Snippet by ID
**Endpoint:** `DELETE /api/snippets/{id}`

**Headers:**
`Authorization: Bearer <token>`

**Responses:**
- `204 No Content`: Snippet deleted.
- `400 Bad Request`: Invalid snippet ID.
- `401 Unauthorized`: User is not authorized to delete this snippet.
- `404 Not Found`: Snippet not found.

---

This documentation outlines the endpoints and expected request/response formats for your REST API.

