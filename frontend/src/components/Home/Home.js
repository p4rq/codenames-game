import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';
import { createGame, joinGame } from '../../services/gameService';

const Home = () => {
  const [username, setUsername] = useState('');
  const [gameId, setGameId] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  // Generate a unique user ID for this session
  const userId = localStorage.getItem('userId') || `user-${Math.random().toString(36).substring(2, 9)}`;

  // Store user ID in localStorage
  if (!localStorage.getItem('userId')) {
    localStorage.setItem('userId', userId);
  }

  const handleCreateGame = async () => {
    if (!username.trim()) {
      setError('Please enter a username');
      return;
    }

    try {
      setError('');
      const game = await createGame(userId, username);
      
      // Store username in localStorage
      localStorage.setItem('username', username);

      // Navigate to the game page
      navigate(`/game/${game.id}`);
    } catch (err) {
      setError('Failed to create game. Please try again.');
      console.error('Error creating game:', err);
    }
  };

  const handleJoinGame = async () => {
    if (!username.trim()) {
      setError('Please enter a username');
      return;
    }

    if (!gameId.trim()) {
      setError('Please enter a game ID');
      return;
    }

    try {
      setError('');
      await joinGame(gameId, userId, username, 'red'); // Default to red team
      
      // Store username in localStorage
      localStorage.setItem('username', username);
      
      // Navigate to the game page
      navigate(`/game/${gameId}`);
    } catch (err) {
      setError('Failed to join game. Please check the game ID and try again.');
      console.error('Error joining game:', err);
    }
  };

  return (
    <div className="home-container">
      <h1>Codenames</h1>
      
      {error && <div className="error-message">{error}</div>}
      
      <div className="form-group">
        <input
          type="text"
          placeholder="Enter your username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>
      
      <div className="button-group">
        <button onClick={handleCreateGame}>Create New Game</button>
        
        <div className="or-divider">
          <span>OR</span>
        </div>
        
        <div className="form-group">
          <input
            type="text"
            placeholder="Enter game ID"
            value={gameId}
            onChange={(e) => setGameId(e.target.value)}
          />
        </div>
        
        <button onClick={handleJoinGame}>Join Game</button>
      </div>
    </div>
  );
};

export default Home;