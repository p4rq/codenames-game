import React, { createContext, useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    // Check local storage for user data
    const storedUser = localStorage.getItem('codenames_user');
    
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    } else {
      // Create a new user if none exists
      const newUser = {
        id: uuidv4(),
        username: `Player_${Math.floor(Math.random() * 1000)}`,
      };
      
      localStorage.setItem('codenames_user', JSON.stringify(newUser));
      setUser(newUser);
    }
  }, []);

  const updateUsername = (newUsername) => {
    if (!newUsername || newUsername.trim() === '') return;
    
    const updatedUser = {
      ...user,
      username: newUsername,
    };
    
    localStorage.setItem('codenames_user', JSON.stringify(updatedUser));
    setUser(updatedUser);
  };

  return (
    <UserContext.Provider value={{ user, updateUsername }}>
      {children}
    </UserContext.Provider>
  );
};