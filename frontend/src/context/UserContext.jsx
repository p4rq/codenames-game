import React, { createContext, useState, useEffect } from 'react';

// Create context
export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  // Initialize state with data from localStorage if available
  const [user, setUser] = useState(() => {
    const savedUser = localStorage.getItem('codenames_user');
    if (savedUser) {
      return JSON.parse(savedUser);
    }
    return { 
      id: `user-${Math.random().toString(36).substring(2, 9)}`,
      username: '' 
    };
  });

  // Save user to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('codenames_user', JSON.stringify(user));
  }, [user]);

  // Update username
  const updateUsername = (newUsername) => {
    setUser(prevUser => ({
      ...prevUser,
      username: newUsername
    }));
  };

  return (
    <UserContext.Provider value={{ user, updateUsername }}>
      {children}
    </UserContext.Provider>
  );
};