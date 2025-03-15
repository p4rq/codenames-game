import React, { useContext, useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { GameContext } from '../../context/GameContext';
import { UserContext } from '../../context/UserContext';
import Navbar from '../../components/Navbar';
import './style.css';

const GamePage = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const { user } = useContext(UserContext);
  const { getGameState, revealCard, setSpymaster, endTurn, error: contextError } = useContext(GameContext);
  
  const [gameState, setGameState] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [darkMode, setDarkMode] = useState(() => {
    return localStorage.getItem('darkMode') === 'true';
  });
  
  // Add WebSocket connection ref
  const socketRef = useRef(null);
  const [isConnected, setIsConnected] = useState(false);

  const toggleDarkMode = () => {
    const newDarkMode = !darkMode;
    setDarkMode(newDarkMode);
    localStorage.setItem('darkMode', newDarkMode);
    if (newDarkMode) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  };

  // Redirect if no game ID
  useEffect(() => {
    if (!gameId || gameId === 'undefined') {
      console.error('Invalid game ID:', gameId);
      navigate('/');
    }
  }, [gameId, navigate]);

  // Set up WebSocket connection
  useEffect(() => {
    if (!gameId || !user || gameId === 'undefined') return;

    // Initial game state load
    const fetchInitialGameState = async () => {
      setLoading(true);
      const data = await getGameState(gameId);
      if (data) {
        setGameState(data);
      }
      setLoading(false);
    };

    fetchInitialGameState();

    // Create WebSocket connection
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const wsUrl = `${protocol}://${host}/ws/game/${gameId}?client_id=${user.id}`;
    
    const socket = new WebSocket(wsUrl);
    socketRef.current = socket;
    
    // Connection opened
    socket.addEventListener('open', (event) => {
      console.log('WebSocket connection established');
      setIsConnected(true);
    });

    // Listen for messages
    socket.addEventListener('message', (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('WebSocket message received:', data);
        setGameState(data);
      } catch (err) {
        console.error('Error parsing WebSocket message:', err);
      }
    });

    // Connection closed
    socket.addEventListener('close', (event) => {
      console.log('WebSocket connection closed');
      setIsConnected(false);
      
      // Try to reconnect after a delay
      setTimeout(() => {
        if (socketRef.current === socket) { // Only reconnect if this is still the current socket
          console.log('Attempting to reconnect WebSocket...');
          // The effect will run again and create a new connection
          setIsConnected(false);
        }
      }, 3000);
    });

    // Connection error
    socket.addEventListener('error', (event) => {
      console.error('WebSocket error:', event);
      setError('Lost connection to the game server. Trying to reconnect...');
    });

    // Clean up
    return () => {
      console.log('Closing WebSocket connection');
      if (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING) {
        socket.close();
      }
      socketRef.current = null;
    };
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
    
    // Optimistically update the UI
    const optimisticUpdate = {
      ...gameState,
      cards: gameState.cards.map(c => 
        c.id === cardId ? { ...c, revealed: true } : c
      )
    };
    setGameState(optimisticUpdate);
    
    // Call the API
    const updatedGame = await revealCard(gameId, cardId, user.id);
    if (updatedGame) {
      // WebSocket will handle the update, but just in case:
      setGameState(updatedGame);
    }
  };
  
  // Same optimistic updates for other actions
  const handleSetSpymaster = async () => {
    if (!user || !gameState) return;
    
    // Optimistic update
    const optimisticUpdate = {
      ...gameState,
      players: gameState.players.map(p => 
        p.id === user.id ? { ...p, is_spymaster: true } : p
      )
    };
    setGameState(optimisticUpdate);
    
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
    
    // Optimistic update
    const newTurn = gameState.current_turn === 'red' ? 'blue' : 'red';
    const optimisticUpdate = {
      ...gameState,
      current_turn: newTurn
    };
    setGameState(optimisticUpdate);
    
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
  
  // Rest of your component remains the same
  return (
    <>
      <Navbar 
        darkMode={darkMode} 
        toggleDarkMode={toggleDarkMode} 
        gameId={gameId}
      />
      <div className={`game-container ${darkMode ? 'dark-mode' : ''}`}>
        {/* Connection status indicator */}
        <div className={`connection-status ${isConnected ? 'connected' : 'disconnected'}`}>
          {isConnected ? (
            <span>🟢 Connected</span>
          ) : (
            <span>🔴 Connecting...</span>
          )}
        </div>
        
        {/* Rest of your existing UI */}
        <div className="game-header">
          <h1>Codenames - Game {gameId}</h1>
          <div className="game-info">
            <div className="teams-info">
              <div className={`team blue ${gameState.current_turn === 'blue' ? 'current-turn' : ''}`}>
                Blue Team: {gameState.blue_cards_left} cards left
              </div>
              <div className={`team red ${gameState.current_turn === 'red' ? 'current-turn' : ''}`}>
                Red Team: {gameState.red_cards_left} cards left
              </div>
            </div>
            
            {gameState.winning_team && (
              <div className={`winner ${gameState.winning_team}`}>
                {gameState.winning_team.toUpperCase()} TEAM WINS!
              </div>
            )}
          </div>
        </div>
        
        <div className="three-column-layout">
          {/* BLUE TEAM PANEL */}
          <div className="team-panel blue-panel">
            <div className="team-header blue">
              <h2>Blue Team</h2>
              <div className={`team-status ${gameState.current_turn === 'blue' ? 'active' : ''}`}>
                {gameState.blue_cards_left} cards left
                {gameState.current_turn === 'blue' && <span className="turn-indicator">Current Turn</span>}
              </div>
            </div>
            
            <div className="team-players">
              <ul>
                {gameState?.players
                  ?.filter(p => p.team === 'blue')
                  ?.map(p => (
                    <li key={p.id} className={`
                      ${p.id === user?.id ? 'current-player' : ''}
                      ${p.is_spymaster ? 'spymaster' : ''}
                    `}>
                      {p.username} 
                      {p.is_spymaster && <span className="role-badge">Spymaster</span>}
                      {p.id === user?.id && <span className="you-badge">You</span>}
                    </li>
                  ))}
                {gameState?.players?.filter(p => p.team === 'blue').length === 0 && (
                  <li className="empty-team">No players yet</li>
                )}
              </ul>
            </div>
            
            {currentPlayer?.team === 'blue' && (
              <div className="team-actions">
                {!currentPlayer.is_spymaster && !isGameOver && (
                  <button className="spymaster-btn" onClick={handleSetSpymaster}>Become Spymaster</button>
                )}
                
                {isCurrentPlayerTurn && !isGameOver && (
                  <button className="end-turn-btn" onClick={handleEndTurn}>End Turn</button>
                )}
              </div>
            )}
          </div>

          {/* MIDDLE SECTION - CARD GRID */}
          <div className="middle-section">
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
            
            <div className="game-controls">
              {isGameOver && (
                <button className="new-game-btn" onClick={() => window.location.href = "/"}>New Game</button>
              )}
              
              {(error || contextError) && <div className="error-message">{error || contextError}</div>}
            </div>
          </div>

          {/* RED TEAM PANEL */}
          <div className="team-panel red-panel">
            <div className="team-header red">
              <h2>Red Team</h2>
              <div className={`team-status ${gameState.current_turn === 'red' ? 'active' : ''}`}>
                {gameState.red_cards_left} cards left
                {gameState.current_turn === 'red' && <span className="turn-indicator">Current Turn</span>}
              </div>
            </div>
            
            <div className="team-players">
              <ul>
                {gameState?.players
                  ?.filter(p => p.team === 'red')
                  ?.map(p => (
                    <li key={p.id} className={`
                      ${p.id === user?.id ? 'current-player' : ''}
                      ${p.is_spymaster ? 'spymaster' : ''}
                    `}>
                      {p.username}
                      {p.is_spymaster && <span className="role-badge">Spymaster</span>}
                      {p.id === user?.id && <span className="you-badge">You</span>}
                    </li>
                  ))}
                {gameState?.players?.filter(p => p.team === 'red').length === 0 && (
                  <li className="empty-team">No players yet</li>
                )}
              </ul>
            </div>
            
            {currentPlayer?.team === 'red' && (
              <div className="team-actions">
                {!currentPlayer.is_spymaster && !isGameOver && (
                  <button className="spymaster-btn" onClick={handleSetSpymaster}>Become Spymaster</button>
                )}
                
                {isCurrentPlayerTurn && !isGameOver && (
                  <button className="end-turn-btn" onClick={handleEndTurn}>End Turn</button>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
};

export default GamePage;