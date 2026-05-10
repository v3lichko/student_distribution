ALTER TABLE students
ADD COLUMN group_number INTEGER REFERENCES groups(number);