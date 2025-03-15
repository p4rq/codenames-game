// Mock data to simulate chat messages
let mockMessages = {
  red: [],
  blue: []
};

// Generate a unique ID for messages
const generateId = () => {
  return Math.random().toString(36).substring(2) + Date.now().toString(36);
};

export const getMessages = async (gameId, team = null) => {
  console.log(`Mock: Getting messages for game ${gameId}, team: ${team}`);
  
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 300));
  
  if (team === 'red') {
    return [...mockMessages.red];
  } else if (team === 'blue') {
    return [...mockMessages.blue];
  } else {
    // Return all messages if no team specified
    return [...mockMessages.red, ...mockMessages.blue].sort(
      (a, b) => new Date(a.timestamp) - new Date(b.timestamp)
    );
  }
};

export const sendMessage = async (content, senderId, username, gameId, team = null) => {
  console.log(`Mock: Sending message to game ${gameId}, team: ${team}`);
  
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 300));
  
  const newMessage = {
    id: generateId(),
    content,
    sender_id: senderId,
    username,
    chat_id: gameId,
    team,
    timestamp: new Date().toISOString()
  };
  
  // Add message to the appropriate team's array
  if (team === 'red') {
    mockMessages.red.push(newMessage);
  } else if (team === 'blue') {
    mockMessages.blue.push(newMessage);
  } else {
    // Default to red if no team specified (shouldn't happen with your UI)
    mockMessages.red.push(newMessage);
  }
  
  return newMessage;
};

// Helper function to reset mock data (useful for testing)
export const resetMockData = () => {
  mockMessages = { red: [], blue: [] };
};