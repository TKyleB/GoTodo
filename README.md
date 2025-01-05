# GoTodo REST API
An API for tracking todos using implemented using golang

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

