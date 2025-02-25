# **Frontend Workflow - Library Management System (LibMS)**

## **Overview**
The frontend of the **Library Management System (LibMS)** provides an intuitive user interface for **owners, library admins, and readers** to interact with the system. Built using **ReactJS** with **React Router v6**, it ensures smooth navigation, state management using **Context API**, and a clean UI using **vanilla CSS**.

---

## **Technology Stack**
- **Framework:** ReactJS (Functional Components & Hooks)
- **State Management:** Context API (AuthContext)
- **Routing:** React Router v6 for navigation
- **Styling:** Vanilla CSS with Google Font **Poppins**
- **API Handling:** Fetch API to interact with the backend
- **Form Validation:** Controlled components with validation checks
- **Error Handling:** User-friendly error messages and UI feedback

---

## **Frontend Project Structure**
```
frontend
├── src
│   ├── components      # UI components
│   │   ├── Auth        # Login & Registration
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── OwnerRegister.jsx
│   │   ├── Admin       # Admin functionalities
│   │   │   ├── AddBookForm.jsx
│   │   │   ├── RemoveBookForm.jsx
│   │   │   ├── UpdateBookForm.jsx
│   │   │   ├── IssueRequestList.jsx
│   │   ├── Owner       # Owner functionalities
│   │   │   ├── AssignAdmin.jsx
│   │   │   ├── BookStatus.jsx
│   │   ├── User        # Reader functionalities
│   │   │   ├── BookSearch.jsx
│   │   │   ├── IssueRequestForm.jsx
│   │   ├── Dashboard.jsx  # Role-based Dashboard
│   │   ├── NavBar.jsx      # Navigation Bar
│   ├── context
│   │   ├── AuthContext.jsx  # Authentication & State Management
│   ├── services
│   │   ├── apiService.js  # API Requests
│   ├── styles
│   │   ├── main.css
│   ├── App.jsx
│   ├── index.jsx
├── package.json
├── public
│   ├── index.html
```

---

## **Authentication Workflow**
### **1️⃣ User Login (`Login.jsx`)**
1. User enters **email** and **password**.
2. Data is sent to `POST /api/auth/login`.
3. If credentials are valid:
   - JWT token is **stored in localStorage**.
   - User is redirected to the dashboard (`/dashboard`).
4. If invalid, an error message is displayed.

### **2️⃣ User Registration (`Register.jsx`)**
1. User submits **name, email, password, contact number, and library ID**.
2. Data is sent to `POST /api/auth/register`.
3. If successful, the user is redirected to the login page.

### **3️⃣ Authentication Context (`AuthContext.jsx`)**
- Manages authentication state using `useContext`.
- Stores user data and token in `localStorage`.
- Provides **login, logout, and auth state** to components.

```jsx
export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(() => {
        return JSON.parse(localStorage.getItem("authData")) || null;
    });

    const login = (data) => {
        setUser(data);
        localStorage.setItem("authData", JSON.stringify(data));
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem("authData");
    };

    return (
        <AuthContext.Provider value={{ user, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};
```

---

## **Book Management Workflow**
### **1️⃣ Add Book (`AddBookForm.jsx`)**
1. Admin enters **ISBN, Title, Author, Publisher, Language, Copies**.
2. Data is sent to `POST /api/books`.
3. If successful, confirmation is displayed and book list is updated.

### **2️⃣ Remove Book (`RemoveBookForm.jsx`)**
1. Admin searches for a book using **ISBN**.
2. Selects how many copies to remove.
3. Sends request to `POST /api/books/remove`.
4. If all copies are removed, book entry is deleted.

### **3️⃣ Update Book (`UpdateBookForm.jsx`)**
1. Admin searches for book by **ISBN**.
2. Updates required fields and submits the form.
3. Sends request to `PUT /api/books/:isbn`.

---

## **Request Handling Workflow**
### **1️⃣ Raise Book Request (`IssueRequestForm.jsx`)**
1. Reader searches for a book.
2. Clicks **Request Issue** button.
3. Sends request to `POST /api/requestEvents`.
4. If successful, request is logged and displayed in user dashboard.

### **2️⃣ Admin Approve/Reject Requests (`IssueRequestList.jsx`)**
1. Admin views pending requests.
2. Can **Approve** or **Reject**.
3. Updates request status in `PUT /api/issueRequests/:id`.

---

## **Navigation & Role-Based Access**
### **1️⃣ Role-Based Dashboard (`Dashboard.jsx`)**
- **Owner:** Sees user management, book status.
- **Library Admin:** Manages books and issue requests.
- **Reader:** Can search for books, request issues.

```jsx
return (
    <div>
        {user.role === "Owner" && <OwnerDashboard />}
        {user.role === "LibraryAdmin" && <AdminDashboard />}
        {user.role === "Reader" && <ReaderDashboard />}
    </div>
);
```

### **2️⃣ Navigation Bar (`NavBar.jsx`)**
- Displays menu items **based on user role**.
- **Logout button** removes token from `localStorage`.

```jsx
{user && user.role === "LibraryAdmin" && (
    <Link to="/admin/add-book">Add Books</Link>
)}
```

---

## **API Endpoints Summary**
### **Authentication**
- `POST /api/auth/register` → Register new user
- `POST /api/auth/login` → User login & JWT token generation

### **Books**
- `POST /api/books` → Add/increment book copies
- `GET /api/books` → Retrieve all books
- `POST /api/books/remove` → Remove book copies
- `PUT /api/books/:isbn` → Update book details

### **Book Requests**
- `POST /api/requestEvents` → Request book issue
- `GET /api/issueRequests` → Get all book requests
- `PUT /api/issueRequests/:id` → Approve/reject issue request
