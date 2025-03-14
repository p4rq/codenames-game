import React, { useState, useContext } from 'react';
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
          <label htmlFor="username">Your Name:</label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter your name"
          />
        </div>
        
        {!isJoining ? (
          <div className="create-game-section">
            <button 
              className="create-btn"
              onClick={handleCreateGame}
              disabled={!username}
            >
              Create New Game
            </button>
            <p>
              Already have a game code? <button onClick={() => setIsJoining(true)}>Join Game</button>
            </p>
          </div>
        ) : (
          <div className="join-game-section">
            <div className="form-section">
              <label htmlFor="gameId">Game Code:</label>
              <input
                id="gameId"
                type="text"
                value={gameId}
                onChange={(e) => setGameId(e.target.value)}
                placeholder="Enter game code"
              />
            </div>
            
            <div className="form-section">
              <label htmlFor="team">Team:</label>
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
              Want to create a new game? <button onClick={() => setIsJoining(false)}>Create Game</button>
            </p>
          </div>
        )}
        
        {error && <p className="error-message">{error}</p>}
      </div>
    </div>
  );
};

export default HomePage;