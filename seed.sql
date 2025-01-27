\c appdb

-- Create a table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "from" VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);

-- Insert sample data
INSERT INTO messages ("from", content) VALUES ('User1', 'Hello, World!');
INSERT INTO messages ("from", content) VALUES ('User2', 'Welcome to PostgreSQL!');
