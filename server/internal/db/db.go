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
	"os"

	"github.com/jackc/pgx/v5"
)

type Db struct {
	conn *pgx.Conn
}

func (d *Db) ConnectToDatabase() error {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")

	dbURL := "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName
	var err error
	d.conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) CloseConnection() {
	d.conn.Close(context.Background())
}
