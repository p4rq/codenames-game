import React, { useState, useEffect, useRef } from 'react';
import api from '../../services/api';
import './Chat.css';

const Chat = ({ gameId, username }) => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [error, setError] = useState('');
  const chatContainerRef = useRef(null);

  const fetchMessages = async () => {
    try {
      const data = await api.getMessages(gameId);
      setMessages(data);
      setError('');
    } catch (err) {
      setError('Error loading messages');
      console.error(err);
    }
  };

  useEffect(() => {
    fetchMessages();
    
    // Set up polling for new messages
    const intervalId = setInterval(fetchMessages, 3000);
    
    return () => clearInterval(intervalId);
  }, [gameId]);

  useEffect(() => {
    // Scroll to bottom of chat when new messages arrive
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;
    
    try {
      await api.sendMessage(gameId, username, newMessage);
      setNewMessage('');
      await fetchMessages();
    } catch (err) {
      setError('Failed to send message');
      console.error(err);
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-header">
        <h3>Game Chat</h3>
      </div>
      
      <div className="chat-messages" ref={chatContainerRef}>
        {messages.length === 0 ? (
          <p className="no-messages">No messages yet</p>
        ) : (
          messages.map((message, index) => (
            <div 
              key={index} 
              className={`message ${api.getUserId() === message.sender_id ? 'my-message' : ''}`}
            >
              <div className="message-header">
                <span className="message-username">{message.username}</span>
                <span className="message-time">
                  {new Date(message.timestamp).toLocaleTimeString()}
                </span>
              </div>
              <div className="message-content">{message.content}</div>
            </div>
          ))
        )}
      </div>
      
      <form className="chat-input" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type a message..."
          disabled={!username}
        />
        <button type="submit" disabled={!newMessage.trim() || !username}>
          Send
        </button>
      </form>
      
      {error && <div className="chat-error">{error}</div>}
    </div>
  );
};

export default Chat;