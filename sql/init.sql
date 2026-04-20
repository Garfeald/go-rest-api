
DROP TABLE IF EXISTS tasks;

CREATE TABLE
    tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO tasks (title, description, completed) VALUES
    ('Sample Task learn GO', 'This is a sample task', FALSE),
    ('Sample Task learn SQL', 'This is a sample task', FALSE),
    ('Sample Task learn Docker', 'This is a sample task', FALSE);