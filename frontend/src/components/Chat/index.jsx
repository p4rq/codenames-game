import React, { useState, useEffect, useContext, useRef } from 'react';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';  // Add this import

import { getMessages, sendMessage } from '../../services/chatService';
import './Chat.css';

const Chat = ({ gameId, team }) => {
  const { user } = useContext(UserContext);
  const { game } = useContext(GameContext); // Add this to get game context
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const messagesEndRef = useRef(null);
  
  // Debug log to see the user and team values
  useEffect(() => {
    console.log("Chat component user:", user);
    console.log("Chat component team:", team);
  }, [user, team]);
  
  // TEMPORARY: Force team access for debugging
  useEffect(() => {
    if (user && !user.team) {
      console.log("TEMPORARY: Forcing user team for debug");
      localStorage.setItem('user', JSON.stringify({...user, team: 'red'}));
    }
  }, [user]);
  
  // Fetch messages on mount and periodically
  useEffect(() => {
    if (!gameId) return;
    
    const fetchMessages = async () => {
      try {
        const data = await getMessages(gameId, team);
        setMessages(data || []);
        setError(null);
      } catch (err) {
        console.error('Failed to load messages:', err);
        setError('Failed to load messages');
      }
    };
    
    // Initial fetch
    fetchMessages();
    
    // Set up polling
    const interval = setInterval(fetchMessages, 3000);
    
    return () => clearInterval(interval);
  }, [gameId, team]);
  
  // Scroll to bottom when messages update
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);
  
  const handleSendMessage = async (e) => {
    e.preventDefault();
    
    if (!newMessage.trim() || !user || !gameId) return;
    
    setLoading(true);
    
    try {
      await sendMessage(newMessage, user.id, user.username, gameId, team);
      setNewMessage('');
      
      // Fetch updated messages
      const data = await getMessages(gameId, team);
      setMessages(data || []);
    } catch (err) {
      setError('Failed to send message');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Update the canAccessChat function

  const canAccessChat = () => {
    // Debug logging
    console.log("canAccessChat check:", {
      user,
      userTeam: user?.team,
      chatTeam: team,
      role: user?.role
    });
    
    if (!team) {
      // Everyone can access the general chat
      return true;
    }
    
    if (!user) {
      // No user means no access to team chats
      return false;
    }
    
    // If user doesn't have a team directly, check for team in game players
    let userTeam = user.team;
    if (!userTeam && game && game.players) {
      // Find the user in the game's players
      const playerInGame = game.players.find(p => p.id === user.id);
      if (playerInGame) {
        userTeam = playerInGame.team;
      }
    }
    
    if (!userTeam) {
      console.log("User has no team - denying access to team chat");
      return false;
    }
    
    // Convert to lowercase for case-insensitive comparison
    const normalizedUserTeam = userTeam.toLowerCase();
    const normalizedChatTeam = team.toLowerCase();
    
    // Spymasters can access all team chats
    if (user.role === 'SPYMASTER') {
      return true;
    }
    
    // Users can only access their own team's chat
    const hasAccess = normalizedUserTeam === normalizedChatTeam;
    console.log(`Team comparison: user=${normalizedUserTeam}, chat=${normalizedChatTeam}, access=${hasAccess}`);
    return hasAccess;
  };
  
  return (
    <div className={`chat-container ${team ? `team-${team.toLowerCase()}` : ''}`}>
      <h3>{team ? `${team} Team Chat` : 'Game Chat'}</h3>
      
      {!canAccessChat() ? (
        <p className="access-denied">You don't have access to this team's chat</p>
      ) : (
        <>
          <div className="messages-container">
            {messages.length === 0 ? (
              <p className="no-messages">No messages yet</p>
            ) : (
              messages.map((msg) => (
                <div 
                  key={msg.id || `${msg.sender_id}-${msg.timestamp}`} 
                  className={`message ${msg.sender_id === user?.id ? 'own-message' : ''}`}
                >
                  <span className="message-username">{msg.username}:</span>
                  <span className="message-content">{msg.content}</span>
                  <span className="message-time">
                    {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </span>
                </div>
              ))
            )}
            <div ref={messagesEndRef} />
            
            {error && <p className="error-message">{error}</p>}
          </div>
          
          <form className="message-form" onSubmit={handleSendMessage}>
            <input
              type="text"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Type a message..."
              disabled={loading || !user}
            />
            <button 
              type="submit" 
              disabled={loading || !newMessage.trim() || !user}
            >
              Send
            </button>
          </form>
        </>
      )}
    </div>
  );
};

export default Chat;