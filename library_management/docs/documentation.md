# 📚 Library Management System (Go)

[![Go](https://img.shields.io/badge/Language-Go-blue?logo=go)](https://golang.org)
[![Status](https://img.shields.io/badge/Project-In%20Progress-yellow)]()

A simple console-based Library Management System built in Go as part of an internship learning project.

This system demonstrates foundational Go concepts including:

- Interfaces
- Struct composition
- Maps and slices
- Method receivers
- Clean modular architecture
- Console input/output with `bufio.Reader`

---

## 🚀 Features

| Feature                            | Status    |
| ---------------------------------- | --------- |
| ✅ Add new books                   | Completed |
| ✅ Add new members                 | Completed |
| ✅ Borrow books                    | Completed |
| ✅ Return books                    | Completed |
| ✅ Remove books                    | Completed |
| ✅ List available books            | Completed |
| ✅ List borrowed books (by member) | Completed |
| ✅ Auto-generate unique IDs        | Completed |
| ✅ Console-based interaction       | Completed |
| 🔧 Input validation                | Partial   |

---

## 📂 Folder Structure

library_management/\
├── main.go # Entry point of the app\
├── controllers/ # Handles user input and actions\
│ └── library_controller.go\
│ └── helper.go\
├── models/ # Book and Member struct definitions\
│ ├── book.go\
│ └── member.go\
├── services/ # Core logic (implements LibraryManager interface)\
│ └── library_service.go\
├── docs/ # Documentation and usage notes\
│ └── documentation.md\
└── go.mod # Go module file

---

## 📦 Technologies & Concepts

- ✅ Structs and custom types
- ✅ Interfaces and method implementation
- ✅ Maps and slices for in-memory storage
- ✅ Auto-incrementing ID generators
- ✅ Terminal input/output with `fmt` and `bufio`
- ✅ Separation of concerns (Controller ↔ Service ↔ Models)
- 🔧 Basic error handling and validation

---

## 🧪 How to Run

Make sure Go is installed on your system.

```bash
# Clone the repository
git clone https://github.com/A-selam/a2sv-go.git
cd a2sv-go/library_management/

# Run the application
go run main.go
```

## 📌 Notes

- This system is in-memory only — data resets each time the app restarts.
- Designed for learning purposes — not production-grade.
- Modular structure makes it easy to expand with file/database persistence later.

## 🙌 Acknowledgements

Built as part of the Go track in the A2SV internship learning phase. Inspired by practical system design exercises for learning backend fundamentals.
