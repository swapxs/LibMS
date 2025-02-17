// frontend/src/components/Admin/IssueRequestList.jsx
import React, { useState, useEffect } from 'react';
import apiService from '../../services/apiService';

function IssueRequestList() {
  const [issues, setIssues] = useState([]);
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);

  // Fetch both issue requests and books when the component mounts.
  useEffect(() => {
    async function fetchData() {
      try {
        const issuesResponse = await apiService.getIssueRequests();
        const booksResponse = await apiService.getBooks();
        console.log("Issues response:", issuesResponse);
        console.log("Books response:", booksResponse);
        if (issuesResponse.success && Array.isArray(issuesResponse.data)) {
          setIssues(issuesResponse.data);
        } else {
          setIssues([]);
        }
        if (booksResponse.success && Array.isArray(booksResponse.data)) {
          setBooks(booksResponse.data);
        } else {
          setBooks([]);
        }
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  // Merge each issue request with its corresponding book details.
  const mergedData = issues.map(req => {
    const reqBookID = String(req.bookID);
    const matchedBook = books.find(
      book => String(book.id) === reqBookID || String(book.ID) === reqBookID
    );
    return { ...req, bookDetails: matchedBook || {} };
  });

  // Filter duplicates based on a unique key: combination of bookID and userID.
  const uniqueData = [];
  const seenKeys = new Set();
  mergedData.forEach((req) => {
    const key = String(req.bookID) + "-" + String(req.userID || req.UserID);
    if (!seenKeys.has(key)) {
      seenKeys.add(key);
      uniqueData.push(req);
    }
  });

  // Update the status of an issue request when the dropdown value changes.
  const handleStatusChange = async (reqId, newStatus) => {
    const response = await apiService.updateIssueRequest(reqId, { status: newStatus });
    if (response.success) {
      setIssues(prevIssues =>
        prevIssues.map(req => {
          const currentId = String(req.id || req.ID);
          if (currentId === String(reqId)) {
            return { ...req, status: newStatus };
          }
          return req;
        })
      );
    } else {
      alert("Failed to update status");
    }
  };

  return (
    <div className="card">
      <h2>Issue Requests</h2>
      {loading ? (
        <p>Loading issue requests...</p>
      ) : uniqueData.length === 0 ? (
        <p>No issue requests found.</p>
      ) : (
        <table className="books-table">
          <thead>
            <tr>
              <th>Book ID</th>
              <th>User ID</th>
              <th>Book Title</th>
              <th>Language</th>
              <th>Copies Left</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {uniqueData.map(req => {
              const reqId = req.id || req.ID;
              const status = req.status || req.Status || 'pending';
              return (
                <tr key={reqId + "-" + req.bookID + "-" + req.userID}>
                  <td>{req.bookID}</td>
                  <td>{req.userID || req.UserID}</td>
                  <td>{req.bookDetails && req.bookDetails.title ? req.bookDetails.title : "N/A"}</td>
                  <td>{req.bookDetails && req.bookDetails.language ? req.bookDetails.language : "N/A"}</td>
                  <td>
                    {typeof req.bookDetails.availableCopies !== 'undefined'
                      ? req.bookDetails.availableCopies
                      : "N/A"}
                  </td>
                  <td>
                    <select
                      value={status}
                      onChange={(e) => handleStatusChange(reqId, e.target.value)}
                    >
                      <option value="pending">Pending</option>
                      <option value="approved">Approved</option>
                      <option value="rejected">Rejected</option>
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
}

export default IssueRequestList;

