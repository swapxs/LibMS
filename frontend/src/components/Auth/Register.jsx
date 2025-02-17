import React, { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import apiService from "../../services/apiService";

function Register() {
  const navigate = useNavigate();
  const [libraries, setLibraries] = useState([]);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
    contact_number: "",
    library_id: "",
  });
  const [error, setError] = useState("");

  // Fetch libraries on component mount.
  useEffect(() => {
    async function fetchLibraries() {
      try {
        const response = await apiService.getLibraries();
        // Check for different response formats.
        if (response.success && Array.isArray(response.data)) {
          setLibraries(response.data);
        } else if (Array.isArray(response.libraries)) {
          setLibraries(response.libraries);
        } else if (Array.isArray(response)) {
          setLibraries(response);
        } else {
          setLibraries([]);
        }
      } catch (err) {
        console.error("Error fetching libraries:", err);
        setLibraries([]);
      }
    }
    fetchLibraries();
  }, []);

  // Handle input field changes.
  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  // Validate and submit registration form.
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (formData.password !== formData.confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      const payload = {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        contact_number: formData.contact_number,
        // Convert library_id to number.
        library_id: Number(formData.library_id),
      };
      const response = await apiService.register(payload);
      // Check if the response message indicates success.
      if (
        response.message &&
        response.message.toLowerCase().includes("successful")
      ) {
        navigate("/login");
      } else {
        setError(response.error || "Registration failed");
      }
    } catch (err) {
      console.error("Registration error:", err);
      setError("An error occurred during registration");
    }
  };

  return (
    <div className="container">
      <h2>Register</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Full Name:</label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </div>
        
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
          <label htmlFor="contact_number">Contact Number:</label>
          <input
            type="text"
            id="contact_number"
            name="contact_number"
            value={formData.contact_number}
            onChange={handleChange}
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="library_id">Select Library:</label>
          <select
            id="library_id"
            name="library_id"
            value={formData.library_id}
            onChange={handleChange}
            required
          >
            <option value="">-- Select a Library --</option>
            {libraries.map((lib) => (
              <option key={lib.ID} value={lib.ID}>
                {lib.Name}
              </option>
            ))}
          </select>
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
        
        <div className="form-group">
          <label htmlFor="confirmPassword">Confirm Password:</label>
          <input
            type="password"
            id="confirmPassword"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
            required
          />
        </div>
        
        <button type="submit">Register</button>
      </form>
      <p>
        Already have an account? <Link to="/login">Login here</Link>
      </p>
    </div>
  );
}

export default Register;
