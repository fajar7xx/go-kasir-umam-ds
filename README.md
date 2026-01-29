# Go Kasir Umam DS

This is a simple cashier application written in Go. It provides a REST API for managing products and categories.

## Developer Docs

### Project Structure

The project is structured using a layered architecture:

*   **`main.go`**: The entry point of the application. It initializes the server and routes.
*   **`handlers`**: Contains the HTTP handlers for the API endpoints.
*   **`services`**: Contains the business logic of the application.
*   **`repositories`**: Contains the data access logic.
*   **`models`**: Contains the data structures.
*   **`internal/database`**: Contains the database connection logic.
*   **`config`**: Contains the configuration logic.

### How to Run

1.  Clone the repository.
2.  Run `go run main.go` to start the server.
3.  The server will be running on `http://localhost:8080`.

## API Endpoints

### Health Check

*   **GET /health**: Checks the health of the application.

### Categories

*   **GET /api/v1/categories**: Get all categories.
*   **GET /api/v1/categories/{id}**: Get a category by ID.
*   **POST /api/v1/categories**: Create a new category.
*   **PUT /api/v1/categories/{id}**: Update a category.
*   **DELETE /api/v1/categories/{id}**: Delete a category.

### Products

*   **GET /api/v1/products**: Get all products.
*   **GET /api/v1/products/{id}**: Get a product by ID.
*   **POST /api/v1/products**: Create a new product.
*   **PUT /api/v1/products/{id}**: Update a product.
*   **DELETE /api/v1/products/{id}**: Delete a product.
