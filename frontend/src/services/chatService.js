import axios from 'axios';

// Standalone implementation not dependent on your api.js
const chatApi = axios.create({
  baseURL: '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
});

export const getMessages = async (gameId, team = null) => {
  try {
    let endpoint = `/api/games/${gameId}/messages`;
    if (team) {
      endpoint += `?team=${team}`;
    }
    
    const response = await chatApi.get(endpoint);
    return response.data;
  } catch (error) {
    console.error('Error fetching messages:', error);
    
    // Return empty array as fallback
    return [];
  }
};

export const sendMessage = async (content, senderId, username, gameId, team = null) => {
  try {
    const endpoint = `/api/games/${gameId}/messages`;
    const payload = {
      content,
      sender_id: senderId,
      username,
      team,
    };
    
    const response = await chatApi.post(endpoint, payload);
    return response.data;
  } catch (error) {
    console.error('Error sending message:', error);
    return null;
  }
};