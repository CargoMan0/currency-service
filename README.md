# Currency Exchange Service

This project provides a suite of services to deliver daily currency exchange rates to users. It includes the `Gateway` and `Currency` services, each with specific functionalities.

## Overview

This system is designed to:
- Provide users with daily currency exchange rates.
- Validate user permissions and convert user requests.
- Fetch and store currency rates from a public API.

## Services

### Gateway Service

The `Gateway` service is responsible for the initial handling of user requests, validating user permissions, and converting requests to a format suitable for further processing by the `Currency` service. This service masks the implementation details of the `Currency` service from the end user.

#### Key Features
1. **HTTP REST API**: Processes user requests via a standard HTTP REST API.
2. **Authorization Requests**: Interfaces with an authentication service to generate and validate access tokens.
3. **Data Request Forwarding**: Sends transformed requests to the `Currency` service using REST/RPC.

#### Assumptions
- Does not maintain a separate database for user storage; opts for simple in-memory structures like slices or maps.
- User registration is simplified; user data may be hardcoded or configured through files.
- Minimal user data is utilized: only `login` and `password` are necessary.

### Currency Service

The `Currency` service is tasked with fetching currency exchange rates from a public API and storing this data in a database. It also processes requests from the `Gateway` service to retrieve stored rates and historical data.

#### Key Features
1. **Automated Worker**: Runs daily to fetch the current RUB exchange rate against one foreign currency.
2. **Data Storage**: Saves fetched exchange rates in a database.
3. **Data Retrieval**: Responds to `Gateway` service requests for specific dates and historical exchange data over time.

## Deployment and Setup

The services can be orchestrated using Docker Compose, with dependencies on a PostgreSQL database for data storage.

### Prerequisites
- Docker and Docker Compose should be installed on the host machine.

### Running the Services

1. **Start Services**:
   Use Docker Compose to start the services:

   ```sh
   docker compose up --build

2. **Generate test data**:
   You can generate some test data using script
   ```sh
   go run ./currency/internal/scripts/generate_test_data.go

# Testing the API with `curl`

This guide outlines how to interact with the API using `curl` for registering a user, logging in, and fetching currency rates.

## Step 1: Register User

Register a new user by sending a `POST` request with the username and password in JSON format.

    curl -X POST http://localhost:8080/api/v1/register \
      -H "Content-Type: application/json" \
      -d '{
            "Username": "test",
            "Password": "test"
          }'

## Step 2: Login User

Authenticate the user and retrieve a token by sending a `POST` request with the username and password.

    curl -X POST http://localhost:8080/api/v1/login \
      -H "Content-Type: application/json" \
      -d '{
            "Username": "test",
            "Password": "test"
          }'

**Note**: The response will contain a JSON object with a `"token"` field. Extract the token value from the response to use it in the next request.

## Step 3: Get Currency Rates

Fetch currency rates by sending a `GET` request. Include the `Authorization` header with the token obtained from the login step.

    curl -X GET "http://localhost:8080/api/v1/rate?currency=USD&date_from=2024-10-01&date_to=2024-10-21" \
      -H "Authorization: Bearer <YOUR_TOKEN>"

**Replace `<YOUR_TOKEN>`** with the actual token obtained in step 2.

## Notes

- **Authorization Token**: Ensure that the token from the login step is used in the `Authorization` header of the third request.
- **Content-Type Header**: The `Content-Type` header should be set to `application/json` for `POST` requests to specify that the request body contains JSON data.

These steps assume the API is running locally on `localhost:8080`. Adjust the host and port as necessary if your API is hosted elsewhere.