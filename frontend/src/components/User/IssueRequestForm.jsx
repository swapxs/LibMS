import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";

const IssueRequestForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.REACT_APP_API_URL || "http://localhost:5000"}/api/issueRequests`,
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
      setMessage("Error submitting issue request: " + err.message);
    }
  };

  return (
    <div className="container">
      <h3>Issue Request</h3>
      {message && <p style={{ color: "red" }}>{message}</p>}
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter Book ISBN"
          value={isbn}
          onChange={(e) => setIsbn(e.target.value)}
          required
        />
        <button type="submit">Request Issue</button>
      </form>
    </div>
  );
};

export default IssueRequestForm;
