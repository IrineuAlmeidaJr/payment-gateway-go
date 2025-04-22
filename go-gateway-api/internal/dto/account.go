package dto

import (
	"time"

	"github.com/irineualmeidajr/imersao22/go-gateway/internal/domain"
)

type CreatedAccountInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	APIKey    string    `json:"api_key,omitempty"` // se o campo estiver vazio ele não será incluído no JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}


func ToAccount(input CreatedAccountInput) *domain.Account {
	return domain.NewAccount(input.Name, input.Email)	
}

func FromAccount(account *domain.Account) AccountOutput {
	return AccountOutput{
		ID: account.ID,
		Name: account.Name,
		Email: account.Email,
		Balance: account.Balance,
		APIKey: account.APIKey,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}