import React, { createContext, useState, useCallback } from 'react';
import { createGame, joinGame, getGame, revealCard, setSpymaster, endTurn } from '../services/gameService';

export const GameContext = createContext();

export const GameProvider = ({ children }) => {
  const [gameState, setGameState] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const startNewGame = useCallback(async (userId, username) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await createGame(userId, username);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to create game');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const joinExistingGame = useCallback(async (gameId, userId, username, team) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await joinGame(gameId, userId, username, team);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to join game');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const refreshGameState = useCallback(async (gameId) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await getGame(gameId);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to get game state');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const revealGameCard = useCallback(async (gameId, cardId, playerId) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await revealCard(gameId, cardId, playerId);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to reveal card');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const becomeSpymaster = useCallback(async (gameId, playerId) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await setSpymaster(gameId, playerId);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to set spymaster');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  const finishTurn = useCallback(async (gameId, playerId) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await endTurn(gameId, playerId);
      setGameState(response);
      return response;
    } catch (err) {
      setError(err.message || 'Failed to end turn');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <GameContext.Provider value={{
      gameState,
      loading,
      error,
      startNewGame,
      joinExistingGame,
      refreshGameState,
      revealGameCard,
      becomeSpymaster,
      finishTurn
    }}>
      {children}
    </GameContext.Provider>
  );
};