# ✅ Task Manager REST API (Go + Gin)

[![Go](https://img.shields.io/badge/Language-Go-blue?logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-red?logo=go)](https://gin-gonic.com)
[![Status](https://img.shields.io/badge/Project-Completed-brightgreen)]()

A lightweight RESTful API built using **Go** and the **Gin Framework**, developed during the A2SV internship learning phase.

The API allows users to manage tasks with full CRUD operations and in-memory data storage. It demonstrates foundational concepts in Go web backend development.

---

## 🔗 Live API Docs

📄 **Full API Documentation (via Postman)**:  
👉 [View Postman Docs](https://documenter.getpostman.com/view/33813578/2sB34ijeaC)

---

## 🚀 Features

| Feature                                  | Status    |
| ---------------------------------------- | --------- |
| ✅ Create new task (POST)                | Completed |
| ✅ Get all tasks (GET)                   | Completed |
| ✅ Get task by ID (GET)                  | Completed |
| ✅ Update task by ID (PUT)               | Completed |
| ✅ Delete task by ID (DELETE)            | Completed |
| ✅ In-memory task store                  | Completed |
| ✅ Custom error types                    | Completed |
| ✅ Input validation                      | Completed |
| ✅ Status constraint (Pending/Completed) | Completed |
| ✅ Postman collection documentation      | Completed |

---

## 🧪 How to Run

Make sure Go is installed.

```bash
# Clone the repo
git clone https://github.com/A-selam/a2sv-go.git
cd a2sv-go/task_manager/

# Run the application
go run main.go
```

## 📌 Notes

- This API uses in-memory data — all tasks reset when the app restarts.
- Built for learning purposes, focusing on Go fundamentals and REST APIs.
- You can easily extend it with persistent storage (PostgreSQL, MongoDB, etc.) in the future.
