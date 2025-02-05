# Chirpy API

Chirpy is a simple yet powerful platform for sharing short posts called "chirps". This API allows users to create, read, and manage chirps, along with features like authentication, user management, and webhooks. It uses PostgreSQL as the database and provides a RESTful interface.

## Table of Contents

1.  [Features](#features)
2.  [Tech Stack](#tech-stack)
3.  [Setup Instructions](#setup-instructions)
4.  [API Documentation](#api-documentation)
    -   [Health Check](#health-check)
    -   [User  Endpoints](#user-endpoints)
    -   [Chirp Endpoints](#chirp-endpoints)
    -   [Webhook Endpoints](#webhook-endpoints)
5.  [Contributing](#contributing)
6.  [License](#license)

## Features

-   **User  Authentication**: Secure login and JWT-based authentication.
-   **Chirps**: Create, read, update, and delete chirps.
-   **Webhooks**: Integrate with third-party services for event-driven architecture.
-   **Health Checks & Metrics**: Monitor API health and performance.

## Tech Stack

-   **Backend**: Go 1.22+ (Golang)
-   **Database**: PostgreSQL
-   **Authentication**: JWT (JSON Web Tokens)
-   **Webhooks**: Polka Integration
-   **Environment Variables**: `.env` for managing configurations

## Setup Instructions

### Prerequisites

Before setting up Chirpy, ensure that you have the following installed:

-   Go 1.22+ (Download from [golang.org](https://golang.org/dl/))
-   PostgreSQL (Download from [postgresql.org](https://www.postgresql.org/download/))
-   Git (Download from [git-scm.com](https://git-scm.com/downloads))

### Clone the Repository

First, clone the repository using Git:

```bash
git clone https://github.com/half-blood-prince-2710/chirpy.git
cd chirpy
```

### Install Dependencies

Install the Go dependencies:

```bash
go mod tidy
```

### Create a `.env` File

Create a `.env` file in the root of the project with the following environment variables:

```ini
DB_URL=your_postgres_database_url
JWT_SECRET=your_jwt_secret
PLATFORM=development_or_production_mode
POLKA_KEY=your_polka_api_key`
```

### Run the Server

Start the server with the following command:

```bash
go run main.go

The server will start on `http://localhost:8080`.
```

* * * * *

API Documentation
-----------------

### Health Check

#### GET /api/healthz

Check the health of the Chirpy API.

Response:

```json
{  "status":  "OK"  }
```

### User Endpoints

#### POST /api/users

Create a new user.

Request Body:

```json
{
  "email":  "name@example.com",
  "password":  "secretpassword"
}
```

Response:

```json
{
  "id":  "a uuid",  
  "email":  "name@example.com",
  "is_chirpy_red":  false
}
```

#### POST /api/login

Login and obtain JWT tokens.

Request Body:

```json
{
   "email":  "name@example.com",
   "password":  "secretpassword"
}
```

Response:

```json
{  "token":  "your_jwt_token",  "refresh_token":  "your_refresh_token"  }
```

#### PUT /api/users

Update user information. Requires authentication via JWT.

Request Body:

```json
{   "email":  "name@example.com",   "password":  "newpassword"  }
```

* * * * *

### Chirp Endpoints

#### POST /api/chirps

Create a new chirp.

Request Body:

```json
{   "content":  "This is my first chirp!",    "author_id":  "a uuid"  }
```

Response:

```json

{    "id":  "chirp_id",    "content":  "This is my first chirp!",    "created_at":  "2025-02-05T14:42:41.780234Z"  }
```

#### GET /api/chirps

Get a list of chirps, optionally filtered by `authorid` and sorted by `created_at`.

Request Parameters:

-   `authorid`: Optional. Filter chirps by the author's user ID.
-   `sort`: Optional. Sort chirps by `created_at` in ascending or descending order.[asc|desc]

Response:

```json
[   {   "id":  "chirp_id",   "content":  "This is my first chirp!",   "created_at":  "2025-02-05T14:42:41.780234Z"   },  ...   ]
```
#### GET /api/chirps/{id}

Get a single chirp by ID.

Response:

```json

{    "id":  "chirp_id",  "content":  "This is my first chirp!",   "created_at":  "2025-02-05T14:42:41.780234Z"  }
```

#### DELETE /api/chirps/{id}

Delete a chirp by ID. Requires authentication.

* * * * *

### Webhook Endpoints

#### POST /api/polka/webhooks

Handle webhooks from Polka.

Response:

```header
HTTP Status: 204 No Content
```

Contributing
------------

We welcome contributions to the Chirpy project! Please fork the repository, create a new branch, and submit a pull request. Ensure that your code adheres to the project's coding style and includes appropriate tests.

License
-------

This project is licensed under the MIT License - see the LICENSE file for details.
