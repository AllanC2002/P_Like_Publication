# Like publication Microservice

## Project Overview

This service provides functionality for users to "like" publications. It is a Go-based microservice that interacts with a MongoDB database to store and manage likes. The service exposes a single API endpoint and uses JWT for authentication.

## Folder Structure

The project is organized with the following folder structure:

-   **`/` (Root Directory):**
    -   `main.go`: The main application entry point. Initializes the Gin router, sets up the API endpoint, and handles incoming requests.
    -   `go.mod`, `go.sum`: Go module files for managing project dependencies.
    -   `dockerfile`: Instructions for building the Docker image for the service.
    -   `.env`: (Typically present for local development) Used to store environment variables such as database credentials and JWT secret keys. Loaded by `godotenv`.
    -   `README.md`: This file.
-   **`connection/`**:
    -   `mongo.go`: Contains the logic for establishing and managing the connection to the MongoDB database.
-   **`functions/`**:
    -   `Like.go`: Implements the core business logic for the "like" functionality, including interacting with the database to record a like.
-   **`.github/`**:
    -   `workflows/`: Contains GitHub Actions workflow files for Continuous Integration and Continuous Deployment (CI/CD).
        -   `docker-publish.yml`: Defines the workflow for building the Docker image, pushing it to a registry, and deploying to the production environment.
        -   `docker-publish_qa.yml`: (Likely) Defines a similar workflow for a QA or staging environment.

## Backend Design Pattern

The service utilizes a **Layered Architecture**:

1.  **Presentation Layer (`main.go`):** Handles HTTP request/response interactions using the Gin web framework. This layer is responsible for API routing, request parsing, authentication (JWT validation), and formatting responses.
2.  **Application/Service Layer (`functions/Like.go`):** Contains the core business logic. For instance, the `LikePublication` function orchestrates the steps involved in liking a publication, including data validation and interaction with the data access layer.
3.  **Data Access Layer (`connection/mongo.go` and MongoDB interactions in `functions/Like.go`):** Manages data persistence. `connection/mongo.go` centralizes database connection setup. The actual database operations (CRUD) are performed within the application layer functions using the MongoDB driver.

This separation of concerns promotes modularity and maintainability.

## Communication Architecture

-   **Client-Server Model:** The service acts as a server, responding to requests from clients (e.g., frontend applications, other microservices).
-   **RESTful API:** Communication is primarily through a RESTful API. The service exposes HTTP endpoints for its functionalities.
-   **Synchronous Communication:** API calls are synchronous; the client sends a request and waits for the server's response.
-   **JWT Authentication:** The `/like` endpoint is secured using JSON Web Tokens. Clients must provide a valid JWT Bearer token in the `Authorization` header to access the endpoint. The server validates this token before processing requests.
-   **Internal Function Calls:** Within the application, different modules and components communicate via direct Go function calls.

## Endpoint Instructions

### Like a Publication

Allows an authenticated user to add a "like" to a specific publication.

-   **Endpoint:** `/like`
-   **HTTP Method:** `POST`
-   **Authentication:**
    -   Required: JWT Bearer Token.
    -   Header: `Authorization: Bearer <your_jwt_token>`
    -   The JWT must include a `user_id` claim (numeric).
-   **Request Headers:**
    -   `Content-Type: application/json`
-   **Request Body (JSON):**

    ```json
    {
      "_id": "string" // The MongoDB ObjectID of the publication
    }
    ```
    -   `_id` (string, required): The unique identifier of the publication to be liked.

-   **Responses:**

    -   **`200 OK` (Success):**
        ```json
        {
          "message": "Like Added like"
        }
        ```
    -   **`400 Bad Request` (Client Error):**
        -   Invalid JSON payload: `{"error": "Invalid JSON"}`
        -   Invalid publication ID format: `{"error": "Invalid publication ID"}`
        -   User has already liked the publication: `{"error": "Yo have already liked this publication"}`
    -   **`401 Unauthorized` (Authentication Error):**
        -   Token missing or malformed: `{"error": "Token missing or invalid"}`
        -   Invalid or expired token: `{"error": "Invalid or expired token"}`
        -   Invalid token claims (e.g., `user_id` missing): `{"error": "Invalid token claims"}` or `{"error": "user_id not found in token"}`
    -   **`500 Internal Server Error` (Server Error):**
        -   Database operation failure: `{"error": "Database error: <details>"}`
