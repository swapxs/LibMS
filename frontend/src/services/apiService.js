const API_BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:5000/api";

const apiService = {
  // Login a user
  login: async (email, password) => {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    return response.json();
  },

  // Register a new user
  register: async (userData) => {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userData),
    });
    return response.json();
  },

  // Get the list of libraries
  getLibraries: async () => {
    const response = await fetch(`${API_BASE_URL}/libraries`);
    return response.json();
  },

  // Add a new book or increment copies
  addBook: async (bookData, token) => {
    const response = await fetch(`${API_BASE_URL}/books`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(bookData),
    });
    return response.json();
  },

  // Remove copies of a book (or delete the book if copies become 0)
  removeBook: async (isbn, copies, token) => {
    const response = await fetch(`${API_BASE_URL}/books/remove`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ isbn, copies }),
    });
    return response.json();
  },

  // Update book details
  updateBook: async (isbn, bookData, token) => {
    const response = await fetch(`${API_BASE_URL}/books/${isbn}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(bookData),
    });
    return response.json();
  },

  // Assign admin rights to a user
  assignAdmin: async (data, token) => {
    const response = await fetch(`${API_BASE_URL}/owner/assign-admin`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Revoke admin rights from a user
  revokeAdmin: async (data, token) => {
    const response = await fetch(`${API_BASE_URL}/owner/revoke-admin`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Get the list of users
  getUsers: async (token) => {
    const response = await fetch(`${API_BASE_URL}/users`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.json();
  },

  // Get the list of issue requests
  getIssueRequests: async (token) => {
    const response = await fetch(`${API_BASE_URL}/issueRequests`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.json();
  },

  // Update the status of an issue request
  updateIssueRequest: async (reqId, data, token) => {
    const response = await fetch(`${API_BASE_URL}/issueRequests/${reqId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    return response.json();
  },

  // Create an issue request (or register an issued book)
  createIssueRequest: async (issueData, token) => {
    const response = await fetch(`${API_BASE_URL}/issueRegistry`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(issueData),
    });
    return response.json();
  },

  // Get the list of books
  getBooks: async (token) => {
    const response = await fetch(`${API_BASE_URL}/books`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.json();
  }
};

export default apiService;
