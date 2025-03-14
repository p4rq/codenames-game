import React from 'react';
import Card from '../Card';
import './style.css';

const GameBoard = ({ cards, isSpymaster, currentTeam, playerTeam, gameOver }) => {
  if (!cards || cards.length === 0) {
    return <div className="empty-board">No cards available</div>;
  }
  
  return (
    <div className="game-board">
      {cards.map((card) => (
        <Card
          key={card.id}
          card={card}
          isSpymaster={isSpymaster}
          currentTeam={currentTeam}
          playerTeam={playerTeam}
          gameOver={gameOver}
        />
      ))}
    </div>
  );
};

export default GameBoard;