import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import NavBar from './components/NavBar';

// Import Auth components from the Auth directory
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';

// Import Admin components
import AddBookForm from './components/Admin/AddBookForm';
import RemoveBookForm from './components/Admin/RemoveBookForm';
import UpdateBookForm from './components/Admin/UpdateBookForm';
import IssueRequestList from './components/Admin/IssueRequestList';

// Import User components
import IssueRequestForm from './components/User/IssueRequestForm';
import Books from './components/Books';

// Import Owner pages
import OwnerRegister from './components/Auth/OwnerRegister';
import AssignAdmin from './components/Owner/AssignAdmin';
import BookStatus from './components/Owner/BookStatus';

import './styles/main.css';

function App() {
  return (
    <AuthProvider>
      <div>
        <header>
          <NavBar />
        </header>
        <div className="container">
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/owner/register" element={<OwnerRegister />} />
            <Route path="/owner/assign-admin" element={<AssignAdmin />} />
            <Route path="/owner/book-status" element={<BookStatus />} />
            <Route path="/admin/add-book" element={<AddBookForm />} />
            <Route path="/admin/remove-book" element={<RemoveBookForm />} />
            <Route path="/admin/update-book" element={<UpdateBookForm />} />
            <Route path="/admin/issue-requests" element={<IssueRequestList />} />
            <Route path="/user/issue-request" element={<IssueRequestForm />} />
            <Route path="/admin/all-books" element={<Books />} />
            <Route path="/user/all-books" element={<Books />} />
            <Route
              path="/"
              element={
                <div className="card">
                  <h2>Welcome to the Library Management System</h2>
                  <p>Manage your library with ease.</p>
                </div>
              }
            />
            {/* <Route path="*" element={<Navigate to="/login" replace />} /> */}
          </Routes>
        </div>
      </div>
    </AuthProvider>
  );
}

export default App;
