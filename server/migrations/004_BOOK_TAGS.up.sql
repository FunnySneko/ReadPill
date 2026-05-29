CREATE TABLE
    "book_tag" (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        book_id INT NOT NULL REFERENCES "book" (id),
        tag_id INT NOT NULL REFERENCES "tag" (id)
    );