import React, { useContext } from 'react';
import { GameContext } from '../../context/GameContext';
import { UserContext } from '../../context/UserContext';
import './style.css';

const Card = ({ card, isSpymaster, currentTeam, playerTeam, gameOver }) => {
  const { revealGameCard } = useContext(GameContext);
  const { user } = useContext(UserContext);
  
  const handleCardClick = () => {
    // Only allow revealing if:
    // 1. Card is not revealed yet
    // 2. Game is not over
    // 3. It's the player's team's turn
    // 4. Player is not a spymaster
    if (
      !card.revealed &&
      !gameOver &&
      currentTeam === playerTeam &&
      !isSpymaster &&
      user
    ) {
      revealGameCard(card.game_id, card.id, user.id);
    }
  };
  
  // Determine CSS classes for the card
  const cardClasses = ['game-card'];
  
  if (card.revealed) {
    cardClasses.push('revealed');
    cardClasses.push(`${card.type}-card`);
  } else if (isSpymaster) {
    // Show the card type to spymaster even if not revealed
    cardClasses.push(`${card.type}-card-spy`);
  }
  
  // Determine if card is clickable
  const isClickable = 
    !card.revealed && 
    !gameOver && 
    currentTeam === playerTeam && 
    !isSpymaster;
  
  if (isClickable) {
    cardClasses.push('clickable');
  }
  
  return (
    <div 
      className={cardClasses.join(' ')} 
      onClick={handleCardClick}
    >
      <span className="card-word">{card.word}</span>
    </div>
  );
};

export default Card;