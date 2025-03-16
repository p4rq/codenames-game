import React, { useState, useContext, useEffect } from 'react';
import { Link, useLocation, useParams } from 'react-router-dom';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './style.css';

const Navbar = ({ darkMode, toggleDarkMode }) => {
  const location = useLocation();
  const params = useParams();
  const { user, updateUsername } = useContext(UserContext);
  const { game, changeTeam } = useContext(GameContext);
  
  const [showSettings, setShowSettings] = useState(false);
  const [newUsername, setNewUsername] = useState('');
  // Force re-render when user changes
  const [currentTeam, setCurrentTeam] = useState('');
  
  // Get the gameId from the URL if we're on a game page
  const gameId = params.gameId || (location.pathname.startsWith('/game/') 
    ? location.pathname.split('/')[2] 
    : null);

  // Update username form when user changes
  useEffect(() => {
    if (user?.username) {
      setNewUsername(user.username);
    }
  }, [user?.username]);
  
  // Update current team whenever user or game changes
  useEffect(() => {
    const team = getCurrentTeam();
    console.log("Navbar - Current team detected:", team);
    setCurrentTeam(team);
  }, [user, game, location.pathname]);
  
  // Debug - log user and team info
  useEffect(() => {
    console.log("Navbar - User updated:", user);
    console.log("Navbar - Current team displayed:", currentTeam);
    console.log("Navbar - Game ID from URL:", gameId);
  }, [user, currentTeam, gameId]);
  
  const handleUsernameChange = (e) => {
    e.preventDefault();
    if (newUsername.trim()) {
      updateUsername(newUsername.trim());
      setShowSettings(false);
    }
  };
  
  const handleTeamChange = async (team) => {
    if (gameId && user) {
      console.log(`Attempting to change team to ${team} for user ${user.id} in game ${gameId}`);
      
      // Force update the current team in local state BEFORE the API call
      setCurrentTeam(team);
      
      // If you have a way to directly update the user context, do it here:
      // Example (assuming updateUser is available from UserContext):
      // updateUser({ ...user, team: team });
      
      await changeTeam(gameId, user.id, team);
      setShowSettings(false);
      
      // Force another team update after API call to be doubly sure
      setTimeout(() => {
        const freshTeam = getCurrentTeam();
        console.log("Post-API team check:", freshTeam);
        setCurrentTeam(freshTeam);
      }, 100);
    } else {
      console.log("Can't change team - missing gameId or user:", { gameId, userId: user?.id });
    }
  };
  
  // Get current team from different sources
  const getCurrentTeam = () => {
    // Try to get user's team directly from user object first
    if (user?.team) {
      console.log("Team found in user object:", user.team);
      return user.team;
    }
    
    // If not found there, try to find the player in the current game
    if (game?.players && user?.id) {
      const player = game.players.find(p => p.id === user.id);
      if (player?.team) {
        console.log("Team found in game players:", player.team);
        return player.team;
      }
    }
    
    // Check localStorage as a last resort
    try {
      const savedUser = JSON.parse(localStorage.getItem('user'));
      if (savedUser?.team) {
        console.log("Team found in localStorage:", savedUser.team);
        return savedUser.team;
      }
    } catch (e) {
      console.error("Error reading from localStorage:", e);
    }
    
    // Default fallback
    return '';
  };
  
  // Get display name for team
  const getTeamDisplay = () => {
    if (!currentTeam) return 'Observer';
    if (currentTeam === 'observer') return 'Observer';
    return currentTeam.charAt(0).toUpperCase() + currentTeam.slice(1);
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
                className={`user-menu-button ${currentTeam ? `${currentTeam}-user` : ''}`}
                onClick={() => setShowSettings(!showSettings)}
              >
                <div className="user-info">
                  <span className="username">{user.username}</span>
                  <span className="team-label">{getTeamDisplay()}</span>
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
                  
                  {/* Change team - only show if in a game */}
                  {gameId && (
                    <div className="team-selection">
                      <label>Change Team</label>
                      <div className="team-buttons">
                        <button 
                          className={`team-btn red ${currentTeam === 'red' ? 'active' : ''}`}
                          onClick={() => handleTeamChange('red')}
                        >
                          Red Team
                        </button>
                        <button 
                          className={`team-btn blue ${currentTeam === 'blue' ? 'active' : ''}`}
                          onClick={() => handleTeamChange('blue')}
                        >
                          Blue Team
                        </button>
                      </div>
                    </div>
                  )}
                  
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