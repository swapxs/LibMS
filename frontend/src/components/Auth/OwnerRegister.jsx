import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";

const RegisterOwner = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    contact_number: "",
    library_name: "",
  });
  const [error, setError] = useState("");

  const handleChange = (e) =>
    setFormData({ ...formData, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.REACT_APP_API_URL || "http://localhost:5000"}/api/owner/registration`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(formData),
        }
      );
      const result = await response.json();
      if (result.message) {
        navigate("/login");
      } else {
        setError(result.error || "Registration failed");
      }
    } catch (err) {
      setError("An error occurred during registration");
    }
  };

  return (
    <div className="container">
      <h2>Owner Registration</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="name"
          placeholder="Full Name"
          value={formData.name}
          onChange={handleChange}
          required
        />
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={formData.email}
          onChange={handleChange}
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          value={formData.password}
          onChange={handleChange}
          required
        />
        <input
          type="text"
          name="contact_number"
          placeholder="Contact Number"
          value={formData.contact_number}
          onChange={handleChange}
          required
        />
        <input
          type="text"
          name="library_name"
          placeholder="Library Name"
          value={formData.library_name}
          onChange={handleChange}
          required
        />
        <button type="submit">Register as Owner</button>
      </form>
      <p>
        Already have an account? <Link to="/login">Login here</Link>
      </p>
    </div>
  );
};

export default RegisterOwner;
