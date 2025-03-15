import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { UserProvider } from './context/UserContext';
import { GameProvider } from './context/GameContext';
import HomePage from './pages/Home';
import GamePage from './pages/Game';
import NotFoundPage from './pages/NotFound';
import Navbar from './components/Navbar';
import './App.css';

function App() {
  const [darkMode, setDarkMode] = useState(
    () => localStorage.getItem('darkMode') === 'true'
  );

  useEffect(() => {
    localStorage.setItem('darkMode', darkMode);
    if (darkMode) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  }, [darkMode]);

  const toggleDarkMode = () => {
    setDarkMode(prevMode => !prevMode);
  };

  return (
    <UserProvider>
      <GameProvider>
        <Router>
          <div className={`app ${darkMode ? 'dark-mode' : ''}`}>
            <Routes>
              <Route 
                path="/" 
                element={<>
                  <Navbar darkMode={darkMode} toggleDarkMode={toggleDarkMode} />
                  <HomePage />
                </>} 
              />
              <Route 
                path="/game/:gameId" 
                element={({ match }) => <>
                  <Navbar 
                    darkMode={darkMode} 
                    toggleDarkMode={toggleDarkMode} 
                    gameId={match?.params?.gameId} 
                  />
                  <GamePage />
                </>} 
              />
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </div>
        </Router>
      </GameProvider>
    </UserProvider>
  );
}

export default App;