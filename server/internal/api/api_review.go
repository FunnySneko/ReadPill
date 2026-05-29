package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (s *ServerHandler) GetReviewRules(w http.ResponseWriter, r *http.Request) {
	reviewRules := mapReviewRulesToAPI(s.Agg.ReviewRules)
	response, err := json.Marshal(reviewRules)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (s *ServerHandler) PostBooksIdReviews(w http.ResponseWriter, r *http.Request, bookId int) {
	userId := int(r.Context().Value("userID").(float64))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	rev := ReviewPost{}
	err = json.Unmarshal([]byte(body), &rev)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
		return
	}

	ratings := []review.Rating{}
	for _, rating := range rev.Ratings {
		ratings = append(ratings, mapRatingFromAPI(rating))
	}
	err = s.Agg.ValidateRatings(ratings)
	if err != nil {
		ErrorOut(w, err, "invalid ratings", http.StatusBadRequest)
		return
	}
	err = s.Ah.PostReview(userId, bookId, ratings)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}
}
