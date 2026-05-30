package aggregation

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func NewAggregator(database *db.Db) (*Aggregator, error) {
	agg := &Aggregator{
		Db: database,
	}

	if err := agg.readReviewRules(); err != nil {
		return nil, err
	}

	return agg, nil
}

type Aggregator struct {
	Db          *db.Db
	ReviewRules review.ReviewRules
}

func (a *Aggregator) readReviewRules() error {
	data, err := os.ReadFile("config/review_rules.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &a.ReviewRules)
	return err
}

func (a *Aggregator) FormReview(ctx context.Context, reviewId int) (review.Review, error) {
	review := review.Review{}
	review, err := a.Db.GetReview(ctx, reviewId)
	if err != nil {
		return review, err
	}

	ratings, err := a.Db.GetReviewRatings(ctx, reviewId)
	if err != nil {
		return review, err
	}

	// need to convert uncompatible rating conventions here

	review.Ratings = ratings

	return review, nil
}

func (a *Aggregator) CollectBookReviews(ctx context.Context, bookId int) ([]review.Review, error) {
	reviewIDs, err := a.Db.GetBookReviewIDs(ctx, bookId)
	if err != nil {
		return nil, err
	}

	reviews := []review.Review{}
	for _, reviewId := range reviewIDs {
		review, err := a.FormReview(ctx, reviewId)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		reviews = append(reviews, review)
	}

	return reviews, err
}
