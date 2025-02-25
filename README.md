# **Library Management System**

## **📌 Overview**
The **Library Management System (LibMS)** is a **full-stack application** that enables libraries to manage book inventories, process issue requests, and oversee borrowing history. It supports three roles:

- **Owner**: Manages the library, assigns admin roles, and oversees administrative functions, including policy settings, performance monitoring, and access control.
- **Library Admin**: Adds/updates books, processes book requests, maintains the library database, and manages issue and return requests in an efficient manner.
- **Reader**: Searches for books, requests book issues, keeps track of borrowing history, and interacts with library resources.

This system is designed to streamline library operations by reducing manual work and increasing efficiency. It offers features such as automated book tracking, real-time request handling, and secured user authentication, ensuring a seamless user experience. Additionally, it provides **advanced reporting** features, analytics on book usage, and logs for auditing purposes.

---

## **🛠️ Technology Stack**

### **Backend:**
- **Language**: Go (Gin framework) for a lightweight and high-performance API.
- **Database**: PostgreSQL (GORM ORM) to manage relational data with efficient querying.
- **Authentication**: JWT-based authentication for secure user access and role-based control.
- **Routing & Middleware**: Gin framework with JWT Middleware for request handling, session management, and authentication.
- **Configuration**: `.env` files for managing environment-specific variables and configurations.
- **Logging**: Implementing structured logging for better debugging and monitoring.
- **Testing**: Unit tests with Go’s built-in testing framework and `sqlmock` for database simulation.

### **Frontend:**
- **Framework**: ReactJS (Functional Components & Hooks) for a modern, component-based UI.
- **State Management**: Context API (AuthContext) to manage authentication states and role-based access.
- **Styling**: Vanilla CSS (Poppins font, library-themed color scheme) to maintain a clean and user-friendly UI.
- **Routing**: React Router v6 for efficient navigation between pages.
- **Responsive Design**: The frontend is optimized for different screen sizes, ensuring accessibility across devices.
- **Form Validation**: Uses controlled form elements and validation checks to prevent erroneous data entry.
- **Error Handling**: Implements user-friendly error messages and fallbacks for network or API failures.

---

## **💁️ Project Structure**

### **Backend (`/backend`)**
```
backend/
├── gen_report.sh
├── go.mod
├── go.sum
├── README.md
├── src
│   ├── db
│   │   └── db.go
│   ├── handlers
│   │   ├── auth_handler.go
│   │   ├── book_handler.go
│   │   ├── claims_handler.go
│   │   ├── issue_handler.go
│   │   ├── library_handler.go
│   │   ├── owner_handler.go
│   │   ├── request_events_handler.go
│   │   └── user_handler.go
│   ├── main.go
│   ├── middleware
│   │   └── jwt.go
│   ├── models
│   │   ├── book_inventory_model.go
│   │   ├── issue_registry_model.go
│   │   ├── library_model.go
│   │   ├── request_events_model.go
│   │   └── user_model.go
│   └── routes
│       └── routes.go
└── test
    ├── books_test.go
    ├── db_setup_test.go
    ├── issue_request_test.go
    ├── jwt_test.go
    ├── library_test.go
    ├── login_user_test.go
    ├── negative_test.go
    ├── owner_operations_test.go
    ├── raise_request_test.go
    ├── register_user_test.go
    └── user_test.go
```

### **Frontend (`/frontend`)**
```
frontend/
├── package.json
├── public
│   └── index.html
├── README.md
└── src
    ├── App.jsx
    ├── assets
    │   ├── AdminBanner.gif
    │   ├── MainBanner.gif
    │   ├── OwnerBanner.gif
    │   └── ReaderBanner.gif
    ├── components
    │   ├── AccountDetails.jsx
    │   ├── Admin
    │   │   ├── AddBookForm.jsx
    │   │   ├── IssueRequestList.jsx
    │   │   ├── RemoveBookForm.jsx
    │   │   └── UpdateBookForm.jsx
    │   ├── Auth
    │   │   ├── Login.jsx
    │   │   ├── OwnerRegister.jsx
    │   │   └── Register.jsx
    │   ├── Books.jsx
    │   ├── Dashboard.jsx
    │   ├── NavBar.jsx
    │   ├── Owner
    │   │   ├── AssignAdmin.jsx
    │   │   └── BookStatus.jsx
    │   └── User
    │       ├── BookSearch.jsx
    │       └── IssueRequestForm.jsx
    ├── context
    │   └── AuthContext.jsx
    ├── index.jsx
    ├── services
    │   └── apiService.js
    └── styles
        └── main.css
```

---

## **📄 Database Schema**
![ERD](./img/er.png)

The database schema follows a relational structure that efficiently manages books, users, and transactions. It ensures data integrity and minimizes redundancy by implementing foreign key constraints where necessary. The design allows easy scalability for adding more libraries, book categories, and multi-location tracking.

---

## **🔒 Authentication & Authorization**
- **JWT-based authentication** ensures secure access control.
- **Middleware (`jwt.go`)**:
  - Extracts token from the request header for validation.
  - Verifies token authenticity using a secret key.
  - Sets user claims in context for role-based access, preventing unauthorized actions.
  - Ensures session expiration policies are enforced.

---

## **📌 API Endpoints**

### **Authentication (`auth.go`)**
- `POST /api/auth/register` → Register a user with encrypted credentials.
- `POST /api/auth/login` → Authenticate and return a JWT token for secure access.

### **Library (`library.go`)**
- `POST /api/library` → Create a new library with a unique name.
- `GET /api/libraries` → Fetch all registered libraries.

### **Books (`book.go`)**
- `POST /api/books` → Add a new book / Increment book copies with metadata.
- `GET /api/books` → Retrieve all books available in a library.
- `POST /api/books/remove` → Remove book copies, updating the inventory.
- `PUT /api/books/:isbn` → Modify book details such as author, publisher, and available copies.

### **Requests (`request_events.go`)**
- `POST /api/requestEvents` → Raise a book issue request with validation.
- `GET /api/issueRequests` → List all issue requests with statuses (Pending, Approved, Rejected).

### **Issues (`issue.go`)**
- `POST /api/issueRegistry` → Issue a book to a reader, updating the inventory.
- `POST /api/issueRegistry/return` → Mark book as returned, updating inventory and return logs.

### **Admin & Owner Actions (`owner.go`)**
- `POST /api/owner/assign-admin` → Promote a user to LibraryAdmin for enhanced privileges.
- `POST /api/owner/revoke-admin` → Revoke admin privileges, demoting the user to Reader status.
- `GET /api/owner/audit-logs` → Retrieve action logs for auditing purposes.

