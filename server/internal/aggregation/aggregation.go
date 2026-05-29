package aggregation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

var ErrInvalidRatings = errors.New("invalid review")

func NewAggregator(database *db.Db) (*Aggregator, error) {
	agg := &Aggregator{}

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

func (a *Aggregator) ValidateRatings(ratings []review.Rating) error {
	for _, rating := range ratings {
		invalid := true
		for _, ratingRule := range a.ReviewRules.RatingRules {
			if rating.Name == ratingRule.Name && rating.Value <= ratingRule.ValueCeiling {
				invalid = false
				break
			}
		}
		if invalid {
			return ErrInvalidRatings
		}
	}
	return nil
}

func (a *Aggregator) FormReview(reviewId int) (review.Review, error) {
	review := review.Review{}
	review, err := a.Db.GetReview(reviewId)
	if err != nil {
		return review, err
	}

	ratings, err := a.Db.GetReviewRatings(reviewId)
	if err != nil {
		return review, err
	}

	err = a.ValidateRatings(ratings)
	if err != nil {

	}

	return review, nil
}

func (a *Aggregator) CollectBookReviews(bookId int) ([]review.Review, error) {
	reviewIDs, err := a.Db.GetBookReviewIDs(bookId)
	if err != nil {
		return nil, err
	}

	reviews := []review.Review{}
	for _, reviewId := range reviewIDs {
		review, err := a.FormReview(reviewId)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		reviews = append(reviews, review)
	}

	return reviews, err
}
