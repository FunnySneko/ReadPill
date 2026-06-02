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
	"fmt"
	"log/slog"
	"math"

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

const USER_OPINION_RATING_COUNT_THRESHOLD = 15
const USER_OPINION_MAD_THRESHOLD = 0.1

const USER_BIAS_RATING_COUNT_THRESHOLD = 30
const USER_BIAS_MAD_THRESHOLD = 0.15

func getNonZeroAverage(values []float32) (float32, int) {
	var sum float32 = 0
	count := 0
	for _, value := range values {
		if value == 0 {
			continue
		}
		sum += value
		count++
	}
	if count == 0 {
		return 0, 0
	}
	return sum / float32(count), count
}

func getMeanAbsoluteDeviation(values []float32, average float32) float32 {
	var mad float64 = 0
	count := 0
	for _, value := range values {
		if value == 0 {
			continue
		}
		mad += math.Abs(float64(value - average))
		count++
	}
	if count == 0 {
		return 0
	}
	return float32(mad / float64(count))
}

func (a *Aggregator) CalculateContributeRating(ratings []review.Rating) float32 {
	values := []float32{}
	for _, rating := range ratings {
		if rating.Contribute {
			var value float32 = rating.Value / rating.ValueCeiling
			values = append(values, value)
		}
	}
	contributeRating, _ := getNonZeroAverage(values)
	return contributeRating
}

func getAverageContributeRating(reviews []review.Review) (float32, int, float32) {
	contributeRatings := []float32{}
	for _, review := range reviews {
		contributeRatings = append(contributeRatings, review.ContributeRating)
	}
	average, count := getNonZeroAverage(contributeRatings)
	mad := getMeanAbsoluteDeviation(contributeRatings, average)
	return average, count, mad
}

func (a *Aggregator) formAverageBookContributeRating(ctx context.Context, bookId int) (float32, int, float32, error) {
	reviews, err := a.Db.GetBookReviews(ctx, bookId)
	if err != nil {
		return 0, 0, 0, trash.WrapError("AGG FORM AVERAGE BOOK CONTRIBUTE RATING", err)
	}
	average, count, mad := getAverageContributeRating(reviews)
	return average, count, mad, nil
}

func (a *Aggregator) formAverageUserContributeRating(ctx context.Context, userId int) (float32, int, float32, error) {
	reviews, err := a.Db.GetUserReviews(ctx, userId)
	if err != nil {
		return 0, 0, 0, trash.WrapError("AGG FORM AVERAGE USER CONTRIBUTE RATING", err)
	}
	average, count, mad := getAverageContributeRating(reviews)
	return average, count, mad, nil
}

func (a *Aggregator) FormUserOpinion(ctx context.Context, contributeRating float32, userId int) (float32, float32, error) {
	userAverageContributeRating, ratingCount, mad, err := a.formAverageUserContributeRating(ctx, userId)
	if err != nil {
		return 0, 0, trash.WrapError("AGG FORM USER OPINION", err)
	}
	userOpinion := contributeRating - userAverageContributeRating

	quantityConfidence := float32(math.Min(float64(ratingCount)/USER_OPINION_RATING_COUNT_THRESHOLD, 1))
	var qualityCondidence float32
	if mad != 0 {
		qualityCondidence = float32(math.Min(float64(USER_OPINION_MAD_THRESHOLD/mad), 1))
	} else {
		quantityConfidence = 0
	}
	confidence := quantityConfidence * qualityCondidence
	slog.Info(fmt.Sprintf("RATING COUNT: %d MAD: %f QUANTITY CONFIDENCE: %f QUALITY CONFIDENCE: %f CONFIDENCE: %f", ratingCount, mad, quantityConfidence, qualityCondidence, confidence))

	return userOpinion, confidence, nil
}

func (a *Aggregator) FormUserBias(ctx context.Context, contributeRating float32, bookId int) (float32, float32, error) {
	bookAverageContributeRating, ratingCount, mad, err := a.formAverageBookContributeRating(ctx, bookId)
	if err != nil {
		return 0, 0, trash.WrapError("AGG FORM USER BIAS", err)
	}
	userBias := contributeRating - bookAverageContributeRating
	quantityConfidence := float32(math.Min(float64(ratingCount)/USER_BIAS_RATING_COUNT_THRESHOLD, 1))
	var qualityCondidence float32
	if mad != 0 {
		qualityCondidence = float32(math.Min(float64(USER_BIAS_MAD_THRESHOLD/mad), 1))
	} else {
		quantityConfidence = 0
	}
	confidence := quantityConfidence * qualityCondidence
	slog.Info(fmt.Sprintf("RATING COUNT: %d MAD: %f QUANTITY CONFIDENCE: %f QUALITY CONFIDENCE: %f CONFIDENCE: %f", ratingCount, mad, quantityConfidence, qualityCondidence, confidence))

	return userBias, confidence, nil
}
