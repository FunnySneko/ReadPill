CREATE TABLE
    "rating" (
        "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        "review_id" INT NOT NULL REFERENCES "review" ("id"),
        "name" TEXT NOT NULL,
        "value" REAL NOT NULL,
        "value_ceiling" REAL NOT NULL,
        "contribute" BOOLEAN NOT NULL
    );