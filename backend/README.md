# **Backend Workflow - Library Management System**

## **Overview**
The backend of the **Library Management System (LibMS)** is responsible for handling user authentication, managing book inventories, processing issue requests, and maintaining user roles and permissions. Built with **Go (Gin framework)** and **PostgreSQL (GORM ORM)**, it ensures high performance and secure operations.

---

## **Technology Stack**
- **Language:** Go (Gin framework)
- **Database:** PostgreSQL (GORM ORM), SQLite (testing)
- **Authentication:** JWT-based authentication
- **Middleware:** JWT for authorization, CORS for cross-origin requests
- **Configuration:** `.env` files for environment management
- **Logging & Error Handling:** Structured logging for debugging and monitoring
- **Testing:** Unit testing using Go's built-in framework and `sqlmock` for database interactions

## NOTE
The test report is saved as a csv in `./.blob/report`. You can generate the report via running the shellscipt `./gen_report.sh`


---

## **Backend Project Structure**
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

---

## **Database Schema**
![ERD](./img/er.png)

The database schema follows a relational structure that efficiently manages books, users, and transactions. It ensures data integrity and minimizes redundancy by implementing foreign key constraints where necessary. The design allows easy scalability for adding more libraries, book categories, and multi-location tracking.

## **Authentication Workflow**
### **User Registration (`POST /api/auth/register`)**
1. User submits registration details (name, email, password, contact, role, library ID).
2. Password is **hashed** using `bcrypt`.
3. Data is saved in the `users` table with a default role of `Reader`.
4. Response: `201 Created` with success message.

### **User Login (`POST /api/auth/login`)**
1. User submits credentials (email & password).
2. Credentials are **validated** against the database.
3. If valid, a **JWT token** is generated and returned.
4. Token contains user ID, role, and library ID for authorization.

### **JWT Authentication Middleware (`jwt.go`)**
- Extracts JWT token from the `Authorization` header.
- Validates the token.
- Sets user claims in context for role-based access.

---

## **Book Management Workflow**
### **Add Book (`POST /api/books`)**
1. Admin submits book details (ISBN, title, author, copies, etc.).
2. If book exists, copies are incremented.
3. If book is new, a **new record** is created in `book_inventory`.

### **Remove Book (`POST /api/books/remove`)**
1. Admin selects a book via ISBN.
2. Requested number of copies are deducted from the inventory.
3. If all copies are removed, the **book record is deleted**.

### **Update Book (`PUT /api/books/:isbn`)**
1. Admin submits updated details (title, author, language, etc.).
2. Changes are applied to the `book_inventory` table.

---

## **Request Handling Workflow**
### **Raise Book Request (`POST /api/requestEvents`)**
1. Readers can request up to **4 active book requests**.
2. System checks **book availability** before processing request.
3. Request is stored in `request_events` with status `Issue`.

### **Approve / Reject Request (`PUT /api/issueRequests/:id`)**
1. Admin reviews the request.
2. If approved:
   - Available book copies are decremented.
   - Request status is updated to `Approved`.
3. If rejected, request status is updated to `Rejected`.

---

## **Book Issue & Return Workflow**
### **Issue Book (`POST /api/issueRegistry`)**
1. Approved requests result in book issuance.
2. Entry is created in `issue_registry` with **expected return date**.
3. Book’s **available copies** are reduced in `book_inventory`.

### **Return Book (`POST /api/issueRegistry/return`)**
1. User returns the book.
2. Admin marks **return date** in `issue_registry`.
3. Available copies are incremented back in `book_inventory`.

---

## **Admin & Owner Management Workflow**
### **Assign Admin (`POST /api/owner/assign-admin`)**
1. Owner selects a user via email.
2. User’s role is updated to `LibraryAdmin`.

### **Revoke Admin (`POST /api/owner/revoke-admin`)**
1. Admin privileges are revoked.
2. User is demoted to `Reader`.

### **Audit Logs (`GET /api/owner/audit-logs`)**
1. Tracks all admin actions (book additions, role changes, request approvals).

---

## **API Endpoints Summary**
### **Authentication**
- `POST /api/auth/register` → Register new user
- `POST /api/auth/login` → User login & JWT token generation

### **Library Management**
- `POST /api/library` → Create a new library
- `GET /api/libraries` → Get all libraries

### **Book Inventory**
- `POST /api/books` → Add/increment book copies
- `GET /api/books` → Retrieve all books
- `POST /api/books/remove` → Remove book copies
- `PUT /api/books/:isbn` → Update book details

### **Book Requests**
- `POST /api/requestEvents` → Request book issue
- `GET /api/issueRequests` → Get all book requests
- `PUT /api/issueRequests/:id` → Approve/reject issue request

### **Issue & Return**
- `POST /api/issueRegistry` → Issue a book
- `POST /api/issueRegistry/return` → Return a book

### **Admin Actions**
- `POST /api/owner/assign-admin` → Assign admin role
- `POST /api/owner/revoke-admin` → Revoke admin role
- `GET /api/owner/audit-logs` → Retrieve audit logs
