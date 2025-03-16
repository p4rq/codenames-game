import React, { useState, useEffect, useContext, useRef } from 'react';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import { getMessages, sendMessage } from '../../services/chatService';
import './Chat.css';

const Chat = ({ gameId, team }) => {
  const { user } = useContext(UserContext);
  const { game } = useContext(GameContext);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [canAccess, setCanAccess] = useState(false); // New state to track access
  const messagesEndRef = useRef(null);
  
  // Debug log to see the user and team values
  useEffect(() => {
    console.log("Chat component user:", user);
    console.log("Chat component team:", team);
  }, [user, team]);
  
  // Evaluate access when user or team changes
  useEffect(() => {
    // Calculate if user can access this chat
    const hasAccess = checkChatAccess();
    setCanAccess(hasAccess);
    
    // If access changed from false to true, fetch messages
    if (hasAccess && gameId) {
      fetchChatMessages();
    }
  }, [user, team, game]); // Re-evaluate when user, team, or game changes
  
  // Fetch messages on mount and periodically
  useEffect(() => {
    if (!gameId || !canAccess) return;
    
    fetchChatMessages();
    
    // Set up polling only if user has access
    const interval = setInterval(fetchChatMessages, 3000);
    
    return () => clearInterval(interval);
  }, [gameId, team, canAccess]); // Add canAccess as dependency
  
  const fetchChatMessages = async () => {
    try {
      const data = await getMessages(gameId, team);
      setMessages(data || []);
      setError(null);
    } catch (err) {
      console.error('Failed to load messages:', err);
      setError('Failed to load messages');
    }
  };
  
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
      fetchChatMessages();
    } catch (err) {
      setError('Failed to send message');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Update the canAccessChat function to be more reliable
  const checkChatAccess = () => {
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
    
    // Try to get user's team from different sources
    let userTeam = user.team;
    if (!userTeam && game && game.players) {
      // Find the user in the game's players
      const playerInGame = game.players.find(p => p.id === user.id);
      if (playerInGame) {
        userTeam = playerInGame.team;
      }
    }
    
    // Another fallback - check localStorage directly
    if (!userTeam) {
      try {
        const storedUser = JSON.parse(localStorage.getItem('user'));
        if (storedUser && storedUser.team) {
          userTeam = storedUser.team;
        }
      } catch (e) {
        console.error("Error parsing user from localStorage:", e);
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
      
      {!canAccess ? (
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