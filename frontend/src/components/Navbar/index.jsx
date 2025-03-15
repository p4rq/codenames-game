import React, { useState, useContext, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './style.css';

const Navbar = ({ darkMode, toggleDarkMode, gameId }) => {
  const { user, updateUsername } = useContext(UserContext);
  const { changeTeam, game } = useContext(GameContext);
  
  const [showSettings, setShowSettings] = useState(false);
  const [newUsername, setNewUsername] = useState(user?.username || '');
  const [teamChangeLoading, setTeamChangeLoading] = useState(false);
  const [error, setError] = useState(null);
  
  // Update username in state when user context changes
  useEffect(() => {
    if (user?.username) {
      setNewUsername(user.username);
    }
  }, [user?.username]);
  
  const handleUsernameChange = (e) => {
    e.preventDefault();
    setError(null);
    if (newUsername.trim()) {
      updateUsername(newUsername.trim());
      setShowSettings(false);
    } else {
      setError("Username cannot be empty");
    }
  };
  
  const handleTeamChange = async (team) => {
    if (gameId && user) {
      setTeamChangeLoading(true);
      setError(null);
      try {
        // For debugging, log all parameters
        console.log("Changing team with params:", {
          gameId,
          userId: user.id,
          team
        });
        
        const result = await changeTeam(gameId, user.id, team);
        if (result) {
          console.log("Team changed successfully");
          setShowSettings(false);
        }
      } catch (err) {
        console.error("Error changing team:", err);
        setError("Failed to change team. Try again later.");
      } finally {
        setTeamChangeLoading(false);
      }
    } else {
      console.error("Cannot change team - missing game ID or user", { gameId, user });
      setError("Cannot change team - game or user information is missing");
    }
  };
  
  // Determine text color class based on dark mode
  const textColorClass = darkMode ? 'text-light' : 'text-dark';
  
  return (
    <nav className={`navbar ${darkMode ? 'dark' : 'light'}`}>
      <div className="navbar-content">
        <div className="navbar-brand">
          <Link to="/" className={textColorClass}>
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
                className={`user-menu-button ${textColorClass}`}
                onClick={() => setShowSettings(!showSettings)}
              >
                <span className="username">{user.username}</span>
                {user.team && (
                  <span className={`team-indicator ${user.team}`}></span>
                )}
                <span className="dropdown-icon">▼</span>
              </button>
              
              {showSettings && (
                <div className="settings-dropdown">
                  {error && <div className="error-message">{error}</div>}
                  
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
                          disabled={teamChangeLoading || user?.team === 'red'}
                        >
                          {teamChangeLoading && user?.team !== 'red' ? 'Loading...' : 'Red Team'}
                        </button>
                        <button 
                          className={`team-btn blue ${user?.team === 'blue' ? 'active' : ''}`}
                          onClick={() => handleTeamChange('blue')}
                          disabled={teamChangeLoading || user?.team === 'blue'}
                        >
                          {teamChangeLoading && user?.team !== 'blue' ? 'Loading...' : 'Blue Team'}
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