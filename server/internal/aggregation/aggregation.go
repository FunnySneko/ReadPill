package aggregation

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

var ErrInvalidRatings = errors.New("invalid review")

func NewAggregator() (*Aggregator, error) {
	agg := &Aggregator{}

	if err := agg.readReviewRules(); err != nil {
		return nil, err
	}

	return agg, nil
}

type Aggregator struct {
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

func (a *Aggregator) LogRules() {
	for _, rule := range a.ReviewRules.RatingRules {
		slog.Info(rule.Name)
	}
}
