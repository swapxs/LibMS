// frontend/src/components/Dashboard.jsx
import React from "react";
import { useAuth } from "../context/AuthContext";

function Dashboard() {
  const { user } = useAuth();

  if (!user) {
    return <div>Please log in to view your dashboard.</div>;
  }

  return (
    <div className="card">
      <h2>Hello {user.role}!</h2>
      <p>Welcome to the online library.</p>
      <p>You can do a lot of things in this portal. Here are the list of things you can do.</p>
    </div>
  );
}

export default Dashboard;

