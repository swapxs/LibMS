import React, { useState } from "react";
import { useAuth } from "../context/AuthContext";

function AccountDetails() {
  const storedUser = localStorage.getItem("authData");
  const user = storedUser ? JSON.parse(storedUser) : null;
const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="account-details">
      <h4 onClick={() => setShowDetails(!showDetails)} style={{ cursor: "pointer" }}>
        Account Details {showDetails ? "▲" : "▼"}
      </h4>

      {showDetails && user && (
        <div className="details-content">
          <p><strong>Name:</strong> {user.name}</p>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Contact Number:</strong> {user.contact_number}</p>
          <p><strong>Library:</strong> {user.library_name}</p>
        </div>
      )}
    </div>
  );
}

export default AccountDetails;
