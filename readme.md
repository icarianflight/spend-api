# Spend Transaction Management API

This project is a **REST API** built using **Golang** to manage bank accounts and transactions to track spending over time. It follows the **Hexagonal Architecture** for clean, modular, and maintainable code. The database interaction is implemented with **MariaDB** (with intended future support for other SQL databases), and the API allows for basic account and transaction management operations.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
- [Running the API](#running-the-api)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- Create and manage bank accounts.
- Record transactions for bank accounts.
- Hexagonal architecture following **Domain-Driven Design** (DDD) principles.
- TLS-enabled database connection for secure data storage.
- Configurable via environment variables for database connection details.

## Project Structure

The project is structured using **Hexagonal Architecture**, ensuring separation of concerns between the core business logic and infrastructure (adapters for persistence, APIs, etc.).

```text
/internal/
    /domain/
        /accounts/
            models.go        # Domain models for accounts
            service.go       # Business logic for accounts
        /transactions/
            models.go        # Domain models for transactions
            service.go       # Business logic for transactions
    /app/
        /adapters/
            /db/
                accounts/    # Persistence logic for accounts in MariaDB
                transactions/ # Persistence logic for transactions in MariaDB
            /rest/
                accounts/    # API handler for accounts
                transactions/ # API handler for transactions
        /infra/
            db/              # Database interfaces and execution logic
            config/          # Configuration management
/cmd/
    /api/
        main.go          # Entry point for starting the API
```

## Setup Instructions

### Prerequisites
* Golang version 1.23.1 or higher
* MariaDB

### Environment Variables
The application requires environment variables for the database connection and TLS settings.

Create a .env file or set these variables in your environment:

```text
DB_USER=<your_db_user>
DB_PASSWORD=<your_db_password>
DB_HOST=<db_host>
DB_PORT=<db_port>
DB_NAME=<db_name>
CACERT_PATH=<path_to_ca_cert>
CLIENT_CERT_PATH=<path_to_client_cert>
CLIENT_KEY_PATH=<path_to_client_key>
```

### Installing Dependencies
Clone the repository:

```text
git clone https://github.com/yourusername/transaction-api.git
cd transaction-api
```

### Install dependencies:

```text
go mod tidy
```

## Running the API
After setting up the environment variables and the database:

Start the API:
```test
go run ./cmd/api/main.go
```

## Testing
The project follows Test-Driven Development (TDD) principles and includes comprehensive unit tests.

Run all unit tests:

```test
go test ./...
```

SQLMock is used to mock database operations, ensuring that all database interactions are unit tested without requiring a live database.

## Contributing
We welcome contributions! Please follow these steps:

* Fork the repository.
* Create a feature branch (git checkout -b feature-branch).
* Commit your changes (git commit -m 'Add some feature').
* Push to the branch (git push origin feature-branch).
* Create a pull request.
* Please ensure that your code adheres to the project’s coding standards and that all tests pass before submitting a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details