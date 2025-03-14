import React, { useContext, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import GameBoard from '../../components/GameBoard';
import PlayerList from '../../components/PlayerList';
import Chat from '../../components/Chat';
import './style.css';

const GamePage = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  
  const { user } = useContext(UserContext);
  const { 
    gameState, 
    loading, 
    error, 
    refreshGameState, 
    becomeSpymaster, 
    finishTurn 
  } = useContext(GameContext);
  
  const [refreshInterval, setRefreshInterval] = useState(null);

  // Get initial game state and set up polling
  useEffect(() => {
    if (!user || !gameId) return;
    
    // Get initial game state
    refreshGameState(gameId);
    
    // Set up polling every 3 seconds
    const interval = setInterval(() => {
      refreshGameState(gameId);
    }, 3000);
    
    setRefreshInterval(interval);
    
    // Clean up on unmount
    return () => {
      if (interval) clearInterval(interval);
    };
  }, [gameId, user, refreshGameState]);

  // Handle user actions
  const handleSetSpymaster = () => {
    if (!user || !gameId) return;
    becomeSpymaster(gameId, user.id);
  };
  
  const handleEndTurn = () => {
    if (!user || !gameId) return;
    finishTurn(gameId, user.id);
  };
  
  const handleBackToLobby = () => {
    navigate('/');
  };
  
  // Get player's team and role
  const currentPlayer = gameState?.players?.find(p => p.id === user?.id);
  const isSpymaster = currentPlayer?.is_spymaster || false;
  const playerTeam = currentPlayer?.team || 'spectator';
  
  // Check if it's this player's team's turn
  const isPlayerTurn = gameState?.current_turn === playerTeam;
  
  return (
    <div className="game-page">
      <header className="game-header">
        <h1>Codenames</h1>
        <div className="game-info">
          {gameState && (
            <>
              <span className="game-code">Game Code: {gameId}</span>
              <span className={`current-turn ${gameState.current_turn}-turn`}>
                Current Turn: {gameState.current_turn} Team
              </span>
              <div className="score">
                <span className="red-score">Red: {gameState.red_cards_left}</span>
                <span className="blue-score">Blue: {gameState.blue_cards_left}</span>
              </div>
            </>
          )}
          {gameState?.winning_team && (
            <div className="winner-announcement">
              {gameState.winning_team} Team Wins!
            </div>
          )}
        </div>
        <button className="back-btn" onClick={handleBackToLobby}>Back to Lobby</button>
      </header>
      
      <div className="game-content">
        <aside className="game-sidebar">
          <PlayerList players={gameState?.players || []} currentPlayerId={user?.id} />
          
          <div className="game-actions">
            {!isSpymaster && (
              <button onClick={handleSetSpymaster}>
                Become Spymaster
              </button>
            )}
            {isPlayerTurn && !gameState?.winning_team && (
              <button 
                onClick={handleEndTurn}
                className="end-turn-btn"
              >
                End Turn
              </button>
            )}
          </div>
          
          <Chat gameId={gameId} />
        </aside>
        
        <main className="game-board-container">
          {loading && !gameState ? (
            <div className="loading">Loading game...</div>
          ) : error ? (
            <div className="error-message">{error}</div>
          ) : gameState ? (
            <GameBoard 
              cards={gameState.cards} 
              isSpymaster={isSpymaster}
              currentTeam={gameState.current_turn}
              playerTeam={playerTeam}
              gameOver={!!gameState.winning_team}
            />
          ) : (
            <div className="error-message">Game not found</div>
          )}
        </main>
      </div>
    </div>
  );
};

export default GamePage;