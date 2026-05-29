package db

import "context"

func (d *Db) CreateReview(bookId int, userId int) (int, error) {
	reviewId := 0
	err := d.conn.QueryRow(context.Background(), `INSERT INTO "review" (book_id, user_id) VALUES ($1, $2) RETURNING id`, bookId, userId).Scan(&reviewId)
	return reviewId, err
}
