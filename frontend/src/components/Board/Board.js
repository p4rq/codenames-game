import React from 'react';
import Card from '../Card/Card';
import './Board.css';

const Board = ({ cards, onCardClick, isSpymaster }) => {
  if (!cards || cards.length === 0) {
    return <div className="board-loading">Loading cards...</div>;
  }

  return (
    <div className="board">
      {cards.map(card => (
        <Card
          key={card.id}
          card={card}
          onClick={() => !card.revealed && onCardClick(card.id)}
          isSpymaster={isSpymaster}
        />
      ))}
    </div>
  );
};

export default Board;