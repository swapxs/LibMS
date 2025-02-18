import React, { useState, useEffect } from "react";
import apiService from "../../services/apiService";
import { useAuth } from "../../context/AuthContext";

const IssuedBooks = () => {
  const { user } = useAuth();
  const [issuedBooks, setIssuedBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchIssuedBooks = async () => {
    try {
      const response = await apiService.getIssuedBooks(user.token);
      console.log("Issued books response:", response);
      let fetchedBooks = [];
      if (response.success && Array.isArray(response.data)) {
        fetchedBooks = response.data;
      } else if (Array.isArray(response)) {
        fetchedBooks = response;
      } else {
        setError("Invalid response format from getIssuedBooks");
        fetchedBooks = [];
      }
      setIssuedBooks(fetchedBooks);
      setError("");
    } catch (err) {
      console.error("Error fetching issued books:", err);
      setError("Error fetching issued books");
      setIssuedBooks([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user && user.token) {
      fetchIssuedBooks();
    }
  }, [user]);

  return (
    <div className="card">
      <h2>Issued Books</h2>
      {loading ? (
        <p>Loading issued books...</p>
      ) : error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : issuedBooks.length === 0 ? (
        <p>No issued books found.</p>
      ) : (
        <table className="books-table">
          <thead>
            <tr>
              <th>Issue ID</th>
              <th>ISBN</th>
              <th>Reader ID</th>
              <th>Issue Date</th>
              <th>Expected Return Date</th>
              <th>Return Date</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {issuedBooks.map((issue) => (
              <tr key={issue.issue_id}>
                <td>{issue.issue_id}</td>
                <td>{issue.isbn}</td>
                <td>{issue.reader_id}</td>
                <td>{new Date(issue.issue_date).toLocaleString()}</td>
                <td>{new Date(issue.expected_return_date).toLocaleString()}</td>
                <td>
                  {issue.return_date
                    ? new Date(issue.return_date).toLocaleString()
                    : "Not Returned"}
                </td>
                <td>{issue.issue_status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      <button onClick={fetchIssuedBooks}>Refresh</button>
    </div>
  );
};

export default IssuedBooks;
