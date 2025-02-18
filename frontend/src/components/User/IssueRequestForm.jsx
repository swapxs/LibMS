import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

const IssueRequestForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.REACT_APP_API_URL || "http://localhost:5000"}/api/requestEvents`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${user.token}`,
          },
          body: JSON.stringify({ bookID: isbn }),
        }
      );
      const data = await response.json();
      if (data.message) {
        setMessage(data.message);
      } else {
        setMessage(data.error || "Issue request failed");
      }
    } catch (err) {
      setMessage("Error submitting issue request: " + err.message);
    }
  };

  return (
    <div className="card">
      <h3>Issue Request</h3>
      {message && <p style={{ color: "red" }}>{message}</p>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
                    <label for="isbn">Enter Book's ISBN: </label>
        <input
          type="text"
        name="isbn"
          value={isbn}
          onChange={(e) => setIsbn(e.target.value)}
          required
        />
        </div>
        <button type="submit">Request Issue</button>
      </form>
    </div>
  );
};

export default IssueRequestForm;
