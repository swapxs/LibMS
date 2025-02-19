import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

const RemoveBookForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [bookFound, setBookFound] = useState(false);
  const [bookData, setBookData] = useState(null);
  const [copiesToRemove, setCopiesToRemove] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  // Function to reset the component state
  const resetForm = () => {
    setIsbn("");
    setBookFound(false);
    setBookData(null);
    setCopiesToRemove("");
    setMessage("");
    setError("");
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

      // Find the book with the matching ISBN.
      const book = books.find((b) => b.ISBN === isbn);
      if (book) {
        setBookData({
          title: book.Title || book.title || "",
          author: book.Author || book.author || "",
          publisher: book.Publisher || book.publisher || "",
          language: book.Language || book.language || "",
          version: book.Version || book.version || "",
          total_copies: book.TotalCopies || book.totalCopies || 0,
          available_copies: book.AvailableCopies || book.availableCopies || 0,
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

  // Handler for removing copies.
  const handleRemove = async (e) => {
    e.preventDefault();
    setMessage("");
    setError("");

    if (!isbn) {
      setError("Please provide the ISBN of the book to update.");
      return;
    }
    if (!copiesToRemove) {
      setError("Please provide the number of copies to remove.");
      return;
    }

    const numCopies = Number(copiesToRemove);
    if (isNaN(numCopies) || numCopies <= 0) {
      setError("Copies to remove must be a number greater than 0.");
      return;
    }

    try {
      const response = await apiService.removeBook(isbn, numCopies, user.token);
      console.log("Remove book response:", response);

      if (
        response.success === true ||
        (response.message &&
          (response.message.toLowerCase().includes("deleted") ||
            response.message.toLowerCase().includes("removed")))
      ) {
        setMessage(response.message || "Book copies removed successfully.");

        // Reset form after a short delay to show the success message
        setTimeout(() => {
          resetForm();
          window.location.reload(); // Reload the page
        }, 2000);
      } else {
        setError(response.error || "Failed to remove book copies.");
      }
    } catch (err) {
      console.error("Error removing book copies:", err);
      setError("An error occurred while removing book copies.");
    }
  };

  return (
    <div className="card">
      <h3>Remove Book Copies</h3>

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

      {bookFound && bookData && (
        <div>
          <div className="card">
            <h4>Book Details</h4>
            <p><strong>Title:</strong> {bookData.title}</p>
            <p><strong>Author:</strong> {bookData.author}</p>
            <p><strong>Publisher:</strong> {bookData.publisher}</p>
            <p><strong>Language:</strong> {bookData.language}</p>
            <p><strong>Version:</strong> {bookData.version}</p>
            <p><strong>Total Copies:</strong> {bookData.total_copies}</p>
            <p><strong>Available Copies:</strong> {bookData.available_copies}</p>
          </div>
          <form onSubmit={handleRemove}>
            <div className="form-group">
              <label htmlFor="copiesToRemove">Number of Copies to Remove:</label>
              <input
                type="number"
                name="copiesToRemove"
                value={copiesToRemove}
                onChange={(e) => setCopiesToRemove(e.target.value)}
                required
              />
            </div>
            <button type="submit">Remove Copies</button>
          </form>
        </div>
      )}
    </div>
  );
};

export default RemoveBookForm;
