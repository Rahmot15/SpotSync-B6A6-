# 🚗 SpotSync – Smart Parking & EV Charging Reservation API

A high-performance, concurrent RESTful API built with Go and Echo. SpotSync provides a centralized platform for airports and malls to manage parking zones, specifically resolving race conditions for high-demand, limited EV charging spots through database row-level locking.

---

## 🌐 Live Deployment & Submission Links

* **Live API URL:** `https://spotsync-api.onrender.com` *(Replace with your deployed URL)*
* **GitHub Repository:** `https://github.com/Rahmot15/SpotSync-B6A6-`
* **Interview Video:** `https://youtu.be/your-video-link` *(Replace with your video link)*

---

## 🛠️ Tech Stack

* **Language:** Go (Golang 1.22+)
* **Web Framework:** [Echo v4](https://github.com/labstack/echo) (Minimalist, high-performance HTTP framework)
* **ORM:** [GORM](https://gorm.io/) with PostgreSQL Driver
* **Database:** PostgreSQL (Cloud hosted on NeonDB)
* **Validation:** [validator v10](https://github.com/go-playground/validator)
* **Authentication:** JWT (`golang-jwt/jwt/v5`) & Password Hashing (`golang.org/x/crypto/bcrypt`)

---

## 🏛️ Clean Architecture

SpotSync strictly adheres to **Clean Architecture** principles to separate concerns and ensure maintainability and testability.

```
          +-------------------------------------------------------+
          |                    HTTP Client                        |
          +-------------------------------------------------------+
                                      |
                                      v
          +-------------------------------------------------------+
          |                   Echo HTTP Layer                     |
          |  - Handlers (Bind & Validate DTOs, Extract JWT Claims) |
          +-------------------------------------------------------+
                                      |
                                      v
          +-------------------------------------------------------+
          |                    Service Layer                      |
          |  - Business Logic, Hashing, Auth & Capacity Checks    |
          +-------------------------------------------------------+
                                      |
                                      v
          +-------------------------------------------------------+
          |                   Repository Layer                    |
          |  - GORM DB Operations, Transactions, Row Locks        |
          +-------------------------------------------------------+
                                      |
                                      v
          +-------------------------------------------------------+
          |                  PostgreSQL Database                  |
          +-------------------------------------------------------+
```

### Layer Responsibilities
* **DTO (`internal/domain/*/dto`)**: Data Transfer Objects defining request payloads and response structures. Exposed models are decoupled from internal database entities.
* **Handler (`internal/domain/*/handler.go`)**: Manages HTTP requests, validates DTOs using `go-playground/validator`, extracts claims, and invokes services.
* **Service (`internal/domain/*/service.go`)**: Encapsulates core business logic, password verification, token generation, and capacity validations.
* **Repository (`internal/domain/*/repository.go`)**: Executes database queries, handles transactions, and enforces row-level locking.
* **Entity (`internal/domain/*/entity.go`)**: Defines GORM structs representing database tables (`users`, `parking_zones`, `reservations`).

---

## ✨ Key Features

1. **Role-Based Access Control (RBAC)**: Supports `driver` and `admin` roles with custom Echo middleware authorization.
2. **Dynamic Spot Availability**: Parking zones calculate `available_spots` dynamically (`total_capacity` minus active reservations).
3. **Concurrency Control (EV Bottleneck Prevention)**: Utilizes GORM database transactions with pessimistic row locking (`FOR UPDATE`) on the parking zone record during reservation creation to prevent over-capacity bookings during simultaneous requests.
4. **Automated Schema Migration**: GORM AutoMigrate seamlessly initializes database tables upon application startup.

---

## 🚀 Local Setup & Installation

### Prerequisites
* [Go 1.22+](https://golang.org/dl/) installed.
* [Git](https://git-scm.com/) installed.
* A valid PostgreSQL connection string (e.g., from NeonDB or local PostgreSQL).

### Setup Instructions

1. **Clone the Repository**
   ```bash
   git clone https://github.com/Rahmot15/SpotSync-B6A6-.git
   cd SpotSync(B6A6)
   ```

2. **Configure Environment Variables**
   Create a `.env` file in the root directory and define the following variables:
   ```env
   PORT=8080
   DATABASE_URL=postgresql://user:password@host/neondb?sslmode=require
   JWT_SECRET=your-super-secret-jwt-key
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run the Application**
   ```bash
   go run ./cmd
   ```
   The server will start on `http://localhost:8080`.

---

## 🌐 API Endpoints Specification

### 🔑 Authentication Module

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Public | Register a new driver or admin user |
| `POST` | `/api/v1/auth/login` | Public | Authenticate user and return JWT token |

### 🅿️ Parking Zones Module

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/zones` | Public | View all parking zones with dynamic `available_spots` |
| `GET` | `/api/v1/zones/:id` | Public | View details of a specific parking zone |
| `POST` | `/api/v1/zones` | Admin Only | Create a new parking zone |
| `PATCH` | `/api/v1/zones/:id` | Admin Only | Update parking zone details |
| `DELETE` | `/api/v1/zones/:id` | Admin Only | Delete a parking zone |

### 🚗 Reservations Module (Core Logic)

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/reservations` | Authenticated | Reserve a parking spot (Concurrency Safe with Row Locking) |
| `GET` | `/api/v1/reservations/my-reservations` | Authenticated | View current user's reservations with zone details |
| `DELETE` | `/api/v1/reservations/:id` | Authenticated | Cancel user's own active reservation |
| `GET` | `/api/v1/reservations` | Admin Only | View all reservations across the system |

---

## 📋 Response Formats

### Standard Success Response
```json
{
  "success": true,
  "message": "Operation description",
  "data": {}
}
```

### Standard Error Response
```json
{
  "success": false,
  "message": "Error description",
  "errors": "Error details"
}
```
