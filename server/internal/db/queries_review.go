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
	"fmt"

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/review"
	"github.com/jackc/pgx/v5"
)

func getReviewSelectQuery(column string, op string) string {
	return fmt.Sprintf(`SELECT id, book_id, user_id, contribute_rating, user_opinion, user_opinion_confidence, user_bias, user_bias_confidence FROM "review" WHERE %s %s $1`, column, op)
}

func readReviews(rows pgx.Rows) ([]review.Review, error) {
	defer rows.Close()
	reviews := []review.Review{}
	for rows.Next() {
		review := review.Review{}
		err := rows.Scan(&review.Id, &review.BookId, &review.UserId, &review.ContributeRating, &review.UserOpinion, &review.UserOpinionConfidence, &review.UserBias, &review.UserBiasConfidence)
		if err != nil {
			return nil, trash.WrapError("DB READ REVIEWS", err)
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (d *Db) CreateReview(ctx context.Context, bookId int, userId int, contributeRating float32, userOpinion float32, userOpinionConfidence float32, userBias float32, userBiasConfidence float32) (int, error) {
	reviewId := 0
	err := d.conn.QueryRow(ctx, `INSERT INTO "review" (book_id, user_id, contribute_rating, user_opinion, user_opinion_confidence, user_bias, user_bias_confidence) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, bookId, userId, contributeRating, userOpinion, userOpinionConfidence, userBias, userBiasConfidence).Scan(&reviewId)
	return reviewId, trash.WrapError("DB CREATE REVIEW", err)
}

func (d *Db) GetReviews(ctx context.Context, reviewIDs []int) ([]review.Review, error) {
	rows, err := d.conn.Query(ctx, getReviewSelectQuery("id", "IN"), reviewIDs)
	if err != nil {
		return nil, trash.WrapError("DB GET REVIEWS", err)
	}
	defer rows.Close()
	return readReviews(rows)
}

func (d *Db) GetBookReviews(ctx context.Context, bookId int) ([]review.Review, error) {
	rows, err := d.conn.Query(ctx, getReviewSelectQuery("book_id", "="), bookId)
	if err != nil {
		return nil, trash.WrapError("DB GET BOOK REVIEWS", err)
	}
	defer rows.Close()
	return readReviews(rows)
}

func (d *Db) GetUserReviews(ctx context.Context, userId int) ([]review.Review, error) {
	rows, err := d.conn.Query(ctx, getReviewSelectQuery("user_id", "="), userId)
	if err != nil {
		return nil, trash.WrapError("DB GET USER REVIEWS", err)
	}
	defer rows.Close()
	return readReviews(rows)
}
