# 📋 Task Manager API

A secure and efficient **Task Manager REST API** built with **Go**, **Gin**, and **MongoDB**, designed during the A2SV internship. This API provides authentication with JWT, role-based access control, and complete CRUD functionality for managing tasks.

---

## 🔧 Tech Stack

- **Language:** Go
- **Framework:** Gin Gonic
- **Database:** MongoDB
- **Auth:** JWT (JSON Web Tokens)

---

## ✅ Features

| Feature                       | Description                                           |
| ----------------------------- | ----------------------------------------------------- |
| 🔐 Authentication             | Secure login with JWT token                           |
| 👤 Role-Based Access          | Admin vs. Regular User permissions                    |
| ➕ Create Task                | Any authenticated user                                |
| 📄 Get All Tasks (admin only) | Admins can view all tasks on the system               |
| 🔎 Get Own Tasks              | Users can view tasks they created                     |
| 🛠️ Update Own Task            | Users can update only their own tasks                 |
| ❌ Delete Own Task            | Users can delete only their own tasks                 |
| ⏱️ Status & DueDate           | Tasks must have status (Pending/Completed) & due date |
| 🧪 Environment Config         | `.env` file for DB URI and credentials                |

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/A-selam/a2sv-go.git
cd a2sv-go/task_manager/
```

### 2. Set up environment

Create a `.env` file at the root:

```env
# For local MongoDB
MONGODB_URI=mongodb://localhost:27017
DB_NAME=taskmanager

# For MongoDB Atlas (optional)
# MONGODB_URI=mongodb+srv://<username>:<password>@cluster0.mongodb.net/?retryWrites=true&w=majority
# DB_NAME=taskmanager
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the application

```bash
go run main.go
```

### 5. Seed Data & Admin User

When you run the server for the first time, the application automatically checks if the database collections are empty. If empty, it seeds:

- A default admin user with:

  - Username: Mr. Admin
  - Password: 12345678

- A sample task owned by this admin user.

Use this admin account to immediately test protected endpoints and manage tasks.

---

## 🔐 Authentication & Authorization

### 🔑 Token

- Uses **JWT** for user sessions
- Token is valid for **24 hours**
- Must be passed in request headers for protected routes

### 🔒 Header format

```http
Authorization: Bearer <your_jwt_token_here>
```

### 👤 Roles

- **Admin**: Can manage all tasks
- **User**: Can only manage tasks they created

---

## 📬 API Endpoints

> **Base URL**: `http://localhost:3000`

### 👥 Auth

#### 🔐 Register User

```http
POST /register
```

**Body:**

```json
{
  "username": "exampleuser",
  "password": "securepassword",
  "role": "User" // or "Admin"
}
```

#### 🔐 Login User

```http
POST /login
```

**Body:**

```json
{
  "username": "exampleuser",
  "password": "securepassword"
}
```

**Returns:** JWT token

---

### 📋 Tasks (Protected)

> **All routes below require Authorization header**:

```http
Authorization: Bearer <jwt_token>
```

#### ➕ Create Task

```http
POST /tasks
```

**Body:**

```json
{
  "title": "Finish assignment",
  "description": "Complete by end of week",
  "due_date": "2025-07-30",
  "status": "Pending"
}
```

#### 📄 Get All Tasks (Admin Only)

```http
GET /tasks
```

- Admins get all tasks
- Regular users will be denied

#### 📄 Get My Tasks (User Only)

```http
GET /tasks
```

- Returns tasks created by the current user

#### 🔄 Update Task by ID

```http
PUT /tasks/:id
```

**Body:** (Partial updates supported)

```json
{
  "title": "Updated title",
  "status": "Completed"
}
```

- Only the creator of the task (or admin) can update

#### ❌ Delete Task by ID

```http
DELETE /tasks/:id
```

- Only the creator (or admin) can delete the task

---

## 🗃️ Models Overview

### 🔐 User

```go
type User struct {
  ID       int    `json:"id"`
  Username string `json:"username"`
  Password string `json:"password"`
  Role     role   `json:"role"` // "Admin" or "User"
}
```

### 📋 Task

```go
type Task struct {
  ID          int    `json:"id"`
  Title       string `json:"title"`
  Description string `json:"description"`
  DueDate     string `json:"due_date"`
  Status      status `json:"status"` // "Pending" or "Completed"
  CreatedBy   string `json:"created_by"`
}
```

---

## 🧠 Notes

- Passwords are securely hashed using `bcrypt`
- Task IDs are managed manually – consider switching to MongoDB ObjectIDs in future

---

## ⚠️ Error Handling

- All error responses include meaningful messages with appropriate HTTP status codes:
  - 400 Bad Request for validation or malformed requests
  - 401 Unauthorized for missing/invalid auth token
  - 404 Not Found when resource is not found or access is forbidden
  - 500 Internal Server Error for unexpected server errors

---

## 🧪 Testing & Usage Tips

- Use the seeded admin account (Mr. Admin / 12345678) for initial testing.
- Remember to pass the JWT token as Authorization: Bearer <token> header for protected endpoints.
- JWT tokens expire after 24 hours — re-login to obtain new tokens.
- Role must always be either Admin or User (during registration).
- Status must always be either Pending or Completed.
