package db

import (
	"context"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (d *Db) CreateRating(ctx context.Context, reviewId int, name string, value float32, valueCeiling float32, contribute bool) error {
	_, err := d.conn.Exec(ctx, `INSERT INTO "rating" (review_id, name, value, value_ceiling, contribute) VALUES ($1, $2, $3, $4, $5)`, reviewId, name, value, valueCeiling, contribute)
	return err
}

func (d *Db) GetReviewRatings(ctx context.Context, reviewId int) ([]review.Rating, error) {
	rows, err := d.conn.Query(ctx, `SELECT name, value FROM "rating" WHERE review_id = $1`, reviewId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
