package account

import (
	"github.com/baking-bad/bcdhub/internal/models/types"
	"github.com/go-pg/pg/v10"
)

// Account -
type Account struct {
	// nolint
	tableName struct{} `pg:"accounts"`

	ID      int64
	Network types.Network     `pg:",type:SMALLINT,unique:account"`
	Type    types.AccountType `pg:",type:SMALLINT"`
	Address string            `pg:",type:varchar(36),unique:account"`
	Alias   string
}

// GetID -
func (a *Account) GetID() int64 {
	return a.ID
}

// GetIndex -
func (a *Account) GetIndex() string {
	return "accounts"
}

// Save -
func (a *Account) Save(tx pg.DBI) error {
	_, err := tx.Model(a).
		Where("network = ?", a.Network).
		Where("address = ?", a.Address).
		Returning("id").
		SelectOrInsert(a)
	return err
}

// IsEmpty -
func (a *Account) IsEmpty() bool {
	return a.Address == "" || a.Network == types.Empty || a.Type == types.AccountTypeUnknown
}
