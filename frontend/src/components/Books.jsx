// frontend/src/components/Books.jsx
import React, { useState, useEffect } from 'react';
import apiService from '../services/apiService';
import { useAuth } from '../context/AuthContext';

function Books() {
  const { user } = useAuth();
  const [books, setBooks] = useState([]);
  const [filteredBooks, setFilteredBooks] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);

  // Fetch books when the component mounts.
  useEffect(() => {
    async function fetchBooks() {
      try {
        const response = await apiService.getBooks(user.token);
        console.log("Books response:", response);
        let fetchedBooks = [];
        // Check for different response formats.
        if (response.success && Array.isArray(response.data)) {
          fetchedBooks = response.data;
        } else if (Array.isArray(response.books)) {
          fetchedBooks = response.books;
        } else if (Array.isArray(response)) {
          fetchedBooks = response;
        } else {
          fetchedBooks = [];
        }
        setBooks(fetchedBooks);
        setFilteredBooks(fetchedBooks);
      } catch (error) {
        console.error('Error fetching books:', error);
      } finally {
        setLoading(false);
      }
    }
    if (user && user.token) {
      fetchBooks();
    }
  }, [user]);

  // Filter books based on search query.
  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredBooks(books);
    } else {
      const query = searchQuery.toLowerCase();
      setFilteredBooks(
        books.filter(book =>
          ((book.isbn || book.ISBN ) && (book.isbn || book.ISBN).toLowerCase().includes(query)) ||
          ((book.title || book.Title) && (book.title || book.Title).toLowerCase().includes(query)) ||
          ((book.author || book.Author) && (book.author || book.Author).toLowerCase().includes(query)) ||
          ((book.publisher || book.Publisher) && (book.publisher || book.Publisher).toLowerCase().includes(query))
        )
      );
    }
  }, [searchQuery, books]);

  return (
    <div className="card">
      <h2>All Books</h2>
      <div className="form-group">
        <label htmlFor="search">Search by Title, Author, or Publisher:</label>
        <input
          id="search"
          type="text"
          placeholder="Enter search query..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
      {loading ? (
        <p>Loading books...</p>
      ) : filteredBooks.length === 0 ? (
        <p>No books found.</p>
      ) : (
        <table className="books-table">
          <thead>
            <tr>
              {/* <th>Book ID</th> */}
              <th>ISBN</th>
              <th>Title</th>
              <th>Author</th>
              <th>Publisher</th>
              <th>Language</th>
              <th>Total Copies</th>
              <th>Available Copies</th>
            </tr>
          </thead>
          <tbody>
            {filteredBooks.map((book) => (
              <tr key={book.id || book.ID}>
                {/* <td>{book.id || book.ID}</td> */}
                <td>{book.isbn || book.ISBN}</td>
                <td>{book.title || book.Title}</td>
                <td>{book.author || book.Author}</td>
                <td>{book.publisher || book.Publisher}</td>
                <td>{book.language || book.Language}</td>
                <td>{book.totalCopies || book.TotalCopies}</td>
                <td>{book.availableCopies || book.AvailableCopies}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default Books;
