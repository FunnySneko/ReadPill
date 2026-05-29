CREATE TABLE
    "review" (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        book_id INT NOT NULL REFERENCES "book" (id),
        user_id INT NOT NULL REFERENCES "user" (id)
    );