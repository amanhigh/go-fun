-- Active: 1702221591394@@docker@5432@play

-- Create Sample Student Table

CREATE TABLE students (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        age INTEGER NOT NULL,
        gender VARCHAR(255) NOT NULL
    );

-- Insert 5 Sample Records

-- Create Index on Name
CREATE INDEX idx_name ON students (name);