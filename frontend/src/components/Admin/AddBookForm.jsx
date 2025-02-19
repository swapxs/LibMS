// frontend/src/components/Admin/AddBookForm.jsx
import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

function AddBookForm() {
  const { user } = useAuth();
  const [isNewBook, setIsNewBook] = useState(true);
  const [formData, setFormData] = useState({
    isbn: "",
    title: "",
    author: "",       // Single author field
    publisher: "",
    language: "",     // New input field for language
    version: "",
    copies: "",
  });
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleToggle = () => {
    setIsNewBook(!isNewBook);
    setFormData({
      isbn: "",
      title: "",
      author: "",
      publisher: "",
      language: "",
      version: "",
      copies: "",
    });
    setMessage("");
    setError("");
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    // If adding a new book, require Title, Author, and Language.
    if (isNewBook && (!formData.title || !formData.author || !formData.language)) {
      setError("Title, Author, and Language are required for a new book.");
      return;
    }
    try {
      const payload = {
        isbn: formData.isbn,
        title: formData.title,
        author: formData.author,
        publisher: formData.publisher,
        language: formData.language,
        version: formData.version,
        copies: Number(formData.copies),
        increment_only: !isNewBook,
      };
      const result = await apiService.addBook(payload, user.token);
      // Check if the result message indicates a successful operation.
      if (
        result.message &&
        (result.message.toLowerCase().includes("added") ||
          result.message.toLowerCase().includes("incremented"))
      ) {
        setMessage(result.message);
        setError("");
      } else {
        setError(result.error || "Operation failed");
        setMessage("");
      }
    } catch (err) {
      console.error("Error adding book:", err);
      setError("An error occurred while adding the book.");
      setMessage("");
    }
  };

  return (
    <div className="card">
    <div className="toggle-header">
        <button
          className={isNewBook ? "active" : "inactive"}
          onClick={() => { if (!isNewBook) handleToggle(true);}}
        >
          Add New Book
        </button>
        <button
          className={!isNewBook ? "active" : "inactive"}
          onClick={() => { if (isNewBook) handleToggle(false);}}
        >
          Add Copies
        </button>
      </div>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="isbn">ISBN:</label>
          <input
            type="text"
            id="isbn"
            name="isbn"
            value={formData.isbn}
            onChange={handleChange}
            required
          />
        </div>
        {isNewBook && (
          <>
            <div className="form-group">
              <label htmlFor="title">Title:</label>
              <input
                type="text"
                id="title"
                name="title"
                value={formData.title}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="author">Author:</label>
              <input
                type="text"
                id="author"
                name="author"
                value={formData.author}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="publisher">Publisher:</label>
              <input
                type="text"
                id="publisher"
                name="publisher"
                value={formData.publisher}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="language">Language:</label>
              <input
                type="text"
                id="language"
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
                id="version"
                name="version"
                value={formData.version}
                onChange={handleChange}
              />
            </div>
          </>
        )}
        <div className="form-group">
          <label htmlFor="copies">Number of Copies:</label>
          <input
            type="number"
            id="copies"
            name="copies"
            value={formData.copies}
            onChange={handleChange}
            required
          />
        </div>
      {error && <p className="error">{error}</p>}
      {message && <p className="success">{message}</p>}
        <button type="submit">
          {isNewBook ? "Add New Book" : "Add Copies"}
        </button>
      </form>
    </div>
  );
}

export default AddBookForm;
