import { get, post } from './api';

/**
 * Creates a new game
 * @param {string} userId - User ID
 * @param {string} username - Username
 * @returns {Promise<any>} - Game data
 */
export const createGame = async (userId, username) => {
  return post('/game/start', {
    creator_id: userId,
    username: username
  });
};

/**
 * Joins an existing game
 * @param {string} gameId - Game ID
 * @param {string} userId - User ID
 * @param {string} username - Username
 * @param {string} team - Team ('red' or 'blue')
 * @returns {Promise<any>} - Game data
 */
export const joinGame = async (gameId, userId, username, team) => {
  return post('/game/join', {
    game_id: gameId,
    player_id: userId,
    username: username,
    team: team
  });
};

/**
 * Gets the current state of a game
 * @param {string} gameId - Game ID
 * @returns {Promise<any>} - Game data
 */
export const getGame = async (gameId) => {
  return get('/game/state', { id: gameId });
};

/**
 * Reveals a card
 * @param {string} gameId - Game ID
 * @param {string} cardId - Card ID
 * @param {string} userId - User ID
 * @returns {Promise<any>} - Updated game data
 */
export const revealCard = async (gameId, cardId, userId) => {
  return post('/game/reveal', {
    game_id: gameId,
    card_id: cardId,
    player_id: userId
  });
};

/**
 * Sets the player as a spymaster
 * @param {string} gameId - Game ID
 * @param {string} userId - User ID
 * @returns {Promise<any>} - Updated game data
 */
export const setSpymaster = async (gameId, userId) => {
  return post('/game/set-spymaster', null, {
    game_id: gameId,
    player_id: userId
  });
};

/**
 * Ends the current team's turn
 * @param {string} gameId - Game ID
 * @param {string} userId - User ID
 * @returns {Promise<any>} - Updated game data
 */
export const endTurn = async (gameId, userId) => {
  return post('/game/end-turn', null, {
    game_id: gameId,
    player_id: userId
  });
};