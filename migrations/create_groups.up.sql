CREATE TABLE groups(
    number INTEGER PRIMARY KEY,
    capacity INTEGER NOT NULL CHECK (capacity > 0)
);