// █▀▀▀▀▀█ ▀ █ ▄ ▀ ▄ █▀▀▀▀▀█    PROJECT: ReadPill
// █ ███ █  ▀█▄▄█  ▀ █ ███ █    AUTHOR:  FunnySneko
// █ ▀▀▀ █  ▄ ▄  █▀▀ █ ▀▀▀ █
// ▀▀▀▀▀▀▀ █▄▀ █▄▀ ▀ ▀▀▀▀▀▀▀    © 2026
// ▀█▄█▀ ▀▄▄▀▄█▀▄█▄  ▀▄▄▄▄▄▀
// ▀▀▄█▀█▀▀ ▀ ██▀██  ▄█▄█▄ █
// █▄█▀ █▀▀▄▀█▀  ▄▄█▀▀▀▄▀▄▀
// █▀▀▄▄▀▀▄▀▀▀ ▀ ▀█▀█  ████▀
// ▀ ▀  ▀▀▀█▄ ▀▀█▀ █▀▀▀█▄█▄▀
// █▀▀▀▀▀█  ██▀▄▄ ▄█ ▀ █▄▀ ▀
// █ ███ █ █▀█▄▀▀▀████▀▀▀▄█▄
// █ ▀▀▀ █ ▄▀▀█ ▀▀▄ ▀█▀█▀███
// ▀▀▀▀▀▀▀ ▀▀▀▀▀   ▀▀   ▀  ▀

package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrTagNotFound = errors.New("tag not found")

func (d *Db) FindTag(ctx context.Context, tagName string) (int, error) {
	tagId := 0
	err := d.conn.QueryRow(ctx, "SELECT id FROM tag WHERE name = $1", tagName).Scan(&tagId)
	if errors.Is(err, pgx.ErrNoRows) {
		return tagId, ErrTagNotFound
	}
	return tagId, err
}

func (d *Db) CreateTag(ctx context.Context, tagName string) (int, error) {
	tagId := 0
	err := d.conn.QueryRow(ctx, "INSERT INTO tag (name) VALUES ($1) RETURNING id", tagName).Scan(&tagId)
	return tagId, err
}

func (d *Db) CreateBookTag(ctx context.Context, bookId int, tagId int) error {
	_, err := d.conn.Exec(ctx, "INSERT INTO book_tag (book_id, tag_id) VALUES ($1, $2)", bookId, tagId)
	return err
}
