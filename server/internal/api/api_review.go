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

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	trash "github.com/FunnySneko/ReadPill/server/internal"
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

func (s *ServerHandler) GetReviewRules(w http.ResponseWriter, r *http.Request) {
	reviewRules := mapReviewRulesToAPI(s.Agg.ReviewRules)
	response, err := json.Marshal(reviewRules)
	if err != nil {
		ErrorOut(w, trash.WrapError("API GET REVIEW RULES", err))
		return
	}
	w.Write(response)
}

func (s *ServerHandler) PostBooksIdReviews(w http.ResponseWriter, r *http.Request, bookId int) {
	userId := int(r.Context().Value("userID").(float64))

	// NO CHECK IF THE BOOK EVEN EXISTS LOOOOOOOOOOOOL

	rev := ReviewPost{}
	err := json.NewDecoder(r.Body).Decode(&rev)
	if err != nil {
		ErrorOut(w, trash.WrapError("API POST BOOK ID REVIEWS", err))
		return
	}

	err = s.validateRatings(rev.Ratings)
	if err != nil {
		ErrorOut(w, trash.WrapError("API POST BOOK ID REVIEWS", err))
		return
	}

	ratings := []review.Rating{}
	for _, rating := range rev.Ratings {
		ratings = append(ratings, s.Agg.FormRating(rating.Name, rating.Value))
	}

	contributeRating := s.Agg.CalculateContributeRating(ratings)
	userOpinion, userOpinionConfidence, err := s.Agg.FormUserOpinion(r.Context(), contributeRating, userId)
	if err != nil {
		ErrorOut(w, trash.WrapError("API POST BOOK ID REVIEWS", err))
		return
	}
	userBias, userBiasConfidence, err := s.Agg.FormUserBias(r.Context(), contributeRating, bookId)
	if err != nil {
		ErrorOut(w, trash.WrapError("API POST BOOK ID REVIEWS", err))
		return
	}

	err = s.Ah.PostReview(r.Context(), bookId, userId, contributeRating, userOpinion, userOpinionConfidence, userBias, userBiasConfidence, ratings)
	if err != nil {
		ErrorOut(w, trash.WrapError("API POST BOOK ID REVIEWS", err))
	}
}

func (s *ServerHandler) GetBooksIdReviews(w http.ResponseWriter, r *http.Request, bookId int) {
	rev, err := s.Agg.CollectBookReviews(r.Context(), bookId)
	if err != nil {
		ErrorOut(w, trash.WrapError("API GET BOOKS ID REVIEWS", err))
		return
	}
	reviews := []Review{}
	for _, review := range rev {
		reviews = append(reviews, mapReviewToAPI(review))
	}

	response, err := json.Marshal(reviews)
	if err != nil {
		ErrorOut(w, trash.WrapError("API GET BOOKS ID REVIEWS", err))
	}
	w.Write(response)
}
