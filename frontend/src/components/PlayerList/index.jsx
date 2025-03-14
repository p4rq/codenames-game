import React from 'react';
import './style.css';

const PlayerList = ({ players, currentPlayerId }) => {
  if (!players || players.length === 0) {
    return <div className="player-list-empty">No players</div>;
  }
  
  // Organize players by team
  const redTeam = players.filter(player => player.team === 'red');
  const blueTeam = players.filter(player => player.team === 'blue');
  
  return (
    <div className="player-list">
      <h3>Players</h3>
      
      <div className="team red-team">
        <h4>Red Team</h4>
        <ul>
          {redTeam.map(player => (
            <li 
              key={player.id} 
              className={`
                player 
                ${player.id === currentPlayerId ? 'current-player' : ''}
                ${player.is_spymaster ? 'spymaster' : ''}
              `}
            >
              {player.username}
              {player.is_spymaster && <span className="role-tag">Spymaster</span>}
              {player.id === currentPlayerId && <span className="you-tag">You</span>}
            </li>
          ))}
          {redTeam.length === 0 && <li className="no-players">No players</li>}
        </ul>
      </div>
      
      <div className="team blue-team">
        <h4>Blue Team</h4>
        <ul>
          {blueTeam.map(player => (
            <li 
              key={player.id} 
              className={`
                player 
                ${player.id === currentPlayerId ? 'current-player' : ''}
                ${player.is_spymaster ? 'spymaster' : ''}
              `}
            >
              {player.username}
              {player.is_spymaster && <span className="role-tag">Spymaster</span>}
              {player.id === currentPlayerId && <span className="you-tag">You</span>}
            </li>
          ))}
          {blueTeam.length === 0 && <li className="no-players">No players</li>}
        </ul>
      </div>
    </div>
  );
};

export default PlayerList;