import React, { useState, useEffect } from 'react';
import apiService from '../../services/apiService';
import { useAuth } from '../../context/AuthContext';

function IssueRequestList() {
  const { user } = useAuth();
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Fetch issue requests from backend.
  const fetchRequests = async () => {
    try {
      const response = await apiService.getIssueRequests(user.token);
      console.log("Issue requests response:", response);
      let fetchedRequests = [];
      if (response.requests && Array.isArray(response.requests)) {
        fetchedRequests = response.requests;
      } else if (response.success && Array.isArray(response.data)) {
        fetchedRequests = response.data;
      } else if (Array.isArray(response)) {
        fetchedRequests = response;
      } else {
        setError("Invalid response format from getIssueRequests");
        fetchedRequests = [];
      }
      setRequests(fetchedRequests);
      setError("");
    } catch (err) {
      console.error("Error fetching issue requests:", err);
      setError("Error fetching issue requests");
      setRequests([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user && user.token) {
      fetchRequests();
    }
  }, [user]);

  // Handler for changing the status via the dropdown.
  const handleStatusChange = async (reqId, newStatus) => {
    let updatePayload = { request_type: newStatus };
    if (newStatus === "Approve") {
      // Calculate expected return date 7 days from now.
      const now = new Date();
      const expectedReturnDate = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);
      updatePayload.expectedReturnDate = expectedReturnDate.toISOString();
    } else {
      // For non-approved statuses, clear expectedReturnDate.
      updatePayload.expectedReturnDate = null;
    }
    try {
      const response = await apiService.updateIssueRequest(reqId, updatePayload, user.token);
      console.log("Update issue request response:", response);
      // Check if response.message indicates success.
      if (
        response.message &&
        response.message.toLowerCase().includes("updated")
      ) {
        // Update local state: update the request's status and expectedReturnDate.
        setRequests(prev =>
          prev.map(r =>
            r.ReqID === reqId
              ? { ...r, RequestType: newStatus, expectedReturnDate: updatePayload.expectedReturnDate }
              : r
          )
        );
        setError("");
      } else {
        setError(response.error || "Failed to update issue request status");
      }
    } catch (err) {
      console.error("Error updating issue request status:", err);
      setError("An error occurred while updating issue request status");
    }
  };

  // Calculate expected return date from ApprovalDate (7 days after approval).
  const calculateExpectedReturnDate = (approvalDate) => {
    if (!approvalDate) return "Invalid Date";
    const date = new Date(approvalDate);
    date.setDate(date.getDate() + 7);
    return date.toLocaleString();
  };

  return (
    <div className="card">
      <h2>Issue Requests</h2>
      {loading ? (
        <p>Loading issue requests...</p>
      ) : error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : requests.length === 0 ? (
        <p>No issue requests found.</p>
      ) : (
        <table className="books-table">
          <thead>
            <tr>
              <th>Issue ID</th>
              <th>Book ID</th>
              <th>Reader ID</th>
              <th>Request Date</th>
              <th>Status</th>
              <th>Expected Return Date</th>
            </tr>
          </thead>
          <tbody>
            {requests.map((req) => {
              const reqId = req.ReqID;
              const bookId = req.BookID;
              const readerId = req.ReaderID;
              const requestDate = req.RequestDate;
              const status = req.RequestType || "Pending";
              return (
                <tr key={reqId}>
                  <td>{reqId}</td>
                  <td>{bookId}</td>
                  <td>{readerId}</td>
                  <td>{new Date(requestDate).toLocaleString()}</td>
                  <td>
                    <select
                      value={status}
                      onChange={(e) => handleStatusChange(reqId, e.target.value)}
                    >
                      <option value="Pending">Pending</option>
                      <option value="Reject">Reject</option>
                      <option value="Approve">Approve</option>
                    </select>
                  </td>
                  <td>
                    {status === "Approve" && req.ApprovalDate
                      ? calculateExpectedReturnDate(req.ApprovalDate)
                      : "Invalid Date"}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
      <button onClick={fetchRequests}>Refresh</button>
    </div>
  );
}

export default IssueRequestList;
