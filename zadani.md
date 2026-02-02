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

2. Message Management

The system must:


3. Search & Filtering

The system must:

* Support case-insensitive search
* Use only built-in PostgreSQL functionality for search (no external search engines)



Non-Functional Requirements

1. Configuration

The system must:

* Use environment variables for all configuration

2. Database Migrations

The system must:

3. Logging

The system must:

* Use structured logging (JSON format)

4. Security

The system must:

* Validate all input data

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

