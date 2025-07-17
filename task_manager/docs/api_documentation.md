# ✅ Task Manager REST API (Go + Gin + MongoDB)

[![Go](https://img.shields.io/badge/Language-Go-blue?logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-red?logo=go)](https://gin-gonic.com)
[![Database](https://img.shields.io/badge/Database-MongoDB-green?logo=mongodb)](https://www.mongodb.com)
[![Status](https://img.shields.io/badge/Project-Completed-brightgreen)]()

A lightweight RESTful API built using **Go**, the **Gin Framework**, and **MongoDB**, developed during the A2SV internship learning phase.

## 🔗 Live API Docs

📄 **Full API Documentation (via Postman)**:  
👉 [View Postman Docs](https://documenter.getpostman.com/view/33813578/2sB34ijeaC)

## 🚀 Features

| Feature                                  | Status    |
| ---------------------------------------- | --------- |
| ✅ Create new task (POST)                | Completed |
| ✅ Get all tasks (GET)                   | Completed |
| ✅ Get task by ID (GET)                  | Completed |
| ✅ Update task by ID (PUT)               | Completed |
| ✅ Delete task by ID (DELETE)            | Completed |
| ✅ MongoDB persistent storage            | Completed |
| ✅ Environment configuration             | Completed |
| ✅ Custom error types                    | Completed |
| ✅ Input validation                      | Completed |
| ✅ Status constraint (Pending/Completed) | Completed |
| ✅ Postman collection documentation      | Completed |

## 🧰 Prerequisites

- **Go** (v1.20+) - [Download](https://golang.org/dl/)
- **MongoDB** (v6.0+) running either:
  - **Local**: [MongoDB Community Server](https://www.mongodb.com/try/download/community)
  - **Cloud**: [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) (free tier available)
- **Git** for cloning the repository

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/A-selam/a2sv-go.git
cd a2sv-go/task_manager/
```

### 2. Set up environment

create .env file:

```env
# For local MongoDB
MONGODB_URI=mongodb://localhost:27017
DB_NAME=taskmanager

# For MongoDB Atlas (uncomment and replace)
# MONGODB_URI=mongodb+srv://<username>:<password>@cluster0.example.mongodb.net/?retryWrites=true&w=majority
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
