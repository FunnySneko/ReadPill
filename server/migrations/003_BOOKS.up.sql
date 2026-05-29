CREATE TABLE
    "book" (
        id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        title VARCHAR(50) NOT NULL,
        author_id INT NOT NULL REFERENCES "writer" (id),
        description TEXT,
        year_of_release INT NOT NULL,
        cover_image TEXT
    );