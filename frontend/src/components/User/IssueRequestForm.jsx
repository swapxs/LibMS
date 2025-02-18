// frontend/src/components/User/IssueRequestForm.jsx
import React, { useState, useEffect } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

const IssueRequestForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [message, setMessage] = useState("");
  const [activeRequestCount, setActiveRequestCount] = useState(0);
  const [isLimitReached, setIsLimitReached] = useState(false);

  const fetchActiveRequestCount = async () => {
    try {
      const response = await apiService.getUserIssueInfo(user.token);
      if (response.success && Array.isArray(response.data)) {
        const active = response.data.filter(
          (issue) => issue.issue_status === "Issue" || issue.issue_status === "Approve"
        );
        setActiveRequestCount(active.length);
        setIsLimitReached(active.length >= 4);
      } else {
        setMessage(response.error || "Failed to fetch active requests");
      }
    } catch (err) {
      setMessage("Failed to fetch active requests: " + err.message);
    }
  };

  useEffect(() => {
    fetchActiveRequestCount();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isLimitReached) {
      setMessage("You have reached the maximum of 4 active requests.");
      return;
    }
    try {
      const response = await fetch(
        "http://localhost:5000/api/requestEvents",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${user.token}`,
          },
          body: JSON.stringify({ bookID: isbn })
        }
      );
      if (!response.ok) {
        const text = await response.text();
        let errorMessage = "An unknown error occurred";
        try {
          const json = JSON.parse(text);
          errorMessage = json.message || json.error || errorMessage;
        } catch (parseError) {
          errorMessage = text;
        }
        throw new Error(errorMessage);
      }
      let data;
      try {
        data = await response.json();
      } catch (jsonError) {
        throw new Error("Failed to parse server response");
      }
      setMessage(data.message || data.error || "Request processed successfully");
      fetchActiveRequestCount();
    } catch (err) {
      setMessage(`Error submitting issue request: ${err.message}`);
    }
  };

  return (
    <div className="card">
      <h3>Issue Request</h3>
      {message && <p style={{ color: "red" }}>{message}</p>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="isbn">Enter Book's ISBN:</label>
          <input
            type="text"
            name="isbn"
            value={isbn}
            onChange={(e) => setIsbn(e.target.value)}
            required
            disabled={isLimitReached}
          />
        </div>
        <button type="submit" disabled={isLimitReached}>Request Issue</button>
      </form>
    </div>
  );
};

export default IssueRequestForm;
