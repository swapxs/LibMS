import React, { useState, useEffect } from 'react';
import apiService from '../../services/apiService';
import { useAuth } from '../../context/AuthContext';

function AssignAdmin() {
  const { user } = useAuth();
  const [users, setUsers] = useState([]);
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');

  // Fetch users (excluding owners) from the backend.
  const fetchUsers = async () => {
    try {
      const response = await apiService.getUsers(user.token);
      console.log("Fetched users response:", response);
      let allUsers = [];
      if (response.success && Array.isArray(response.data)) {
        allUsers = response.data;
      } else if (Array.isArray(response.users)) {
        allUsers = response.users;
      } else if (Array.isArray(response)) {
        allUsers = response;
      } else {
        setError("No users found or invalid response format");
        return;
      }
      // Exclude users with Role "Owner"
      const nonOwners = allUsers.filter(u => u.Role !== 'Owner');
      // Optionally, sort non-owners with admin users first.
      nonOwners.sort((a, b) => {
        if (a.Role === 'LibraryAdmin' && b.Role !== 'LibraryAdmin') return -1;
        if (a.Role !== 'LibraryAdmin' && b.Role === 'LibraryAdmin') return 1;
        return 0;
      });
      setUsers(nonOwners);
      setFilteredUsers(nonOwners);
      setError('');
    } catch (err) {
      console.error("Error fetching users:", err);
      setUsers([]);
      setFilteredUsers([]);
      setError("Error fetching users");
    }
  };

  useEffect(() => {
    if (user && user.token) {
      fetchUsers();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user]);

  // Filter users by search query.
  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredUsers(users);
    } else {
      const query = searchQuery.toLowerCase();
      setFilteredUsers(
        users.filter(u =>
          u.Name.toLowerCase().includes(query) ||
          u.Email.toLowerCase().includes(query) ||
          (u.ContactNumber && u.ContactNumber.toLowerCase().includes(query))
        )
      );
    }
  }, [searchQuery, users]);

  // Handler to assign admin rights.
  const handleAssignAdmin = async (u) => {
    const confirmAssign = window.confirm(
      `Are you sure you want to assign admin rights to ${u.Name} (${u.Email})?`
    );
    if (!confirmAssign) return;
    try {
      const response = await apiService.assignAdmin({ email: u.Email }, user.token);
      console.log("Assign admin response:", response);
      // Check the response message instead of a "success" field.
      if (response.message && response.message.toLowerCase().includes("promoted")) {
        setMessage(`Successfully promoted ${u.Name} to admin.`);
        // Update local state.
        const updated = users.map(item =>
          item.Email === u.Email ? { ...item, Role: 'LibraryAdmin' } : item
        );
        setUsers(updated);
        setFilteredUsers(updated);
      } else {
        setError(`Failed to assign admin rights: ${response.error || ""}`);
      }
    } catch (err) {
      console.error(err);
      setError("An error occurred while assigning admin rights.");
    }
  };

  // Handler to revoke admin rights.
  const handleRevokeAdmin = async (u) => {
    const confirmRevoke = window.confirm(
      `Are you sure you want to revoke admin rights from ${u.Name} (${u.Email})?`
    );
    if (!confirmRevoke) return;
    try {
      const response = await apiService.revokeAdmin({ email: u.Email }, user.token);
      console.log("Revoke admin response:", response);
      // Check response message for revoke confirmation.
      if (response.message && response.message.toLowerCase().includes("revoked")) {
        setMessage(`Successfully revoked admin rights from ${u.Name}.`);
        const updated = users.map(item =>
          item.Email === u.Email ? { ...item, Role: 'Reader' } : item
        );
        setUsers(updated);
        setFilteredUsers(updated);
      } else {
        setError(`Failed to revoke admin rights: ${response.error || ""}`);
      }
    } catch (err) {
      console.error(err);
      setError("An error occurred while revoking admin rights.");
    }
  };

  return (
    <div className="card">
      <h2>Assign/Revoke Admin Rights</h2>
      <div className="form-group">
        <label htmlFor="search">Search Users (by Name or Email):</label>
        <input
          id="search"
          type="text"
          placeholder="Search name or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
      {error && <div className="error">{error}</div>}
      {message && <div className="success">{message}</div>}
      <button onClick={fetchUsers} style={{ marginBottom: '1rem' }}>
        Refresh
      </button>
      {filteredUsers.length === 0 ? (
        <p>No users to display.</p>
      ) : (
        <div className="cards-container">
          {filteredUsers.map(u => (
            <div key={u.ID || u.Email} className="card profile-card">
              <div className="profile-card-content">
                <div className="profile-info">
                  <h3 className="profile-name">
                    {u.Role === 'LibraryAdmin' ? (
                      <i className="fa-solid fa-user-secret"></i>
                    ) : (
                      <i className="fa-solid fa-user"></i>
                    )} {u.Name}
                  </h3>
                  <p><i class="fa-solid fa-envelope"></i> {u.Email}</p>
                  <p><i class="fa-solid fa-phone"></i> {u.ContactNumber}</p>
                  {/* <p>{u.Role}</p> */}
                  {/* {u.LibraryID && <p>Library: {u.LibraryID}</p>} */}
                </div>
                <div className="profile-actions">
                  {u.Role !== 'LibraryAdmin' ? (
                    <button className="assign" onClick={() => handleAssignAdmin(u)}>
                      <span>Promote</span>
                    </button>
                  ) : (
                    <button className="revoke" onClick={() => handleRevokeAdmin(u)}>
                      <span>Revoke</span>
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default AssignAdmin;
