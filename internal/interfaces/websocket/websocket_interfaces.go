package websocket

// UpdateBroadcaster defines the interface for broadcasting game updates
type UpdateBroadcaster interface {
	// BroadcastGameUpdate sends a game update to all clients in a game
	BroadcastGameUpdate(gameID string, data []byte)
}
