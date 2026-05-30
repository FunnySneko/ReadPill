package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

var ErrInvalidRatings = errors.New("invalid review")

func (s *ServerHandler) validateRatings(ratings []Rating) error {
	for _, rating := range ratings {
		invalid := true
		for _, ratingRule := range s.Agg.ReviewRules.RatingRules {
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

func (s *ServerHandler) writeFromRuleToRating(rating review.Rating) review.Rating {
	r := rating
	ratingRule := review.RatingRule{}
	for _, rule := range s.Agg.ReviewRules.RatingRules {
		if rating.Name == rule.Name {
			ratingRule = rule
			break
		}
	}
	r.Contribute = ratingRule.Contribute
	r.ValueCeiling = ratingRule.ValueCeiling
	return r
}

func (s *ServerHandler) GetReviewRules(w http.ResponseWriter, r *http.Request) {
	reviewRules := mapReviewRulesToAPI(s.Agg.ReviewRules)
	response, err := json.Marshal(reviewRules)
	if err != nil {
		ErrorOut(w, err)
		return
	}
	w.Write(response)
}

func (s *ServerHandler) PostBooksIdReviews(w http.ResponseWriter, r *http.Request, bookId int) {
	userId := int(r.Context().Value("userID").(float64))

	rev := ReviewPost{}
	err := json.NewDecoder(r.Body).Decode(&rev)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	err = s.validateRatings(rev.Ratings)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	ratings := []review.Rating{}
	for _, rating := range rev.Ratings {
		ratings = append(ratings, s.writeFromRuleToRating(mapRatingFromAPI(rating)))
	}

	err = s.Ah.PostReview(r.Context(), userId, bookId, ratings)
	if err != nil {
		ErrorOut(w, err)
	}
}

func (s *ServerHandler) GetBooksIdReviews(w http.ResponseWriter, r *http.Request, bookId int) {
	rev, err := s.Agg.CollectBookReviews(r.Context(), bookId)
	if err != nil {
		ErrorOut(w, err)
		return
	}
	reviews := []Review{}
	for _, review := range rev {
		reviews = append(reviews, mapReviewToAPI(review))
	}

	response, err := json.Marshal(reviews)
	if err != nil {
		ErrorOut(w, err)
	}
	w.Write(response)
}
