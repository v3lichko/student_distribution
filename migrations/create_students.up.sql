CREATE TABLE students (
    isu INTEGER PRIMARY KEY,
    full_name TEXT NOT NULL,
    telegram TEXT NOT NULL UNIQUE,
    score INTEGER NOT NULL CHECK (score >= 0)
);