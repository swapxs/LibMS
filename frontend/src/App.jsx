import React from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import NavBar from "./components/NavBar";

// Auth pages
import Login from "./components/Auth/Login";
import Register from "./components/Auth/Register";
import OwnerRegister from "./components/Auth/OwnerRegister";

// Other pages (admin, reader, etc.)
import AddBookForm from "./components/Admin/AddBookForm";
import RemoveBookForm from "./components/Admin/RemoveBookForm";
import UpdateBookForm from "./components/Admin/UpdateBookForm";
import IssueRequestList from "./components/Admin/IssueRequestList";
import IssueRequestForm from "./components/User/IssueRequestForm";
import Books from "./components/Books";
import UserBorrowHistory from "./components/User/UserBorrowHistory";

import "./styles/main.css";

function App() {
  return (
    <AuthProvider>
      <NavBar />
      <div className="container">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/owner/register" element={<OwnerRegister />} />
          <Route path="/admin/add-book" element={<AddBookForm />} />
          <Route path="/admin/remove-book" element={<RemoveBookForm />} />
          <Route path="/admin/update-book" element={<UpdateBookForm />} />
          <Route path="/admin/issue-requests" element={<IssueRequestList />} />
          <Route path="/user/issue-request" element={<IssueRequestForm />} />
          <Route path="/admin/all-books" element={<Books />} />
          <Route path="/user/all-books" element={<Books />} />
          <Route path="/user/borrow-history" element={<UserBorrowHistory />} />
          {/* <Route path="*" element={<Navigate to="/login" replace />} /> */}
        </Routes>
      </div>
    </AuthProvider>
  );
}

export default App;
