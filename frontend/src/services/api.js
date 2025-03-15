import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

/**
 * Makes a GET request to the API
 * @param {string} endpoint - API endpoint
 * @param {object} params - Query parameters
 * @returns {Promise<any>} - Response data
 */
export const get = async (endpoint, params = {}) => {
  const url = new URL(`${API_URL}${endpoint}`);
  
  // Add query parameters
  Object.keys(params).forEach(key => {
    if (params[key] !== undefined && params[key] !== null) {
      url.searchParams.append(key, params[key]);
    }
  });
  
  const response = await fetch(url.toString());
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  
  return response.json();
};

/**
 * Makes a POST request to the API
 * @param {string} endpoint - API endpoint
 * @param {object} data - Request body data
 * @returns {Promise<any>} - Response data
 */
export const post = async (endpoint, data = {}) => {
  const response = await fetch(`${API_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  
  return response.json();
};

// Generate random user ID for this session
const USER_ID = `user-${Math.floor(Math.random() * 10000)}`;

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Intercept responses to handle common errors
api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

const apiMethods = {
  // Game operations
  createGame: async (username) => {
    const response = await api.post(`/game/start`, {
      creator_id: USER_ID,
      username: username || 'Anonymous'
    });
    return response.data;
  },

  joinGame: async (gameId, username, team) => {
    const response = await api.post(`/game/join`, {
      game_id: gameId,
      player_id: USER_ID,
      username: username || 'Anonymous',
      team: team || 'red'
    });
    return response.data;
  },

  getGameState: async (gameId) => {
    const response = await api.get(`/game/state?id=${gameId}`);
    return response.data;
  },

  setSpymaster: async (gameId) => {
    const response = await api.post(`/game/set-spymaster?game_id=${gameId}&player_id=${USER_ID}`);
    return response.data;
  },

  revealCard: async (gameId, cardId) => {
    const response = await api.post(`/game/reveal`, {
      game_id: gameId,
      card_id: cardId,
      player_id: USER_ID
    });
    return response.data;
  },

  endTurn: async (gameId) => {
    const response = await api.post(`/game/end-turn?game_id=${gameId}&player_id=${USER_ID}`);
    return response.data;
  },

  // Chat operations
  sendMessage: async (gameId, username, content) => {
    const response = await api.post(`/chat/send`, {
      content: content,
      sender_id: USER_ID,
      username: username || 'Anonymous',
      chat_id: gameId
    });
    return response.data;
  },

  getMessages: async (gameId) => {
    const response = await api.get(`/chat/messages?chat_id=${gameId}`);
    return response.data;
  },

  // Helper method to get current user ID
  getUserId: () => USER_ID
};

export default apiMethods;