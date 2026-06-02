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

package aggregation

import (
	"context"
	"encoding/json"
	"os"

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func NewAggregator(database *db.Db) (*Aggregator, error) {
	agg := &Aggregator{
		Db: database,
	}

	if err := agg.readReviewRules(); err != nil {
		return nil, trash.WrapError("AGG INIT AGGREGATOR", err)
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
		return trash.WrapError("AGG READ REVIEW RULES", err)
	}
	err = json.Unmarshal(data, &a.ReviewRules)
	return trash.WrapError("AGG READ REVIEW RULES", err)
}

func (a *Aggregator) FormRating(name string, value float32) review.Rating {
	ratingRule := review.RatingRule{}
	for _, rule := range a.ReviewRules.RatingRules {
		if name == rule.Name {
			ratingRule = rule
			break
		}
	}
	return review.Rating{
		Name:         name,
		Value:        value,
		ValueCeiling: ratingRule.ValueCeiling,
		Contribute:   ratingRule.Contribute,
	}
}

func (a *Aggregator) CollectBookReviews(ctx context.Context, bookId int) ([]review.Review, error) {
	reviews, err := a.Db.GetBookReviews(ctx, bookId)
	if err != nil {
		return nil, trash.WrapError("AGG COLLECT BOOK REVIEWS", err)
	}
	r := []review.Review{}
	for _, review := range reviews {
		ratings, err := a.Db.GetReviewRatings(ctx, review.Id)
		if err != nil {
			return nil, trash.WrapError("AGG COLLECT BOOK REVIEWS", err)
		}
		review.Ratings = ratings
		r = append(r, review)
	}
	return r, nil
}
