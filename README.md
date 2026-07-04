# Library API

A simple RESTful Library Management API built with **Go**, **Fiber** and **MongoDB**.

---

## 🚀 About The Project

Library API is a backend project developed to manage basic library operations such as user registration, login, book creation, book listing and book borrowing.

This project was created as an early backend development practice project and demonstrates REST API development, MongoDB usage and password hashing with bcrypt.

---

## 🛠 Tech Stack

- **Go**
- **Fiber**
- **MongoDB**
- **bcrypt**
- **REST API**

---

## ✨ Features

- User registration
- User login
- Password hashing with bcrypt
- Add books
- List library books
- Borrow books
- User profile with borrowed books
- MongoDB persistence
- Borrow limit rule: each user can borrow up to 2 books

---

## 📁 Project Structure

```text
Library_api/
├── go.mod
├── go.sum
├── server.go
└── README.md
```

---

## ⚙️ Installation

Clone the repository:

```bash
git clone https://github.com/SerhatErbil/Library_api.git
cd Library_api
```

Install dependencies:

```bash
go mod tidy
```

Make sure MongoDB is running locally:

```bash
mongodb://localhost:27017
```

Run the application:

```bash
go run server.go
```

The server will start on:

```text
http://localhost:3000
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/register` | Register a new user |
| POST | `/login` | Login user |
| POST | `/books` | Add a new book |
| GET | `/library` | List all books |
| GET | `/profile` | Get user profile and borrowed books |
| POST | `/borrow` | Borrow a book |

---

## 📦 Example Request

### Register User

```json
{
  "username": "serhat",
  "password": "123456"
}
```

### Add Book

```json
{
  "id": 1,
  "name": "Clean Code",
  "author": "Robert C. Martin"
}
```

### Borrow Book

```json
{
  "book_id": 1,
  "username": "serhat"
}
```

---

## 🧠 What I Practiced

- Building REST APIs with Go Fiber
- Connecting Go applications to MongoDB
- Structuring request and response models
- Hashing passwords with bcrypt
- Implementing simple business rules
- Handling JSON requests and API responses

---

## 📌 Notes

This is an early backend practice project.  
The goal of the project is to demonstrate basic API development concepts with Go, Fiber and MongoDB.

---

## 📄 License

This project is licensed under the MIT License.