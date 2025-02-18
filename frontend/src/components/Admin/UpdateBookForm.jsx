import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

const UpdateBookForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [bookFound, setBookFound] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    author: "",
    publisher: "",
    language: "",
    version: "",
    total_copies: "",
  });
  const [message, setMessage] = useState("");

  // Handler for input changes for book details.
  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  // Handler to search for a book by ISBN.
  const handleSearch = async () => {
    if (!isbn) {
      setMessage("Please provide an ISBN.");
      return;
    }
    try {
      // Fetch books from the backend.
      const response = await apiService.getBooks(user.token);
      console.log("Books response:", response);
      let books = [];
      // Check for different response formats.
      if (response.books && Array.isArray(response.books)) {
        books = response.books;
      } else if (response.success && Array.isArray(response.data)) {
        books = response.data;
      } else if (Array.isArray(response)) {
        books = response;
      } else {
        setMessage("Invalid response format.");
        return;
      }
      // Find the book with a matching ISBN.
      const book = books.find((b) => b.ISBN === isbn);
      if (book) {
        // Populate formData using the values from the matched book.
        setFormData({
          title: book.Title || "",
          author: book.Author || "",
          publisher: book.Publisher || "",
          language: book.Language || "",
          version: book.Version || "",
          total_copies: book.TotalCopies ? String(book.TotalCopies) : "",
        });
        setBookFound(true);
        setMessage("");
      } else {
        setMessage("Book not found.");
        setBookFound(false);
      }
    } catch (error) {
      console.error("Error searching for book:", error);
      setMessage("Error searching for book.");
      setBookFound(false);
    }
  };

  // Handler for updating the book details.
  const handleUpdate = async (e) => {
    e.preventDefault();
    if (!isbn) {
      setMessage("Please provide the ISBN of the book to update.");
      return;
    }
    try {
      const payload = {};
      // Only include provided fields.
      Object.keys(formData).forEach((key) => {
        if (formData[key] !== "") {
          // For total_copies, convert to number.
          payload[key] = key === "total_copies" ? Number(formData[key]) : formData[key];
        }
      });
      const result = await apiService.updateBook(isbn, payload, user.token);
      setMessage(result.message || result.error || "Update failed");
    } catch (err) {
      console.error("Error updating book:", err);
      setMessage("An error occurred.");
    }
  };

  return (
    <div className="card">
      <h3>Update Book Details</h3>
      {message && <p>{message}</p>}
      {/* If no book has been found yet, show only the ISBN search field */}
      {!bookFound && (
                <>
        <div className="form-group">
          <label htmlFor="isbn">Enter ISBN of Book:</label>
          <input
            type="text"
            name="isbn"
            value={isbn}
            onChange={(e) => setIsbn(e.target.value)}
            required
          />
                    </div>
          <button type="button" onClick={handleSearch}>
            Search
          </button>
</>
      )}
      {/* If the book is found, display the form populated with its details */}
      {bookFound && (
        <form onSubmit={handleUpdate}>
          <div className="form-group">
            <label htmlFor="isbn">ISBN:</label>
            <input type="text" name="isbn" value={isbn} disabled />
          </div>
          <div className="form-group">
            <label htmlFor="title">Title:</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="author">Author:</label>
            <input
              type="text"
              name="author"
              value={formData.author}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="publisher">Publisher:</label>
            <input
              type="text"
              name="publisher"
              value={formData.publisher}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="language">Language:</label>
            <input
              type="text"
              name="language"
              value={formData.language}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="version">Version:</label>
            <input
              type="text"
              name="version"
              value={formData.version}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="total_copies">Total Copies:</label>
            <input
              type="number"
              name="total_copies"
              value={formData.total_copies}
              onChange={handleChange}
            />
          </div>
          <button type="submit">Update Book</button>
        </form>
      )}
    </div>
  );
};

export default UpdateBookForm;
