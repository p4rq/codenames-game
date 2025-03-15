import React, { createContext, useState, useCallback, useRef, useEffect } from 'react';
import axios from 'axios';
import API_URL from '../config/api';

export const GameContext = createContext({});

export const GameProvider = ({ children }) => {
  const [game, setGame] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [wsConnected, setWsConnected] = useState(false);
  const socketRef = useRef(null);
  const pollingIntervalRef = useRef(null);
  const pollingTimeoutRef = useRef(null);
  
  // For debugging
  useEffect(() => {
    console.log("WebSocket connected state:", wsConnected);
  }, [wsConnected]);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Function to fetch game state (polling)
  const fetchGameState = useCallback(async (gameId) => {
    if (!gameId) return;
    
    try {
      console.log("Polling game state for:", gameId);
      const response = await axios.get(`${API_URL}/api/game/state?game_id=${gameId}`);
      setGame(response.data);
    } catch (err) {
      console.error("Error fetching game state:", err);
    }
  }, []);

  // Stop polling function
  const stopPolling = useCallback(() => {
    if (pollingIntervalRef.current) {
      console.log("Stopping polling");
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
  }, []);

  // Start polling with rate limiting to prevent excessive requests
  const startPolling = useCallback((gameId) => {
    if (wsConnected) {
      console.log("WebSocket connected, not starting polling");
      return;
    }
    
    if (pollingIntervalRef.current) {
      stopPolling();
    }
    
    console.log("Starting polling for game:", gameId);
    // Use a longer interval (10 seconds) to reduce server load
    pollingIntervalRef.current = setInterval(() => {
      // Additional check to avoid polling when WebSocket is connected
      if (!wsConnected && gameId) {
        fetchGameState(gameId);
      } else if (wsConnected) {
        console.log("WebSocket now connected, stopping polling");
        stopPolling();
      }
    }, 10000); // 10 seconds
  }, [fetchGameState, stopPolling, wsConnected]);

  // Setup WebSocket connection with reconnection logic
  const setupWebSocket = useCallback((gameId) => {
    if (!gameId) return null;
    
    // Clean up any existing socket
    if (socketRef.current) {
      console.log("Closing existing WebSocket connection");
      socketRef.current.close();
      socketRef.current = null;
    }
    
    const userId = localStorage.getItem('userId');
    // Use relative path for WebSocket URL to match the current host
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${wsProtocol}//${window.location.host}/ws/game/${gameId}?client_id=${userId}`;
    
    console.log("Setting up WebSocket connection to:", wsUrl);
    
    try {
      const socket = new WebSocket(wsUrl);
      
      socket.onopen = () => {
        console.log("WebSocket connection established");
        setWsConnected(true);
        stopPolling();
      };
      
      socket.onmessage = (event) => {
        try {
          const gameUpdate = JSON.parse(event.data);
          console.log("Received game update via WebSocket");
          setGame(gameUpdate);
        } catch (err) {
          console.error("Error parsing WebSocket message:", err);
        }
      };
      
      socket.onclose = (event) => {
        console.log("WebSocket connection closed:", event);
        setWsConnected(false);
        
        // Start fallback polling with a delay
        if (gameId && !pollingIntervalRef.current) {
          console.log("WebSocket closed, setting up fallback polling");
          startPolling(gameId);
        }
      };
      
      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
        setWsConnected(false);
      };
      
      socketRef.current = socket;
      return socket;
    } catch (err) {
      console.error("Error setting up WebSocket:", err);
      return null;
    }
  }, [stopPolling, startPolling]);

  // Function to create a new game
  const createGame = useCallback(async (playerId, username, language = 'en') => {
    try {
      clearError();
      setLoading(true);
      
      console.log("Creating game with:", { playerId, username, language });
      
      const response = await axios.post(`${API_URL}/api/game/start`, {
        player_id: playerId,
        username,
        language
      });
      
      console.log("Server response:", response.data);
      setGame(response.data);
      
      // Set up WebSocket connection after creating the game
      setupWebSocket(response.data.id);
      
      return response.data;
    } catch (err) {
      console.error("Error creating game:", err);
      setError(err.response?.data?.message || "Failed to create game");
      return null;
    } finally {
      setLoading(false);
    }
  }, [setupWebSocket, clearError]);

  // Function to join or load a game
  const joinOrLoadGame = useCallback(async (gameId, playerId, username, team) => {
    try {
      clearError();
      setLoading(true);
      
      // First try to join the game
      const joinResponse = await axios.post(`${API_URL}/api/game/join`, {
        game_id: gameId,
        player_id: playerId,
        username,
        team
      });
      
      setGame(joinResponse.data);
      
      // Set up WebSocket connection after successfully joining
      setupWebSocket(gameId);
      
      // Set up fallback polling only if WebSocket isn't connected yet
      if (!wsConnected) {
        console.log("WebSocket not yet connected, setting up initial fallback polling");
        // Delay the start of polling by a few seconds to give WebSocket a chance to connect
        pollingTimeoutRef.current = setTimeout(() => {
          if (!wsConnected) {
            startPolling(gameId);
          }
        }, 3000); // Wait 3 seconds before starting polling
      }
      
      return joinResponse.data;
    } catch (err) {
      console.error("Error joining game:", err);
      setError(err.response?.data?.message || 'Failed to join game');
      return null;
    } finally {
      setLoading(false);
    }
  }, [setupWebSocket, wsConnected, startPolling, clearError]);

  // Function to reveal a card
  const revealCard = useCallback(async (gameId, cardIndex) => {
    try {
      clearError();
      setLoading(true);
      
      console.log(`Revealing card ${cardIndex} in game ${gameId}`);
      
      const response = await axios.post(`${API_URL}/api/game/reveal`, {
        game_id: gameId,
        card_index: cardIndex
      });
      
      setGame(response.data);
      return response.data;
    } catch (err) {
      console.error("Error revealing card:", err);
      setError(err.response?.data?.message || "Failed to reveal card");
      return null;
    } finally {
      setLoading(false);
    }
  }, [clearError]);

  // Function to set a player as spymaster
  const setSpymaster = useCallback(async (gameId, playerId) => {
    try {
      clearError();
      setLoading(true);
      
      console.log(`Setting player ${playerId} as spymaster in game ${gameId}`);
      
      const response = await axios.post(`${API_URL}/api/game/set-spymaster`, {
        game_id: gameId,
        player_id: playerId
      });
      
      setGame(response.data);
      return response.data;
    } catch (err) {
      console.error("Error setting spymaster:", err);
      setError(err.response?.data?.message || "Failed to set spymaster");
      return null;
    } finally {
      setLoading(false);
    }
  }, [clearError]);

  // Function to end the current team's turn
  const endTurn = useCallback(async (gameId) => {
    try {
      clearError();
      setLoading(true);
      
      console.log(`Ending turn in game ${gameId}`);
      
      const response = await axios.post(`${API_URL}/api/game/end-turn`, {
        game_id: gameId
      });
      
      setGame(response.data);
      return response.data;
    } catch (err) {
      console.error("Error ending turn:", err);
      setError(err.response?.data?.message || "Failed to end turn");
      return null;
    } finally {
      setLoading(false);
    }
  }, [clearError]);

  // Function to change team with WebSocket reconnection
  const changeTeam = async (gameId, playerId, team) => {
    try {
      clearError();
      console.log(`Changing team for player ${playerId} to ${team} in game ${gameId}`);
      
      const response = await axios.post(`${API_URL}/api/game/change-team`, {
        game_id: gameId,
        player_id: playerId,
        team: team
      });
      
      console.log("Change team response:", response.data);
      
      if (!response.data || !response.data.id) {
        console.error("Invalid game response:", response.data);
        setError("Server returned an invalid game. Please try again.");
        return null;
      }
      
      const updatedGame = response.data;
      setGame(updatedGame);
      
      // Re-establish WebSocket connection after team change
      // Give the server a moment to process the team change before reconnecting
      setTimeout(() => {
        console.log("Re-establishing WebSocket connection after team change");
        setupWebSocket(gameId);
      }, 500);
      
      return updatedGame;
    } catch (err) {
      console.error('Error changing team:', err);
      setError(err.response?.data?.message || 'Failed to change team. Please try again.');
      return null;
    }
  };

  // Clean up all resources on unmount
  useEffect(() => {
    return () => {
      console.log("Cleaning up GameContext resources");
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
      
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
        pollingIntervalRef.current = null;
      }
      
      if (pollingTimeoutRef.current) {
        clearTimeout(pollingTimeoutRef.current);
        pollingTimeoutRef.current = null;
      }
    };
  }, []);

  // Value object for the context provider
  const contextValue = {
    game,
    loading,
    error,
    clearError,
    createGame,
    joinOrLoadGame,
    revealCard,
    setSpymaster,
    endTurn,
    changeTeam,
    wsConnected
  };

  return (
    <GameContext.Provider value={contextValue}>
      {children}
    </GameContext.Provider>
  );
};