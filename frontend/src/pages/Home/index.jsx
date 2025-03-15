import React, { useState, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './style.css';

const HomePage = () => {
  const { user, createUser, updateUsername } = useContext(UserContext);
  const { startNewGame, joinExistingGame, error } = useContext(GameContext);
  
  const [username, setUsername] = useState(user?.username || '');
  const [gameId, setGameId] = useState('');
  const [team, setTeam] = useState('red');
  const [isJoining, setIsJoining] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  
  const navigate = useNavigate();

  // Update username state when user changes
  useEffect(() => {
    if (user?.username) {
      setUsername(user.username);
    }
  }, [user]);

  const handleCreateGame = async () => {
    if (!username.trim()) return;
    setIsProcessing(true);
    
    try {
      // Ensure we have a user
      let currentUser = user;
      if (!currentUser) {
        currentUser = createUser(username);
      } else if (username !== currentUser.username) {
        currentUser = updateUsername(username);
      }
      
      // Now create the game
      const game = await startNewGame(currentUser.id, username);
      if (game) {
        navigate(`/game/${game.id}`);
      }
    } catch (err) {
      console.error('Error creating game:', err);
    } finally {
      setIsProcessing(false);
    }
  };

  const handleJoinGame = async () => {
    if (!username.trim() || !gameId) return;
    setIsProcessing(true);
    
    try {
      // Ensure we have a user
      let currentUser = user;
      if (!currentUser) {
        currentUser = createUser(username);
      } else if (username !== currentUser.username) {
        currentUser = updateUsername(username);
      }
      
      // Now join the game
      const game = await joinExistingGame(gameId, currentUser.id, username, team);
      if (game) {
        navigate(`/game/${game.id}`);
      }
    } catch (err) {
      console.error('Error joining game:', err);
    } finally {
      setIsProcessing(false);
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
              disabled={!username || !gameId || isProcessing}
            >
              {isProcessing ? 'Joining...' : 'Join Game'}
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
              disabled={!username || isProcessing}
            >
              {isProcessing ? 'Creating...' : 'Create New Game'}
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