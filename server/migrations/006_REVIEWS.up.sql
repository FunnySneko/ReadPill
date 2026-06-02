CREATE TABLE
    "review" (
        "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        "book_id" INT NOT NULL REFERENCES "book" ("id"),
        "user_id" INT NOT NULL REFERENCES "user" ("id"),
        "contribute_rating" REAL NOT NULL,
        "user_opinion" REAL NOT NULL,
        "user_opinion_confidence" REAL NOT NULL,
        "user_bias" REAL NOT NULL,
        "user_bias_confidence" REAL NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT NOW ()
    );