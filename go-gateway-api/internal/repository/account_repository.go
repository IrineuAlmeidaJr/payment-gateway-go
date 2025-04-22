package repository

import (
	"database/sql"
	"time"

	"github.com/irineualmeidajr/imersao22/go-gateway/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Save salva uma fatura no banco de dados
func (r *AccountRepository) Save(account *domain.Account) error {
	// stantment
	stmt, err := r.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`)

	if err != nil {
		return err
	}

	// Depois que executar fecha a conexão
	defer stmt.Close()

	// O Exec, retorna a quantidade de linhas afetadas e o erro, como não
	// precisamos saber a quantidade de linhas, colocamos _ para descartar o valor.
	// É o blank identifier
	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// FindByID busca uma fatura pelo ID
func (r *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE api_key = $1;
	`, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt
	return &account, nil
}

// FindByAccountID busca todas as faturas de um determinado accountID
func (r *AccountRepository) FindByID(id string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, UpdatedAt time.Time

	err := r.db.QueryRow(`
		SELECT id, name, email, api_key, created_at, updated_at
		FROM accounts
		WHERE id = $1;
	`, id).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = UpdatedAt
	return &account, nil
}

// UpdateStatus atualiza o status de uma fatura
func (r *AccountRepository) UpdateBalance(account *domain.Account) error {
	// Aqui devemos tomar cuidado no caso de concorrência, por essa razão iremos
	// dar lock enquanto estiver atualizando o dado

	// tx: de transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	// Quando utiliza o FOR UPDATE, ninguem irá conseguir fazer uma outra atualização nesse ID
	err = tx.QueryRow(`
		SELECT balance 
		FROM accounts 
		WHERE id = $1 FOR UPDATE;
	`, account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3;
	`, account.Balance, time.Now(), account.ID)

	if err != nil {
		return err
	}

	return tx.Commit()

}
