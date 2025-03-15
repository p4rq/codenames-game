package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Constants for connection behavior
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a connected WebSocket client
type Client struct {
	// Unique client identifier
	ID string

	// The websocket connection
	Conn *Connection

	// Reference to the hub
	hub *Hub

	// Game ID this client is connected to
	gameID string

	// Flag to track if the client is being unregistered
	unregistering bool

	// Mutex to protect the unregistering flag
	mutex sync.Mutex
}

// Connection wraps a websocket connection
type Connection struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Mutex to protect operations on the connection
	mutex sync.Mutex

	// Flag to track if the connection is already closed
	closed bool
}

// Hub maintains the set of active clients per game
type Hub struct {
	// Game ID to client mapping
	gameClients map[string]map[*Client]bool

	// Mutex for concurrent access
	mutex sync.RWMutex

	// Channels for client registration/unregistration
	register   chan *clientRegistration
	unregister chan *Client
}

// clientRegistration holds registration data
type clientRegistration struct {
	Client *Client
	GameID string
}

// NewConnection creates a new connection
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		conn:   conn,
		send:   make(chan []byte, 256),
		closed: false,
	}
}

// NewClient creates a new client
func NewClient(id string, conn *Connection, hub *Hub, gameID string) *Client {
	return &Client{
		ID:            id,
		Conn:          conn,
		hub:           hub,
		gameID:        gameID,
		unregistering: false,
	}
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		gameClients: make(map[string]map[*Client]bool),
		register:    make(chan *clientRegistration),
		unregister:  make(chan *Client),
		mutex:       sync.RWMutex{},
	}
}

// RegisterClient registers a client with the hub
func (h *Hub) RegisterClient(client *Client, gameID string) {
	h.register <- &clientRegistration{
		Client: client,
		GameID: gameID,
	}
}

// UnregisterClient unregisters a client from the hub (exported method)
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.mutex.Lock()
		if !c.unregistering {
			c.unregistering = true
			c.mutex.Unlock()
			c.hub.unregister <- c
		} else {
			c.mutex.Unlock()
		}
	}()

	c.Conn.conn.SetReadLimit(maxMessageSize)
	c.Conn.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.conn.SetPongHandler(func(string) error {
		c.Conn.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// We're not expecting client messages in this implementation
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.Conn.send:
			c.Conn.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Conn.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Conn.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WriteMessage sends a message to the client
func (c *Connection) WriteMessage(message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return &MessageBufferFullError{}
	}

	select {
	case c.send <- message:
		return nil
	default:
		return &MessageBufferFullError{}
	}
}

// Update the Close method to avoid closing an already closed channel:
func (c *Connection) Close() {
	// Use a mutex to protect access to the send channel
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Only close if not already closed
	if !c.closed {
		close(c.send)
		c.closed = true
		c.conn.Close()
	}
}

// MessageBufferFullError is returned when the message buffer is full
type MessageBufferFullError struct{}

func (e *MessageBufferFullError) Error() string {
	return "message buffer full"
}

// Run starts the hub and handles client connections
func (h *Hub) Run() {
	for {
		select {
		case registration := <-h.register:
			// Register client to a game
			h.mutex.Lock()
			gameID := registration.GameID
			client := registration.Client

			if _, ok := h.gameClients[gameID]; !ok {
				h.gameClients[gameID] = make(map[*Client]bool)
			}
			h.gameClients[gameID][client] = true
			h.mutex.Unlock()

			log.Printf("Client %s registered for game %s", client.ID, gameID)

		case client := <-h.unregister:
			// Unregister client from all games
			h.mutex.Lock()
			clientFound := false
			for gameID, clients := range h.gameClients {
				if _, ok := clients[client]; ok {
					delete(h.gameClients[gameID], client)
					log.Printf("Client %s unregistered from game %s", client.ID, client.gameID)
					clientFound = true

					// Clean up empty game rooms
					if len(h.gameClients[gameID]) == 0 {
						delete(h.gameClients, gameID)
						log.Printf("Game room %s removed (no clients left)", gameID)
					}
				}
			}
			h.mutex.Unlock()

			// Only close the connection if the client was found and unregistered
			if clientFound {
				client.Conn.Close()
			}
		}
	}
}

// Broadcast sends a message to all clients in a specific game
func (h *Hub) Broadcast(gameID string, message []byte) {
	h.mutex.RLock()
	clients := h.gameClients[gameID]
	h.mutex.RUnlock()

	if clients == nil {
		return
	}

	log.Printf("Broadcasting to %d clients in game %s", len(clients), gameID)

	for client := range clients {
		err := client.Conn.WriteMessage(message)
		if err != nil {
			log.Printf("Error broadcasting to client %s: %v", client.ID, err)
			client.mutex.Lock()
			if !client.unregistering {
				client.unregistering = true
				client.mutex.Unlock()
				h.unregister <- client
			} else {
				client.mutex.Unlock()
			}
			client.Conn.Close()
		}
	}
}
