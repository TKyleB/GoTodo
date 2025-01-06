# Snippetz REST API
An API for creating code snippets for your favorite programming languages

## API
The REST API to the app is described below.

### Auth

POST /api/users/register

**Request**
```
{
    "email": string,
    "password": string
}
```
**Response**
```
HTTP 201
{
    "id": number,
    "created_at": date,
    "updated_at": date,
    "email": string,
}
```

POST /api/users/login

**Request**
```
{
    "email": string,
    "password": string
}
```
**Response**
```
HTTP 200
{
    "id": number,
    "created_at": date,
    "updated_at": date,
    "email": string,
    "token": string
    "refresh_token": string
}
```

POST /api/users/refresh

**Request**
```
Authorization: Bearer TOKEN_STRING
```
**Response**
```
HTTP 200
{
    "token": string
}
```

