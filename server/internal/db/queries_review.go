package db

import (
	"context"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (d *Db) CreateReview(ctx context.Context, bookId int, userId int) (int, error) {
	reviewId := 0
	err := d.conn.QueryRow(ctx, `INSERT INTO "review" (book_id, user_id) VALUES ($1, $2) RETURNING id`, bookId, userId).Scan(&reviewId)
	return reviewId, err
}

func (d *Db) GetBookReviewIDs(ctx context.Context, bookId int) ([]int, error) {
	rows, err := d.conn.Query(ctx, `SELECT id from "review" WHERE book_id = $1`, bookId)
	if err != nil {
		return nil, err
	}
	reviewIDs := []int{}
	for rows.Next() {
		reviewId := 0
		err = rows.Scan(&reviewId)
		if err != nil {
			return nil, err
		}
		reviewIDs = append(reviewIDs, reviewId)
	}
	return reviewIDs, nil
}

func (d *Db) GetReview(ctx context.Context, reviewId int) (review.Review, error) {
	review := review.Review{}
	err := d.conn.QueryRow(ctx, `SELECT book_id, user_id FROM "review" WHERE id = $1`, reviewId).Scan(&review.BookId, &review.UserId)
	return review, err
}
