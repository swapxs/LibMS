// frontend/src/components/User/BookSearch.jsx
import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";

const BookSearch = () => {
  const { user } = useAuth();
  const [query, setQuery] = useState("");
  const [books, setBooks] = useState([]);
  const [message, setMessage] = useState("");

  const handleSearch = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.REACT_APP_API_URL || "http://localhost:5000"}/api/books?query=${query}`,
        { headers: { Authorization: `Bearer ${user.token}` } }
      );
      const data = await response.json();
      if (data.books) {
        setBooks(data.books);
      } else {
        setMessage(data.error || "No books found");
      }
    } catch (err) {
      setMessage("Error searching books: " + err.message);
    }
  };

  const handleIssue = async (isbn) => {
    try {
      const response = await fetch(
        `${process.env.REACT_APP_API_URL || "http://localhost:5000"}/api/issueRegistry`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${user.token}`,
          },
          body: JSON.stringify({ isbn }),
        }
      );
      const data = await response.json();
      setMessage(data.message || data.error);
    } catch (err) {
      setMessage("Error issuing request: " + err.message);
    }
  };

  return (
    <div className="container">
      <h3>Search Books</h3>
      {message && <p style={{ color: "red" }}>{message}</p>}
      <form onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search by title, author, or publisher"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          required
        />
        <button type="submit">Search</button>
      </form>
      <div className="card-container">
        {books.map((book) => (
          <div key={book.isbn} className="book-card">
            <div className="book-card-header">
              <h3>{book.title}</h3>
            </div>
            <hr />
            <div className="book-card-body">
              <p><strong>Authors:</strong> {book.authors}</p>
              <p><strong>Publisher:</strong> {book.publisher}</p>
              <button onClick={() => handleIssue(book.isbn)}>Issue Book</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default BookSearch;
