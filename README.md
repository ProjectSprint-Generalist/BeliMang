# Project 3: BeliMang!

A food delivery backend service that user can use to buy food and drinks!

## Setup

1. Install Docker
2. Install sqlc [here](https://docs.sqlc.dev/en/stable/overview/install.html)
3. Make sure `go` is installed and version is 1.25.1
4. Run `go mod tidy` to install dependencies

## Tech Stack

- Golang 1.25.1 + Gin
- pgx + sqlc (type-safe queries)
- PostgreSQL + PostGIS
- MinIO

## Project Structure

```
BeliMang/
├── internal/
├── compose.yaml
├── go.mod
├── .env.example
├── main.go                    # Main function entry point
└── README.md
```
