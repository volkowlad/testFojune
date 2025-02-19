package methodsdb

import (
	"context"
	"fmt"
)

func (s *Storage) DeleteWallet(ctx context.Context, uuid string) error {
	fail := func(err error) error {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM wallet WHERE walletid=$1", uuid)
	if err != nil {
		tx.Rollback()
		return fail(err)
	}

	return tx.Commit()
}
