// frontend/src/components/Books.jsx
import React, { useState, useEffect } from 'react';
import apiService from '../services/apiService';
import { useAuth } from '../context/AuthContext';

function Books() {
  const { user } = useAuth();
  const [books, setBooks] = useState([]);
  const [filteredBooks, setFilteredBooks] = useState([]);
  const [activeRequests, setActiveRequests] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');

  useEffect(() => {
    async function fetchBooks() {
      try {
        const response = await apiService.getBooks(user.token);
        let fetchedBooks = [];
        if (response.success && Array.isArray(response.data)) {
          fetchedBooks = response.data;
        } else if (Array.isArray(response.books)) {
          fetchedBooks = response.books;
        } else if (Array.isArray(response)) {
          fetchedBooks = response;
        }
        setBooks(fetchedBooks);
        setFilteredBooks(fetchedBooks);
      } catch (error) {
        console.error('Error fetching books:', error);
      } finally {
        setLoading(false);
      }
    }
    async function fetchActiveRequests() {
      try {
        const response = await apiService.getUserIssueInfo(user.token);
        if (response.success && Array.isArray(response.data)) {
          setActiveRequests(response.data.map((request) => request.isbn));
        }
      } catch (error) {
        console.error('Error fetching active requests:', error);
      }
    }
    if (user && user.token) {
      fetchBooks();
      fetchActiveRequests();
    }
  }, [user]);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredBooks(books);
    } else {
      const query = searchQuery.toLowerCase();
      setFilteredBooks(
        books.filter(book =>
          ((book.isbn || book.ISBN) && (book.isbn || book.ISBN).toLowerCase().includes(query)) ||
          ((book.title || book.Title) && (book.title || book.Title).toLowerCase().includes(query)) ||
          ((book.author || book.Author) && (book.author || book.Author).toLowerCase().includes(query)) ||
          ((book.publisher || book.Publisher) && (book.publisher || book.Publisher).toLowerCase().includes(query))
        )
      );
    }
  }, [searchQuery, books]);

  const handleIssueBook = async (isbn) => {
    try {
      const response = await fetch("http://localhost:5000/api/requestEvents", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${user.token}`,
        },
        body: JSON.stringify({ bookID: isbn }),
      });
      if (!response.ok) {
        const text = await response.text();
        throw new Error(text);
      }
      setActiveRequests((prev) => [...prev, isbn]);
      setMessage(<p className='success'>Book issue request submitted successfully.</p>);
    } catch (err) {
      setMessage(<p className='error'>Error submitting issue request. You have exhausted you quota for the time being</p>);
            console.log(err.message)
    }
  };

  return (
    <div className="card">
      <h2>All Books</h2>
      {message && <p style={{ color: "red" }}>{message}</p>}
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
              <th>ISBN</th>
              <th>Title</th>
              <th>Author</th>
              <th>Publisher</th>
              <th>Language</th>
              {user && user.role === 'Reader' && <th>Action</th>}
            </tr>
          </thead>
          <tbody>
            {filteredBooks.map((book) => (
              <tr key={book.isbn || book.ISBN}>
                <td>{book.isbn || book.ISBN}</td>
                <td>{book.title || book.Title}</td>
                <td>{book.author || book.Author}</td>
                <td>{book.publisher || book.Publisher}</td>
                <td>{book.language || book.Language}</td>
                {user && user.role === 'Reader' && (
                  <td className="profile-actions">
                    <button
                      onClick={() => handleIssueBook(book.isbn || book.ISBN)}
                      id={`btn-${book.isbn || book.ISBN}`}
                      style={{ display: activeRequests.includes(book.isbn || book.ISBN) ? 'none' : 'inline-block' }}
                      disabled={activeRequests.includes(book.isbn || book.ISBN) || activeRequests.length >= 4}
                    >
                      Issue
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default Books;
