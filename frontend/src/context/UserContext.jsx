import React, { createContext, useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    const savedUser = localStorage.getItem('codenames_user');
    return savedUser ? JSON.parse(savedUser) : {
      id: `user-${uuidv4().substring(0, 8)}`,
      username: '',
      team: null
    };
  });

  useEffect(() => {
    localStorage.setItem('codenames_user', JSON.stringify(user));
  }, [user]);

  const updateUsername = (username) => {
    setUser(prevUser => ({ ...prevUser, username }));
  };

  const updateTeam = (team) => {
    setUser(prevUser => ({ ...prevUser, team }));
  };

  return (
    <UserContext.Provider value={{ user, updateUsername, updateTeam }}>
      {children}
    </UserContext.Provider>
  );
};