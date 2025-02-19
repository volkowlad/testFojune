package methodsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) UpdateWallet(ctx context.Context, uuid string, newBalance int) (int, error) {
	fail := func(err error) (int, error) {
		return 0, fmt.Errorf("failed to update wallet: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallet SET balance=$1 WHERE walletid=$2", newBalance, uuid)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return fail(fmt.Errorf("failed to found wallet: %w", err))
	}
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	return newBalance, tx.Commit()
}
