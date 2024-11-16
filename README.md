# Hotel Guide Project

This project is a simple hotel guide application built using Go. It consists of two microservices communicating with each other via HTTP: `hotel-service` and `report-service`.

## Table of Contents

1. [Overview](#overview)
2. [Technologies](#technologies)
3. [API Endpoints](#api-endpoints)
4. [Setup and Installation](#setup-and-installation)
5. [Testing](#testing)
6. [Environment Variables](#environment-variables)
7. [Unit Test Coverage](#unit-test-coverage)

---

## Overview

The Hotel Guide application is designed to handle hotel-related data and generate reports for specific locations. The project implements the following features:

- Hotel management (create, delete)
- Adding and removing hotel contact information
- Generating location-based reports (asynchronous)
- Viewing report details and statuses

The application follows a microservices architecture, where two services communicate through HTTP:

- `hotel-service`: Manages hotel data.
- `report-service`: Handles report generation and fetching statistics.

---

## Technologies

- **Go** for backend services
- **PostgreSQL** for database management
- **RabbitMQ** for message queuing
- **Docker** for containerization
- **GORM** for ORM-based database interaction
- **Gorilla Mux** for routing

---

## API Endpoints

### Hotel-Service (http://localhost:8081)

#### **POST /hotels**  
Create a new hotel.

- **Request Body**:
    ```json
    {
        "ownerName": "John",
        "ownerSurname": "Doe",
        "companyTitle": "Sample Hotel"
    }
    ```
- **Example**:  
  `curl -X POST http://localhost:8081/hotels -d '{"ownerName":"John","ownerSurname":"Doe","companyTitle":"Sample Hotel"}'`

---

#### **GET /hotels**  
Retrieve all hotels.

- **Example**:  
  `curl http://localhost:8081/hotels`

---

#### **POST /hotels/{id}/contacts**  
Add contact information to a hotel.

- **Request Body**:
    ```json
    {
        "info_type": "location",
        "info_content": "New York"
    }
    ```
- **Example**:  
  `curl -X POST http://localhost:8081/hotels/{hotel_id}/contacts -d '{"info_type":"location","info_content":"New York"}'`

---

#### **DELETE /hotels/{id}**  
Delete a hotel.

- **Example**:  
  `curl -X DELETE http://localhost:8081/hotels/{hotel_id}`

---

#### **DELETE /hotels/{id}/contacts/{contact_id}**  
Delete contact information for a hotel.

- **Example**:  
  `curl -X DELETE http://localhost:8081/hotels/{hotel_id}/contacts/{contact_id}`

---

#### **GET /hotels/officials**  
Retrieve officials associated with hotels.

- **Example**:  
  `curl http://localhost:8081/hotels/officials`

---

#### **GET /hotels/{id}**  
Retrieve a specific hotel by ID.

- **Example**:  
  `curl http://localhost:8081/hotels/{hotel_id}`

---

#### **GET /hotels/stats**  
Retrieve statistics about hotels for a specific location.

- **Query Parameters**:  
  `location` (required) - The location to fetch stats for.
- **Example**:  
  `curl http://localhost:8081/hotels/stats?location=New+York`

---

### Report-Service (http://localhost:8082)

#### **POST /reports**  
Request a new report for a specific location.

- **Request Body**:
    ```json
    {
        "location": "New York"
    }
    ```
- **How it works**:
    When a new report is requested, the request is placed in a RabbitMQ queue, and a worker consumes the task asynchronously. The report includes statistics about hotels and phone numbers for the specified location. 
    The report is processed in the background, and the status will be updated to "Completed" once the task is done.

- **Example**:  
  `curl -X POST http://localhost:8082/reports -d '{"location":"New York"}'`

---

#### **GET /reports**  
Retrieve all reports.

- **Example**:  
  `curl http://localhost:8082/reports`

---

#### **GET /reports/{id}**  
Retrieve details of a specific report.

- **Example**:  
  `curl http://localhost:8082/reports/{report_id}`

---

## Setup and Installation

### Prerequisites

- **Go** (1.18+)
- **Docker** (for containerization)
- **PostgreSQL** (for database)

### Steps

1. **Clone the repository**

    ```bash
    git clone https://github.com/GokhanCagritekin/hotel-guide.git
    cd hotel-guide
    ```

2. **Set up environment variables**

    Create a `.env` file in the root directory and add the following:

    ```
    DB_USER=myuser
    DB_PASSWORD=mysecretpassword
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=hotels
    DB_SSLMODE=disable
    
    MQ_USER=guest
    MQ_PASSWORD=guest
    MQ_HOST=localhost
    MQ_PORT=5672
    
    HOTEL_SERVICE_URL=http://localhost:8081
    REPORT_SERVICE_URL=http://localhost:8082

    ```
    
3. **Development Environment Setup**
 ### RabbitMQ
Run the RabbitMQ container with the management plugin enabled:
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```
**Management UI**: Available at [http://localhost:15672](http://localhost:15672)  

**Default credentials**:  
- **Username**: `guest`  
- **Password**: `guest`

---
### PostgreSQL  
Run the PostgreSQL container for the hotel database:  

```bash
docker run --name hotel-db -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=myuser -e POSTGRES_DB=hotels -p 5432:5432 -d postgres
```

**Database Configuration**:  

- **Username**: `myuser`  
- **Password**: `mysecretpassword`  
- **Database**: `hotels`  

4. **Run the services with Docker**

    You can run both services using Docker. Ensure Docker is installed, then run the following commands to build and start the services:

    ```bash
    docker-compose up --build
    ```

    This will start the `hotel-service` and `report-service` on `localhost` and set up PostgreSQL in a container.

---

## Testing

### Unit Tests

Run unit tests using the following command:

```bash
go test -cover ./...
