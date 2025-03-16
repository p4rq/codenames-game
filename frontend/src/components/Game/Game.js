import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import Board from '../Board/Board';
import Chat from '../Chat/Chat';
import './Game.css';

const Game = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const [gameState, setGameState] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [username, setUsername] = useState(() => {
    return localStorage.getItem('codenames_username') || 'Anonymous';
  });
  const [team, setTeam] = useState('red');

  // Get current player from game state
  const currentPlayer = gameState?.players?.find(p => p.id === api.getUserId());
  
  const refreshGame = async () => {
    try {
      const data = await api.getGameState(gameId);
      setGameState(data);
      setError('');
    } catch (err) {
      setError('Error loading game: ' + (err.response?.data || err.message));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    // Save username to localStorage
    if (username !== 'Anonymous') {
      localStorage.setItem('codenames_username', username);
    }

    // Initial game load
    refreshGame();

    // Set up polling for game state
    const intervalId = setInterval(refreshGame, 5000);
    
    return () => clearInterval(intervalId);
  }, [gameId]);

  const handleJoinTeam = async (selectedTeam) => {
    try {
      setIsLoading(true);
      await api.joinGame(gameId, username, selectedTeam);
      setTeam(selectedTeam);
      await refreshGame();
    } catch (err) {
      setError('Error joining team: ' + (err.response?.data || err.message));
    } finally {
      setIsLoading(false);
    }
  };

  const handleSetSpymaster = async () => {
    try {
      setIsLoading(true);
      await api.setSpymaster(gameId);
      await refreshGame();
    } catch (err) {
      setError('Error setting spymaster: ' + (err.response?.data || err.message));
    } finally {
      setIsLoading(false);
    }
  };

  const handleEndTurn = async () => {
    try {
      setIsLoading(true);
      await api.endTurn(gameId);
      await refreshGame();
    } catch (err) {
      setError('Error ending turn: ' + (err.response?.data || err.message));
    } finally {
      setIsLoading(false);
    }
  };

  const handleCardClick = async (cardId) => {
    if (!currentPlayer || gameState.winning_team) return;
    
    // Only allow clicking if it's your team's turn
    if (currentPlayer.team !== gameState.current_turn) {
      setError("It's not your team's turn!");
      return;
    }
    
    // Don't allow spymasters to click cards
    if (currentPlayer.is_spymaster) {
      setError("Spymasters can't reveal cards!");
      return;
    }

    try {
      setIsLoading(true);
      await api.revealCard(gameId, cardId);
      await refreshGame();
    } catch (err) {
      setError('Error revealing card: ' + (err.response?.data || err.message));
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading && !gameState) {
    return <div className="loading">Loading game...</div>;
  }

  if (!gameState) {
    return (
      <div className="error-container">
        <p>Error: {error || 'Game not found'}</p>
        <button onClick={() => navigate('/')}>Back to Home</button>
      </div>
    );
  }

  return (
    <div className="game-container">
      <div className="game-header">
        <h1>Game #{gameId}</h1>
        <div className="game-info">
          <div className="teams-info">
            <div className={`team-red ${gameState.current_turn === 'red' ? 'active-turn' : ''}`}>
              Red Team: {gameState.red_cards_left} cards left
            </div>
            <div className={`team-blue ${gameState.current_turn === 'blue' ? 'active-turn' : ''}`}>
              Blue Team: {gameState.blue_cards_left} cards left
            </div>
          </div>
          
          <div className="current-player">
            {currentPlayer ? (
              <>
                You: {currentPlayer.username} ({currentPlayer.team} team)
                {currentPlayer.is_spymaster && ' - Spymaster'}
              </>
            ) : (
              <>Not joined yet</>
            )}
          </div>
          
          {gameState.winning_team && (
            <div className="winner-announcement">
              {gameState.winning_team.toUpperCase()} TEAM WINS!
            </div>
          )}
        </div>
      </div>
      
      {user && user.team === 'spectator' && (
        <div className="spectator-notice">
          <p>
            You're currently a <strong>Spectator</strong>. 
            Choose a team from your profile menu to participate in the game.
          </p>
        </div>
      )}
      
      <div className="game-content">
        <div className="game-board">
          <Board 
            cards={gameState.cards} 
            onCardClick={handleCardClick} 
            isSpymaster={currentPlayer?.is_spymaster || false}
          />
          
          <div className="game-controls">
            {!currentPlayer && (
              <div className="join-controls">
                <input
                  type="text"
                  placeholder="Your username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                />
                <select 
                  value={team} 
                  onChange={(e) => setTeam(e.target.value)}
                >
                  <option value="red">Red Team</option>
                  <option value="blue">Blue Team</option>
                </select>
                <button onClick={() => handleJoinTeam(team)}>Join Game</button>
              </div>
            )}
            
            {currentPlayer && !currentPlayer.is_spymaster && (
              <button onClick={handleSetSpymaster}>Become Spymaster</button>
            )}
            
            {currentPlayer && gameState.current_turn === currentPlayer.team && !gameState.winning_team && (
              <button onClick={handleEndTurn}>End Turn</button>
            )}
            
            <button onClick={refreshGame}>Refresh Game</button>
            <button onClick={() => navigate('/')}>Leave Game</button>
          </div>
        </div>
        
        <div className="game-sidebar">
          <Chat 
            gameId={gameId} 
            username={currentPlayer?.username || username}
          />
        </div>
      </div>
      
      {error && <div className="error-message">{error}</div>}
    </div>
  );
};

export default Game;