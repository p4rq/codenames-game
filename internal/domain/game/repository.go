package game

// Repository defines the storage operations for games and words
type GameRepository interface {
	// Game operations
	Save(game *GameState) error
	FindByID(id string) (*GameState, error)
	FindAll() ([]*GameState, error)
	Update(game *GameState) error
	Delete(id string) error

	// Word operations
	GetWords() ([]string, error)
	AddWord(word string) error
	AddWords(words []string) error
	DeleteWord(word string) error
}

// WordRepository can be used separately if you want to split the interfaces
type WordRepository interface {
	GetWords() ([]string, error)
	AddWord(word string) error
	AddWords(words []string) error
	DeleteWord(word string) error
}
