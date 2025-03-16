import React, { useState, useContext } from 'react';
import { Link } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './style.css';

const Navbar = ({ darkMode, toggleDarkMode, gameId }) => {
  const { user, updateUsername } = useContext(UserContext);
  const { changeTeam } = useContext(GameContext);
  
  const [showSettings, setShowSettings] = useState(false);
  const [newUsername, setNewUsername] = useState(user?.username || '');
  
  const handleUsernameChange = (e) => {
    e.preventDefault();
    if (newUsername.trim()) {
      updateUsername(newUsername.trim());
      setShowSettings(false);
    }
  };
  
  const handleTeamChange = async (team) => {
    if (gameId && user) {
      await changeTeam(gameId, user.id, team);
      setShowSettings(false);
    }
  };
  
  // Get display name for team
  const getTeamDisplay = () => {
    if (!user || !user.team) return 'Observer';
    return user.team.charAt(0).toUpperCase() + user.team.slice(1);
  };
  
  return (
    <nav className={`navbar ${darkMode ? 'dark' : ''}`}>
      <div className="navbar-content">
        <div className="navbar-brand">
          <Link to="/">
            <h1>Codenames</h1>
          </Link>
        </div>
        
        <div className="navbar-actions">
          <button 
            className="theme-toggle" 
            onClick={toggleDarkMode} 
            aria-label={darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
          >
            {darkMode ? '☀️' : '🌙'}
          </button>
          
          {user && (
            <div className="user-menu">
              <button 
                className={`user-menu-button ${user.team || 'observer'}-user`}
                onClick={() => setShowSettings(!showSettings)}
              >
                <div className="user-info">
                  <span className="username">{user.username}</span>
                  {user.team && (
                    <span className="team-label">{getTeamDisplay()}</span>
                  )}
                </div>
                <span className="dropdown-icon">▼</span>
              </button>
              
              {showSettings && (
                <div className="settings-dropdown">
                  {/* Change username */}
                  <form onSubmit={handleUsernameChange} className="settings-form">
                    <label htmlFor="username">Change Name</label>
                    <div className="input-group">
                      <input
                        type="text"
                        id="username"
                        value={newUsername}
                        onChange={(e) => setNewUsername(e.target.value)}
                        placeholder="New username"
                      />
                      <button type="submit">Save</button>
                    </div>
                  </form>
                  
                  {/* Change team */}
                  {gameId && (
                    <div className="team-selection">
                      <label>Change Team</label>
                      <div className="team-buttons">
                        <button 
                          className={`team-btn red ${user?.team === 'red' ? 'active' : ''}`}
                          onClick={() => handleTeamChange('red')}
                        >
                          Red Team
                        </button>
                        <button 
                          className={`team-btn blue ${user?.team === 'blue' ? 'active' : ''}`}
                          onClick={() => handleTeamChange('blue')}
                        >
                          Blue Team
                        </button>
                      </div>
                    </div>
                  )}
                  
                  {/* Close button */}
                  <button 
                    className="close-settings"
                    onClick={() => setShowSettings(false)}
                  >
                    Close
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;