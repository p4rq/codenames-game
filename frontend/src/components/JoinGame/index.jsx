const handleJoinGame = async (gameId) => {
  if (!user) return;
  
  try {
    const response = await fetch('/api/game/join', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        game_id: gameId,
        player_id: user.id,
        username: user.username,
        team: selectedTeam || 'red' // Default to red team if none selected
      }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to join game');
    }
    
    const gameData = await response.json();
    
    // Update the user context with the team
    updateUser({ team: selectedTeam || 'red' });
    
    // Navigate to the game
    navigate(`/game/${gameId}`);
  } catch (error) {
    console.error('Error joining game:', error);
    setError('Failed to join game');
  }
};