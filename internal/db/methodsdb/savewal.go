package methodsdb

import (
	"context"
	"fmt"
)

func (s *Storage) SaveWallet(ctx context.Context, amount int) (string, error) {
	fail := func(err error) (string, error) {
		return "", fmt.Errorf("failed to save wallet: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	var uuid string
	query := `
			INSERT INTO wallet (balance)
			VALUES %1
			RETURNING walletid
			`
	err = tx.QueryRowContext(ctx, query, amount).Scan(&uuid)
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	return uuid, tx.Commit()
}
