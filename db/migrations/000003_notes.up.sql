CREATE TABLE IF NOT EXISTS notes (
                                     id SERIAL PRIMARY KEY,
                                     user_name TEXT NOT NULL,
                                     title TEXT NOT NULL,
                                     content TEXT,
                                     metadata TEXT
);