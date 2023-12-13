-- Active: 1702221591394@@docker@5432@play

-- Create Sample Student Table

CREATE TABLE students (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        age INTEGER NOT NULL,
        gender VARCHAR(255) NOT NULL
    );

-- Insert 5 Sample Records

INSERT INTO students (name, age, gender)
VALUES
('John', 20, 'Male'),
('Jane', 21, 'Female'),
('Bob', 22, 'Male'),
('Sarah', 23, 'Female'),
('Mike', 24, 'Male');

-- Create Index on Name
CREATE INDEX idx_name ON students (name);