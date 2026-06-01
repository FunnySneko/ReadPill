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

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/review"
	"github.com/jackc/pgx/v5"
)

func readRatings(rows pgx.Rows) ([]review.Rating, error) {
	defer rows.Close()
	ratings := []review.Rating{}
	for rows.Next() {
		rating := review.Rating{}
		err := rows.Scan(&rating.Name, &rating.Value, &rating.ValueCeiling, &rating.Contribute)
		if err != nil {
			return nil, trash.WrapError("DB READ RATINGS", err)
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

func (d *Db) CreateRating(ctx context.Context, reviewId int, name string, value float32, valueCeiling float32, contribute bool) error {
	_, err := d.conn.Exec(ctx, `INSERT INTO "rating" (review_id, name, value, value_ceiling, contribute) VALUES ($1, $2, $3, $4, $5)`, reviewId, name, value, valueCeiling, contribute)
	return trash.WrapError("DB CREATE RATING", err)
}

func (d *Db) GetReviewRatings(ctx context.Context, reviewId int) ([]review.Rating, error) {
	rows, err := d.conn.Query(ctx, `SELECT (name, value, value_ceiling, contribute) FROM "rating" WHERE review_id = $1`, reviewId)
	if err != nil {
		return nil, trash.WrapError("DB GET REVIEW RATING", err)
	}
	defer rows.Close()

	return readRatings(rows)
}
