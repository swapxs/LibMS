// frontend/src/context/AuthContext.jsx
import React, { createContext, useState, useContext, useEffect } from "react";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  // Rehydrate from localStorage on initial load.
  const [user, setUser] = useState(() => {
    const storedUser = localStorage.getItem("authData");
    return storedUser ? JSON.parse(storedUser) : null;
  });

  const login = (data) => {
    setUser(data);
    localStorage.setItem("authData", JSON.stringify(data));
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem("authData");
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
export { AuthContext };
