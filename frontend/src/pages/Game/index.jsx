import React, { useContext, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { GameContext } from '../../context/GameContext';
import { UserContext } from '../../context/UserContext';
import './style.css';

const GamePage = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(UserContext);
  const { getGameState, revealCard, setSpymaster, endTurn, error: contextError } = useContext(GameContext);
  
  const [gameState, setGameState] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null); // Add state for local errors

  // Redirect if no game ID
  useEffect(() => {
    if (!gameId || gameId === 'undefined') {
      console.error('Invalid game ID:', gameId);
      navigate('/');
    }
  }, [gameId, navigate]);

  // Load game data
  useEffect(() => {
    if (!gameId || !user || gameId === 'undefined') return;

    const fetchGameState = async () => {
      setLoading(true);
      const data = await getGameState(gameId);
      if (data) {
        setGameState(data);
      }
      setLoading(false);
    };

    fetchGameState();
    
    // Poll for updates every 3 seconds
    const interval = setInterval(fetchGameState, 3000);
    return () => clearInterval(interval);
  }, [gameId, user, getGameState]);

  // Clear error after 5 seconds
  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => {
        setError(null);
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [error]);

  const handleCardClick = async (cardId) => {
    if (!user || !gameState) return;
    
    // Find the current player
    const currentPlayer = gameState.players.find(p => p.id === user.id);
    if (!currentPlayer) return;
    
    // Don't allow spymasters to reveal cards
    if (currentPlayer.is_spymaster) {
      setError("Spymasters can't reveal cards!");
      return;
    }
    
    // Only allow revealing cards on your team's turn
    if (currentPlayer.team !== gameState.current_turn) {
      setError("It's not your team's turn!");
      return;
    }
    
    const updatedGame = await revealCard(gameId, cardId, user.id);
    if (updatedGame) {
      setGameState(updatedGame);
    }
  };
  
  const handleSetSpymaster = async () => {
    if (!user || !gameState) return;
    
    const updatedGame = await setSpymaster(gameId, user.id);
    if (updatedGame) {
      setGameState(updatedGame);
    }
  };
  
  const handleEndTurn = async () => {
    if (!user || !gameState) return;
    
    // Find the current player
    const currentPlayer = gameState.players.find(p => p.id === user.id);
    if (!currentPlayer) return;
    
    // Only allow ending turn on your team's turn
    if (currentPlayer.team !== gameState.current_turn) {
      setError("It's not your team's turn!");
      return;
    }
    
    const updatedGame = await endTurn(gameId, user.id);
    if (updatedGame) {
      setGameState(updatedGame);
    }
  };
  
  // Find current player in game state
  const currentPlayer = gameState?.players?.find(p => p.id === user.id);
  const isCurrentPlayerTurn = currentPlayer?.team === gameState?.current_turn;
  const isGameOver = gameState?.winning_team !== null;
  
  if (loading) {
    return <div className="loading">Loading game...</div>;
  }
  
  if (!gameState) {
    return <div className="error">Game not found</div>;
  }
  
  return (
    <div className="game-container">
      <div className="game-header">
        <h1>Codenames - Game {gameId}</h1>
        <div className="game-info">
          <div className="teams-info">
            <div className={`team red ${gameState.current_turn === 'red' ? 'current-turn' : ''}`}>
              Red Team: {gameState.red_cards_left} cards left
            </div>
            <div className={`team blue ${gameState.current_turn === 'blue' ? 'current-turn' : ''}`}>
              Blue Team: {gameState.blue_cards_left} cards left
            </div>
          </div>
          
          {gameState.winning_team && (
            <div className={`winner ${gameState.winning_team}`}>
              {gameState.winning_team.toUpperCase()} TEAM WINS!
            </div>
          )}
        </div>
      </div>
      
      <div className="game-content">
        <div className="card-grid">
          {gameState?.cards?.map(card => (
            <div 
              key={card.id} 
              className={`game-card ${card.revealed ? card.type : ''} ${currentPlayer?.is_spymaster && !card.revealed ? `spymaster-${card.type}` : ''}`}
              onClick={() => !card.revealed && !isGameOver && handleCardClick(card.id)}
            >
              {card.word}
            </div>
          ))}
        </div>
        
        <div className="game-sidebar">
          <div className="players-list">
            <h3>Players</h3>
            <div className="team-players">
              <h4>Red Team</h4>
              <ul>
                {gameState?.players
                  ?.filter(p => p.team === 'red')
                  ?.map(p => (
                    <li key={p.id} className={p.id === user?.id ? 'current-player' : ''}>
                      {p.username} {p.is_spymaster ? '(Spymaster)' : ''}
                    </li>
                  ))}
              </ul>
            </div>
            <div className="team-players">
              <h4>Blue Team</h4>
              <ul>
                {gameState?.players
                  ?.filter(p => p.team === 'blue')
                  ?.map(p => (
                    <li key={p.id} className={p.id === user?.id ? 'current-player' : ''}>
                      {p.username} {p.is_spymaster ? '(Spymaster)' : ''}
                    </li>
                  ))}
              </ul>
            </div>
          </div>
          
          <div className="game-actions">
            {!currentPlayer?.is_spymaster && !isGameOver && (
              <button onClick={handleSetSpymaster}>Become Spymaster</button>
            )}
            
            {isCurrentPlayerTurn && !isGameOver && (
              <button onClick={handleEndTurn}>End Turn</button>
            )}
            
            {isGameOver && (
              <button onClick={() => window.location.href = "/"}>New Game</button>
            )}
          </div>
          
          {(error || contextError) && <div className="error-message">{error || contextError}</div>}
        </div>
      </div>
    </div>
  );
};

export default GamePage;