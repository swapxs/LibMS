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
            <Link to="/owner/assign-admin"><i class="fa-solid fa-users"></i> Manage Users</Link>
            <Link to="/owner/book-status"><i class="fa-solid fa-book-open-reader"></i> Book Status</Link>
          </>
        )}
        {user && (user.role === "admin" || user.role === "LibraryAdmin") && (
          <>
            <Link to="/admin/add-book"><i class="fa-solid fa-square-pen"></i> Add Books</Link>
            <Link to="/admin/remove-book"><i class="fa-solid fa-trash"></i> Remove Books</Link>
            <Link to="/admin/update-book"><i class="fa-solid fa-bookmark"></i> Update Books</Link>
            <Link to="/admin/issue-requests"><i class="fa-solid fa-bell"></i> Issue Requests</Link>
            <Link to="/admin/all-books"><i class="fa-solid fa-search"></i> Search Books</Link>
          </>
        )}
        {user && user.role === "Reader" && (
          <>
            <Link to="/user/all-books"><i class="fa-solid fa-search"></i> Search Books</Link>
            <Link to="/user/issue-request"><i class="fa-solid fa-cart-plus"></i> Issue Books</Link>
          </>
        )}
        {user && (
          <button className="logoutbtn" onClick={handleLogout}>
            <i className="fa-solid fa-right-from-bracket"></i>
          </button>
        )}
      </div>
    </nav>
  );
}

export default NavBar;
