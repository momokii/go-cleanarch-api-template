# Clean Architecture API Development Implementation using Golang Fiber

This project demonstrates the implementation of Clean Architecture principles in building a RESTful API using the Golang Fiber framework. The goal is to create a scalable and maintainable architecture by separating the application into distinct layers, such as domain, use case, repository, and delivery.

The structure promotes independence of frameworks, testability, and flexibility to adapt to various external dependencies. This makes it easier to manage changes over time and ensures the codebase remains clean and efficient as the application grows.

Key features of this project include:
- Fiber as the web framework
- PostgreSQL as the database (easily interchangeable with other databases)
- JWT authentication for secure endpoints
- Clean Architecture structure for better code organization
- Repository pattern for data persistence

**Note**: The current version of this project does not include any unit or integration tests. Future updates will introduce testing to improve code reliability and maintainability.

## Project Structure

The project follows a well-organized folder structure:

- **api/**: Contains API specifications and documentation, typically in OpenAPI/Swagger format. Use this to generate API documentation and share with developers.
  
- **internal/**: Contains code that is specific to this application. It is structured according to the Clean Architecture principles, with distinct layers to separate concerns:
  
  - **domain/**: Defines core entities, business rules, and logic that are central to the application's functionality. This is the heart of the business logic and should not depend on external frameworks.
  
  - **infrastructure/**: Handles external services such as database access, third-party APIs, etc. It serves as the implementation detail for the persistence and other system dependencies.
  
  - **interfaces/**: Defines how the outside world interacts with the application, such as controllers or HTTP handlers. This layer handles the communication between the external environment (e.g., web requests) and the inner layers.
  
  - **service/**: Contains the application's use cases, which define the interactions between the domain and external entities (via the `interfaces` and `infrastructure` layers). It orchestrates the flow of data between the layers and ensures that business rules are enforced.

- **pkg/**: Holds reusable code that can be shared across different applications. Libraries, utilities, and other generic functionality that isn't tied to the specific business logic of this application should go here.

This folder structure helps to maintain a clear separation of concerns, ensures a clean codebase, and follows Clean Architecture principles, allowing for scalability and ease of testing.


Follow the steps below to get started.

## Getting Started

### 1. Configure Environment Variables
Create a `.env` file in the root directory of the project and fill it with your configuration settings using the values from `.example.env`.

### 2. Install Dependencies
Run the following command to ensure all necessary modules are installed:

```bash
go mod tidy
```

### 3. Start the Development Server
To start the development server, run:

```bash
go run main.go
```

This will start the server and automatically reload changes when you rerun the command after making updates.

Alternatively, for live reloading during development, you can use **air** by running:

```bash
air
```

Make sure to configure air according to your project's needs by adjusting the settings in the `.air.toml` file.

### 4. Start the Production Server
To start the server in production mode, you can build the binary and run it:

#### On Windows:
```bash
go build -o lorem.exe
lorem.exe
```

#### On Linux/macOS:
```bash
go build -o lorem
./lorem
```

