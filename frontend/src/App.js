import React from 'react';
import { Routes, Route } from 'react-router-dom';
import HomePage from './pages/Home';
import GamePage from './pages/Game';
import { UserProvider } from './context/UserContext';
import { GameProvider } from './context/GameContext';
import './App.css';

function App() {
  return (
    <UserProvider>
      <GameProvider>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/game/:gameId" element={<GamePage />} />
        </Routes>
      </GameProvider>
    </UserProvider>
  );
}

export default App;