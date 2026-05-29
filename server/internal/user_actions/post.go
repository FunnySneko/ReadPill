package user_actions

import (
	"errors"
	"log/slog"

	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (ah *ActionsHandler) PostBookTag(bookId int, tag string) error {
	tagId, err := ah.db.FindTag(tag)
	if err != nil {
		if errors.Is(err, db.ErrTagNotFound) {
			tagId, err = ah.db.CreateTag(tag)
		} else {
			return err
		}
	}
	return ah.db.CreateBookTag(bookId, tagId)
}

func (ah *ActionsHandler) PostBook(title string, authorName string, description string, yearOfRelease int, tags []string, coverImageURL string) error {
	authorId, err := ah.db.FindWriter(authorName)
	if err != nil {
		if errors.Is(err, db.ErrWriterNotFound) {
			authorId, err = ah.db.CreateWriter(authorName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	bookId, err := ah.db.CreateBook(title, authorId, description, yearOfRelease, coverImageURL)
	for _, tag := range tags {
		err = ah.PostBookTag(bookId, tag)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	return nil
}

func (ah *ActionsHandler) PostRating(reviewId int, rating review.Rating) error {
	err := ah.db.CreateRating(reviewId, rating.Name, rating.Value)
	return err
}

func (ah *ActionsHandler) PostReview(userId int, bookId int, ratings []review.Rating) error {
	reviewId, err := ah.db.CreateReview(bookId, userId)
	if err != nil {
		return err
	}
	for _, rating := range ratings {
		err = ah.PostRating(reviewId, rating)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	return nil
}
