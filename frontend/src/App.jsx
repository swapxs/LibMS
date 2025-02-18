// frontend/src/App.jsx
import React from "react";
import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import NavBar from "./components/NavBar";

// Auth pages
import Login from "./components/Auth/Login";
import Register from "./components/Auth/Register";
import OwnerRegister from "./components/Auth/OwnerRegister";

// Owner pages
import AssignAdmin from "./components/Owner/AssignAdmin";
import BookStatus from "./components/Owner/BookStatus";

// Admin pages
import AddBookForm from "./components/Admin/AddBookForm";
import RemoveBookForm from "./components/Admin/RemoveBookForm";
import UpdateBookForm from "./components/Admin/UpdateBookForm";
import IssueRequestList from "./components/Admin/IssueRequestList";
import Books from "./components/Books";

// Reader pages
import IssueRequestForm from "./components/User/IssueRequestForm";


import "./styles/main.css";

// Unified Dashboard Component
function Dashboard() {
  // Retrieve user data from localStorage (this assumes your Login component stores authData there)
  const storedUser = localStorage.getItem("authData");
  const user = storedUser ? JSON.parse(storedUser) : null;

  return (
    <div className="card">
      {user ? (
        <>
          <h2>Hello {user.role}!</h2>
          <p>Welcome to the online library.</p>
        </>
      ) : (
        <div>
          <h2>Welcome to the online library.</h2>
          <p>Please log in to see your dashboard.</p>
        </div>
      )}
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <NavBar />
      <div className="container">
        <Routes>
          {/* Auth Routes */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/owner/register" element={<OwnerRegister />} />

          {/* Unified Dashboard */}
          <Route path="/dashboard" element={<Dashboard />} />

          {/* Owner Routes */}
          <Route path="/owner/assign-admin" element={<AssignAdmin />} />
          <Route path="/owner/book-status" element={<BookStatus />} />

          {/* Admin Routes */}
          <Route path="/admin/add-book" element={<AddBookForm />} />
          <Route path="/admin/remove-book" element={<RemoveBookForm />} />
          <Route path="/admin/update-book" element={<UpdateBookForm />} />
          <Route path="/admin/issue-requests" element={<IssueRequestList />} />
          <Route path="/admin/all-books" element={<Books />} />

          {/* Reader Routes */}
          <Route path="/user/issue-request" element={<IssueRequestForm />} />
          <Route path="/user/all-books" element={<Books />} />


          {/* Fallback: Display a 404 Not Found page */}
          <Route path="/" element={<Dashboard />} />
          <Route
            path="*"
            element={
              <div className="card">
                <h2>404 Not Found</h2>
                <p>Error 404: Page not found.</p>
              </div>
            }
          />
        </Routes>
      </div>
    </AuthProvider>
  );
}

export default App;
