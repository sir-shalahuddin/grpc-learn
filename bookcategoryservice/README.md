# Category Service

The **Category Service** is a microservice responsible for managing book categories within the library system. It allows for the organization of books into different categories, facilitating easier searching and filtering for users.

## Key Features

- **Category Management**: Enables the addition, updating, and deletion of book categories.
- **Unique Category Names**: Ensures that each category has a unique name to avoid duplication.
- **Category Retrieval**: Provides functionality to retrieve all categories or specific categories as needed.

## Database Structure

The Category Service utilizes a PostgreSQL database to store information about book categories.

### Table: `book_categories`

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE book_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE
);
```


| Column | Data Type   | Description                                      |
|--------|-------------|--------------------------------------------------|
| `id`   | UUID        | Primary key, a unique identifier for each category (auto-generated). |
| `name` | VARCHAR(255)| The name of the category (must be unique).      |


## API Documentation

The API documentation for this project is available and can be accessed through Swagger. It provides a comprehensive overview of all available endpoints, including request and response formats.

You can explore the API documentation at the following URL:

- [Swagger Documentation](http://book-category-rest.sirlearn.my.id/swagger)

### How to Use the API

To use the API, ensure that you have the correct authentication (if required) and follow the structure of requests as outlined in the Swagger documentation. You can test the endpoints directly within the Swagger interface or use tools like `curl` or Postman.