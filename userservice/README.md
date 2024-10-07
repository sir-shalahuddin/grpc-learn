# User Service

The **User Service** is a microservice responsible for managing user-related functionalities within the application. It handles user registration, authentication, and profile management, ensuring secure access to the system.

## Key Features

- **User Registration**: Allows new users to create an account with necessary information such as name, email, and password.
- **Authentication**: Implements JWT-based authentication to ensure secure access to protected resources.
- **Profile Management**: Enables users to view and update their profile information.
- **Role Management**: Supports different user roles (e.g., regular user, librarian, super admin) for managing access to various functionalities within the application.
## Database Setup

### Database Structure

We are using **PostgreSQL** as the primary database for this application. Below is an overview of the database schema, which includes the `users` table and an enum type for user roles.

#### Enum Type: `user_role`

We have defined a custom enum type `user_role` to represent the different roles available in the system:

- `user`: Standard user.
- `super admin`: Admin with full privileges.
- `librarian`: User responsible for managing books and library operations.

```sql
-- Create an enum type for roles
CREATE TYPE user_role AS ENUM ('user', 'super admin', 'librarian');
```

#### Table: `users`

The `users` table stores information about the registered users of the application. Each user has a role defined by the `user_role` enum.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

| Column         | Data Type                     | Description                                                                 |
|----------------|-------------------------------|-----------------------------------------------------------------------------|
| `id`           | UUID                          | Primary key, a unique identifier for each user (auto-generated).            |
| `name`         | VARCHAR(100)                  | The full name of the user (maximum 100 characters).                         |
| `email`        | VARCHAR(100)                  | The user’s email address (must be unique).                                  |
| `password_hash`| VARCHAR(255)                  | The hashed password for the user.                                           |
| `role`         | `user_role`                   | The role of the user, can be `user`, `super admin`, or `librarian`. Defaults to `user`. |
| `created_at`   | TIMESTAMP WITH TIME ZONE      | The timestamp when the user was created (auto-generated).                   |
| `updated_at`   | TIMESTAMP WITH TIME ZONE      | The timestamp when the user's information was last updated (auto-generated).|


## API Documentation

The API documentation for this project is available and can be accessed through Swagger. It provides a comprehensive overview of all available endpoints, including request and response formats.

You can explore the API documentation at the following URL:

- [Swagger Documentation](http://user-rest.sirlearn.my.id/swagger)

### How to Use the API

To use the API, ensure that you have the correct authentication (if required) and follow the structure of requests as outlined in the Swagger documentation. You can test the endpoints directly within the Swagger interface or use tools like `curl` or Postman.
