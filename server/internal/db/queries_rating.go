package db

import "context"

func (d *Db) CreateRating(reviewId int, name string, value float32) error {
	_, err := d.conn.Exec(context.Background(), `INSERT INTO "rating" (review_id, name, value) VALUES ($1, $2, $3)`, reviewId, name, value)
	return err
}
