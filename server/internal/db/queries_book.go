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

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/FunnySneko/ReadPill/server/internal/book"
)

type queryBuilder struct {
	index int
	query string
	args  []any
}

func (qb *queryBuilder) filterByTextMatch(columnName string, text string, isNegative bool) {
	qb.query += " AND ("
	op := "ILIKE"
	if isNegative {
		op = "NOT ILIKE"
	}
	pattern := "%" + text + "%"
	qb.query += fmt.Sprintf("%s %s $%d", columnName, op, qb.index)
	qb.args = append(qb.args, pattern)
	qb.index++
	qb.query += ")"
}

func (qb *queryBuilder) filterByValues(columnName string, values []int, isNegative bool) {
	qb.query += " AND ("
	op := "="
	if isNegative {
		op = "!="
	}
	for i, v := range values {
		if i > 0 {
			qb.query += " AND "
		}
		qb.query += fmt.Sprintf("%s %s $%d", columnName, op, qb.index)
		qb.args = append(qb.args, v)
		qb.index++
	}
	qb.query += ")"
}

func (qb *queryBuilder) filterByRanges(columnName string, ranges map[int]int, isNegative bool) {
	qb.query += " AND ("
	op := "BETWEEN"
	if isNegative {
		op = "NOT BETWEEN"
	}
	i := 0
	for k, v := range ranges {
		if i > 0 {
			qb.query += " AND "
		}
		qb.query += fmt.Sprintf("%s %s $%d AND $%d", columnName, op, qb.index, qb.index+1)
		qb.args = append(qb.args, k, v)
		i++
		qb.index += 2
	}
	qb.query += ")"
}

func (qb *queryBuilder) filterByTags(tags []int, isNegative bool) {
	op := "EXISTS"
	if isNegative {
		op = "NOT EXISTS"
	}
	for _, tag := range tags {
		qb.query += fmt.Sprintf(` AND %s (
            SELECT 1 FROM book_tag bt2
            WHERE bt2.book_id = b.id AND bt2.tag_id = $%d)`, op, qb.index)
		qb.args = append(qb.args, tag)
		qb.index++
	}
}

func (qb *queryBuilder) applyFilter(filter *book.BookFilter, isNegative bool) {
	if isNegative {
		qb.query += " AND ((1 = 1)"
	} else {
		qb.query += " OR ((1 = 1)"
	}
	if len(filter.ID) > 0 {
		qb.filterByValues("b.id", filter.ID, isNegative)
	}
	if filter.Title != nil {
		qb.filterByTextMatch("b.title", *filter.Title, isNegative)
	}
	if len(filter.AuthorID) > 0 {
		qb.filterByValues("b.author_id", filter.AuthorID, isNegative)
	}
	if filter.AuthorFullName != nil {
		qb.filterByTextMatch("w.full_name", *filter.AuthorFullName, isNegative)
	}
	if filter.Description != nil {
		qb.filterByTextMatch("b.description", *filter.Description, isNegative)
	}

	if len(filter.YearOfRelease) > 0 {
		qb.filterByRanges("b.year_of_release", filter.YearOfRelease, isNegative)
	}
	if len(filter.Tags) > 0 {
		qb.filterByTags(filter.Tags, isNegative)
	}
	qb.query += ")"
}

func (d *Db) GetBooks(ctx context.Context, filters []book.BookFilter, negativeFilters []book.BookFilter) ([]book.Book, error) {
	qb := queryBuilder{
		index: 1,
		query: `SELECT
					b.id,
					b.title,
					w.id,
					w.full_name,
					b.description,
					b.year_of_release,
					b.cover_image,
					COALESCE(
						json_agg(
							json_build_object('id', t.id, 'name', t.name)
						) FILTER (WHERE t.id IS NOT NULL),
						'[]'
					) AS tags
				FROM book b
				JOIN writer w ON w.id = b.author_id
				LEFT JOIN book_tag bt ON bt.book_id = b.id
				LEFT JOIN tag t ON t.id = bt.tag_id
				WHERE (1 = 1)`,
		args: []any{},
	}

	if len(filters) > 0 {
		qb.query += " AND ((1 = 0)"
		for _, filter := range filters {
			qb.applyFilter(&filter, false)
		}
		qb.query += ")"
	}

	if len(negativeFilters) > 0 {
		qb.query += " AND ((1 = 1)"
		for _, negativeFilter := range negativeFilters {
			qb.applyFilter(&negativeFilter, true)
		}
		qb.query += ")"
	}

	qb.query += " GROUP BY b.id, w.id"

	rows, err := d.conn.Query(ctx, qb.query, qb.args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []book.Book
	for rows.Next() {
		b := book.Book{}
		var tagsJSON []byte
		err = rows.Scan(&b.Id, &b.Title, &b.AuthorId, &b.AuthorName, &b.Description, &b.YearOfRelease, &b.CoverImageURL, &tagsJSON)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		err = json.Unmarshal(tagsJSON, &b.Tags)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		books = append(books, b)
	}

	return books, nil
}

func (d *Db) CreateBook(ctx context.Context, title string, authorId int, description string, yearOfRelease int, coverImageURL string, userId int) (int, error) {
	bookId := 0
	err := d.conn.QueryRow(ctx, "INSERT INTO book (title, author_id, description, year_of_release, cover_image, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", title, authorId, description, yearOfRelease, coverImageURL, userId).Scan(&bookId)
	return bookId, err
}
