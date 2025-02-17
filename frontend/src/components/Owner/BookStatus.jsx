import React, { useState, useEffect } from "react";
import apiService from "../../services/apiService";
import { useAuth } from "../../context/AuthContext";

function BookStatus() {
  const { user } = useAuth();
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Fetch books from the backend using the token.
  const fetchBooks = async () => {
    try {
      const response = await apiService.getBooks(user.token);
      console.log("BookStatus fetched response:", response);
      let fetchedBooks = [];
      if (response.success && Array.isArray(response.data)) {
        fetchedBooks = response.data;
      } else if (Array.isArray(response.books)) {
        fetchedBooks = response.books;
      } else if (Array.isArray(response)) {
        fetchedBooks = response;
      } else {
        setError("Invalid response format from getBooks");
        fetchedBooks = [];
      }
      setBooks(fetchedBooks);
      setError("");
    } catch (err) {
      console.error("Error fetching books in BookStatus:", err);
      setError("Error fetching books");
      setBooks([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user && user.token) {
      fetchBooks();
    }
  }, [user]);

  return (
    <div className="card">
      <h2>Book Status</h2>
      <button onClick={fetchBooks}>Refresh</button>
      {loading ? (
        <p>Loading books...</p>
      ) : error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : books.length === 0 ? (
        <p>No books found.</p>
      ) : (
        <div className="cards-container">
          {books.map((book) => {
            // Use uppercase keys if available; fallback to lowercase.
            const bookID = book.ID || book.id;
            const isbn = book.ISBN || book.isbn;
            const title = book.Title || book.title;
            const author = book.Author || book.author;
            const publisher = book.Publisher || book.publisher;
            const language = book.Language || book.language;
            const totalCopies = book.TotalCopies || book.totalCopies;
            const availableCopies = book.AvailableCopies || book.availableCopies;
            const currentlyIssued = totalCopies - availableCopies;

            return (
              <div key={bookID} className="card book-card">
                <div className="book-card-header">
                  <h3>{title}</h3>
                </div>
                <div className="book-card-body">
                  <p><strong>Book ID:</strong> {bookID}</p>
                  <p><strong>ISBN:</strong> {isbn}</p>
                  <p><strong>Author:</strong> {author}</p>
                  <p><strong>Publisher:</strong> {publisher}</p>
                  <p><strong>Language:</strong> {language}</p>
                  <p><strong>Total Copies:</strong> {totalCopies}</p>
                  <p><strong>Available Copies:</strong> {availableCopies}</p>
                  <p><strong>Currently Issued:</strong> {currentlyIssued}</p>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

export default BookStatus;
