import React, { createContext, useState, useEffect } from 'react';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Load user from localStorage on mount
  useEffect(() => {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      try {
        setUser(JSON.parse(savedUser));
      } catch (e) {
        console.error('Failed to parse saved user', e);
        localStorage.removeItem('user');
      }
    }
    setLoading(false);
  }, []);

  // Create a user if one doesn't exist
  const createUser = (username) => {
    const newUser = {
      id: `user-${Math.random().toString(36).substring(2, 10)}`,
      username,
      createdAt: new Date().toISOString()
    };
    setUser(newUser);
    localStorage.setItem('user', JSON.stringify(newUser));
    return newUser;
  };

  // Update username only
  const updateUsername = (username) => {
    if (!user) {
      return createUser(username);
    }
    
    const updatedUser = { ...user, username };
    setUser(updatedUser);
    localStorage.setItem('user', JSON.stringify(updatedUser));
    return updatedUser;
  };

  // Update any user properties
  const updateUser = (updates) => {
    if (!user) return null;
    
    const updatedUser = { ...user, ...updates };
    setUser(updatedUser);
    localStorage.setItem('user', JSON.stringify(updatedUser));
    return updatedUser;
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('user');
  };

  return (
    <UserContext.Provider
      value={{
        user,
        loading,
        createUser,
        updateUsername,
        updateUser,
        logout
      }}
    >
      {children}
    </UserContext.Provider>
  );
};