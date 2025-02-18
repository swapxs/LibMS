import React, { useState, useEffect } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

function UserBorrowHistory() {
  const { user } = useAuth();
  const [borrowHistory, setBorrowHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Fetch borrowing history for the logged-in user.
  useEffect(() => {
    async function fetchHistory() {
      try {
        // Assuming your API returns borrowings for the current user.
        // Adjust the function name and token usage as needed.
        const response = await apiService.getUserIssueInfo(user.token);
        console.log("User borrow history response:", response);
        let history = [];
        if (response.success && Array.isArray(response.data)) {
          history = response.data;
        } else if (Array.isArray(response)) {
          history = response;
        } else {
          setError("Invalid response format for borrow history");
        }
        setBorrowHistory(history);
      } catch (err) {
        console.error("Error fetching borrow history:", err);
        setError("Error fetching borrow history");
      } finally {
        setLoading(false);
      }
    }
    if (user && user.token) {
      fetchHistory();
    }
  }, [user]);

  // Helper function to calculate "Days Left" as (DueDate - IssueDate) in days.
  const calculateDaysLeft = (issueDate, dueDate) => {
    if (!issueDate || !dueDate) return "N/A";
    const issue = new Date(issueDate);
    const due = new Date(dueDate);
    const diffTime = due - issue;
    const diffDays = Math.round(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  return (
    <div className="card">
      <h2>Your Borrowing History</h2>
      {loading ? (
        <p>Loading your borrowing history...</p>
      ) : error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : borrowHistory.length === 0 ? (
        <p>No borrowing records found.</p>
      ) : (
        <table className="books-table">
          <thead>
            <tr>
              <th>Book ISBN</th>
              <th>Book Name</th>
              <th>Issue Date</th>
              <th>Due Date</th>
              <th>Days Left</th>
              <th>Issue Approver's Email</th>
            </tr>
          </thead>
          <tbody>
            {borrowHistory.map((record, index) => {
              // Adjust keys according to your backend's response.
              const isbn = record.BookISBN || record.isbn || "N/A";
              const bookName = record.BookName || record.bookName || "N/A";
              const issueDateRaw = record.IssueDate || record.issueDate;
              const dueDateRaw = record.DueDate || record.dueDate;
              const issueDate = issueDateRaw ? new Date(issueDateRaw).toLocaleString() : "N/A";
              const dueDate = dueDateRaw ? new Date(dueDateRaw).toLocaleString() : "N/A";
              const daysLeft = (issueDateRaw && dueDateRaw)
                ? calculateDaysLeft(issueDateRaw, dueDateRaw)
                : "N/A";
              const approverEmail = record.IssueApproverEmail || record.issueApproverEmail || "N/A";

              return (
                <tr key={index}>
                  <td>{isbn}</td>
                  <td>{bookName}</td>
                  <td>{issueDate}</td>
                  <td>{dueDate}</td>
                  <td>{daysLeft}</td>
                  <td>{approverEmail}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default UserBorrowHistory;
