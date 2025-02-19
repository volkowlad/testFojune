package methodsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) GetWallet(ctx context.Context, uuid string) (int, error) {
	// create a helper function for preparing failure results
	fail := func(err error) (int, error) {
		return 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	var balance int
	err = tx.QueryRowContext(ctx, "SELECT balance FROM wallet WHERE walletid=$1", uuid).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return fail(fmt.Errorf("wallet not found: %w", err))
	}
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	return balance, tx.Commit()
}
