import React, { createContext, useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    // Try to load user from localStorage on initial render
    try {
      const savedUser = localStorage.getItem('user');
      if (savedUser) {
        return JSON.parse(savedUser);
      }
    } catch (error) {
      console.error('Error loading user from localStorage:', error);
    }
    
    // Default user with random ID
    return {
      id: `user-${uuidv4().substring(0, 8)}`,
      username: 'Guest',
    };
  });

  // Save user to localStorage whenever it changes
  useEffect(() => {
    try {
      localStorage.setItem('user', JSON.stringify(user));
      console.log("UserContext - User saved to localStorage:", user);
    } catch (error) {
      console.error('Error saving user to localStorage:', error);
    }
  }, [user]);

  // Function to update user information
  const updateUser = (updatedData) => {
    setUser(prevUser => {
      const newUser = { ...prevUser, ...updatedData };
      console.log("UserContext - Updating user:", prevUser, "→", newUser);
      return newUser;
    });
  };

  // Update username specifically
  const updateUsername = (username) => {
    updateUser({ username });
  };

  return (
    <UserContext.Provider value={{ 
      user, 
      setUser, 
      updateUser, 
      updateUsername 
    }}>
      {children}
    </UserContext.Provider>
  );
};