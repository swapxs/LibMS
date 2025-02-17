import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

function RemoveBookForm() {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [copies, setCopies] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!isbn || !copies) {
      setError("Please provide ISBN and number of copies to remove.");
      return;
    }
    const numCopies = Number(copies);
    if (isNaN(numCopies) || numCopies <= 0) {
      setError("Copies must be a number greater than 0.");
      return;
    }
    try {
      const response = await apiService.removeBook(isbn, numCopies, user.token);
      // Check if response.message contains "deleted" or "removed"
      if (
        (response.success && response.success === true) ||
        (response.message &&
          (response.message.toLowerCase().includes("deleted") ||
            response.message.toLowerCase().includes("removed")))
      ) {
        setMessage(response.message || "Operation successful.");
        setError("");
      } else {
        setError(response.error || "Failed to remove book copies.");
      }
    } catch (err) {
      console.error("Remove book error:", err);
      setError("An error occurred while removing book copies.");
    }
  };

  return (
    <div className="card">
      <h3>Remove Book Copies</h3>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {message && <p style={{ color: "green" }}>{message}</p>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="isbn">ISBN:</label>
          <input
            type="text"
            id="isbn"
            name="isbn"
            value={isbn}
            onChange={(e) => setIsbn(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="copies">Number of Copies to Remove:</label>
          <input
            type="number"
            id="copies"
            name="copies"
            value={copies}
            onChange={(e) => setCopies(e.target.value)}
            required
          />
        </div>
        <button type="submit">Remove Copies</button>
      </form>
    </div>
  );
}

export default RemoveBookForm;
