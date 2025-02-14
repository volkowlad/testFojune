package initdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"strings"
	"testFojune/internal/config"
)

type Storage struct {
	db *sql.DB
}

func InitDB(cfg config.DB) (*Storage, error) {
	dsn := fmt.Sprintf(`
					user=%s 
					password=%s
					host=%s
					port=%s		
					dbname=%s 
					sslmode=%s`, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetWallet(uuid uuid.UUID) (int, error) {
	stmt, err := s.db.Prepare("SELECT balance FROM wallet WHERE walletid=$1")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare query to get wallet: %w", err)
	}

	var resBalance int
	err = stmt.QueryRow(uuid).Scan(&resBalance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("wallet not found")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	return resBalance, nil
}

func (s *Storage) ChangeWallet(amount int, uuid uuid.UUID, action string) (int, error) {
	stmt, err := s.db.Prepare("SELECT balance FROM wallet WHERE walletid=$1")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare query to change wallet: %w", err)
	}

	var balance int
	err = stmt.QueryRow(uuid).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, errors.New("wallet not found")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	switch strings.ToLower(action) {
	case "deposit":
		balance = balance + amount
	case "withdraw":
		balance = balance - amount
	default:
		return 0, fmt.Errorf("invalid action provided: %s", action)
	}

	stmt, err = s.db.Prepare("UPDATE wallet SET balance=$1 WHERE walletid=$2")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare query to insert wallet: %w", err)
	}

	_, err = stmt.Exec(balance, uuid)
	if err != nil {
		return 0, fmt.Errorf("failed to %s to wallet: %w", strings.ToUpper(action), err)
	}

	return balance, nil
}

func (s *Storage) DeleteWallet(uuid uuid.UUID) error {
	stmt, err := s.db.Prepare("DELETE FROM wallet WHERE walletid=$1")
	if err != nil {
		return fmt.Errorf("failed to prepare query to delete wallet: %w", err)
	}

	_, err = stmt.Exec(uuid)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

func (s *Storage) SaveWallet(amount int) (string, error) {
	stmt, err := s.db.Prepare("INSERT INTO wallet (balance) VALUES ($1) RETURNING walletid")
	if err != nil {
		return "", fmt.Errorf("failed to prepare query to save wallet: %w", err)
	}

	var uuid string
	err = stmt.QueryRow(amount).Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("failed to save wallet: %w", err)
	}

	return uuid, nil
}
