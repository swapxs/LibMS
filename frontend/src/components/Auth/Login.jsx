// frontend/src/components/Auth/Login.jsx
import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import apiService from "../../services/apiService";
import { useAuth } from "../../context/AuthContext";

function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({ email: "", password: "" });
  const [error, setError] = useState("");

  const handleChange = (e) =>
    setFormData({ ...formData, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      const response = await apiService.login(formData.email, formData.password);
      if (response.token) {
        // If the response does not include a name, derive one from the email.
        if (!response.name && response.email) {
          response.name = response.email.split("@")[0];
        }
        // Save the full response in AuthContext.
        login(response);
        // Navigate to the unified dashboard.
        navigate("/dashboard");
      } else {
        setError(response.error || "Login failed. Please check your credentials.");
      }
    } catch (err) {
      console.error("Login error:", err);
      setError("An error occurred during login.");
    }
  };

  return (
        <div className="container">
    <div className="card">
      <h1>Login</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email:</label>
          <input 
            type="email" 
            id="email" 
            name="email"
            value={formData.email} 
            onChange={handleChange} 
            required 
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password:</label>
          <input 
            type="password" 
            id="password" 
            name="password"
            value={formData.password} 
            onChange={handleChange} 
            required 
          />
        </div>
        <button type="submit">Login</button>
      </form>
      <p className="bottom">
        Don't have an account? <Link to="/register">Register here</Link>
      </p>
    </div>
</div>
  );
}

export default Login;
