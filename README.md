# KV-Store

A simple, concurrent-safe, and persistent key-value store implemented in Go.

## Overview

This project provides a standalone key-value store server that supports basic CRUD operations. It features concurrent access handling, time-to-live (TTL) for keys, and periodic data persistence to a JSON file on disk. It also maintains a log of all operations.

## Features

-   **CRUD Operations**: Basic Create, Read, Update, and Delete operations for key-value pairs.
-   **Concurrent Safe**: Handles multiple requests concurrently using mutexes.
-   **Data Persistence**: Periodically saves the in-memory store to a `kv-data.json` file to prevent data loss.
-   **Time-To-Live (TTL)**: Supports automatic expiration of keys.
-   **Logging**: Logs all operations with timestamps to `kv-logs.log`.

## Project Structure

```
KV-Store/
├── cmd/
│   └── main.go         # Main application entry point
├── pkg/
│   ├── controllers/
│   │   └── controller.go # Request handlers for API endpoints
│   ├── store/
│   │   └── store.go      # Core key-value store logic
│   └── utils/
│       └── util.go       # Utility functions (e.g., logging)
├── go.mod
├── data.json             # Data persistence file (auto-generated)
└── kv-logs.log           # Log file (auto-generated)
```

## Getting Started

### Prerequisites

-   [Go](https://golang.org/doc/install) installed on your machine.

### Installation & Running

1.  Clone the repository:
    ```bash
    git clone https://github.com/patnaikankit/KV-Store.git
    cd KV-Store
    ```

2.  Run the server:
    ```bash
    go run cmd/main.go
    ```
    The server will start on port `4000`.

## API Endpoints

The API provides endpoints for managing key-value pairs.

### 1. Get a Key

Retrieves the value for a given key.

-   **Endpoint**: `/get`
-   **Method**: `GET`
-   **Query Parameter**: `key`
-   **Example**:
    ```bash
    curl "http://localhost:4000/get?key=mykey"
    ```

### 2. Set a Key

Creates a new key-value pair. Fails if the key already exists.

-   **Endpoint**: `/set`
-   **Method**: `POST`
-   **Body**:
    ```json
    {
        "Key": "mykey",
        "value": "myvalue",
        "ttl": "1h"
    }
    ```
-   **Example**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"Key": "mykey", "value": "myvalue", "ttl": "2h"}' http://localhost:4000/set
    ```

### 3. Update a Key

Updates the value of an existing key.

-   **Endpoint**: `/update`
-   **Method**: `PUT` or `PATCH`
-   **Body**:
    ```json
    {
        "Key": "mykey",
        "value": "newvalue"
    }
    ```
-   **Example**:
    ```bash
    curl -X PUT -H "Content-Type: application/json" -d '{"Key": "mykey", "value": "newvalue"}' http://localhost:4000/update
    ```

### 4. Delete a Key

Deletes a key-value pair.

-   **Endpoint**: `/delete`
-   **Method**: `DELETE`
-   **Query Parameter**: `Key`
-   **Example**:
    ```bash
    curl -X DELETE "http://localhost:4000/delete?Key=mykey"
    ```

## Postman Collection

You can also use the following Postman collection to test the API:

-   [KV-Store Collection](https://documenter.getpostman.com/view/25374920/2sB2xBDq2X) 