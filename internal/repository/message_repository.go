package repository

import (
	"context"
	"database/sql"
	"message-processor/internal/domain"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO messages (content) VALUES ($1)", message.Content)
	return err
}

func (r *PostgresMessageRepository) MarkAsProcessed(ctx context.Context, content string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE messages SET processed = true WHERE content = $1", content)
	return err
}

func (r *PostgresMessageRepository) GetStatistics(ctx context.Context) (int, int, error) {
	var total, processed int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*), COUNT(*) FILTER (WHERE processed = true) FROM messages").Scan(&total, &processed)
	return total, processed, err
}
