// frontend/src/components/NavBar.jsx
import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function NavBar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    if (window.confirm("Are you sure you want to logout?")) {
      logout();
      navigate("/login");
    }
  };

  return (
    <nav className="navbar">
      <span className="title">
        <Link to="/">LIBRARY MANAGEMENT SYSTEM</Link>
      </span>
      <div className="nav-links">
        {!user && (
          <>
            <Link to="/login">
              <i className="fa-solid fa-right-to-bracket"></i> Login
            </Link>
            <Link to="/register">
              <i className="fa-solid fa-user-plus"></i> Register
            </Link>
          </>
        )}
        {user && user.role === "Owner" && (
          <>
            <Link to="/owner/assign-admin">Manage Users</Link>
            <Link to="/owner/book-status">Book Status</Link>
          </>
        )}
        {user && (user.role === "admin" || user.role === "LibraryAdmin") && (
          <>
            <Link to="/admin/add-book">Add Books</Link>
            <Link to="/admin/remove-book">Remove Books</Link>
            <Link to="/admin/update-book">Update Books</Link>
            <Link to="/admin/issue-requests">Issue Requests</Link>
            <Link to="/admin/all-books">All Books</Link>
          </>
        )}
        {user && user.role === "Reader" && (
          <>
            <Link to="/user/all-books">Search &amp; Issue Books</Link>
            <Link to="/user/issue-request">Manual Issue Request</Link>
          </>
        )}
        {user && (
          <button className="logoutbtn" onClick={handleLogout}>
            <i className="fa-solid fa-right-from-bracket"></i> Logout
          </button>
        )}
      </div>
    </nav>
  );
}

export default NavBar;
