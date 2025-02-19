# Online Library Management System

## Introduction
This is a full-stack webdevelopment/UI Engineering project that enables libraries to manage book inventories, process issue requests and oversee borrowing history.

This system is comprised of three key roles or entities:
- **Owner:** Manges the library and assign admin roles.
- **Library Admin:** Perform CRUD operations on books and processes book requests.
- **Reader:** Searches for books and requests issue for particular books.

## Technologies Used
### Backend
- **Language:** Go (Gin framework)
- **Database:** PostgreSQL (GORM ORM)
- **Authentication:** JWT-based authentication
- **Routing & Middleware:** Gin, JWT Middleware
- **Testing:** Go’s testing package, sqlmock
- **Configuration:** .env files, environment variables

### Frontend
- **Framework:** ReactJS
- **Styling:** Vanilla CSS
- **Routing:** React Router v6
- **Markup:** HTML

## Project Structure
### Backend (/backend)
```
backend/
 ├── main.go
 ├── db/
 │    └── db.go
 ├── models/
 │    ├── library.go
 │    ├── user.go
 │    ├── book_inventory.go
 │    ├── request_events.go
 │    └── issue_registry.go
 ├── controllers/
 │    ├── auth.go
 │    ├── book.go
 │    ├── owner.go
 │    ├── request_events.go
 │    ├── issue.go
 │    └── user.go
 ├── middleware/
 │    └── jwt.go
 ├── routes/
 │    └── routes.go
 ├── test/
 │    ├── auth_test.go
 │    ├── book_test.go
 │    ├── request_events_test.go
 │    ├── issue_test.go
 │    └── user_test.go
 └── seed.go
```


### Frontend (/frontend)

```
frontend/
 ├── src/
 │    ├── components/
 │    │    ├── Auth/
 │    │    │    ├── Login.jsx
 │    │    │    ├── Register.jsx
 │    │    │    └── OwnerRegister.jsx
 │    │    ├── Admin/
 │    │    │    ├── AddBookForm.jsx
 │    │    │    ├── RemoveBookForm.jsx
 │    │    │    ├── UpdateBookForm.jsx
 │    │    │    └── IssueRequestList.jsx
 │    │    ├── Owner/
 │    │    │    ├── AssignAdmin.jsx
 │    │    │    └── BookStatus.jsx
 │    │    ├── User/
 │    │    │    ├── IssueRequestForm.jsx
 │    │    │    └── BookSearch.jsx
 │    │    ├── AccountDetails.jsx
 │    │    ├── Books.jsx
 │    │    ├── Books.jsx
 │    │    └── NavBar.jsx
 │    ├── context/
 │    │    └── AuthContext.jsx
 │    ├── services/
 │    │    └── apiService.js
 │    ├── styles/
 │    │    └── main.css
 │    ├── App.jsx
 │    └── index.jsx
 ├── public/
 │    └── index.html
 └── package.json
```

## Authentication & Authorization
- Used JWT for authentication 
- Package provided by `github.com/golang-jwt/jwt/v4`
- This is implemented using a middleware called `jwt.go`. 
    - It extracts token from the request header
    - Verifies token authenticity
    - Sets user claims in context for role-based access

## Database Schema and ERD

![ERD](./img/er.png)

## API Endpoints
### Authentication (`auth.go`)
- `POST /api/auth/register`: Register a user
- `POST /api/auth/login`: Login & get JWT token

### Library (`library.go`)
- `POST /api/library`: Create a new library
- `GET /api/libraries`: Fetch all libraries

### Books (`book.go`)
- `POST /api/books`: Add a new book / Increment book copies
- `GET /api/books`: Get books in a library
- `POST /api/books/remove`: Remove book copies
- `PUT /api/books/:isbn`: Update book details

### Requests (`request_events.go`)
- `POST /api/requestEvents`: Raise book issue request (Max 4 active requests)
- `GET /api/issueRequests`: List all issue requests
- `PUT /api/issueRequests/:id`: Approve/Reject requests

### Issues (`issue.go`)
- `POST /api/issueRegistry`: Issue a book

### Admin & Owner Actions (`owner.go`)
- `POST /api/owner/assign-admin`: Promote user to LibraryAdmin
- `POST /api/owner/revoke-admin`: Revoke admin privileges


## Frontend Walkthrough
### Key Components
- Auth Pages:
  - `Login.jsx`
  - `Register.jsx`
  - `OwnerRegister.jsx`

- Admin Pages:
  - `AddBookForm.jsx` → Add a new book
  - `RemoveBookForm.jsx` → Remove book copies
  - `UpdateBookForm.jsx` → Modify book details
  - `IssueRequestList.jsx` → Process book issue requests

- Owner Pages:
  - `AssignAdmin.jsx` → Assign/revoke admin role
  - `BookStatus.jsx` → View book availability

- Reader Pages:
  - `IssueRequestForm.jsx` → Request a book
  - `BookSearch.jsx` → Search books

### Routing (`App.jsx`)
- Uses React Router v6 to navigate between:
  - `/dashboard`
  - `/admin/*`
  - `/owner/*`
  - `/user/*`


## **📌 Key Features**
- **Secure Role-Based Access** (Owner, Admin, Reader)  
- **JWT Authentication & Authorization**  
- **RESTful API with Gin Framework**  
- **Book Inventory Management**  
- **Book Issue & Return System**  
- **Admin Dashboard for Requests**  
- **Seamless User Experience with ReactJS**  
