import React, { createContext, useState, useContext } from 'react';
import axios from 'axios';
import { UserContext } from './UserContext';

export const GameContext = createContext();

// Define API base URL with the /api prefix
const API_URL = '/api';

// Add request/response interceptors for debugging
axios.interceptors.request.use(
  config => {
    console.log('API Request:', {
      method: config.method,
      url: config.url,
      data: config.data
    });
    return config;
  },
  error => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  response => {
    console.log('API Response:', {
      status: response.status,
      data: response.data
    });
    return response;
  },
  error => {
    console.error('API Response Error:', {
      status: error.response?.status,
      data: error.response?.data
    });
    return Promise.reject(error);
  }
);

export const GameProvider = ({ children }) => {
  const { user, updateUser } = useContext(UserContext);
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

  const startNewGame = async (playerId, username) => {
    clearError();
    
    try {
      console.log(`Starting new game for player: ${playerId}, ${username}`);
      
      const response = await axios.post(`${API_URL}/game/start`, {
        creator_id: playerId,  // Changed from player_id to creator_id
        username: username
      });
      
      if (response.status === 200 || response.status === 201) {
        console.log("Server response:", response.data);
        setGame(response.data);
        return response.data;
      } else {
        throw new Error(`Failed to start game: ${response.statusText}`);
      }
    } catch (err) {
      console.error("Error starting game:", err);
      setError(`Failed to start game: ${err.message || 'Unknown error'}`);
      return null;
    }
  };

  // Add the joinExistingGame function
  const joinExistingGame = async (gameId, userId, username, team) => {
    try {
      clearError();
      console.log("Joining game:", { gameId, userId, username, team });
      
      const response = await axios.post(`${API_URL}/game/join`, {
        game_id: gameId,
        player_id: userId,
        username: username,
        team: team
      });
      
      console.log("Join game response:", response.data);
      
      if (!response.data || !response.data.id) {
        console.error("Invalid game response:", response.data);
        setError("Server returned an invalid game. Please try again.");
        return null;
      }
      
      const joinedGame = response.data;
      setGame(joinedGame);
      return joinedGame;
    } catch (err) {
      console.error("Error joining game:", err);
      setError(err.response?.data || 'Failed to join game. Please try again.');
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

  const changeTeam = async (gameId, playerId, team) => {
    try {
      clearError();
      const response = await axios.post(`${API_URL}/game/change-team`, {
        game_id: gameId,
        player_id: playerId,
        team: team
      });
      return response.data;
    } catch (err) {
      console.error("Error changing team:", err);
      setError(err.response?.data || 'Failed to change team.');
      return null;
    }
  };

  // Set user team both in backend and in user context
  const setUserTeam = async (gameId, team) => {
    if (!user) return null;
    
    try {
      // Update team on server
      const gameResponse = await changeTeam(gameId, user.id, team);
      if (!gameResponse) {
        throw new Error('Failed to update team on server');
      }
      
      // Update user in local context/storage
      const updatedUser = { ...user, team };
      updateUser(updatedUser);
      return updatedUser;
    } catch (error) {
      console.error('Error setting user team:', error);
      return null;
    }
  };

  // Fixed handleTeamChange function
  const handleTeamChange = async (gameId, teamColor) => {
    if (!user || !game) return;

    try {
      // Update user team on server and in context
      const updatedUser = await setUserTeam(gameId, teamColor);
      
      if (!updatedUser) {
        console.error('Failed to update user team');
        return;
      }

      // Update local game state too
      const updatedGame = { ...game };
      
      // Find and update the player in the game state
      if (updatedGame.players) {
        const playerIndex = updatedGame.players.findIndex(p => p.id === user.id);
        if (playerIndex >= 0) {
          updatedGame.players[playerIndex].team = teamColor;
          setGame(updatedGame);
        }
      }
      
      console.log(`Team changed to ${teamColor} for user ${user.username}`);
    } catch (error) {
      console.error('Error changing team:', error);
    }
  };
  
  // Update the provider value with all functions
  return (
    <GameContext.Provider 
      value={{ 
        game, 
        error, 
        startNewGame, 
        joinExistingGame,
        getGameState, 
        revealCard, 
        setSpymaster, 
        endTurn,
        changeTeam,
        setUserTeam,
        handleTeamChange
      }}
    >
      {children}
    </GameContext.Provider>
  );
};