package db

import (
	"context"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (d *Db) CreateRating(reviewId int, name string, value float32) error {
	_, err := d.conn.Exec(context.Background(), `INSERT INTO "rating" (review_id, name, value) VALUES ($1, $2, $3)`, reviewId, name, value)
	return err
}

func (d *Db) GetReviewRatings(reviewId int) ([]review.Rating, error) {
	rows, err := d.conn.Query(context.Background(), `SELECT name, value FROM "rating" WHERE review_id = $1`, reviewId)
	if err != nil {
		return nil, err
	}
	ratings := []review.Rating{}
	for rows.Next() {
		rating := review.Rating{}
		err = rows.Scan(&rating.Name, &rating.Value)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}
