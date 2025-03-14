import React, { useState, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './style.css';

const HomePage = () => {
  const { user, updateUsername } = useContext(UserContext);
  const { startNewGame, joinExistingGame, error } = useContext(GameContext);
  
  const [username, setUsername] = useState(user?.username || '');
  const [gameId, setGameId] = useState('');
  const [team, setTeam] = useState('red');
  const [isJoining, setIsJoining] = useState(false);
  
  const navigate = useNavigate();

  const handleCreateGame = async () => {
    if (!user) return;
    
    // Update user's name if changed
    if (username !== user.username) {
      updateUsername(username);
    }
    
    const game = await startNewGame(user.id, username);
    if (game) {
      navigate(`/game/${game.id}`);
    }
  };

  const handleJoinGame = async () => {
    if (!user || !gameId) return;
    
    // Update user's name if changed
    if (username !== user.username) {
      updateUsername(username);
    }
    
    const game = await joinExistingGame(gameId, user.id, username, team);
    if (game) {
      navigate(`/game/${game.id}`);
    }
  };

  return (
    <div className="home-container">
      <div className="home-content">
        <h1>Codenames</h1>
        
        <div className="form-section">
          <label htmlFor="username">Your Name</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter your name"
            required
          />
        </div>
        
        {isJoining ? (
          <>
            <div className="form-section">
              <label htmlFor="game-code">Game Code</label>
              <input
                type="text"
                id="game-code"
                value={gameId}
                onChange={(e) => setGameId(e.target.value)}
                placeholder="Enter game code"
                required
              />
            </div>
            
            <div className="form-section">
              <label htmlFor="team">Select Team</label>
              <select
                id="team"
                value={team}
                onChange={(e) => setTeam(e.target.value)}
              >
                <option value="red">Red Team</option>
                <option value="blue">Blue Team</option>
              </select>
            </div>
            
            <button
              className="join-btn"
              onClick={handleJoinGame}
              disabled={!username || !gameId}
            >
              Join Game
            </button>
            
            <p>
              Want to create a new game instead?{' '}
              <button onClick={() => setIsJoining(false)}>Create Game</button>
            </p>
          </>
        ) : (
          <>
            <button
              className="create-btn"
              onClick={handleCreateGame}
              disabled={!username}
            >
              Create New Game
            </button>
            
            <p>
              Have a game code?{' '}
              <button onClick={() => setIsJoining(true)}>Join Existing Game</button>
            </p>
          </>
        )}
        
        {error && <div className="error-message">{error}</div>}
      </div>
    </div>
  );
};

export default HomePage;