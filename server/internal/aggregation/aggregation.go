package aggregation

import (
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
