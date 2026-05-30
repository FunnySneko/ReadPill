CREATE TABLE
    "user" (
        "id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        "username" TEXT NOT NULL,
        "email" TEXT UNIQUE NOT NULL,
        "password_hash" TEXT NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT NOW ()
    );