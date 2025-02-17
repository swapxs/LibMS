import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import apiService from "../../services/apiService";

const UpdateBookForm = () => {
  const { user } = useAuth();
  const [isbn, setIsbn] = useState("");
  const [formData, setFormData] = useState({
    title: "",
    authors: "",
    publisher: "",
    version: "",
    total_copies: "",
  });
  const [message, setMessage] = useState("");

  const handleChange = (e) =>
    setFormData({ ...formData, [e.target.name]: e.target.value });

  const handleUpdate = async (e) => {
    e.preventDefault();
    if (!isbn) {
      setMessage("Please provide the ISBN of the book to update.");
      return;
    }
    try {
      const payload = {};
      Object.keys(formData).forEach((key) => {
        if (formData[key] !== "") {
          payload[key] = key === "total_copies" ? Number(formData[key]) : formData[key];
        }
      });
      const result = await apiService.updateBook(isbn, payload, user.token);
      setMessage(result.message || result.error || "Update failed");
    } catch (err) {
      setMessage("An error occurred.");
    }
  };

  return (
    <div className="container">
      <h3>Update Book Details</h3>
      {message && <p>{message}</p>}
      <form onSubmit={handleUpdate}>
        <input
          type="text"
          placeholder="ISBN of Book"
          value={isbn}
          onChange={(e) => setIsbn(e.target.value)}
          required
        />
        <input
          type="text"
          name="title"
          placeholder="New Title"
          value={formData.title}
          onChange={handleChange}
        />
        <input
          type="text"
          name="authors"
          placeholder="New Authors"
          value={formData.authors}
          onChange={handleChange}
        />
        <input
          type="text"
          name="publisher"
          placeholder="New Publisher"
          value={formData.publisher}
          onChange={handleChange}
        />
        <input
          type="text"
          name="version"
          placeholder="New Version"
          value={formData.version}
          onChange={handleChange}
        />
        <input
          type="number"
          name="total_copies"
          placeholder="New Total Copies"
          value={formData.total_copies}
          onChange={handleChange}
        />
        <button type="submit">Update Book</button>
      </form>
    </div>
  );
};

export default UpdateBookForm;
