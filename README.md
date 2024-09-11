# gRPC-Learn

## Overview

gRPC-Learn is a microservices-based application designed to demonstrate the use of gRPC in a multi-service environment. The application is structured into three main services:

- **Book Service:** Manages the books available in the system.
- **Book Category Service:** Manages the categories of the books.
- **User Service:** Manages user information and authentication.

## Architecture

The application is built using the following technologies:

- **Go:** The programming language used for implementing the services.
- **gRPC:** For efficient, low-latency communication between services.
- **Docker:** Containerization of services to ensure consistency across environments.
- **PostgreSQL:** Database used for persistent data storage.

### Prerequisites

- [Docker](https://www.docker.com/) installed on your machine.
- [Go](https://golang.org/) installed for running and building services.

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/sir-shalahuddin/grpc-learn.git
    cd grpc-learn
    ```

2. **Run the application using Docker Compose:**

    ```bash
    docker-compose up --build
    ```

   This command will build and start all services in the project.

## Usage

### Accessing the Services

- **Book Service:** Runs on `localhost:<book-service-port>`
- **Book Category Service:** Runs on `localhost:<book-category-service-port>`
- **User Service:** Runs on `localhost:<user-service-port>`

## Live Project

The gRPC-Learn project is currently deployed and accessible via Cloud Run. Below are the details to access and interact with the live services.

### Available Services

- **Book Service:**  
  - **REST API:** [book-rest.sirlearn.my.id](http://book-rest.sirlearn.my.id)

- **Book Category Service:**  
  - **REST API:** [book-category-rest.sirlearn.my.id](http://book-category-rest.sirlearn.my.id)  
  - **gRPC API:** [book-category-grpc.sirlearn.my.id](http://book-category-grpc.sirlearn.my.id)

- **User Service:**  
  - **REST API:** [user-rest.sirlearn.my.id](http://user-rest.sirlearn.my.id)  
  - **gRPC API:** [user-grpc.sirlearn.my.id](http://user-grpc.sirlearn.my.id)
