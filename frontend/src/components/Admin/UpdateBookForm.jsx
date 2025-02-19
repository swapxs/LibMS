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
  const [error, setError] = useState("");

  // Handler for input changes for book details.
  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  // Handler to search for a book by ISBN.
  const handleSearch = async () => {
    setMessage("");
    setError("");

    if (!isbn) {
      setError("Please provide an ISBN.");
      return;
    }
    try {
      const response = await apiService.getBooks(user.token);
      console.log("Books response:", response);
      let books = [];

      if (response.books && Array.isArray(response.books)) {
        books = response.books;
      } else if (response.success && Array.isArray(response.data)) {
        books = response.data;
      } else if (Array.isArray(response)) {
        books = response;
      } else {
        setError("Invalid response format.");
        return;
      }

      const book = books.find((b) => b.ISBN === isbn);
      if (book) {
        setFormData({
          title: book.Title || "",
          author: book.Author || "",
          publisher: book.Publisher || "",
          language: book.Language || "",
          version: book.Version || "",
          total_copies: book.TotalCopies ? String(book.TotalCopies) : "",
        });
        setBookFound(true);
        setMessage("Book found!");
      } else {
        setError("Book not found.");
        setBookFound(false);
      }
    } catch (error) {
      console.error("Error searching for book:", error);
      setError("Error searching for book.");
      setBookFound(false);
    }
  };

  // Handler for updating the book details.
  const handleUpdate = async (e) => {
    e.preventDefault();
    setMessage("");
    setError("");

    if (!isbn) {
      setError("Please provide the ISBN of the book to update.");
      return;
    }
    try {
      const payload = {};
      Object.keys(formData).forEach((key) => {
        if (formData[key] !== "") {
          payload[key] = key === "total_copies" ? Number(formData[key]) : formData[key];
        }
      });

      const result = await apiService.updateBook(isbn, payload, user.token);
      
      if (result.success || result.message.toLowerCase().includes("updated")) {
        setMessage(result.message || "Book details updated successfully.");

        // Reset form after success
        setTimeout(() => {
          setIsbn("");
          setBookFound(false);
          setFormData({
            title: "",
            author: "",
            publisher: "",
            language: "",
            version: "",
            total_copies: "",
          });
          setMessage("");
          setError("");
          window.location.reload(); // Reload the page
        }, 2000);
      } else {
        setError(result.error || "Failed to update book details.");
      }
    } catch (err) {
      console.error("Error updating book:", err);
      setError("An error occurred while updating the book.");
    }
  };

  return (
    <div className="card">
      <h3>Update Book Details</h3>

      {error && <p className="error">{error}</p>}
      {message && <p className="success">{message}</p>}

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
