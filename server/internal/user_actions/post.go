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

package user_actions

import (
	"context"
	"errors"
	"log/slog"

	trash "github.com/FunnySneko/ReadPill/server/internal"
	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/review"
)

func (ah *ActionsHandler) PostBookTag(ctx context.Context, bookId int, tag string) error {
	tagId, err := ah.db.FindTag(ctx, tag)
	if err != nil {
		if errors.Is(err, db.ErrTagNotFound) {
			tagId, err = ah.db.CreateTag(ctx, tag)
		} else {
			return trash.WrapError("ACT POST BOOK TAG", err)
		}
	}
	return ah.db.CreateBookTag(ctx, bookId, tagId)
}

func (ah *ActionsHandler) PostBook(ctx context.Context, title string, authorName string, description string, yearOfRelease int, tags []string, coverImageURL string, userId int) error {
	authorId, err := ah.db.FindWriter(ctx, authorName)
	if err != nil {
		if errors.Is(err, db.ErrWriterNotFound) {
			authorId, err = ah.db.CreateWriter(ctx, authorName)
			if err != nil {
				return trash.WrapError("ACT POST BOOK", err)
			}
		} else {
			return trash.WrapError("ACT POST BOOK", err)
		}
	}

	bookId, err := ah.db.CreateBook(ctx, title, authorId, description, yearOfRelease, coverImageURL, userId)
	for _, tag := range tags {
		err = ah.PostBookTag(ctx, bookId, tag)
		if err != nil {
			slog.Error(trash.WrapError("ACT POST BOOK", err).Error())
		}
	}
	return nil
}

func (ah *ActionsHandler) PostRating(ctx context.Context, reviewId int, rating review.Rating) error {
	err := ah.db.CreateRating(ctx, reviewId, rating.Name, rating.Value, rating.ValueCeiling, rating.Contribute)
	return trash.WrapError("ACT POST RATING", err)
}

func (ah *ActionsHandler) PostReview(ctx context.Context, bookId int, userId int, contributeRating float32, userOpinion float32, userOpinionConfidence float32, userBias float32, userBiasConfidence float32, ratings []review.Rating) error {
	reviewId, err := ah.db.CreateReview(ctx, bookId, userId, contributeRating, userOpinion, userOpinionConfidence, userBias, userBiasConfidence)
	if err != nil {
		return trash.WrapError("ACT POST REVIEW", err)
	}
	for _, rating := range ratings {
		err = ah.PostRating(ctx, reviewId, rating)
		if err != nil {
			slog.Error(trash.WrapError("ACT POST REVIEW", err).Error())
		}
	}
	return nil
}
