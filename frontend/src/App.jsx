import React, { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { UserProvider } from './context/UserContext';
import { GameProvider } from './context/GameContext';
import Home from './components/Home';
import Game from './components/Game/Game';
import Navbar from './components/Navbar';
import './App.css';

function App() {
  const [darkMode, setDarkMode] = useState(() => {
    const saved = localStorage.getItem('darkMode');
    return saved ? JSON.parse(saved) : false;
  });

  useEffect(() => {
    localStorage.setItem('darkMode', JSON.stringify(darkMode));
    if (darkMode) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  }, [darkMode]);

  const toggleDarkMode = () => {
    setDarkMode(!darkMode);
  };

  return (
    <UserProvider>
      <GameProvider>
        <BrowserRouter>
          <Navbar darkMode={darkMode} toggleDarkMode={toggleDarkMode} />
          <div className={`app-container ${darkMode ? 'dark-mode' : ''}`}>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/game/:gameId" element={<Game />} />
            </Routes>
          </div>
        </BrowserRouter>
      </GameProvider>
    </UserProvider>
  );
}

export default App;