# Book Service

The **Book Service** is a microservice designed to manage the library’s book inventory and borrowing activities. It facilitates interactions between users and the book database, allowing users to view, borrow, and return books.

## Key Features

- **Book Management**: Enables the addition, updating, and deletion of books in the library catalog.
- **Borrowing Records**: Tracks borrowing activities, including which user borrowed a book, the borrowing date, and due dates for returns.
- **Search Functionality**: Allows users to search for books based on various criteria such as title, author, or category.

## Database Setup

### Database Structure

We are using **PostgreSQL** as the primary database for this application. Below is an overview of the database schema, which includes the `books` and `borrowing_records` tables.

#### Table: `books`

The `books` table stores information about the books available in the library. Each book has a unique ID and contains relevant metadata, such as title, author, and publication details.

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13),
    published_date DATE,
    category_id UUID,
    stock INT,
    added_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version INT DEFAULT 0
);
```

| Column         | Data Type        | Description                                                           |
|----------------|------------------|-----------------------------------------------------------------------|
| `id`           | UUID             | Primary key, a unique identifier for each book (auto-generated).      |
| `title`        | VARCHAR(255)      | The title of the book.                                               |
| `author`       | VARCHAR(255)      | The author of the book.                                              |
| `isbn`         | VARCHAR(13)      | The ISBN of the book (optional).                                     |
| `published_date`| DATE            | The date when the book was published (optional).                     |
| `category_id`  | UUID             | Foreign key referencing the category of the book (optional).         |
| `stock`        | INT              | The number of copies of the book available in the library.           |
| `added_by`     | UUID             | Foreign key, representing the user who added the book (optional).    |
| `created_at`   | TIMESTAMP        | The timestamp when the book was added (auto-generated).              |
| `updated_at`   | TIMESTAMP        | The timestamp when the book details were last updated (auto-generated).|
| `version`      | INT              | Versioning field, useful for optimistic concurrency control.         |

#### Table: `borrowing_record`
The `borrowing_records` table keeps track of all book borrowing activities. It records when a book was borrowed, who borrowed it, and when it is due for return.

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE borrowing_records (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  book_id UUID REFERENCES books(id) ON DELETE CASCADE,
  user_id UUID,
  borrowed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  due_date TIMESTAMP NOT NULL,
  returned_at TIMESTAMP
);

```

## API Documentation

The API documentation for this project is available and can be accessed through Swagger. It provides a comprehensive overview of all available endpoints, including request and response formats.

You can explore the API documentation at the following URL:

- [Swagger Documentation](http://book-rest.sirlearn.my.id/swagger)

### How to Use the API

To use the API, ensure that you have the correct authentication (if required) and follow the structure of requests as outlined in the Swagger documentation. You can test the endpoints directly within the Swagger interface or use tools like `curl` or Postman.