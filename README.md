# Codenames Game

Codenames is a word-based party game where players split into two teams and compete to identify their team's secret agents based on clues given by their spymaster. This project implements a web-based version of the game, allowing players to interact through a web interface.

## Project Structure

```
codenames-game
├── cmd
│   └── main.go                # Entry point of the application
├── internal
│   ├── domain
│   │   ├── game
│   │   │   ├── entity.go      # Defines the Game entity
│   │   │   └── repository.go  # GameRepository interface
│   │   └── chat
│   │       ├── entity.go      # Defines the Chat entity
│   │       └── repository.go  # ChatRepository interface
│   ├── usecase
│   │   ├── game
│   │   │   ├── service.go     # Business logic for game operations
│   │   │   └── service_test.go # Unit tests for game service
│   │   └── chat
│   │       ├── service.go     # Business logic for chat operations
│   │       └── service_test.go # Unit tests for chat service
│   ├── infrastructure
│   │   ├── persistence
│   │   │   ├── game_repository.go # Implementation of GameRepository
│   │   │   └── chat_repository.go # Implementation of ChatRepository
│   │   └── websocket
│   │       └── hub.go         # WebSocket hub for managing connections
│   └── interfaces
│       ├── api
│       │   ├── game_handler.go # HTTP handlers for game API
│       │   └── chat_handler.go # HTTP handlers for chat API
│       └── websocket
│           └── client.go      # WebSocket client management
├── pkg
│   ├── utils
│   │   └── helpers.go         # Utility functions
│   └── errors
│       └── errors.go          # Custom error types and handling
├── configs
│   └── config.go              # Application configuration settings
├── web
│   ├── templates               # HTML templates for views
│   └── static                  # Static assets (CSS, JS)
├── go.mod                      # Module definition and dependencies
├── go.sum                      # Dependency checksums
└── README.md                   # Project documentation
```

## Getting Started

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd codenames-game
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

4. **Access the game:**
   Open your web browser and navigate to `http://localhost:8080`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.