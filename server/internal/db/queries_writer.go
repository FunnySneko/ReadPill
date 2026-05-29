package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrWriterNotFound = errors.New("writer not found")

func (d *Db) FindWriter(writerName string) (int, error) {
	writerId := 0
	err := d.conn.QueryRow(context.Background(), "SELECT id FROM writer WHERE full_name = $1", writerName).Scan(&writerId)
	if errors.Is(err, pgx.ErrNoRows) {
		return writerId, ErrWriterNotFound
	}
	return writerId, err
}

func (d *Db) CreateWriter(writerName string) (int, error) {
	writerId := 0
	err := d.conn.QueryRow(context.Background(), "INSERT INTO writer (full_name) VALUES ($1) RETURNING id", writerName).Scan(&writerId)
	return writerId, err
}
