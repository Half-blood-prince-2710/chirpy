Chirpy API

Chirpy is a simple backend application built in Go that supports features like user authentication, chirp creation, and management, with an emphasis on secure access and real-time data handling. This application is designed to interact with a PostgreSQL database and provides a RESTful API with JWT authentication.

Features:

- User Management: Register, update, and login users.

- Chirp Management: Create, view, and delete chirps.

- Webhooks: Handle third-party webhooks with security.

- Health Check: Monitor server health and metrics.

Requirements:

- Go 1.22+

- PostgreSQL database

- Environment variables file (.env)

Installation:

1\. Clone the repository:

   git clone https://github.com/half-blood-prince-2710/chirpy.git

   cd chirpy

2\. Install dependencies:

   go mod tidy

3\. Create a .env file in the root directory and add the following environment variables:

   DB_URL=your_database_url

   JWT_SECRET=your_jwt_secret

   PLATFORM=development_or_production

   POLKA_KEY=your_polka_key

4\. Build the project:

   go build -o chirpy

5\. Run the application:

   ./chirpy

API Endpoints:

Public Routes:

- POST /api/login - Log in a user and get JWT tokens.

- POST /api/refresh - Refresh JWT token.

- POST /api/revoke - Revoke a refresh token.

Protected Routes (Require JWT):

- POST /api/users - Create a new user.

- PUT /api/users - Update user information.

- POST /api/chirps - Create a new chirp.

- GET /api/chirps - Get a list of chirps.

- GET /api/chirps/{id} - Get a specific chirp by ID.

- DELETE /api/chirps/{id} - Delete a chirp.

Admin Routes:

- POST /admin/reset - Reset the application (admin use).

- GET /admin/metrics - Get server metrics.

Webhooks:

- POST /api/polka/webhooks - Receive webhook events from third-party services.

Database Setup:

1\. Set up a PostgreSQL database and create the required tables.

2\. Update the DB_URL in the .env file with your database connection string.

Contributing:

Feel free to fork the repository, make changes, and submit pull requests. If you find any issues or have suggestions, open an issue on the GitHub repository.

License:

This project is licensed under the MIT License.
