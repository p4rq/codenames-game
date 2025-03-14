import { get, post } from './api';

/**
 * Sends a chat message
 * @param {string} content - Message content
 * @param {string} userId - User ID
 * @param {string} username - Username
 * @param {string} chatId - Game ID for chat
 * @returns {Promise<any>} - Response data
 */
export const sendMessage = async (content, userId, username, chatId) => {
  return post('/chat/send', {
    content: content,
    sender_id: userId,
    username: username,
    chat_id: chatId
  });
};

/**
 * Gets all messages for a chat
 * @param {string} chatId - Game ID for chat
 * @returns {Promise<any>} - Chat messages
 */
export const getMessages = async (chatId) => {
  return get('/chat/messages', { chat_id: chatId });
};