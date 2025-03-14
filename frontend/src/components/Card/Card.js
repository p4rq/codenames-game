import React from 'react';
import './Card.css';

const Card = ({ card, isSpymaster, onClick }) => {
  const { word, type, revealed } = card;
  
  const getCardClass = () => {
    let classes = 'card';
    
    if (revealed) {
      classes += ` revealed ${type}`;
    } else if (isSpymaster) {
      // Show card types to spymaster
      classes += ` spymaster ${type}`;
    }
    
    return classes;
  };

  return (
    <div 
      className={getCardClass()}
      onClick={revealed ? null : onClick}
    >
      <div className="card-word">{word}</div>
    </div>
  );
};

export default Card;