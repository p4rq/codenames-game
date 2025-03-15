import React, { useContext } from 'react';
import { UserContext } from '../../context/UserContext';
import { GameContext } from '../../context/GameContext';
import './TeamSelection.css';

const TeamSelection = () => {
  const { user, setUserTeam } = useContext(UserContext);
  const { gameState } = useContext(GameContext);

  if (!user || !gameState) return null;

  const handleTeamChange = async (team) => {
    try {
      await setUserTeam(gameState.id, team);
    } catch (error) {
      console.error('Failed to change team:', error);
    }
  };

  return (
    <div className="team-selection">
      <h3>Select Your Team</h3>
      <div className="team-buttons">
        <button
          className={`team-button red ${user.team === 'red' ? 'active' : ''}`}
          onClick={() => handleTeamChange('red')}
        >
          Red Team
        </button>
        <button
          className={`team-button blue ${user.team === 'blue' ? 'active' : ''}`}
          onClick={() => handleTeamChange('blue')}
        >
          Blue Team
        </button>
      </div>
    </div>
  );
};

export default TeamSelection;