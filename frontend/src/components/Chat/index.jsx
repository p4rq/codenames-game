import React, { useState, useEffect, useContext, useRef } from 'react';
import { UserContext } from '../../context/UserContext';
import { sendMessage, getMessages } from '../../services/chatService';
import './style.css';

const Chat = ({ gameId }) => {
  const { user } = useContext(UserContext);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const messagesEndRef = useRef(null);
  
  // Fetch messages on mount and periodically
  useEffect(() => {
    if (!gameId) return;
    
    const fetchMessages = async () => {
      try {
        const data = await getMessages(gameId);
        setMessages(data || []);
        setError(null);
      } catch (err) {
        setError('Failed to load messages');
        console.error(err);
      }
    };
    
    // Initial fetch
    fetchMessages();
    
    // Set up polling
    const interval = setInterval(fetchMessages, 3000);
    
    return () => clearInterval(interval);
  }, [gameId]);
  
  // Scroll to bottom when messages update
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);
  
  const handleSendMessage = async (e) => {
    e.preventDefault();
    
    if (!newMessage.trim() || !user || !gameId) return;
    
    setLoading(true);
    
    try {
      await sendMessage(newMessage, user.id, user.username, gameId);
      setNewMessage('');
      
      // Fetch updated messages
      const data = await getMessages(gameId);
      setMessages(data || []);
    } catch (err) {
      setError('Failed to send message');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <div className="chat-container">
      <h3>Chat</h3>
      
      <div className="messages-container">
        {messages.length === 0 ? (
          <p className="no-messages">No messages yet</p>
        ) : (
          messages.map((msg) => (
            <div 
              key={msg.id} 
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
    </div>
  );
};

export default Chat;