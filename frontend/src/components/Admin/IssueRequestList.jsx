// frontend/src/components/Admin/IssueRequestList.jsx
import React, { useState, useEffect } from "react";
import apiService from "../../services/apiService";
import { useAuth } from "../../context/AuthContext";

const IssueRequestList = () => {
  const { user } = useAuth();
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Helper: Calculate Due Date as 7 days after ApprovalDate.
  const calculateDueDate = (approvalDate) => {
    if (!approvalDate || approvalDate === "Rejected") return "Rejected";
    const date = new Date(approvalDate);
    date.setDate(date.getDate() + 7);
    return date.toLocaleString();
  };

  // Fetch issue requests from the backend.
  const fetchRequests = async () => {
    try {
      const response = await apiService.getIssueRequests(user.token);
      console.log("Issue requests response:", response);
      let fetched = [];
      if (response.requests && Array.isArray(response.requests)) {
        fetched = response.requests;
      } else if (response.success && Array.isArray(response.data)) {
        fetched = response.data;
      } else if (Array.isArray(response)) {
        fetched = response;
      } else {
        setError("Invalid response format from getIssueRequests");
      }
      setRequests(fetched);
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

  // Handler for Issue Action dropdown changes.
  const handleIssueActionChange = async (reqId, action) => {
    if (!action) return;
    const lowerAction = action.toLowerCase();
    if (lowerAction === "approve") {
      // For simplicity, simulate an approval action.
      const now = new Date().toISOString();
      try {
        // Example API call (replace with your real API call):
        const updatePayload = {
          request_type: "Approve",
          // We'll assume the API sets an ApprovalDate.
          expected_return_date: new Date(new Date().getTime() + 7 * 24 * 60 * 60 * 1000).toISOString(),
        };
        const updateResponse = await apiService.updateIssueRequest(reqId, updatePayload, user.token);
        console.log("Update issue request response:", updateResponse);
        if (
          updateResponse.success ||
          (updateResponse.message &&
            updateResponse.message.toLowerCase().includes("updated"))
        ) {
          setRequests((prev) =>
            prev.map((r) =>
              r.ReqID === reqId
                ? { ...r, RequestType: "Approve", ApprovalDate: now }
                : r
            )
          );
          setError("");
        } else {
          setError(updateResponse.error || "Failed to update issue request status");
        }
      } catch (err) {
        console.error("Error processing approval:", err);
        setError("Error processing approval");
      }
    } else if (lowerAction === "reject") {
      // When rejecting, update Issue Date, Due Date, and Return Date to "Rejected"
      const updatePayload = { request_type: "Reject", expected_return_date: null };
      try {
        const updateResponse = await apiService.updateIssueRequest(reqId, updatePayload, user.token);
        console.log("Update issue request response:", updateResponse);
        if (
          updateResponse.success ||
          (updateResponse.message &&
            updateResponse.message.toLowerCase().includes("updated"))
        ) {
          setRequests((prev) =>
            prev.map((r) =>
              r.ReqID === reqId
                ? {
                    ...r,
                    RequestType: "Reject",
                    RequestDate: "Rejected",
                    ApprovalDate: "Rejected",
                    ReturnDate: "Rejected",
                  }
                : r
            )
          );
          setError("");
        } else {
          setError(updateResponse.error || "Failed to update issue request status");
        }
      } catch (err) {
        console.error("Error updating issue request status:", err);
        setError("An error occurred while updating issue request status");
      }
    }
  };

  return (
    <div className="card">
      <h2>Issue Requests</h2>
      <button onClick={fetchRequests}>Refresh</button>
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
              <th>ISBN</th>
              <th>Book Name</th>
              <th>Reader Name</th>
              <th>Issue Date</th>
              <th>Due Date</th>
              <th>Return Date</th>
              <th>Issuing Approval Admin</th>
              <th>Issue Action</th>
            </tr>
          </thead>
          <tbody>
            {requests.map((req) => {
              const isbn = req.BookID || "N/A";
              const bookName = req.BookName || "N/A";
              const readerName = req.ReaderName || "N/A";
              const issueDate =
                req.RequestDate && req.RequestDate !== "Rejected"
                  ? new Date(req.RequestDate).toLocaleString()
                  : req.RequestDate || "N/A";
              const dueDate =
                req.ApprovalDate && req.ApprovalDate !== "Rejected"
                  ? calculateDueDate(req.ApprovalDate)
                  : req.ApprovalDate || "N/A";
              const returnDate =
                req.ReturnDate && req.ReturnDate !== "Rejected"
                  ? new Date(req.ReturnDate).toLocaleString()
                  : req.ReturnDate || "N/A";
              const issueApprovalAdmin = req.IssueApproverEmail || "N/A";
              const currentIssueAction = req.RequestType || "Pending";

              return (
                <tr key={req.ReqID}>
                  <td>{isbn}</td>
                  <td>{bookName}</td>
                  <td>{readerName}</td>
                  <td>{currentIssueAction === "Reject" ? "Rejected" : issueDate}</td>
                  <td>{currentIssueAction === "Reject" ? "Rejected" : dueDate}</td>
                  <td>{currentIssueAction === "Reject" ? "Rejected" : returnDate}</td>
                  <td>{issueApprovalAdmin}</td>
                  <td>
                    <select
                      value={currentIssueAction}
                      onChange={(e) =>
                        handleIssueActionChange(req.ReqID, e.target.value)
                      }
                    >
                      <option value="Pending">Pending</option>
                      <option value="Approve">Approve</option>
                      <option value="Reject">Reject</option>
                    </select>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default IssueRequestList;
