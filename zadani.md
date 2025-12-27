
========Task 2 Chat - Project Requirements
Project Overview

A REST API for a chat application built with Go, PostgreSQL, and clean architecture principles. The system manages authors and their messages with full CRUD operations and search capabilities.
 
Stack

* Go
* PostgreSQL db
* Clean Architecture
* REST API
* JSON

Libraries & Tools

* go-jet - SQL query builder 
* goose - Database migrations
* Gorilla Mux or Chi - HTTP router

Project structure

* domain/
* handlers/
* usecases/
* ports/
* migrations/
* main.go



Functional Requirements

1. Author Management

The system must:
// TODO: Pagination
* Retrieve all authors (with pagination support)
* Update author information (username)
* Delete authors (cascading delete of all associated messages)
* Search authors by username (partial match)

2. Message Management

The system must:

* Retrieve a single message by ID
* Retrieve all messages (with pagination support)
* Update message content (cannot change author)
* Delete individual messages
* Retrieve all messages by a specific author (with pagination)
* Search messages by content (partial text match)
* Filter messages by author ID

3. Search & Filtering

The system must:

* Search messages by content keyword
* Search authors by username keyword
* Filter messages by author ID
* Support case-insensitive search
* Return results with pagination
* Use only built-in PostgreSQL functionality for search (no external search engines)



Non-Functional Requirements

1. Configuration

The system must:

* Use environment variables for all configuration
* Required configuration:
* Database connection (host, port, user, password, dbname)
* Server port

2. Database Migrations

The system must:

* Use goose for database migrations
* Version all database schema changes
* Create migration files for all schema modifications
* Ensure migrations can be applied incrementally

3. Logging

The system must:

* Use structured logging (JSON format)

4. Security

The system must:

* Validate all input data
* Prevent SQL injection (using parameterized queries)
* Sanitize user input before storage

5. Code Quality

The system must:

* Follow Go conventions and idioms
* Pass gofmt formatting checks
* Follow clean architecture principles
* Use meaningful variable and function names
* Include code comments for complex logic
* Keep functions small and focused

6. Containerization

The system must:

* Provide docker-compose.yml for local development
* Include both application and database in compose
* Support running entire stack with docker-compose up

