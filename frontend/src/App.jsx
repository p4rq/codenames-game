import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import HomePage from './pages/Home';
import GamePage from './pages/Game';
import NotFoundPage from './pages/NotFound';
import { UserProvider } from './context/UserContext';
import { GameProvider } from './context/GameContext';

function App() {
  return (
    <UserProvider>
      <GameProvider>
        <Router>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/game/:gameId" element={<GamePage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Router>
      </GameProvider>
    </UserProvider>
  );
}

export default App;