import React, { useState } from "react";

function AccountDetails() {
  const storedUser = localStorage.getItem("authData");
  const user = storedUser ? JSON.parse(storedUser) : null;
const [showDetails, setShowDetails] = useState(false);

  return (
    <div className="account-details">
      <div className="account-head"  onClick={() => setShowDetails(!showDetails)} style={{ cursor: "pointer" }}>
        Account Details {showDetails ? <i class="fa-solid fa-caret-up"></i> : <i class="fa-solid fa-caret-down"></i>}
      </div>

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
