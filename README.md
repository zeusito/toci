# Toci
This is a highly opinionated blueprint for building Go applications.
It serves as a foundation for new projects, much like a Lego set where you can pick and choose the pieces you need for your specific use case.


## Features
- Chi Router for HTTP-based endpoints
- Zerolog for logging capabilities
- Koanf for configuration, supports files and env vars
- PGX and Bun for PostgreSQL database access
- DBMate for database migrations
- Session management (using Opaque tokens) and HTTP filter to protect endpoints
- Hashing algorithms, including argon2id
- Makefile with the most common tasks
- Multi-stage Dockerfile for building and running the application
- A basic authentication module

## Getting Started
- Start a new repository from this template, or clone the repository or download the code
- Change the module name in the go.mod file
- Run `make run` to start the application
- Check out the `Makefile` for more information about available commands
- Start customizing the application

## Folder Structure
- cmd/main.go - main entry point
- internal - application-specific business logic
- pkg – shared packages that might be used in multiple modules (follows the Unix philosophy of simple tools that make one thing)
- resources - application-specific resources, such as config files, databases, etc.

