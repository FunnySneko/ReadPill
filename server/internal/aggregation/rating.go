package aggregation

import (
	"context"
	"errors"

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

var ErrNotEnoughUserData = errors.New("not enough user data")

const USER_AVERAGE_RATING_THRESHOLD = 5
const USER_OPINION_THRESHOLD = 5
const USER_BIAS_THRESHOLD = 5

func getAverage(values []float32) float32 {
	var sum float32 = 0
	count := 0
	for _, value := range values {
		sum += value
		count++
	}
	if count == 0 {
		return 0
	}
	return sum / float32(count)
}

func (a *Aggregator) CalculateContributeRating(ratings []review.Rating) float32 {
	values := []float32{}
	for _, rating := range ratings {
		if rating.Contribute {
			var value float32 = rating.Value / rating.ValueCeiling
			values = append(values, value)
		}
	}
	return getAverage(values)
}

func getAverageContributeRating(reviews []review.Review) float32 {
	contributeRatings := []float32{}
	for _, review := range reviews {
		contributeRatings = append(contributeRatings, review.ContributeRating)
	}
	return getAverage(contributeRatings)
}

func (a *Aggregator) formAverageBookContributeRating(ctx context.Context, bookId int) (float32, error) {
	reviews, err := a.Db.GetBookReviews(ctx, bookId)
	if err != nil {
		return 0, trash.WrapError("AGG FORM AVERAGE BOOK CONTRIBUTE RATING", err)
	}
	return getAverageContributeRating(reviews), nil
}

func (a *Aggregator) formAverageUserContributeRating(ctx context.Context, userId int) (float32, error) {
	reviews, err := a.Db.GetUserReviews(ctx, userId)
	if err != nil {
		return 0, trash.WrapError("AGG FORM AVERAGE USER CONTRIBUTE RATING", err)
	}
	return getAverageContributeRating(reviews), nil
}

func (a *Aggregator) FormUserOpinion(ctx context.Context, contributeRating float32, userId int) (float32, error) {
	userAverageContributeRating, err := a.formAverageUserContributeRating(ctx, userId)
	if err != nil {
		return 0, trash.WrapError("AGG FORM USER OPINION", err)
	}
	if userAverageContributeRating == 0 {
		return 0, nil
	}

	return contributeRating - userAverageContributeRating, nil
}

func (a *Aggregator) FormUserBias(ctx context.Context, contributeRating float32, bookId int) (float32, error) {
	bookAverageContributeRating, err := a.formAverageBookContributeRating(ctx, bookId)
	if err != nil {
		return 0, trash.WrapError("AGG FORM USER BIAS", err)
	}
	if bookAverageContributeRating == 0 {
		return 0, nil
	}

	return contributeRating - bookAverageContributeRating, nil
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
