CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    age INT,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE NOT NULL
);

INSERT INTO users (age, first_name, last_name, email)
VALUES 
(25, 'Jane', 'Doe', 'jane.doe@example.com'),
(30, 'Bob', 'Johnson', 'bob.johnson@example.com'),
(22, 'Alice', 'Williams', 'alice.williams@example.com'),
(35, 'Michael', 'Brown', 'michael.brown@example.com'),
(28, 'Emily', 'Davis', 'emily.davis@example.com');

UPDATE users
SET age = 26
WHERE first_name = 'Jane' AND last_name = 'Doe';

UPDATE users
SET email = 'bob.johnson.updated@example.com'
WHERE first_name = 'Bob' AND last_name = 'Johnson';

DELETE FROM users
WHERE first_name = 'Alice' AND last_name = 'Williams';

INSERT INTO users (age, first_name, last_name, email)
VALUES (40, 'David', 'Smith', 'david.smith@example.com');

SELECT * FROM users ORDER BY age DESC;

