package methodsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (s *Storage) ChangeWallet(ctx context.Context, amount int, uuid string, action string) (int, error) {
	fail := func(err error) (int, error) {
		return 0, fmt.Errorf("failed to change wallet: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	var balance int
	querySelect := `
			SELECT balance FROM wallet 
			WHERE walletid=$1
			`
	err = tx.QueryRowContext(ctx, querySelect, uuid).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return fail(fmt.Errorf("failed to found wallet: %w", err))
	}
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	switch strings.ToLower(action) {
	case "deposit":

		balance = balance + amount
	case "withdraw":
		balance = balance - amount
	default:
		tx.Rollback()
		return 0, fmt.Errorf("invalid action provided: %s", action)
	}

	queryUpdate := `
			UPDATE wallet SET balance=$1 
			WHERE walletid=$2
			`
	_, err = tx.ExecContext(ctx, queryUpdate, balance, uuid)
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	return balance, tx.Commit()
}
