package api

import (
	"github.com/FunnySneko/ReadPill/server/internal/book"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func mapBookToAPI(book book.Book) Book {
	b := Book{
		Id:            book.Id,
		Title:         book.Title,
		AuthorId:      book.AuthorId,
		AuthorName:    book.AuthorName,
		Description:   book.Description,
		YearOfRelease: book.YearOfRelease,
		CoverImage:    book.CoverImageURL,
	}
	for _, tag := range book.Tags {
		b.Tags = append(b.Tags, mapTagToAPI(tag))
	}
	return b
}

func mapTagToAPI(tag book.BookTag) Tag {
	return Tag{
		Id:   tag.ID,
		Name: tag.Name,
	}
}

func mapRatingRuleToAPI(ratingRule review.RatingRule) RatingRule {
	return RatingRule{
		Name:         ratingRule.Name,
		ValueCeiling: ratingRule.ValueCeiling,
	}
}

func mapReviewRulesToAPI(reviewRules review.ReviewRules) ReviewRules {
	r := ReviewRules{}
	for _, ratingRule := range reviewRules.RatingRules {
		r.RatingRules = append(r.RatingRules, mapRatingRuleToAPI(ratingRule))
	}
	return r
}

func mapRatingFromAPI(rating Rating) review.Rating {
	return review.Rating{
		Name:  rating.Name,
		Value: rating.Value,
	}
}
