package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex // Serve para travar a transão e não conseguir mudar de valor
	CreatedAt time.Time
	UpdatedAt time.Time
}

// No GO não temos uma método contrutor, mas, por convensão cria um
// função construtora, para criar uma novo objeto
func NewAccount(name, email string) *Account {
	account := &Account{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		APIKey:    generateAPIKey(),
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account
}

func (a *Account) AddBalance(amount float64) {
	a.mu.Lock()
	// posso utilizar o defer que é como finaly que ao fim do método
	// irá realizar o unlock, sem ter que colocar na mão. Defer espera tudo executar
	// e só executa ao fim, isso impede "racing condition" (condição de corrida)
	defer a.mu.Unlock()

	a.Balance += amount
	a.UpdatedAt = time.Now()
	// a.mu.Unlock()
}

func generateAPIKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
