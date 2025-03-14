import React, { createContext, useState } from 'react';
import axios from 'axios';

export const GameContext = createContext();

// Define API base URL with the /api prefix
const API_URL = '/api';

export const GameProvider = ({ children }) => {
  const [error, setError] = useState(null);
  const [game, setGame] = useState({
    id: null,
    cards: [],
    players: [],
    current_turn: null,
    red_cards_left: 0,
    blue_cards_left: 0,
    winning_team: null
  });
  
  const clearError = () => setError(null);

  const startNewGame = async (userId, username) => {
    try {
      clearError();
      console.log("Creating game with:", { userId, username });
      
      // Make sure we use /api/game/start not just /game/start
      const response = await axios.post(`${API_URL}/game/start`, {
        creator_id: userId,
        username: username
      });
      
      console.log("Server response:", response.data);
      
      // Check for valid response
      if (!response.data || !response.data.id) {
        console.error("Invalid game response:", response.data);
        setError("Server returned an invalid game. Please try again.");
        return null;
      }
      
      const newGame = response.data;
      setGame(newGame);
      return newGame;
    } catch (err) {
      console.error("Error creating game:", err);
      setError(err.response?.data || 'Failed to create game. Please try again.');
      return null;
    }
  };

  // Other methods like getGameState also need the API prefix
  const getGameState = async (gameId) => {
    try {
      clearError();
      const response = await axios.get(`${API_URL}/game/state?id=${gameId}`);
      return response.data;
    } catch (err) {
      console.error("Error fetching game state:", err);
      setError(err.response?.data || 'Failed to load game.');
      return null;
    }
  };

  // Make sure all other API calls use the same prefix
  const revealCard = async (gameId, cardId, playerId) => {
    try {
      clearError();
      const response = await axios.post(`${API_URL}/game/reveal`, {
        game_id: gameId,
        card_id: cardId,
        player_id: playerId
      });
      return response.data;
    } catch (err) {
      console.error("Error revealing card:", err);
      setError(err.response?.data || 'Failed to reveal card.');
      return null;
    }
  };

  const setSpymaster = async (gameId, playerId) => {
    try {
      clearError();
      const response = await axios.post(`${API_URL}/game/set-spymaster?game_id=${gameId}&player_id=${playerId}`);
      return response.data;
    } catch (err) {
      console.error("Error setting spymaster:", err);
      setError(err.response?.data || 'Failed to become spymaster.');
      return null;
    }
  };

  const endTurn = async (gameId, playerId) => {
    try {
      clearError();
      const response = await axios.post(`${API_URL}/game/end-turn?game_id=${gameId}&player_id=${playerId}`);
      return response.data;
    } catch (err) {
      console.error("Error ending turn:", err);
      setError(err.response?.data || 'Failed to end turn.');
      return null;
    }
  };

  return (
    <GameContext.Provider 
      value={{ 
        game, 
        error, 
        startNewGame, 
        joinExistingGame: () => {}, // Implement this as needed
        getGameState, 
        revealCard, 
        setSpymaster, 
        endTurn 
      }}
    >
      {children}
    </GameContext.Provider>
  );
};