package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"codenames-game/internal/domain/game"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresRepository implements Repository with PostgreSQL storage
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new repository with PostgreSQL
func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

// Initialize tables if they don't exist
func initTables(db *sql.DB) error {
	// Create games table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS games (
            id TEXT PRIMARY KEY,
            data JSONB NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL
        )
    `)
	if err != nil {
		return err
	}

	// Create words table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS words (
            id SERIAL PRIMARY KEY,
            word TEXT UNIQUE NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            active BOOLEAN NOT NULL DEFAULT true
        )
    `)
	if err != nil {
		return err
	}

	// Check if words table is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM words").Scan(&count)
	if err != nil {
		return err
	}

	// If no words exist, add default words
	if count == 0 {
		defaultWords := []string{
			"AFRICA", "AGENT", "AIR", "ALIEN", "ALPS", "AMAZON", "AMBULANCE", "AMERICA", "ANGEL",
			"ANTARCTICA", "APPLE", "ARM", "ATLANTIS", "AUSTRALIA", "AZTEC", "BACK", "BALL", "BAND",
			"BANK", "BAR", "BARK", "BAT", "BATTERY", "BEACH", "BEAR", "BEAT", "BED", "BEIJING",
			"BELL", "BELT", "BERLIN", "BERMUDA", "BERRY", "BILL", "BLOCK", "BOARD", "BOLT", "BOMB",
			"BOND", "BOOM", "BOOT", "BOTTLE", "BOW", "BOX", "BRIDGE", "BRUSH", "BUCK", "BUFFALO",
			"BUG", "BUGLE", "BUTTON", "CALF", "CANADA", "CAP", "CAPITAL", "CAR", "CARD", "CARROT",
			"CASINO", "CAST", "CAT", "CELL", "CENTAUR", "CENTER", "CHAIR", "CHANGE", "CHARGE", "CHECK",
			// More words...
		}

		// Insert default words in a transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		for _, word := range defaultWords {
			_, err := tx.Exec("INSERT INTO words (word, created_at) VALUES ($1, NOW())", word)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	}

	return nil
}

// Save stores a game in the database
func (r *PostgresRepository) Save(gameState *game.GameState) error {
	// Convert game state to JSON
	data, err := json.Marshal(gameState)
	if err != nil {
		return err
	}

	// Insert into database
	_, err = r.db.Exec(
		"INSERT INTO games (id, data, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		gameState.ID, data, gameState.CreatedAt, gameState.UpdatedAt,
	)
	return err
}

// FindByID retrieves a game from the database by ID
func (r *PostgresRepository) FindByID(id string) (*game.GameState, error) {
	var jsonData []byte
	err := r.db.QueryRow("SELECT data FROM games WHERE id = $1", id).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("game not found")
		}
		return nil, err
	}

	// Convert JSON to game state
	var gameState game.GameState
	if err := json.Unmarshal(jsonData, &gameState); err != nil {
		return nil, err
	}

	return &gameState, nil
}

// FindAll retrieves all games from the database
func (r *PostgresRepository) FindAll() ([]*game.GameState, error) {
	rows, err := r.db.Query("SELECT data FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*game.GameState
	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			return nil, err
		}

		var gameState game.GameState
		if err := json.Unmarshal(jsonData, &gameState); err != nil {
			return nil, err
		}

		games = append(games, &gameState)
	}

	return games, nil
}

// Update modifies a game in the database
func (r *PostgresRepository) Update(gameState *game.GameState) error {
	// Convert game state to JSON
	data, err := json.Marshal(gameState)
	if err != nil {
		return err
	}

	gameState.UpdatedAt = time.Now()

	// Update in database
	result, err := r.db.Exec(
		"UPDATE games SET data = $1, updated_at = $2 WHERE id = $3",
		data, gameState.UpdatedAt, gameState.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("game not found")
	}

	return nil
}

// Delete removes a game from the database
func (r *PostgresRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM games WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("game not found")
	}

	return nil
}

// GetWords retrieves all active words from the database
func (r *PostgresRepository) GetWords() ([]string, error) {
	rows, err := r.db.Query("SELECT word FROM words WHERE active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// AddWord adds a word to the database
func (r *PostgresRepository) AddWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))
	if word == "" {
		return errors.New("word cannot be empty")
	}

	_, err := r.db.Exec(
		"INSERT INTO words (word, created_at) VALUES ($1, NOW()) ON CONFLICT (word) DO UPDATE SET active = true",
		word,
	)
	return err
}

// AddWords adds multiple words to the database
func (r *PostgresRepository) AddWords(words []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, word := range words {
		word = strings.TrimSpace(strings.ToUpper(word))
		if word == "" {
			continue
		}

		_, err := tx.Exec(
			"INSERT INTO words (word, created_at) VALUES ($1, NOW()) ON CONFLICT (word) DO UPDATE SET active = true",
			word,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeleteWord deactivates a word in the database
func (r *PostgresRepository) DeleteWord(word string) error {
	word = strings.TrimSpace(strings.ToUpper(word))

	_, err := r.db.Exec("UPDATE words SET active = false WHERE word = $1", word)
	return err
}
